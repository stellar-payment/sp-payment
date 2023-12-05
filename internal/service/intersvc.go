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

func (s *service) findUserByID(ctx context.Context, id string) (res *indto.User, err error) {
	logger := zerolog.Ctx(ctx)
	invoker := apiutil.NewRequester[indto.User]()
	conf := config.Get()

	apires, err := invoker.SendRequest(ctx, &apiutil.SendRequestParams{
		Endpoint: fmt.Sprintf("%s%s%s", conf.AuthServiceAddr, inconst.ACCOUNT_USRID, id),
		Method:   http.MethodGet,
		Headers: map[string]string{
			"authorization": fmt.Sprintf("Bearer %s", ctxutil.GetTokenCtx(ctx)),
		},
		Body: "",
	})

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return apires.Payload, nil
}

func (s *service) findUserMe(ctx context.Context) (res *indto.User, err error) {
	logger := zerolog.Ctx(ctx)
	invoker := apiutil.NewRequester[indto.User]()
	conf := config.Get()

	apires, err := invoker.SendRequest(ctx, &apiutil.SendRequestParams{
		Endpoint: fmt.Sprintf("%s%s", conf.AuthServiceAddr, inconst.ACCOUNT_ME),
		Method:   http.MethodGet,
		Headers: map[string]string{
			"authorization": fmt.Sprintf("Bearer %s", ctxutil.GetTokenCtx(ctx)),
		},
		Body: "",
	})

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return apires.Payload, nil
}
