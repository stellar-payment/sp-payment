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

func (r *repository) FindMerchants(ctx context.Context, params *indto.MerchantParams) (res []*indto.Merchant, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"m.deleted_at": nil},
	}

	baseStmt := pgSquirrel.Select("m.id", "m.user_id", "m.name", "m.address", "m.phone", "m.email", "m.pic_name", "m.pic_email", "m.pic_phone", "m.photo_profile", "m.row_hash").
		From("merchants m").Where(cond)

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

	res = []*indto.Merchant{}
	for rows.Next() {
		temp := &indto.Merchant{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountMerchants(ctx context.Context, params *indto.MerchantParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"m.deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("merchants m").Where(cond).ToSql()
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

func (r *repository) FindMerchant(ctx context.Context, params *indto.MerchantParams) (res *indto.Merchant, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"m.deleted_at": nil},
	}
	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"id": params.MerchantID})
	}

	if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"user_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Select("m.id", "m.user_id", "m.name", "m.address", "m.phone", "m.email", "m.pic_name", "m.pic_email", "m.pic_phone", "m.photo_profile", "m.row_hash").
		From("merchants m").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.Merchant{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) CreateMerchant(ctx context.Context, payload *model.Merchant) (res *model.Merchant, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("merchants").Columns("id", "user_id", "name", "phone", "email", "address", "pic_name", "pic_email", "pic_phone", "photo_profile", "row_hash").
		Values(payload.ID, payload.UserID, payload.Name, payload.Phone, payload.Email, payload.Address, payload.PICName, payload.PICEmail, payload.PICPhone, payload.PhotoProfile, payload.RowHash).ToSql()
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

func (r *repository) UpdateMerchant(ctx context.Context, payload *model.Merchant) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("merchants").SetMap(map[string]interface{}{
		"name":          payload.Name,
		"phone":         payload.Phone,
		"email":         payload.Email,
		"address":       payload.Address,
		"pic_name":      payload.PICName,
		"pic_phone":     payload.PICPhone,
		"pic_email":     payload.PICEmail,
		"photo_profile": payload.PhotoProfile,
		"row_hash":      payload.RowHash,
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

func (r *repository) DeleteMerchant(ctx context.Context, params *indto.MerchantParams) (err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"deleted_at": nil},
	}

	if params.MerchantID != "" {
		cond = append(cond, squirrel.Eq{"id": params.MerchantID})
	} else if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"user_id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Update("merchants").SetMap(map[string]interface{}{
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
