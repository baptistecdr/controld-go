package controld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Filter struct {
	PK          string           `json:"PK"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Additional  *string          `json:"additional,omitempty"`
	Sources     []string         `json:"sources"`
	Levels      []FilterLevel    `json:"levels,omitempty"`
	Status      IntBool          `json:"status"`
	Resolvers   *FilterResolvers `json:"resolvers,omitempty"`
}

type FilterLevel struct {
	Title  string  `json:"title"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Status IntBool `json:"status"`
	Opt    []Opt   `json:"opt,omitempty"`
}

type Opt struct {
	PK    string `json:"PK"`
	Value any    `json:"value"`
}

type FilterResolvers struct {
	V4 []net.IP `json:"v4"`
	V6 []net.IP `json:"v6"`
}

type ListProfileFiltersParams struct {
	ProfileID string `json:"profile_id"`
}

type ListProfileFiltersBody struct {
	Filters []Filter `json:"filters"`
}

type ListProfileFiltersResponse struct {
	Body ListProfileFiltersBody `json:"body"`
	Response
}

type UpdateProfileFilterParams struct {
	ProfileID string  `json:"profile_id"`
	Filter    string  `json:"filter"`
	Status    IntBool `json:"status"`
}

type UpdateProfileFilterBody struct {
	Filters any `json:"filters"`
}

type UpdateProfileFilterResponse struct {
	Body UpdateProfileFilterBody `json:"body"`
	Response
}

type UpdateProfileFilterEntry struct {
	Filter string  `json:"filter"`
	Status IntBool `json:"status"`
}

type UpdateProfileFiltersParams struct {
	ProfileID string                     `json:"profile_id"`
	Filters   []UpdateProfileFilterEntry `json:"filters"`
}

type FilterState struct {
	Do     DoType  `json:"do"`
	Status IntBool `json:"status"`
}

type UpdateProfileFiltersBody struct {
	// Filters is a map[string]FilterState when at least one filter is
	// reported. ControlD returns an empty array instead of an empty object
	// when there is nothing to report, so this is decoded manually.
	Filters any `json:"filters"`
}

type UpdateProfileFiltersResponse struct {
	Body UpdateProfileFiltersBody `json:"body"`
	Response
}

func (api *API) ListProfileNativeFilters(ctx context.Context, params ListProfileFiltersParams) ([]Filter, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("list: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/filters", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Filter{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileFiltersResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Filter{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Filters, nil
}

func (api *API) ListProfileExternalFilters(ctx context.Context, params ListProfileFiltersParams) ([]Filter, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("list: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/filters/external", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Filter{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileFiltersResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Filter{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Filters, nil
}

func (api *API) UpdateProfileFilter(ctx context.Context, params UpdateProfileFilterParams) (any, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/filters/filter/%s", params.ProfileID, params.Filter)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileFilterResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Filters, nil
}

func (api *API) UpdateProfileFilters(ctx context.Context, params UpdateProfileFiltersParams) (map[string]FilterState, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/filters", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileFiltersResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	switch filters := r.Body.Filters.(type) {
	case map[string]any:
		jsonFilters, err := json.Marshal(filters)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errMarshalError, err)
		}
		var result map[string]FilterState
		if err := json.Unmarshal(jsonFilters, &result); err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
		}
		return result, nil
	// ControlD returns an empty array when no filters are enabled
	case []any:
		return map[string]FilterState{}, nil
	default:
		return nil, fmt.Errorf("%s: %w", errTypeError, errors.New("type: unknown field type"))
	}
}
