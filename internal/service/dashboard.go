package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/util/cryptoutil"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/internal/util/scopeutil"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (s *service) GetAdminDashboard(ctx context.Context) (res *dto.AdminDashboard, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	reports, err := s.repository.FindAdminDashboard(ctx)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	res = &dto.AdminDashboard{
		PeerTrxCount:     reports.PeerTrxCount,
		MerchantTrxCount: reports.MerchantTrxCount,
		SystemTrxCount:   reports.SystemTrxCount,
		TotalCustomers:   reports.TotalCustomers,
		TotalMerchants:   reports.TotalMerchants,
		TrxTraffic:       []dto.GenericDashboardGraph{},
	}

	for _, v := range reports.TrxTraffic {
		res.TrxTraffic = append(res.TrxTraffic, dto.GenericDashboardGraph{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return
}

func (s *service) GetMerchantDashboard(ctx context.Context) (res *dto.MerchantDashboard, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_MERCHANT); !ok {
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

	reports, err := s.repository.FindMerchantDashboard(ctx, &indto.MerchantDashboardParams{
		AccountID:  data.ID,
		MerchantID: data.OwnerID,
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	trx, err := s.repository.FindTransactions(ctx, &indto.TransactionParams{
		RecipientID: data.ID,
		TrxTypes:    []int64{inconst.TRX_TYPE_P2B, inconst.TRX_TYPE_BENEFICIARY, inconst.TRX_TYPE_MERCHANT_SYSTEM},
		Limit:       5,
		Page:        1,
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	res = &dto.MerchantDashboard{
		AccountID:          cryptoutil.DecryptField(data.AccountNo, conf.DBKey),
		AccountBalance:     data.Balance,
		TrxCount:           reports.TrxCount,
		TrxNominal:         reports.TrxNominal,
		SettlementNominal:  reports.SettlementNominal,
		BeneficiaryNominal: reports.BeneficiaryNominal,
		LastTrx:            []dto.TransactionMetaDashboard{},
	}

	for _, v := range trx {
		temp := dto.TransactionMetaDashboard{
			SenderName:    "",
			RecipientName: "",
			Nominal:       v.Nominal,
			TrxDate:       timeutil.FormatDate(v.TrxDatetime),
		}

		if v.AccountName != nil {
			temp.SenderName = cryptoutil.DecryptField(v.AccountName, conf.DBKey)
		}

		if v.RecipientName != nil {
			if v.TrxType == 1 || v.TrxType == 9 {
				temp.RecipientName = cryptoutil.DecryptField(v.RecipientName, conf.DBKey)
			} else {
				temp.RecipientName = string(v.RecipientName)
			}
		}

		res.LastTrx = append(res.LastTrx, temp)
	}

	return
}

func (s *service) GetCustomerDashboard(ctx context.Context) (res *dto.CustomerDashboard, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_CUSTOMER); !ok {
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

	reports, err := s.repository.FindCustomerDashboard(ctx, &indto.CustomerDashboardParams{
		AccountID:  data.ID,
		CustomerID: data.OwnerID,
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	trx, err := s.repository.FindTransactions(ctx, &indto.TransactionParams{
		RecipientID: data.ID,
		TrxTypes:    []int64{inconst.TRX_TYPE_P2P, inconst.TRX_TYPE_P2B, inconst.TRX_TYPE_CUST_SYSTEM},
		Limit:       5,
		Page:        1,
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	res = &dto.CustomerDashboard{
		AccountID:          cryptoutil.DecryptField(data.AccountNo, conf.DBKey),
		AccountBalance:     data.Balance,
		PeerTrxCount:       reports.PeerTrxCount,
		PeerTrxNominal:     reports.PeerTrxNominal,
		MerchantTrxCount:   reports.MerchantTrxCount,
		MerchantTrxNominal: reports.MerchantTrxNominal,
		LastTrx:            []dto.TransactionMetaDashboard{},
	}

	for _, v := range trx {
		temp := dto.TransactionMetaDashboard{
			SenderName:    "",
			RecipientName: "",
			Nominal:       v.Nominal,
			TrxDate:       timeutil.FormatDate(v.TrxDatetime),
		}

		if v.AccountName != nil {
			temp.SenderName = cryptoutil.DecryptField(v.AccountName, conf.DBKey)
		}

		if v.RecipientName != nil {
			if v.TrxType == 1 || v.TrxType == 9 {
				temp.RecipientName = cryptoutil.DecryptField(v.RecipientName, conf.DBKey)
			} else {
				temp.RecipientName = string(v.RecipientName)
			}
		}

		res.LastTrx = append(res.LastTrx, temp)
	}

	return
}
