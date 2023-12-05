package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/component"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"github.com/stellar-payment/sp-payment/internal/util/cryptoutil"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/internal/util/namegen"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/internal/util/structutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) GetAllAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.ListAccountResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER); !ok {
		return nil, errs.ErrNoAccess
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 100
	}

	repoParams := &indto.AccountParams{
		Keyword:     params.Keyword,
		AccountType: params.AccountType,
		Limit:       params.Limit,
		Page:        params.Page,
	}

	res = &dto.ListAccountResponse{
		Accounts: []*dto.AccountResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountAccounts(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindAccounts(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.AccountResponse{
			ID:          v.ID,
			OwnerID:     v.OwnerID,
			AccountType: v.AccountType,
			AccountNo:   cryptoutil.DecryptField(v.AccountNo, conf.DBKey),
		}

		if v.OwnerName != nil {
			if v.AccountType == inconst.ACCOUNT_TYPE_CUST {
				temp.OwnerName = cryptoutil.DecryptField(v.OwnerName, conf.DBKey)
			} else {
				temp.OwnerName = string(v.OwnerName)
			}
		}

		hash := v.AccountNo
		if len(v.RowHash) != 0 {
			if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, v.RowHash) {
				logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "accounts")).Send()
			}
		} else {
			logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "accounts")).Str("merchant-id", v.ID).Msg("row hash not found")
		}

		res.Accounts = append(res.Accounts, temp)
	}

	return
}

func (s *service) GetAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.AccountResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	repoParams := &indto.AccountParams{AccountID: params.AccountID}

	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_CUSTOMER || usrmeta.RoleID == inconst.ROLE_MERCHANT {
		repoParams.UserID = usrmeta.UserID
	}

	data, err := s.repository.FindAccount(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.AccountResponse{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		AccountType: data.AccountType,
		Balance:     data.Balance,
		AccountNo:   cryptoutil.DecryptField(data.AccountNo, conf.DBKey),
	}

	if data.OwnerName != nil {
		if data.AccountType == inconst.ACCOUNT_TYPE_CUST {
			res.OwnerName = cryptoutil.DecryptField(data.OwnerName, conf.DBKey)
		} else {
			res.OwnerName = string(data.OwnerName)
		}
	}

	hash := data.AccountNo
	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "accounts")
			logger.Error().Err(err).Send()
			return
		}
	} else {
		err = errs.New(errs.ErrDataIntegrity, "merchant")
		logger.Error().Err(err).Str("merchant-id", data.ID).Msg("row hash not found")
		return
	}

	return
}

func (s *service) GetAccountByNo(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.AccountResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	repoParams := &indto.AccountParams{AccountNoHash: cryptoutil.HMACSHA512([]byte(params.AccountNo), conf.HashKey)}

	data, err := s.repository.FindAccount(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.AccountResponse{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		AccountType: data.AccountType,
		AccountNo:   cryptoutil.DecryptField(data.AccountNo, conf.DBKey),
	}

	if data.OwnerName != nil {
		if data.AccountType == inconst.ACCOUNT_TYPE_CUST {
			res.OwnerName = cryptoutil.DecryptField(data.OwnerName, conf.DBKey)
		} else {
			res.OwnerName = string(data.OwnerName)
		}
	}

	hash := data.AccountNo
	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "accounts")
			logger.Error().Err(err).Send()
			return
		}
	} else {
		err = errs.New(errs.ErrDataIntegrity, "merchant")
		logger.Error().Err(err).Str("merchant-id", data.ID).Msg("row hash not found")
		return
	}

	return
}

func (s *service) GetAccountMe(ctx context.Context) (res *dto.AccountResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	repoParams := &indto.AccountParams{}

	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_CUSTOMER || usrmeta.RoleID == inconst.ROLE_MERCHANT {
		repoParams.UserID = usrmeta.UserID
	}

	data, err := s.repository.FindAccount(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.AccountResponse{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		AccountType: data.AccountType,
		Balance:     data.Balance,
		AccountNo:   cryptoutil.DecryptField(data.AccountNo, conf.DBKey),
	}

	hash := data.AccountNo
	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "accounts")
			logger.Error().Err(err).Send()
			return
		}
	} else {
		err = errs.New(errs.ErrDataIntegrity, "merchant")
		logger.Error().Err(err).Str("merchant-id", data.ID).Msg("row hash not found")
		return
	}

	return
}

func (s *service) CreateAccount(ctx context.Context, payload *dto.AccountPayload) (err error) {
	logger := component.GetLogger()
	conf := config.Get()

	usrctx := ctxutil.GetUserCTX(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	if val := structutil.CheckMandatoryField(payload); err != nil {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	var usermeta *indto.User

	if usrctx.RoleID <= inconst.ROLE_ADMIN {
		usermeta, err = s.findUserByID(ctx, payload.OwnerID)
	} else {
		usermeta, err = s.findUserMe(ctx)
	}

	if err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if usermeta == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("userID: %s not found", payload.OwnerID)
		return errs.ErrBadRequest
	}

	if exist, err := s.repository.FindAccount(ctx, &indto.AccountParams{UserID: payload.OwnerID}); err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if exist != nil {
		return errs.ErrDuplicatedResources
	}

	for {
		num, err := namegen.GenerateRandomNumber(8)
		if err != nil {
			logger.Warn().Err(err).Msg("failed to generate random account no")
			continue
		}

		if exists, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountNoHash: cryptoutil.HMACSHA512([]byte(fmt.Sprint(num)), conf.HashKey)}); err != nil {
			logger.Warn().Err(err).Msg("failed to check account no duplication")
			continue
		} else if exists != nil {
			continue
		}

		payload.AccountNo = fmt.Sprint(num)
		break
	}

	rowHash := []byte{}
	accModel := &model.Account{
		ID:            uuid.NewString(),
		OwnerID:       payload.OwnerID,
		Balance:       0,
		AccountNo:     cryptoutil.EncryptField([]byte(payload.AccountNo), conf.DBKey, &rowHash),
		AccountNoHash: cryptoutil.HMACSHA512([]byte(payload.AccountNo), conf.HashKey),
		RowHash:       rowHash,
	}

	if usermeta.RoleID == inconst.ROLE_CUSTOMER {
		accModel.AccountType = inconst.ACCOUNT_TYPE_CUST
	} else {
		accModel.AccountType = inconst.ACCOUNT_TYPE_MERCHANT
	}

	if enc, err := bcrypt.GenerateFromPassword([]byte(payload.PIN), bcrypt.DefaultCost); err != nil {
		logger.Error().Err(err).Msgf("failed to hash PIN")
		return err
	} else {
		accModel.PIN = string(enc)
	}

	accModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)
	if _, err = s.repository.CreateAccount(ctx, accModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) UpdateAccount(ctx context.Context, params *dto.AccountsQueryParams, payload *dto.AccountPayload) (err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	meta, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: params.AccountID})
	if err != nil {
		logger.Error().Err(err).Msgf("failed to fetch account meta")
		return err
	}

	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_CUSTOMER || usrmeta.RoleID == inconst.ROLE_MERCHANT {
		if meta.OwnerID != usrmeta.UserID {
			return errs.ErrNoAccess
		}
	}

	rowHash := []byte{}
	accModel := &model.Account{
		ID:            params.AccountID,
		AccountType:   payload.AccountType,
		Balance:       payload.Balance,
		AccountNo:     cryptoutil.EncryptField([]byte(payload.AccountNo), conf.DBKey, &rowHash),
		AccountNoHash: cryptoutil.HMACSHA512([]byte(payload.AccountNo), conf.HashKey),
		RowHash:       rowHash,
	}

	if payload.PIN != "" {
		if enc, err := bcrypt.GenerateFromPassword([]byte(payload.PIN), bcrypt.DefaultCost); err != nil {
			logger.Error().Err(err).Msgf("failed to hash PIN")
			return err
		} else {
			accModel.PIN = string(enc)
		}
	}

	accModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)

	if err = s.repository.UpdateAccount(ctx, accModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) DeleteAccount(ctx context.Context, params *dto.AccountsQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteAccount(ctx, &indto.AccountParams{AccountID: params.AccountID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) AuthenticateAccountMe(ctx context.Context, payload *dto.AccountPayload) (err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	usrmeta := ctxutil.GetUserCTX(ctx)
	repoParams := &indto.AccountParams{UserID: usrmeta.UserID}

	data, err := s.repository.FindAccount(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil || cryptoutil.VerifyHMACSHA512([]byte(payload.AccountNo), conf.HashKey, []byte(data.AccountNo)) {
		logger.Error().Str("reason", "no data found").Msg("failed to authenticate PIN")
		return errs.ErrNoAccess
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.PIN), []byte(payload.PIN))
	if err != nil {
		logger.Error().Err(err).Msg("failed to authenticate PIN")
		return errs.ErrNoAccess
	}

	return
}
