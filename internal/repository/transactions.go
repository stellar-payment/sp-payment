package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/godruoyi/go-snowflake"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
)

func (r *repository) FindTransactions(ctx context.Context, params *indto.TransactionParams) (res []*indto.Transaction, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"t.deleted_at": nil},
	}

	if !params.DateStart.IsZero() && !params.DateEnd.IsZero() {
		cond = append(cond,
			squirrel.Expr("date(t.trx_date) >= date(?)", params.DateStart),
			squirrel.Expr("date(t.trx_date) <= date(?)", params.DateEnd),
		)
	}

	if params.AccountID != "" {
		cond = append(cond, squirrel.Or{
			squirrel.Eq{"t.account_id": params.AccountID},
			squirrel.Eq{"t.recipient_id": params.AccountID},
		})
	}

	if params.TrxType != 0 {
		cond = append(cond, squirrel.Eq{"t.trx_type": params.TrxType})
	} else if len(params.TrxTypes) != 0 {
		cond = append(cond, squirrel.Eq{"t.trx_type": params.TrxTypes})
	}

	baseStmt := pgSquirrel.Select(
		"t.id", "t.account_id", "c1.legal_name account_name", "t.recipient_id", "coalesce(c2.legal_name, convert_to(m2.name, 'utf-8')) recipient_name",
		"t.trx_type", "t.trx_datetime", "t.trx_status", "t.trx_fee", "t.nominal", "t.description").
		From("transactions t").
		LeftJoin("accounts a1 on t.account_id = a1.id and t.trx_type not in (3, 9)").
		LeftJoin("customers c1 on a1.owner_id = c1.user_id").
		LeftJoin("accounts a2 on t.recipient_id = a2.id").
		LeftJoin("customers c2 on a2.owner_id = c2.user_id and t.trx_type in (1, 9)").
		LeftJoin("merchants m2 on a2.owner_id = m2.user_id and t.trx_type in (2, 3, 8)").
		Where(cond).OrderBy("t.created_at desc")

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

	res = []*indto.Transaction{}
	for rows.Next() {
		temp := &indto.Transaction{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountTransactions(ctx context.Context, params *indto.TransactionParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"t.deleted_at": nil},
	}

	if !params.DateStart.IsZero() && !params.DateEnd.IsZero() {
		cond = append(cond,
			squirrel.Expr("date(t.trx_date) >= date(?)", params.DateStart),
			squirrel.Expr("date(t.trx_date) <= date(?)", params.DateEnd),
		)
	}

	if params.AccountID != "" {
		cond = append(cond, squirrel.Or{
			squirrel.Eq{"t.account_id": params.AccountID},
			squirrel.Eq{"t.recipient_id": params.AccountID},
		})
	}

	if params.TrxType != 0 {
		cond = append(cond, squirrel.Eq{"t.trx_type": params.TrxType})
	} else if len(params.TrxTypes) != 0 {
		cond = append(cond, squirrel.Eq{"t.trx_type": params.TrxTypes})
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("transactions t").
		LeftJoin("accounts a1 on t.account_id = a1.id and t.trx_type not in (3, 9)").
		LeftJoin("customers c1 on a1.owner_id = c1.user_id").
		LeftJoin("accounts a2 on t.recipient_id = a2.id").
		LeftJoin("customers c2 on a2.owner_id = c2.user_id and t.trx_type in (1, 9)").
		LeftJoin("merchants m2 on a2.owner_id = m2.user_id and t.trx_type in (2, 3, 9)").
		Where(cond).ToSql()
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

func (r *repository) FindTransaction(ctx context.Context, params *indto.TransactionParams) (res *indto.Transaction, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"t.id": params.TransactionID},
		squirrel.Eq{"t.deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Select(
		"t.id", "t.account_id", "c1.legal_name account_name", "t.recipient_id", "coalesce(c2.legal_name, m2.name::bytea) recipient_name",
		"t.trx_type", "t.trx_datetime", "t.trx_status", "t.trx_fee", "t.nominal", "t.description").
		From("transactions t").
		LeftJoin("accounts a1 on t.account_id = a1.id and t.trx_type not in (3, 9)").
		LeftJoin("customers c1 on a1.owner_id = c1.user_id").
		LeftJoin("accounts a2 on t.recipient_id = a2.id").
		LeftJoin("customers c2 on a2.owner_id = c2.user_id and t.trx_type in (1, 9)").
		LeftJoin("merchants m2 on a2.owner_id = m2.user_id and t.trx_type in (2, 3, 9)").
		Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Transaction{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) CreateTransactionP2P(ctx context.Context, payload *model.Transaction) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	defer tx.Rollback()

	// substract sender's fund
	err = r.updateAccountBalanceTx(ctx, tx, &model.Account{ID: payload.AccountID, Balance: payload.Nominal + payload.TrxFee})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	// add receiver's fund
	// minus value used to denote addition
	err = r.updateAccountBalanceTx(ctx, tx, &model.Account{ID: payload.RecipientID, Balance: -payload.Nominal})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	// finally record transaction
	_, err = r.CreateTransactionTx(ctx, tx, payload)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}

	return
}

func (r *repository) CreateTransactionP2B(ctx context.Context, payload *model.Transaction) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	defer tx.Rollback()

	// substract sender's fund
	err = r.updateAccountBalanceTx(ctx, tx, &model.Account{ID: payload.AccountID, Balance: payload.Nominal + payload.TrxFee})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	// finally record transaction
	_, err = r.CreateTransactionTx(ctx, tx, payload)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}
	_, err = r.createSettlementTx(ctx, tx, &model.Settlement{
		ID:             snowflake.ID(),
		TransactionID:  payload.ID,
		MerchantID:     payload.MerchantID,
		Amount:         payload.Nominal,
		SettlementDate: time.Now(),
	})
	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	return
}

func (r *repository) CreateTransactionSystem(ctx context.Context, payload *model.Transaction) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	defer tx.Rollback()

	// add sender's fund
	err = r.updateAccountBalanceTx(ctx, tx, &model.Account{ID: payload.RecipientID, Balance: -payload.Nominal})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	// finally record transaction
	_, err = r.CreateTransactionTx(ctx, tx, payload)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("tx err")
		return
	}
	return
}

func (r *repository) CreateTransactionTx(ctx context.Context, tx *sql.Tx, payload *model.Transaction) (res *model.Transaction, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("transactions").Columns("id", "account_id", "recipient_id", "trx_type", "trx_datetime", "trx_status", "trx_fee", "nominal", "description").
		Values(payload.ID, payload.AccountID, payload.RecipientID, payload.TrxType, payload.TrxDatetime, payload.TrxStatus, payload.TrxFee, payload.Nominal, payload.Description).ToSql()
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

func (r *repository) UpdateTransaction(ctx context.Context, payload *model.Transaction) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("transactions").SetMap(map[string]interface{}{
		"trx_status": payload.TrxStatus,
		"updated_at": time.Now(),
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

func (r *repository) DeleteTransaction(ctx context.Context, params *indto.TransactionParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"id": params.TransactionID},
		squirrel.Eq{"deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Update("transactions").SetMap(map[string]interface{}{
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
