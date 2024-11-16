package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogLevel struct {
	PK    AnalyticsLevel `json:"PK"`
	Title string         `json:"title"`
}

type ListLogLevelsBody struct {
	Levels []LogLevel `json:"levels"`
}

type ListLogLevelsResponse struct {
	Body ListLogLevelsBody `json:"body"`
	Response
}

type Endpoint struct {
	PK          string `json:"PK"`
	Title       string `json:"title"`
	CountryCode string `json:"country_code"`
}

type ListStorageRegionsBody struct {
	Endpoint []Endpoint `json:"endpoints"`
}

type ListStorageRegionsResponse struct {
	Body ListStorageRegionsBody `json:"body"`
	Response
}

func (api *API) ListLogLevels(ctx context.Context) ([]LogLevel, error) {
	uri := buildURI("/analytics/levels", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []LogLevel{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListLogLevelsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []LogLevel{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Levels, nil
}

func (api *API) ListStorageRegions(ctx context.Context) ([]Endpoint, error) {
	uri := buildURI("/analytics/endpoints", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Endpoint{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListStorageRegionsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Endpoint{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Endpoint, nil
}
