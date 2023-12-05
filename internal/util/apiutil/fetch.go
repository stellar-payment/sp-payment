package apiutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
	"github.com/stellar-payment/sp-payment/pkg/errs"
	"golang.org/x/time/rate"
)

// Rate Limiter
type Transport struct {
	base    http.RoundTripper
	limiter *rate.Limiter
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	res := t.limiter.Reserve()

	select {
	case <-time.After(res.Delay()):
		return t.base.RoundTrip(r)
	case <-r.Context().Done():
		res.Cancel()
		return nil, r.Context().Err()
	}
}

type SendRequestParams struct {
	Endpoint string
	Method   string
	Headers  map[string]string
	Body     string
}

type APIResponse[T any] struct {
	Status  int
	Payload *T
	Headers map[string][]string
}

type APIError struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type APIWrapper[T any] struct {
	Data  *T       `json:"data"`
	Error APIError `json:"error"`
}

// Actual API
type Requester[T any] struct {
	Client http.Client
}

func NewRequester[T any]() Requester[T] {
	t := Requester[T]{
		Client: http.Client{
			Transport: &Transport{base: http.DefaultTransport, limiter: rate.NewLimiter(rate.Limit(50), 1)},
		},
	}
	return t
}

func (r *Requester[T]) SendRequest(ctx context.Context, params *SendRequestParams) (res *APIResponse[T], err error) {
	logger := zerolog.Ctx(ctx)
	conf := config.Get()

	req, err := http.NewRequestWithContext(ctx, params.Method, params.Endpoint, bytes.NewBuffer([]byte(params.Body)))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", fmt.Sprintf("%s-%s %s", conf.ServiceID, conf.Environment, conf.BuildVer))

	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	if v := req.Header.Get("Authorization"); v == "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctxutil.GetTokenCtx(ctx)))
	}

	data, err := r.Client.Do(req)
	if err != nil {
		return
	}

	res = &APIResponse[T]{}
	res.Headers = data.Header
	res.Status = data.StatusCode

	buf := &APIWrapper[T]{}
	rawBody, _ := ioutil.ReadAll(data.Body)

	defer data.Body.Close()

	err = json.Unmarshal(rawBody, buf)
	if err != nil {
		logger.Error().Err(err).Str("raw", string(rawBody)).Send()
		return nil, fmt.Errorf("failed to read response body")
	}

	if data.StatusCode <= 200 || data.StatusCode >= 299 {
		logger.Info().Str("api-msg", buf.Error.Msg).Send()
		if data.StatusCode == 403 {
			return nil, errs.ErrUserSessionExpired
		}

		return nil, errs.ErrUnknown
	}

	res.Payload = buf.Data

	logger.Info().
		Str("url", params.Endpoint).
		Str("status", data.Status).
		Any("req-header", req.Header).
		Any("req-body", params.Body).
		Any("res-header", data.Header).
		Any("res-body", buf).
		Send()

	return
}
