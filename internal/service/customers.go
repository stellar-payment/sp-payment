package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.ListCustomerResponse, err error) {
	logger := log.Ctx(ctx)

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
			LegalName:    v.LegalName,
			Phone:        v.Phone,
			Email:        v.Email,
			Birthdate:    v.Birthdate,
			Address:      v.Address,
			PhotoProfile: v.PhotoProfile,
		}

		res.Customers = append(res.Customers, temp)
	}

	return
}

func (s *service) GetCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.CustomerResponse, err error) {
	logger := log.Ctx(ctx)

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
		LegalName:    data.LegalName,
		Phone:        data.Phone,
		Email:        data.Email,
		Birthdate:    data.Birthdate,
		Address:      data.Address,
		PhotoProfile: data.PhotoProfile,
	}

	return
}

func (s *service) HandleCreateCustomer(ctx context.Context, payload *indto.Customer) (err error) {
	logger := log.Ctx(ctx)

	custModel := &model.Customer{
		ID:           uuid.NewString(),
		UserID:       payload.UserID,
		LegalName:    payload.LegalName,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Birthdate:    payload.Birthdate,
		Address:      payload.Address,
		PhotoProfile: payload.PhotoProfile,
	}

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

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	custModel := &model.Customer{
		ID:           params.CustomerID,
		LegalName:    payload.LegalName,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Birthdate:    payload.Birthdate,
		Address:      payload.Address,
		PhotoProfile: payload.PhotoProfile,
	}
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

func (s *service) HandleDeleteCustomer(ctx context.Context, params *indto.Customer) (err error) {
	logger := log.Ctx(ctx)

	err = s.repository.DeleteCustomer(ctx, &indto.CustomerParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
