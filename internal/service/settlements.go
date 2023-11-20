package service

import (
	"context"
	"math"

	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllSettlement(ctx context.Context, params *dto.SettlementsQueryParams) (res *dto.ListSettlementResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 100
	}

	repoParams := &indto.SettlementParams{
		Keyword: params.Keyword,
		Limit:   params.Limit,
		Page:    params.Page,
	}

	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_MERCHANT {
		merchantMeta, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{UserID: usrmeta.UserID})
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch merchant meta")
			return nil, err
		} else if merchantMeta == nil {
			err = errs.New(errs.ErrNotFound)
			logger.Error().Err(err).Str("user-id", usrmeta.UserID).Msg("failed to fetch merchant meta")
			return nil, err
		}

		repoParams.MerchantID = merchantMeta.ID
	}

	res = &dto.ListSettlementResponse{
		Settlements: []*dto.SettlementResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountSettlements(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindSettlements(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.SettlementResponse{
			ID:             v.ID,
			TransactionID:  v.TransactionID,
			MerchantID:     v.MerchantID,
			MerchantName:   v.MerchantName,
			BeneficiaryID:  v.BeneficiaryID,
			Amount:         v.Amount,
			Status:         v.Status,
			SettlementDate: timeutil.FormatVerboseTime(v.SettlementDate),
		}

		if v.BeneficiaryID != 0 {
			temp.Status = 1
		}

		res.Settlements = append(res.Settlements, temp)
	}

	return
}

func (s *service) GetSettlement(ctx context.Context, params *dto.SettlementsQueryParams) (res *dto.SettlementResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	repoParams := &indto.SettlementParams{SettlementID: params.SettlementID}

	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_MERCHANT {
		merchantMeta, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{UserID: usrmeta.UserID})
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch merchant meta")
			return nil, err
		} else if merchantMeta == nil {
			err = errs.New(errs.ErrNotFound)
			logger.Error().Err(err).Str("user-id", usrmeta.UserID).Msg("failed to fetch merchant meta")
			return nil, err
		}

		repoParams.MerchantID = merchantMeta.ID
	}

	data, err := s.repository.FindSettlement(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.SettlementResponse{
		ID:             data.ID,
		TransactionID:  data.TransactionID,
		MerchantID:     data.MerchantID,
		MerchantName:   data.MerchantName,
		BeneficiaryID:  data.BeneficiaryID,
		Amount:         data.Amount,
		Status:         data.Status,
		SettlementDate: timeutil.FormatVerboseTime(data.SettlementDate),
	}

	if data.BeneficiaryID != 0 {
		res.Status = 1
	}

	return
}
