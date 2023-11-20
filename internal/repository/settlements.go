package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
)

func (r *repository) FindSettlements(ctx context.Context, params *indto.SettlementParams) (res []*indto.Settlement, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"s.deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"s.merchant_id": params.MerchantID})
	}

	if params.BeneficiaryID != 0 {
		cond = append(cond, squirrel.Eq{"s.beneficiary_id": params.BeneficiaryID})
	}

	baseStmt := pgSquirrel.Select("s.id", "s.transaction_id", "s.merchant_id", "m.name merchant_name", "s.beneficiary_id", "s.amount", "s.settlement_date").
		From("settlements s").
		LeftJoin("merchants m on s.merchant_id = m.id").
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

	res = []*indto.Settlement{}
	for rows.Next() {
		temp := &indto.Settlement{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountSettlements(ctx context.Context, params *indto.SettlementParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"s.deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"s.merchant_id": params.MerchantID})
	}

	if params.BeneficiaryID != 0 {
		cond = append(cond, squirrel.Eq{"s.beneficiary_id": params.BeneficiaryID})
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("settlements s").Where(cond).ToSql()
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

func (r *repository) FindSettlement(ctx context.Context, params *indto.SettlementParams) (res *indto.Settlement, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"s.deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"s.merchant_id": params.MerchantID})
	}

	if params.BeneficiaryID != 0 {
		cond = append(cond, squirrel.Eq{"s.beneficiary_id": params.BeneficiaryID})
	}

	stmt, args, err := pgSquirrel.Select("s.id", "s.transaction_id", "s.merchant_id", "m.name merchant_name", "s.beneficiary_id", "s.amount", "s.settlement_date").
		From("settlements s").
		LeftJoin("merchants m on m.id = s.merchant_name").
		Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Settlement{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) createSettlementTx(ctx context.Context, tx *sql.Tx, payload *model.Settlement) (res *model.Settlement, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("settlements").Columns("id", "transaction_id", "merchant_id", "beneficiary_id", "amount", "settlement_date").
		Values(payload.ID, payload.TransactionID, payload.MerchantID, payload.BeneficiaryID, payload.Amount, payload.SettlementDate).ToSql()
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

func (r *repository) updateSettlementTx(ctx context.Context, tx *sql.Tx, payload *model.Settlement) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
	}

	if payload.ID != 0 {
		cond = append(cond, squirrel.Eq{"id": payload.ID})
	} else if payload.BeneficiaryID != 0 {
		cond = append(cond, squirrel.Eq{"beneficiary_id": payload.BeneficiaryID})
	} else if payload.TransactionID != 0 {
		cond = append(cond, squirrel.Eq{"transaction_id": payload.TransactionID})
	}

	stmt, args, err := pgSquirrel.Update("settlements").SetMap(map[string]interface{}{
		"beneficiary_id": squirrel.Expr("coalesce(nullif(?, 0), beneficiary_id)", payload.BeneficiaryID),
		"updated_at":     time.Now(),
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

func (r *repository) updateSettlementBeneficiaryTx(ctx context.Context, tx *sql.Tx, payload *model.Settlement) (aff int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
		squirrel.Eq{"beneficiary_id": 0},
		squirrel.Eq{"merchant_id": payload.MerchantID},
	}

	stmt, args, err := pgSquirrel.Update("settlements").SetMap(map[string]interface{}{
		"beneficiary_id": payload.BeneficiaryID,
		"updated_at":     time.Now(),
	}).Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	fmt.Println(stmt)
	fmt.Println(args...)

	res, err := r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	aff, _ = res.RowsAffected()

	return
}

func (r *repository) deleteSettlementTx(ctx context.Context, tx *sql.Tx, params *indto.SettlementParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
	}

	if params.SettlementID != 0 {
		cond = append(cond, squirrel.Eq{"id": params.SettlementID})
	} else if params.BeneficiaryID != 0 {
		cond = append(cond, squirrel.Eq{"beneficiary_id": params.BeneficiaryID})
	} else if params.TransactionID != 0 {
		cond = append(cond, squirrel.Eq{"transaction_id": params.TransactionID})
	}

	stmt, args, err := pgSquirrel.Update("settlements").SetMap(map[string]interface{}{
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

func (r *repository) FindPendingSettlement(ctx context.Context, params *indto.SettlementParams) (res float64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"s.deleted_at": nil},
		squirrel.Eq{"s.beneficiary_id": 0},
		squirrel.Eq{"s.merchant_id": params.MerchantID},
	}

	stmt, args, err := pgSquirrel.Select("coalesce(sum(amount), 0)").From("settlements s").Where(cond).ToSql()
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
