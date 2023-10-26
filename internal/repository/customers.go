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

func (r *repository) FindCustomers(ctx context.Context, params *indto.CustomerParams) (res []*indto.Customer, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"c.deleted_at": nil},
	}

	baseStmt := pgSquirrel.Select("c.id", "c.user_id", "c.legal_name", "c.phone", "c.email", "c.birthdate", "c.address", "c.photo_profile").
		From("customers c").Where(cond)

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

	res = []*indto.Customer{}
	for rows.Next() {
		temp := &indto.Customer{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountCustomers(ctx context.Context, params *indto.CustomerParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"c.deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("customers c").Where(cond).ToSql()
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

func (r *repository) FindCustomer(ctx context.Context, params *indto.CustomerParams) (res *indto.Customer, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"c.id": params.CustomerID},
		squirrel.Eq{"c.deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Select("c.id", "c.user_id", "c.legal_name", "c.phone", "c.email", "c.birthdate", "c.address", "c.photo_profile").From("customers c").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Customer{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) CreateCustomer(ctx context.Context, payload *model.Customer) (res *model.Customer, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("customers").Columns("id", "user_id", "legal_name", "phone", "email", "birthdate", "address", "photo_profile").
		Values(payload.ID, payload.UserID, payload.LegalName, payload.Phone, payload.Email, payload.Birthdate, payload.Address, payload.PhotoProfile).ToSql()
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

func (r *repository) UpdateCustomer(ctx context.Context, payload *model.Customer) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("customers").SetMap(map[string]interface{}{
		"legal_name":    payload.LegalName,
		"phone":         payload.Phone,
		"email":         payload.Email,
		"address":       payload.Address,
		"birthdate":     payload.Birthdate,
		"photo_profile": payload.PhotoProfile,
		"updated_at":    time.Now(),
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

func (r *repository) DeleteCustomer(ctx context.Context, params *indto.CustomerParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
	}

	if params.CustomerID != "" {
		cond = append(cond, squirrel.Eq{"id": params.CustomerID})
	} else if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"user_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Update("customers").SetMap(map[string]interface{}{
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
