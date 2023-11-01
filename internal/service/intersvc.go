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
)

func (s *service) findUserByID(ctx context.Context, id string) (res *indto.User, err error) {
	logger := zerolog.Ctx(ctx)
	invoker := apiutil.NewRequester[indto.User]()
	conf := config.Get()

	apires, err := invoker.SendRequest(ctx, &apiutil.SendRequestParams{
		Endpoint: fmt.Sprintf("%s%s%s", conf.AuthServiceAddr, inconst.ACCOUNT_USRID, id),
		Method:   http.MethodGet,
		Headers:  nil,
		Body:     "",
	})

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return apires.Payload, nil
}