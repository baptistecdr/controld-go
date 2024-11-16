package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Category struct {
	PK          string `json:"PK"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

type ListServiceCategoriesBody struct {
	Categories []Category `json:"categories"`
}

type ListServiceCategoriesResponse struct {
	Body ListServiceCategoriesBody `json:"body"`
	Response
}

type ListServicesParams struct {
	Category string `json:"category"`
}

type Service struct {
	PK             string  `json:"PK"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	UnlockLocation string  `json:"unlock_location"`
	Warning        *string `json:"warning,omitempty"`
}

type ListServicesBody struct {
	Services []Service `json:"services"`
}

type ListServicesResponse struct {
	Body ListServicesBody `json:"body"`
	Response
}

func (api *API) ListServiceCategories(ctx context.Context) ([]Category, error) {
	uri := buildURI("/services/categories", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Category{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListServiceCategoriesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Category{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Categories, nil
}

func (api *API) ListServices(ctx context.Context, params ListServicesParams) ([]Service, error) {
	if params.Category == "" {
		return []Service{}, fmt.Errorf("list: no category provided")
	}
	baseURL := fmt.Sprintf("/services/categories/%s", params.Category)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Service{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListServicesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Service{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Services, nil
}
