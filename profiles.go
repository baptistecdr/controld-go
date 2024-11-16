package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Profile struct {
	PK      string   `json:"PK"`
	Updated UnixTime `json:"updated"`
	Name    string   `json:"name"`
}

type ListProfilesBody struct {
	Profiles []Profile `json:"profiles"`
}

type ListProfilesResponse struct {
	Body ListProfilesBody `json:"body"`
	Response
}

type CreateProfileParams struct {
	Name           string  `json:"name"`
	CloneProfileID *string `json:"clone_profile_id"`
}

type CreateProfileBody struct {
	Profiles []Profile `json:"profiles"`
}

type CreateProfileResponse struct {
	Body ListProfilesBody `json:"body"`
	Response
}

type UpdateProfileParams struct {
	ProfileID   string   `json:"profile_id"`
	Name        *string  `json:"name"`
	DisableTTL  *int     `json:"disable_ttl"`
	LockStatus  *IntBool `json:"lock_status"`
	LockMessage *string  `json:"lock_message"`
	Password    *string  `json:"password"`
}

type UpdateProfileBody struct {
	Profiles []Profile `json:"profiles"`
}

type UpdateProfileResponse struct {
	Body ListProfilesBody `json:"body"`
	Response
}

type DeleteProfileParams struct {
	ProfileID string `json:"profile_id"`
}

type DeleteProfileResponse struct {
	Body    []any  `json:"body"`
	Message string `json:"message"`
	Response
}

type ProfileOptionType string

const (
	Dropdown ProfileOptionType = "dropdown"
	Field    ProfileOptionType = "field"
	Toggle   ProfileOptionType = "toggle"
)

type ProfilesOption struct {
	PK           string            `json:"PK"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Type         ProfileOptionType `json:"type"`
	DefaultValue any               `json:"default_value"`
	InfoURL      string            `json:"info_url"`
}

type ListProfilesOptionsBody struct {
	Options []ProfilesOption `json:"options"`
}

type ListProfilesOptionsResponse struct {
	Body ListProfilesOptionsBody `json:"body"`
	Response
}

type UpdateProfilesOption struct {
	ProfileID string  `json:"profile_id"`
	Name      string  `json:"name"`
	Status    IntBool `json:"status"`
	Value     *string `json:"value"`
}

type UpdateProfilesOptionBody struct {
	Options any `json:"options"`
	Response
}

type UpdateProfilesOptionResponse struct {
	Body UpdateProfilesOptionBody `json:"body"`
	Response
}

func (api *API) ListProfiles(ctx context.Context) ([]Profile, error) {
	uri := buildURI("/profiles", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfilesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Profiles, nil
}

func (api *API) CreateProfile(ctx context.Context, params CreateProfileParams) ([]Profile, error) {
	uri := buildURI("/profiles", nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r CreateProfileResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Profiles, nil
}

func (api *API) UpdateProfile(ctx context.Context, params UpdateProfileParams) ([]Profile, error) {
	if params.ProfileID == "" {
		return []Profile{}, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Profile{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Profiles, nil
}

func (api *API) DeleteProfile(ctx context.Context, params DeleteProfileParams) ([]any, error) {
	if params.ProfileID == "" {
		return []any{}, fmt.Errorf("delete: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, params)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteProfileResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) ListProfilesOptions(ctx context.Context) ([]ProfilesOption, error) {
	uri := buildURI("/profiles/options", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []ProfilesOption{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfilesOptionsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []ProfilesOption{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Options, nil
}

func (api *API) UpdateProfilesOption(ctx context.Context, params UpdateProfilesOption) (any, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("update: no profile ID provided")
	}
	if params.Name == "" {
		return nil, fmt.Errorf("update: no profile options name provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/options/%s", params.ProfileID, params.Name)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfilesOptionResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Options, nil
}
