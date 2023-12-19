package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
)

func (r *repository) FindAdminDashboard(ctx context.Context) (res *indto.AdminDashboard, err error) {
	logger := zerolog.Ctx(ctx)
	var stmt string
	var namedArgs map[string]any

	dateStart, dateEnd := timeutil.GetStartEndMonth(time.Now())

	res = &indto.AdminDashboard{
		PeerTrxCount:     0,
		MerchantTrxCount: 0,
		SystemTrxCount:   0,
		TotalCustomers:   0,
		TotalMerchants:   0,
		TrxTraffic:       []indto.GenericDashboardGraph{},
	}

	stmt = `
		select
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 1 and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) peer_trx_count,
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 2 and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) merchant_trx_count,
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 9 and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) system_trx_count,
			(select coalesce(count(*), 0) from customers where deleted_at is null) total_customers,
			(select coalesce(count(*), 0) from merchants where deleted_at is null) total_merchants;
	`

	namedArgs = map[string]any{
		"date_start": dateStart,
		"date_end":   dateEnd,
	}

	rows, err := r.db.NamedQueryContext(ctx, stmt, namedArgs)
	if err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Send()
		return
	}

	rows.Next()
	if err = rows.StructScan(res); err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Msg("sql map err")
		return
	}

	stmt, args, err := pgSquirrel.Select("to_char(date_trunc('month', trx_datetime), 'MM') trx_mon", "sum(nominal)").From("transactions").
		Where(squirrel.And{
			squirrel.Eq{"deleted_at": nil},
			squirrel.Expr("date_part('year', trx_datetime) = date_part('year', ?::timestamp)", dateStart),
		}).GroupBy("date_trunc('month', trx_datetime)").ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	fmt.Println(stmt)
	fmt.Println(args)

	rows, err = r.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Str("query", "trx_traffic").Send()
		return
	}

	temp := make(map[int64]float64)
	for rows.Next() {
		var mon int64
		var nominal float64

		if err = rows.Scan(&mon, &nominal); err != nil {
			logger.Error().Err(err).Str("query", "trx_traffic").Msg("sql map err")
			return
		}

		temp[mon] = nominal
	}

	for i := 1; i <= 12; i++ {
		key := int64(i)

		res.TrxTraffic = append(res.TrxTraffic, indto.GenericDashboardGraph{
			Key:   key,
			Value: temp[key],
		})
	}

	return
}
func (r *repository) FindMerchantDashboard(ctx context.Context, param *indto.MerchantDashboardParams) (res *indto.MerchantDashboard, err error) {
	logger := zerolog.Ctx(ctx)
	var stmt string
	var namedArgs map[string]any

	dateStart, dateEnd := timeutil.GetStartEndMonth(time.Now())

	res = &indto.MerchantDashboard{
		TrxCount:           0,
		TrxNominal:         0,
		SettlementNominal:  0,
		BeneficiaryNominal: 0,
	}

	stmt = `
		select
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 2 and recipient_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) trx_count,
			(select coalesce(sum(nominal), 0) from transactions where deleted_at is null and trx_type = 2 and recipient_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) trx_nominal,    
			(select coalesce(sum(amount), 0) from settlements where deleted_at is null and beneficiary_id = 0 and merchant_id = :merchant_id and date(settlement_date) >= date(:date_start) and date(settlement_date) <= date(:date_end)) settlement_nominal,
			(select coalesce(sum(amount), 0) from beneficiaries where deleted_at is null and merchant_id = :merchant_id and date(withdrawal_date) >= date(:date_start) and date(withdrawal_date) <= date(:date_end)) beneficiary_nominal
	`

	namedArgs = map[string]any{
		"date_start":  dateStart,
		"date_end":    dateEnd,
		"account_id":  param.AccountID,
		"merchant_id": param.MerchantID,
	}

	rows, err := r.db.NamedQueryContext(ctx, stmt, namedArgs)
	if err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Send()
		return
	}

	rows.Next()
	if err = rows.StructScan(res); err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Msg("sql map err")
		return
	}

	return
}

func (r *repository) FindCustomerDashboard(ctx context.Context, param *indto.CustomerDashboardParams) (res *indto.CustomerDashboard, err error) {
	logger := zerolog.Ctx(ctx)
	var stmt string
	var namedArgs map[string]any

	dateStart, dateEnd := timeutil.GetStartEndMonth(time.Now())

	res = &indto.CustomerDashboard{
		PeerTrxCount:       0,
		PeerTrxNominal:     0,
		MerchantTrxCount:   0,
		MerchantTrxNominal: 0,
	}

	stmt = `
		select
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 1 and account_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) peer_trx_count,
			(select coalesce(sum(nominal), 0) from transactions where deleted_at is null and trx_type = 1 and account_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) peer_trx_nominal,
			(select coalesce(count(*), 0) from transactions where deleted_at is null and trx_type = 2 and account_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) merchant_trx_count,
			(select coalesce(sum(nominal), 0) from transactions where deleted_at is null and trx_type = 2 and account_id = :account_id and date(trx_datetime) >= date(:date_start) and date(trx_datetime) <= date(:date_end)) merchant_trx_nominal
	`

	namedArgs = map[string]any{
		"date_start": dateStart,
		"date_end":   dateEnd,
		"account_id": param.AccountID,
	}

	rows, err := r.db.NamedQueryContext(ctx, stmt, namedArgs)
	if err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Send()
		return
	}

	rows.Next()
	if err = rows.StructScan(res); err != nil {
		logger.Error().Err(err).Str("query", "overall_dashboard").Msg("sql map err")
		return
	}

	return
}
