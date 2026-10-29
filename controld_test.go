package controld

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the API client being tested.
	client *API

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup(opts ...Option) {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// disable rate limits and retries in testing - prepended so any provided value overrides this
	opts = append([]Option{UsingRateLimit(100000), UsingRetryPolicy(0, 0, 0)}, opts...)

	// Cloudflare client configured to use test server
	client, _ = New("api.1377", opts...)
	client.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

func TestRateLimitRetriesExhaustedKeepsTypedError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})

	_, err := client.makeRequestContext(context.Background(), http.MethodGet, "/devices", nil)

	var rateLimitErr *RatelimitError
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &rateLimitErr, "errors.As must reach the typed error once retries are exhausted")
		assert.Equal(t, ErrorTypeRateLimit, rateLimitErr.Type())
		assert.Equal(t, "exceeded available rate limit retries", err.Error(),
			"the message must stay byte-identical for callers matching on text")
	}
}

func TestServerErrorRetriesExhaustedKeepsTypedError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	_, err := client.makeRequestContext(context.Background(), http.MethodGet, "/devices", nil)

	var serviceErr *ServiceError
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &serviceErr, "errors.As must reach the typed error once retries are exhausted")
		assert.Equal(t, "received service unavailable response (HTTP 503), please try again later", err.Error(),
			"the message must stay byte-identical for callers matching on text")
	}
}
