package apiutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

func (r *Requester[T]) SendRequest(ctx context.Context, endpoint string, method string, params map[string]string, body io.Reader) (res *T, err error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Stellar-Microservice by Misaki-chan")
	req.Header.Add("Content-Type", "application/json")

	queryParam := req.URL.Query()
	for k, v := range params {
		queryParam.Add(k, v)
	}

	req.URL.RawQuery = queryParam.Encode()
	data, err := r.Client.Do(req)
	if err != nil {
		return
	}

	defer data.Body.Close()
	if data.StatusCode < 200 || data.StatusCode > 299 {
		return nil, fmt.Errorf("failed to process request, got: %s", data.Status)
	}

	err = json.NewDecoder(data.Body).Decode(res)
	if err != nil {
		return
	}

	return
}
