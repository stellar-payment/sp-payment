package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
)

func (r *repository) FindAccounts(ctx context.Context, params *indto.AccountParams) (res []*indto.Account, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"a.deleted_at": nil},
	}

	if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"a.owner_id": params.UserID})
	}

	baseStmt := pgSquirrel.Select("a.id", "a.owner_id", "a.account_type", "a.balance", "a.account_no", "row_hash").
		From("accounts a").Where(cond)

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

	res = []*indto.Account{}
	for rows.Next() {
		temp := &indto.Account{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountAccounts(ctx context.Context, params *indto.AccountParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"a.deleted_at": nil},
	}

	if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"a.owner_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("accounts a").Where(cond).ToSql()
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

func (r *repository) FindAccount(ctx context.Context, params *indto.AccountParams) (res *indto.Account, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"a.deleted_at": nil},
	}

	if params.AccountID != "" {
		cond = append(cond, squirrel.Eq{"id": params.AccountID})
	} else if params.AccountNoHash != nil {
		cond = append(cond, squirrel.Eq{"account_no_hash": params.AccountNoHash})
	}

	if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"owner_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Select("a.id", "a.owner_id", "a.account_type", "a.balance", "a.account_no", "row_hash").
		From("accounts a").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Account{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) CreateAccount(ctx context.Context, payload *model.Account) (res *model.Account, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("accounts").Columns("id", "owner_id", "account_type", "balance", "account_no", "account_no_hash", "pin", "row_hash").
		Values(payload.ID, payload.OwnerID, payload.AccountType, payload.Balance, payload.AccountNo, payload.AccountNoHash, payload.PIN, payload.RowHash).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return payload, nil
}

func (r *repository) UpdateAccount(ctx context.Context, payload *model.Account) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("accounts").SetMap(map[string]interface{}{
		"account_type":    payload.AccountType,
		"balance":         payload.Balance,
		"account_no":      payload.AccountNo,
		"account_no_hash": payload.AccountNoHash,
		"pin":             squirrel.Expr("coalesce(nullif(?, ''), pin)", payload.PIN),
		"row_hash":        payload.RowHash,
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

func (r *repository) updateAccountBalanceTx(ctx context.Context, tx *sql.Tx, payload *model.Account) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("accounts").SetMap(map[string]interface{}{
		"balance":    squirrel.Expr("(balance - ?)", payload.Balance),
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

func (r *repository) DeleteAccount(ctx context.Context, params *indto.AccountParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
	}

	if params.AccountID != "" {
		cond = append(cond, squirrel.Eq{"id": params.AccountID})
	} else if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"owner_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Update("accounts").SetMap(map[string]interface{}{
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
