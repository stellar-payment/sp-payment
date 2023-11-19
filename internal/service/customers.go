package service

import (
	"context"
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
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.ListCustomerResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 100
	}

	repoParams := &indto.CustomerParams{
		Keyword: params.Keyword,
		Limit:   params.Limit,
		Page:    params.Page,
	}

	res = &dto.ListCustomerResponse{
		Customers: []*dto.CustomerResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountCustomers(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindCustomers(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.CustomerResponse{
			ID:           v.ID,
			UserID:       v.UserID,
			LegalName:    cryptoutil.DecryptField(v.LegalName, conf.DBKey),
			Phone:        cryptoutil.DecryptField(v.Phone, conf.DBKey),
			Email:        cryptoutil.DecryptField(v.Email, conf.DBKey),
			Birthdate:    cryptoutil.DecryptField(v.Birthdate, conf.DBKey),
			Address:      cryptoutil.DecryptField(v.Address, conf.DBKey),
			PhotoProfile: v.PhotoProfile,
		}

		hash := v.LegalName
		hash = append(hash, v.Phone...)
		hash = append(hash, v.Email...)
		hash = append(hash, v.Birthdate...)
		hash = append(hash, v.Address...)

		if len(v.RowHash) != 0 {
			if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, v.RowHash) {
				logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "customer")).Send()
			}
		} else {
			logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "customer")).Str("customer-id", v.ID).Msg("row hash not found")
		}

		res.Customers = append(res.Customers, temp)
	}

	return
}

func (s *service) GetCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.CustomerResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	data, err := s.repository.FindCustomer(ctx, &indto.CustomerParams{CustomerID: params.CustomerID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.CustomerResponse{
		ID:           data.ID,
		UserID:       data.UserID,
		LegalName:    cryptoutil.DecryptField(data.LegalName, conf.DBKey),
		Phone:        cryptoutil.DecryptField(data.Phone, conf.DBKey),
		Email:        cryptoutil.DecryptField(data.Email, conf.DBKey),
		Birthdate:    cryptoutil.DecryptField(data.Birthdate, conf.DBKey),
		Address:      cryptoutil.DecryptField(data.Address, conf.DBKey),
		PhotoProfile: data.PhotoProfile,
	}

	hash := data.LegalName
	hash = append(hash, data.Phone...)
	hash = append(hash, data.Email...)
	hash = append(hash, data.Birthdate...)
	hash = append(hash, data.Address...)

	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "customer")
			logger.Error().Err(err).Send()
			return
		}
	} else {
		err = errs.New(errs.ErrDataIntegrity, "customer")
		logger.Error().Err(err).Str("customer-id", data.ID).Msg("row hash not found")
		return
	}

	return
}

func (s *service) GetCustomerMe(ctx context.Context) (res *dto.CustomerResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER); !ok {
		return nil, errs.ErrNoAccess
	}

	usrmeta := ctxutil.GetUserCTX(ctx)
	data, err := s.repository.FindCustomer(ctx, &indto.CustomerParams{UserID: usrmeta.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.CustomerResponse{
		ID:           data.ID,
		UserID:       data.UserID,
		LegalName:    cryptoutil.DecryptField(data.LegalName, conf.DBKey),
		Phone:        cryptoutil.DecryptField(data.Phone, conf.DBKey),
		Email:        cryptoutil.DecryptField(data.Email, conf.DBKey),
		Birthdate:    cryptoutil.DecryptField(data.Birthdate, conf.DBKey),
		Address:      cryptoutil.DecryptField(data.Address, conf.DBKey),
		PhotoProfile: data.PhotoProfile,
	}

	hash := data.LegalName
	hash = append(hash, data.Phone...)
	hash = append(hash, data.Email...)
	hash = append(hash, data.Birthdate...)
	hash = append(hash, data.Address...)

	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "customer")
			logger.Error().Err(err).Send()
			return
		}
	} else {
		err = errs.New(errs.ErrDataIntegrity, "customer")
		logger.Error().Err(err).Str("customer-id", data.ID).Msg("row hash not found")
		return
	}

	return
}

func (s *service) HandleCreateCustomer(ctx context.Context, payload *indto.EventCustomer) (err error) {
	logger := component.GetLogger()
	conf := config.Get()

	rowHash := []byte{}
	custModel := &model.Customer{
		ID:           uuid.NewString(),
		UserID:       payload.UserID,
		LegalName:    cryptoutil.EncryptField([]byte(payload.LegalName), conf.DBKey, &rowHash),
		Phone:        cryptoutil.EncryptField([]byte(payload.Phone), conf.DBKey, &rowHash),
		Email:        cryptoutil.EncryptField([]byte(payload.Email), conf.DBKey, &rowHash),
		Birthdate:    cryptoutil.EncryptField([]byte(payload.Birthdate), conf.DBKey, &rowHash),
		Address:      cryptoutil.EncryptField([]byte(payload.Address), conf.DBKey, &rowHash),
		PhotoProfile: payload.PhotoProfile,
	}

	custModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)

	if _, err = s.repository.CreateCustomer(ctx, custModel); err != nil {
		logger.Error().Err(err).Send()

		if inerr := s.publishEvent(ctx, inconst.TOPIC_DELETE_USER, &indto.User{UserID: payload.UserID}); inerr != nil {
			logger.Error().Err(inerr).Send()
		}

		return
	}

	return
}

func (s *service) UpdateCustomer(ctx context.Context, params *dto.CustomersQueryParams, payload *dto.CustomerPayload) (err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	rowHash := []byte{}
	custModel := &model.Customer{
		ID:           params.CustomerID,
		LegalName:    cryptoutil.EncryptField([]byte(payload.LegalName), conf.DBKey, &rowHash),
		Phone:        cryptoutil.EncryptField([]byte(payload.Phone), conf.DBKey, &rowHash),
		Email:        cryptoutil.EncryptField([]byte(payload.Email), conf.DBKey, &rowHash),
		Birthdate:    cryptoutil.EncryptField([]byte(payload.Birthdate), conf.DBKey, &rowHash),
		Address:      cryptoutil.EncryptField([]byte(payload.Address), conf.DBKey, &rowHash),
		PhotoProfile: payload.PhotoProfile,
	}

	custModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)
	if err = s.repository.UpdateCustomer(ctx, custModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) DeleteCustomer(ctx context.Context, params *dto.CustomersQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteCustomer(ctx, &indto.CustomerParams{CustomerID: params.CustomerID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) HandleDeleteCustomer(ctx context.Context, params *indto.EventCustomer) (err error) {
	logger := component.GetLogger()

	err = s.repository.DeleteCustomer(ctx, &indto.CustomerParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
