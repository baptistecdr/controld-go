package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DoType int

const (
	Block    = 0
	Bypass   = 1
	Spoof    = 2
	Redirect = 3
)

type ProfileService struct {
	PK             string   `json:"PK"`
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	UnlockLocation string   `json:"unlock_location"`
	Locations      []string `json:"locations,omitempty"`
	Action         Action   `json:"action"`
	Warning        *string  `json:"warning,omitempty"`
}

type Action struct {
	Do     DoType  `json:"do"`
	Status IntBool `json:"status"`
	Via    *string `json:"via,omitempty"`
	ViaV6  *string `json:"via_v6,omitempty"`
	Group  *int    `json:"group,omitempty"`
	Order  *int    `json:"order,omitempty"`
}

type ListProfileServicesParams struct {
	ProfileID string `json:"profile_id"`
}

type ListProfileServicesBody struct {
	Services []ProfileService `json:"services"`
}

type ListProfileServicesResponse struct {
	Body ListProfileServicesBody `json:"body"`
	Response
}

type UpdateProfileServiceParams struct {
	ProfileID string  `json:"profile_id"`
	Service   string  `json:"service"`
	Do        DoType  `json:"do"`
	Status    IntBool `json:"status"`
	Via       *string `json:"via"`
	ViaV6     *string `json:"via_v6"`
}

type UpdateProfileServiceBody struct {
	Services []Action `json:"services"`
}

type UpdateProfileServiceResponse struct {
	Body UpdateProfileServiceBody `json:"body"`
	Response
}

type DeleteProfileServiceParams struct {
	ProfileID string `json:"profile_id"`
	Hostname  string `json:"hostname"`
}

type DeleteProfileServiceBody struct {
	Services []Action `json:"services"`
}

type DeleteProfileServiceResponse struct {
	Body DeleteProfileServiceBody `json:"body"`
	Response
}

func (api *API) ListProfileServices(ctx context.Context, params ListProfileServicesParams) ([]ProfileService, error) {
	if params.ProfileID == "" {
		return []ProfileService{}, fmt.Errorf("list: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/services", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []ProfileService{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileServicesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []ProfileService{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Services, nil
}

func (api *API) UpdateProfileService(ctx context.Context, params UpdateProfileServiceParams) ([]Action, error) {
	if params.ProfileID == "" {
		return []Action{}, fmt.Errorf("update: no profile ID provided")
	}
	if params.Service == "" {
		return []Action{}, fmt.Errorf("update: no service provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/services/%s", params.ProfileID, params.Service)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return []Action{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileServiceResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Action{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Services, nil
}
