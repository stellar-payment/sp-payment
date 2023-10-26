package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/util/apiutil"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
)

func (s *service) AuthorizedAccessCtx(ctx context.Context, token string) (res context.Context, err error) {
	logger := zerolog.Ctx(ctx)
	conf := config.Get()

	requester := apiutil.NewRequester[indto.UserResponse]()
	usrdata, err := requester.SendRequest(ctx, fmt.Sprintf("%s%s", conf.AuthServiceAddr, inconst.ACCOUNT_ME), http.MethodGet, map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	res = ctxutil.WrapCtx(ctx, inconst.AUTH_CTX_KEY, usrdata)
	return
}
