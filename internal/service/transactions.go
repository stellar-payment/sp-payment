package service

import (
	"context"
	"math"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/component"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"github.com/stellar-payment/sp-payment/internal/util/cryptoutil"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/internal/util/structutil"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (res *dto.ListTransactionResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 100
	}

	repoParams := &indto.TransactionParams{
		AccountID: params.AccountID,
		TrxType:   params.TrxType,
		DateStart: timeutil.ParseDate(params.DateStart),
		DateEnd:   timeutil.ParseDate(params.DateEnd),
		Keyword:   params.Keyword,
		Limit:     params.Limit,
		Page:      params.Page,
	}

	res = &dto.ListTransactionResponse{
		Transactions: []*dto.TransactionResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	usermeta := ctxutil.GetUserCTX(ctx)
	if usermeta.RoleID == inconst.ROLE_CUSTOMER || usermeta.RoleID == inconst.ROLE_MERCHANT {
		val := &indto.Account{}
		val, err = s.repository.FindAccount(ctx, &indto.AccountParams{UserID: usermeta.UserID})
		if err != nil {
			logger.Error().Err(err).Send()
			return
		} else if val == nil {
			return
		}

		repoParams.AccountID = val.ID
	}

	count, err := s.repository.CountTransactions(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindTransactions(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, data := range data {
		temp := &dto.TransactionResponse{
			ID:          data.ID,
			TrxType:     data.TrxType,
			TrxDatetime: timeutil.FormatVerboseTime(data.TrxDatetime),
			TrxStatus:   data.TrxStatus,
			TrxFee:      data.TrxFee,
			Nominal:     data.Nominal,
			Description: data.Description,
		}

		if data.AccountName != nil {
			temp.AccountID = data.AccountID
			temp.AccountName = cryptoutil.DecryptField(data.AccountName, conf.DBKey)
		}

		if data.RecipientName != nil {
			temp.RecipientID = data.RecipientID

			if data.TrxType == 1 || data.TrxType == 9 {
				temp.RecipientName = cryptoutil.DecryptField(data.RecipientName, conf.DBKey)
			} else {
				temp.RecipientName = string(data.RecipientName)
			}
		}

		res.Transactions = append(res.Transactions, temp)
	}

	return
}

func (s *service) GetTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (res *dto.TransactionResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	data, err := s.repository.FindTransaction(ctx, &indto.TransactionParams{TransactionID: params.TransactionID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.TransactionResponse{
		ID:          data.ID,
		TrxType:     data.TrxType,
		TrxDatetime: timeutil.FormatVerboseTime(data.TrxDatetime),
		TrxStatus:   data.TrxStatus,
		TrxFee:      data.TrxFee,
		Nominal:     data.Nominal,
		Description: data.Description,
	}

	if data.AccountName != nil {
		res.AccountID = data.AccountID
		res.AccountName = cryptoutil.DecryptField(data.AccountName, conf.DBKey)
	}

	if data.RecipientName != nil {
		res.RecipientID = data.RecipientID

		if data.TrxType == 1 || data.TrxType == 9 {
			res.RecipientName = cryptoutil.DecryptField(data.RecipientName, conf.DBKey)
		} else {
			res.RecipientName = string(data.RecipientName)
		}
	}

	return
}

func (s *service) CreateTransactionP2P(ctx context.Context, payload *dto.TransactionPayload) (err error) {
	logger := component.GetLogger()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	if val := structutil.CheckMandatoryField(payload); err != nil {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	senderMeta, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.AccountID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	} else if senderMeta == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("sender accountID: %s not found", payload.AccountID)
		return errs.ErrBadRequest
	}

	if exists, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.RecipientID}); err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if exists == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s not found", payload.AccountID)
		return errs.ErrBadRequest
	} else if exists.AccountType != inconst.ACCOUNT_TYPE_CUST {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s is not customer", payload.AccountID)
		return errs.ErrBadRequest
	}

	if senderMeta.Balance < payload.Nominal*1.1 {
		err = errs.ErrInsufficientBalance
		logger.Error().Err(err).Msgf("accountID: %s does not have enough balance. (has=%.2f, need=%.2f)", senderMeta.Balance, payload.Nominal*1.1)
		return
	}

	trxModel := &model.Transaction{
		ID:          snowflake.ID(),
		AccountID:   payload.AccountID,
		RecipientID: payload.RecipientID,
		TrxType:     inconst.TRX_TYPE_P2P,
		TrxDatetime: time.Now(),
		TrxStatus:   inconst.TRX_STATUS_SUCCESS, // always success
		TrxFee:      payload.Nominal * 0.1,
		Nominal:     payload.Nominal,
		Description: payload.Description,
	}

	err = s.repository.CreateTransactionP2P(ctx, trxModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) CreateTransactionSystem(ctx context.Context, payload *dto.TransactionPayload) (err error) {
	logger := component.GetLogger()
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	if val := structutil.CheckMandatoryField(payload); err != nil {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	senderMeta, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.AccountID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	} else if senderMeta == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("accountID: %s not found", payload.AccountID)
		return errs.ErrBadRequest
	}

	if exists, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.RecipientID}); err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if exists == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s not found", payload.RecipientID)
		return errs.ErrBadRequest
	} else if exists.AccountType != inconst.ACCOUNT_TYPE_CUST {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s is not customer", payload.RecipientID)
		return errs.ErrBadRequest
	}

	trxModel := &model.Transaction{
		ID:          snowflake.ID(),
		AccountID:   conf.SystemAccountUUID,
		RecipientID: payload.RecipientID,
		TrxType:     inconst.TRX_TYPE_SYSTEM,
		TrxDatetime: time.Now(),
		TrxStatus:   inconst.TRX_STATUS_SUCCESS, // always success
		TrxFee:      0,
		Nominal:     payload.Nominal,
		Description: payload.Description,
	}

	err = s.repository.CreateTransactionSystem(ctx, trxModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) CreateTransactionP2B(ctx context.Context, payload *dto.TransactionPayload) (err error) {
	logger := component.GetLogger()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	if val := structutil.CheckMandatoryField(payload); err != nil {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	senderMeta, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.AccountID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	} else if senderMeta == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("accountID: %s not found", payload.AccountID)
		return errs.ErrBadRequest
	}

	recipientMeta, err := s.repository.FindAccount(ctx, &indto.AccountParams{AccountID: payload.RecipientID})
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if recipientMeta == nil {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s not found", payload.AccountID)
		return errs.ErrBadRequest
	} else if recipientMeta.AccountType != inconst.ACCOUNT_TYPE_MERCHANT {
		logger.Error().Err(errs.ErrNotFound).Msgf("recepient accountID: %s is not merchant", payload.AccountID)
		return errs.ErrBadRequest
	}

	merchantMeta, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{UserID: recipientMeta.OwnerID})
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch merchant meta")
		return err
	} else if merchantMeta == nil {
		err = errs.New(errs.ErrNotFound)
		logger.Error().Err(err).Str("user-id", recipientMeta.OwnerID).Msg("failed to fetch merchant meta")
		return err
	}

	if senderMeta.Balance < payload.Nominal*1.1 {
		err = errs.ErrInsufficientBalance
		logger.Error().Err(err).Msgf("accountID: %s does not have enough balance. (has=%.2f, need=%.2f)", senderMeta.Balance, payload.Nominal*1.1)
		return
	}

	trxModel := &model.Transaction{
		ID:          snowflake.ID(),
		AccountID:   payload.AccountID,
		RecipientID: payload.RecipientID,
		MerchantID:  merchantMeta.ID,
		TrxType:     inconst.TRX_TYPE_P2B,
		TrxDatetime: time.Now(),
		TrxStatus:   inconst.TRX_STATUS_SUCCESS,
		TrxFee:      payload.Nominal * 0.1,
		Nominal:     payload.Nominal,
		Description: payload.Description,
	}

	err = s.repository.CreateTransactionP2B(ctx, trxModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) UpdateTransaction(ctx context.Context, params *dto.TransactionsQueryParams, payload *dto.TransactionPayload) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	custModel := &model.Transaction{
		ID:        params.TransactionID,
		TrxStatus: payload.TrxStatus,
	}

	if err = s.repository.UpdateTransaction(ctx, custModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) DeleteTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteTransaction(ctx, &indto.TransactionParams{TransactionID: params.TransactionID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
