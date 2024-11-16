package controld

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"golang.org/x/time/rate"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// API holds the configuration for the current API client. A client should not
// be modified concurrently.
type API struct {
	APIToken    string
	BaseURL     string
	UserAgent   string
	headers     http.Header
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	retryPolicy RetryPolicy
	logger      Logger
	Debug       bool
}

// newClient provides shared logic for New and NewWithUserServiceKey.
func newClient(opts ...Option) (*API, error) {
	silentLogger := log.New(io.Discard, "", log.LstdFlags)

	api := &API{
		BaseURL:     fmt.Sprintf("%s://%s%s", defaultScheme, defaultHostname, defaultBasePath),
		UserAgent:   userAgent,
		headers:     make(http.Header),
		rateLimiter: rate.NewLimiter(rate.Limit(4), 1), // 4rps equates to default api limit (1200 req/5 min)
		retryPolicy: RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: 1 * time.Second,
			MaxRetryDelay: 30 * time.Second,
		},
		logger: silentLogger,
	}

	err := api.parseOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("options parsing failed: %w", err)
	}

	// Fall back to http.DefaultClient if the package user does not provide
	// their own.
	if api.httpClient == nil {
		api.httpClient = http.DefaultClient
	}

	return api, nil
}

func New(token string, opts ...Option) (*API, error) {
	if token == "" {
		return nil, errors.New(errEmptyAPIToken)
	}

	api, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	api.APIToken = token

	return api, nil
}

// makeRequest makes a HTTP request and returns the body as a byte slice,
// closing it before returning. params will be serialized to JSON.
//
//nolint:unused
func (api *API) makeRequest(method, uri string, params interface{}) ([]byte, error) {
	return api.makeRequestWithAuthType(context.Background(), method, uri, params)
}

func (api *API) makeRequestContext(ctx context.Context, method, uri string, params interface{}) ([]byte, error) {
	return api.makeRequestWithAuthType(ctx, method, uri, params)
}

func (api *API) makeRequestContextWithHeaders(ctx context.Context, method, uri string, params interface{}, headers http.Header) ([]byte, error) {
	return api.makeRequestWithAuthTypeAndHeaders(ctx, method, uri, params, headers)
}

func (api *API) makeRequestWithAuthType(ctx context.Context, method, uri string, params interface{}) ([]byte, error) {
	return api.makeRequestWithAuthTypeAndHeaders(ctx, method, uri, params, nil)
}

// APIResponse holds the structure for a response from the API. It looks alot
// like `http.Response` however, uses a `[]byte` for the `Body` instead of a
// `io.ReadCloser`.
//
// This may go away in the experimental client in favour of `http.Response`.
type APIResponse struct {
	Body       []byte
	Status     string
	StatusCode int
	Headers    http.Header
}

func (api *API) makeRequestWithAuthTypeAndHeaders(ctx context.Context, method, uri string, params interface{}, headers http.Header) ([]byte, error) {
	res, err := api.makeRequestWithAuthTypeAndHeadersComplete(ctx, method, uri, params, headers)
	if err != nil {
		return nil, err
	}
	return res.Body, err
}

// Use this method if an API response can have different Content-Type headers and different body formats.
//
//nolint:unused
func (api *API) makeRequestContextWithHeadersComplete(ctx context.Context, method, uri string, params interface{}, headers http.Header) (*APIResponse, error) {
	return api.makeRequestWithAuthTypeAndHeadersComplete(ctx, method, uri, params, headers)
}

func (api *API) makeRequestWithAuthTypeAndHeadersComplete(ctx context.Context, method, uri string, params interface{}, headers http.Header) (*APIResponse, error) {
	var err error
	var resp *http.Response
	var respErr error
	var respBody []byte

	for i := 0; i <= api.retryPolicy.MaxRetries; i++ {
		var reqBody io.Reader
		if params != nil {
			if r, ok := params.(io.Reader); ok {
				reqBody = r
			} else if paramBytes, ok := params.([]byte); ok {
				reqBody = bytes.NewReader(paramBytes)
			} else {
				var jsonBody []byte
				jsonBody, err = json.Marshal(params)
				if err != nil {
					return nil, fmt.Errorf("error marshalling params to JSON: %w", err)
				}
				reqBody = bytes.NewReader(jsonBody)
			}
		}

		if i > 0 {
			// expect the backoff introduced here on errored requests to dominate the effect of rate limiting
			// don't need a random component here as the rate limiter should do something similar
			// nb time duration could truncate an arbitrary float. Since our inputs are all ints, we should be ok
			sleepDuration := time.Duration(math.Pow(2, float64(i-1)) * float64(api.retryPolicy.MinRetryDelay))

			if sleepDuration > api.retryPolicy.MaxRetryDelay {
				sleepDuration = api.retryPolicy.MaxRetryDelay
			}
			// useful to do some simple logging here, maybe introduce levels later
			api.logger.Printf("Sleeping %s before retry attempt number %d for request %s %s", sleepDuration.String(), i, method, uri)

			select {
			case <-time.After(sleepDuration):
			case <-ctx.Done():
				return nil, fmt.Errorf("operation aborted during backoff: %w", ctx.Err())
			}
		}

		err = api.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, fmt.Errorf("error caused by request rate limiting: %w", err)
		}

		resp, respErr = api.request(ctx, method, uri, reqBody, headers)

		// short circuit processing on context timeouts
		if respErr != nil && errors.Is(respErr, context.DeadlineExceeded) {
			return nil, respErr
		}

		// retry if the server is rate limiting us or if it failed
		// assumes server operations are rolled back on failure
		if respErr != nil || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
				respErr = errors.New("exceeded available rate limit retries")
			}

			if respErr == nil {
				respErr = fmt.Errorf("received %s response (HTTP %d), please try again later", strings.ToLower(http.StatusText(resp.StatusCode)), resp.StatusCode)
			}
			continue
		} else {
			respBody, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("could not read response body: %w", err)
			}

			break
		}
	}

	// still had an error after all retries
	if respErr != nil {
		return nil, respErr
	}

	if resp.StatusCode >= http.StatusBadRequest {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, &ServiceError{controldError: &Error{
				StatusCode: resp.StatusCode,
				Error: ResponseInfo{
					Message: errInternalServiceError,
				},
			}}
		}

		errBody := &Response{}
		err = json.Unmarshal(respBody, &errBody)
		if err != nil {
			return nil, fmt.Errorf(errUnmarshalErrorBody+": %w", err)
		}

		err := &Error{
			StatusCode: resp.StatusCode,
			Error:      errBody.Error,
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			err.Type = ErrorTypeAuthorization
			return nil, &AuthorizationError{controldError: err}
		case http.StatusForbidden:
			err.Type = ErrorTypeAuthentication
			return nil, &AuthenticationError{controldError: err}
		case http.StatusNotFound:
			err.Type = ErrorTypeNotFound
			return nil, &NotFoundError{controldError: err}
		case http.StatusTooManyRequests:
			err.Type = ErrorTypeRateLimit
			return nil, &RatelimitError{controldError: err}
		default:
			err.Type = ErrorTypeRequest
			return nil, &RequestError{controldError: err}
		}
	}

	return &APIResponse{
		Body:       respBody,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
	}, nil
}

// request makes a HTTP request to the given API endpoint, returning the raw
// *http.Response, or an error if one occurred. The caller is responsible for
// closing the response body.
func (api *API) request(ctx context.Context, method, uri string, reqBody io.Reader, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, api.BaseURL+uri, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed: %w", err)
	}

	combinedHeaders := make(http.Header)
	copyHeader(combinedHeaders, api.headers)
	copyHeader(combinedHeaders, headers)
	req.Header = combinedHeaders

	req.Header.Set("Authorization", "Bearer "+api.APIToken)

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if api.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}

		log.Printf("\n%s", string(dump))
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if api.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}
		log.Printf("\n%s", string(dump))
	}

	return resp, nil
}

// copyHeader copies all headers for `source` and sets them on `target`.
// based on https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}

type DateTime struct {
	time.Time
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(dt.Format(time.RFC3339)), nil
}
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var dateTimeStr = string(data[1 : len(data)-1])
	parse, err := time.Parse(time.RFC1123Z, dateTimeStr)
	if err != nil {
		return err
	}
	dt.Time = parse
	return nil
}

// ResponseInfo contains a code and message returned by the API as errors or
// informational messages inside the response.
type ResponseInfo struct {
	Date    DateTime `json:"date"`
	Message string   `json:"message"`
	Code    int      `json:"code"`
}

// Response is a template.  There will also be a result struct.  There will be a
// unique response type for each response, which will include this type.
type Response struct {
	Success bool         `json:"success"`
	Error   ResponseInfo `json:"error"`
}

// RawResponse keeps the result as JSON form.
type RawResponse struct {
	Response
	Body json.RawMessage `json:"body"`
}

// Raw makes an HTTP request with user provided params and returns the
// result as a RawResponse, which contains the untouched JSON result.
func (api *API) Raw(ctx context.Context, method, endpoint string, data interface{}, headers http.Header) (RawResponse, error) {
	var r RawResponse
	res, err := api.makeRequestContextWithHeaders(ctx, method, endpoint, data, headers)
	if err != nil {
		return r, err
	}

	if err := json.Unmarshal(res, &r); err != nil {
		return r, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r, nil
}

// RetryPolicy specifies number of retries and min/max retry delays
// This config is used when the client exponentially backs off after errored requests.
type RetryPolicy struct {
	MaxRetries    int
	MinRetryDelay time.Duration
	MaxRetryDelay time.Duration
}

// Logger defines the interface this library needs to use logging
// This is a subset of the methods implemented in the log package.
type Logger interface {
	Printf(format string, v ...interface{})
}

// ReqOption is a functional option for configuring API requests.
type ReqOption func(opt *reqOption)

type reqOption struct {
	//nolint:unused
	params url.Values
}
