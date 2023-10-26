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

func (s *service) GetAllMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.ListMerchantResponse, err error) {
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

	repoParams := &indto.MerchantParams{
		Keyword: params.Keyword,
		Limit:   params.Limit,
		Page:    params.Page,
	}

	res = &dto.ListMerchantResponse{
		Merchants: []*dto.MerchantResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountMerchants(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindMerchants(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.MerchantResponse{
			ID:           v.ID,
			UserID:       v.UserID,
			Name:         v.Name,
			Phone:        v.Phone,
			Email:        v.Email,
			Address:      v.Address,
			PICName:      v.PICName,
			PICEmail:     v.PICEmail,
			PICPhone:     v.PICPhone,
			PhotoProfile: v.PhotoProfile,
		}

		res.Merchants = append(res.Merchants, temp)
	}

	return
}

func (s *service) GetMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.MerchantResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	data, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{MerchantID: params.MerchantID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.MerchantResponse{
		ID:           data.ID,
		UserID:       data.UserID,
		Name:         data.Name,
		Phone:        data.Phone,
		Email:        data.Email,
		Address:      data.Address,
		PICName:      data.PICName,
		PICEmail:     data.PICEmail,
		PICPhone:     data.PICPhone,
		PhotoProfile: data.PhotoProfile,
	}

	return
}

func (s *service) HandleCreateMerchant(ctx context.Context, payload *indto.Merchant) (err error) {
	logger := log.Ctx(ctx)

	custModel := &model.Merchant{
		ID:           uuid.NewString(),
		UserID:       payload.UserID,
		Name:         payload.Name,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Address:      payload.Address,
		PICName:      payload.PICName,
		PICEmail:     payload.PICEmail,
		PICPhone:     payload.PICPhone,
		PhotoProfile: payload.PhotoProfile,
	}

	if _, err = s.repository.CreateMerchant(ctx, custModel); err != nil {
		logger.Error().Err(err).Send()

		if inerr := s.publishEvent(ctx, inconst.TOPIC_DELETE_USER, &indto.User{UserID: payload.UserID}); inerr != nil {
			logger.Error().Err(inerr).Send()
		}

		return
	}

	return
}

func (s *service) UpdateMerchant(ctx context.Context, params *dto.MerchantsQueryParams, payload *dto.MerchantPayload) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	custModel := &model.Merchant{
		ID:           params.MerchantID,
		Name:         payload.Name,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Address:      payload.Address,
		PICName:      payload.PICName,
		PICEmail:     payload.PICEmail,
		PICPhone:     payload.PICPhone,
		PhotoProfile: payload.PhotoProfile,
	}
	if err = s.repository.UpdateMerchant(ctx, custModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) DeleteMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteMerchant(ctx, &indto.MerchantParams{MerchantID: params.MerchantID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) HandleDeleteMerchant(ctx context.Context, params *indto.Merchant) (err error) {
	logger := log.Ctx(ctx)

	err = s.repository.DeleteMerchant(ctx, &indto.MerchantParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
