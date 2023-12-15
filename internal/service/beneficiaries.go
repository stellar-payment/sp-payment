package service

import (
	"context"
	"math"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAllBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (res *dto.ListBeneficiaryResponse, err error) {
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

	repoParams := &indto.BeneficiaryParams{
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

	res = &dto.ListBeneficiaryResponse{
		Beneficiaries: []*dto.BeneficiaryResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountBeneficiaries(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindBeneficiaries(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.BeneficiaryResponse{
			ID:           v.ID,
			MerchantID:   v.MerchantID,
			MerchantName: v.MerchantName,
			Amount:       v.Amount,
			Status:       v.Status,
		}

		if v.WithdrawalDate.Valid {
			temp.WithdrawalDate = timeutil.FormatVerboseTime(v.WithdrawalDate.Time)
		}

		res.Beneficiaries = append(res.Beneficiaries, temp)
	}

	return
}

func (s *service) GetBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (res *dto.BeneficiaryResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	repoParams := &indto.BeneficiaryParams{BeneficiaryID: params.BeneficiaryID}

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

	data, err := s.repository.FindBeneficiary(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.BeneficiaryResponse{
		ID:             data.ID,
		MerchantID:     data.MerchantID,
		MerchantName:   data.MerchantName,
		Amount:         data.Amount,
		WithdrawalDate: "",
		Status:         data.Status,
	}

	if data.WithdrawalDate.Valid {
		res.WithdrawalDate = timeutil.FormatVerboseTime(data.WithdrawalDate.Time)
	}

	return
}

func (s *service) GetBeneficiaryPreview(ctx context.Context, params *dto.BeneficiariesQueryParams) (res float64, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_MERCHANT); !ok {
		return 0, errs.ErrNoAccess
	}

	repoParams := &indto.SettlementParams{MerchantID: params.MerchantID}
	usrmeta := ctxutil.GetUserCTX(ctx)
	if usrmeta.RoleID == inconst.ROLE_MERCHANT {
		merchantMeta, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{UserID: usrmeta.UserID})
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch merchant meta")
			return 0, err
		} else if merchantMeta == nil {
			err = errs.New(errs.ErrNotFound)
			logger.Error().Err(err).Str("user-id", usrmeta.UserID).Msg("failed to fetch merchant meta")
			return 0, err
		}

		repoParams.MerchantID = merchantMeta.ID
	}

	res, err = s.repository.FindPendingSettlement(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) CreateBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_MERCHANT); !ok {
		return errs.ErrNoAccess
	}

	repoParams := &indto.SettlementParams{MerchantID: params.MerchantID}

	usrmeta := ctxutil.GetUserCTX(ctx)
	userID := usrmeta.UserID
	if usrmeta.RoleID == inconst.ROLE_MERCHANT {
		merchantMeta, err := s.repository.FindMerchant(ctx, &indto.MerchantParams{UserID: usrmeta.UserID})
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch merchant meta")
			return err
		} else if merchantMeta == nil {
			err = errs.New(errs.ErrNotFound)
			logger.Error().Err(err).Str("user-id", usrmeta.UserID).Msg("failed to fetch merchant meta")
			return err
		}

		userID = usrmeta.UserID
		repoParams.MerchantID = merchantMeta.ID
	}

	accountMeta, err := s.repository.FindAccount(ctx, &indto.AccountParams{UserID: userID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	} else if accountMeta == nil {
		err = errs.ErrNotFound
		logger.Error().Err(err).Send()
		return
	}

	nominal, err := s.repository.FindPendingSettlement(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	beneModel := &model.Beneficiary{
		ID:             snowflake.ID(),
		AccountID:      accountMeta.ID,
		MerchantID:     repoParams.MerchantID,
		Amount:         nominal,
		WithdrawalDate: &time.Time{},
		Status:         inconst.BNF_STATUS_CONFIRM,
	}
	*beneModel.WithdrawalDate = time.Now()

	_, err = s.repository.CreateBeneficiary(ctx, beneModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
