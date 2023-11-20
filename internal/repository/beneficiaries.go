package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func (r *repository) FindBeneficiaries(ctx context.Context, params *indto.BeneficiaryParams) (res []*indto.Beneficiary, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"b.deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"b.merchant_id": params.MerchantID})
	}

	baseStmt := pgSquirrel.Select("b.id", "b.merchant_id", "m.name merchant_name", "b.amount", "b.withdrawal_date", "b.status").
		From("beneficiaries b").
		LeftJoin("merchants m on m.id = b.merchant_id").
		Where(cond)

	if params.Limit != 0 && params.Page >= 1 {
		baseStmt = baseStmt.Limit(params.Limit).Offset((params.Page - 1) * params.Limit)
	}

	stmt, args, err := baseStmt.ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	rows, err := r.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	res = []*indto.Beneficiary{}
	for rows.Next() {
		temp := &indto.Beneficiary{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountBeneficiaries(ctx context.Context, params *indto.BeneficiaryParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"b.deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"b.merchant_id": params.MerchantID})
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("beneficiaries b").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	err = r.db.QueryRowxContext(ctx, stmt, args...).Scan(&res)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}

func (r *repository) FindBeneficiary(ctx context.Context, params *indto.BeneficiaryParams) (res *indto.Beneficiary, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"b.deleted_at": nil},
		squirrel.Eq{"b.id": params.BeneficiaryID},
	}

	stmt, args, err := pgSquirrel.Select("b.id", "b.merchant_id", "m.name merchant_name", "b.amount", "b.withdrawal_date", "b.status").
		From("beneficiaries b").
		LeftJoin("merchants m on m.id = b.merchant_id").
		Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Beneficiary{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) CreateBeneficiary(ctx context.Context, payload *model.Beneficiary) (res *model.Beneficiary, err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	defer tx.Rollback()

	_, err = r.createBeneficiaryTx(ctx, tx, payload)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	aff, err := r.updateSettlementBeneficiaryTx(ctx, tx, &model.Settlement{MerchantID: payload.MerchantID, BeneficiaryID: payload.ID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if aff == 0 {
		err = errs.ErrUnknown
		logger.Error().Err(err).Msg("settlement not properly affected")
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}

	return payload, nil
}

func (r *repository) createBeneficiaryTx(ctx context.Context, tx *sql.Tx, payload *model.Beneficiary) (res *model.Beneficiary, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("beneficiaries").Columns("id", "merchant_id", "amount", "withdrawal_date", "status").
		Values(payload.ID, payload.MerchantID, payload.Amount, payload.WithdrawalDate, payload.Status).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return payload, nil
}

func (r *repository) UpdateBeneficiary(ctx context.Context, payload *model.Beneficiary) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("beneficiaries").SetMap(map[string]interface{}{
		"withdrawal_date": payload.WithdrawalDate,
		"status":          payload.Status,
		"updated_at":      time.Now(),
	}).Where(squirrel.And{
		squirrel.Eq{"id": payload.ID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}

func (r *repository) DeleteBeneficiary(ctx context.Context, params *indto.BeneficiaryParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
		squirrel.Eq{"id": params.BeneficiaryID},
	}

	stmt, args, err := pgSquirrel.Update("beneficiaries").SetMap(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
	}).Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}
