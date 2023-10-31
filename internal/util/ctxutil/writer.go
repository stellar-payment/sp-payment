package ctxutil

import (
	"context"

	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
)

func WrapCtx(ctx context.Context, key inconst.CtxKey, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetCtx[T any](ctx context.Context, key inconst.CtxKey) (res T, ok bool) {
	res, ok = ctx.Value(key).(T)
	return
}

func GetUserCTX(ctx context.Context) (res *indto.UserResponse) {
	res, _ = GetCtx[*indto.UserResponse](ctx, inconst.AUTH_CTX_KEY)

	return
}

func GetCompanyIDCtx(ctx context.Context) (res int64) {
	res, _ = GetCtx[int64](ctx, inconst.MID_CTX_KEY)

	return
}

func GetTokenCtx(ctx context.Context) (res string) {
	res, _ = GetCtx[string](ctx, inconst.TOKEN_CTX_KEY)

	return
}
