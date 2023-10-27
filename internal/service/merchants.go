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
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.ListMerchantResponse, err error) {
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
			PICName:      cryptoutil.DecryptField(v.PICName, conf.DBKey),
			PICEmail:     cryptoutil.DecryptField(v.PICEmail, conf.DBKey),
			PICPhone:     cryptoutil.DecryptField(v.PICPhone, conf.DBKey),
			PhotoProfile: v.PhotoProfile,
		}

		hash := v.PICName
		hash = append(hash, v.PICEmail...)
		hash = append(hash, v.PICPhone...)

		if len(v.RowHash) != 0 {
			if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, v.RowHash) {
				logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "merchant")).Send()
			}
		} else {
			logger.Warn().Err(errs.New(errs.ErrDataIntegrity, "merchant")).Str("merchant-id", v.ID).Msg("row hash not found")
		}

		res.Merchants = append(res.Merchants, temp)
	}

	return
}

func (s *service) GetMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.MerchantResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

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
		PICName:      cryptoutil.DecryptField(data.PICName, conf.DBKey),
		PICEmail:     cryptoutil.DecryptField(data.PICEmail, conf.DBKey),
		PICPhone:     cryptoutil.DecryptField(data.PICPhone, conf.DBKey),
		PhotoProfile: data.PhotoProfile,
	}

	hash := data.PICName
	hash = append(hash, data.PICEmail...)
	hash = append(hash, data.PICPhone...)

	if len(data.RowHash) != 0 {
		if !cryptoutil.VerifyHMACSHA512(hash, conf.HashKey, data.RowHash) {
			err = errs.New(errs.ErrDataIntegrity, "merchant")
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

func (s *service) HandleCreateMerchant(ctx context.Context, payload *indto.EventMerchant) (err error) {
	logger := component.GetLogger()
	conf := config.Get()

	rowHash := []byte{}
	custModel := &model.Merchant{
		ID:           uuid.NewString(),
		UserID:       payload.UserID,
		Name:         payload.Name,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Address:      payload.Address,
		PICName:      cryptoutil.EncryptField([]byte(payload.PICName), conf.DBKey, &rowHash),
		PICEmail:     cryptoutil.EncryptField([]byte(payload.PICEmail), conf.DBKey, &rowHash),
		PICPhone:     cryptoutil.EncryptField([]byte(payload.PICPhone), conf.DBKey, &rowHash),
		PhotoProfile: payload.PhotoProfile,
	}

	custModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)

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
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	rowHash := []byte{}
	custModel := &model.Merchant{
		ID:           params.MerchantID,
		Name:         payload.Name,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Address:      payload.Address,
		PICName:      cryptoutil.EncryptField([]byte(payload.PICName), conf.DBKey, &rowHash),
		PICEmail:     cryptoutil.EncryptField([]byte(payload.PICEmail), conf.DBKey, &rowHash),
		PICPhone:     cryptoutil.EncryptField([]byte(payload.PICPhone), conf.DBKey, &rowHash),
		PhotoProfile: payload.PhotoProfile,
	}

	custModel.RowHash = cryptoutil.HMACSHA512(rowHash, conf.HashKey)

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

func (s *service) HandleDeleteMerchant(ctx context.Context, params *indto.EventMerchant) (err error) {
	logger := component.GetLogger()

	err = s.repository.DeleteMerchant(ctx, &indto.MerchantParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
