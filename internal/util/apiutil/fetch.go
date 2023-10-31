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
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
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

type APIWrapper[T any] struct {
	Data  *T  `json:"data"`
	Error any `json:"error"`
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

	req, err := http.NewRequestWithContext(ctx, params.Method, params.Endpoint, bytes.NewBuffer([]byte(params.Body)))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Stellar-Microservice by Misaki-chan")

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

	if data.StatusCode <= 200 && data.StatusCode >= 299 {
		return
	}

	buf := &APIWrapper[T]{}

	defer data.Body.Close()
	defer logger.Info().Any("req-header", req.Header).Send()

	rawBody, _ := ioutil.ReadAll(data.Body)
	err = json.Unmarshal(rawBody, buf)
	if err != nil {
		logger.Error().Err(err).Str("raw", string(rawBody)).Send()
		return nil, fmt.Errorf("failed to read response body")
	}

	res.Payload = buf.Data
	return
}
