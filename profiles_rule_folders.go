package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GroupAction struct {
	Status IntBool `json:"status"`
	Do     *DoType `json:"do,omitempty"`
}

type Group struct {
	PK     int         `json:"PK"`
	Group  string      `json:"group"`
	Action GroupAction `json:"action"`
	Count  int         `json:"count"`
}

type ListProfileRuleFoldersParams struct {
	ProfileID string `json:"profile_id"`
}

type ListProfileRuleFoldersBody struct {
	Groups []Group `json:"groups"`
}

type ListProfileRuleFoldersResponse struct {
	Body ListProfileRuleFoldersBody `json:"body"`
	Response
}

type CreateProfileRuleFolderParams struct {
	ProfileID string   `json:"profile_id"`
	Name      string   `json:"name"`
	Do        *DoType  `json:"do,omitempty"`
	Via       *string  `json:"via,omitempty"`
	Status    *IntBool `json:"status,omitempty"`
}

type CreateProfileRuleFolderBody struct {
	Groups []Group `json:"groups"`
}

type CreateProfileRuleFolderResponse struct {
	Body CreateProfileRuleFolderBody `json:"body"`
	Response
}

type UpdateProfileRuleFolderParams struct {
	ProfileID string   `json:"profile_id"`
	FolderID  string   `json:"folder"`
	Do        *DoType  `json:"do,omitempty"`
	Via       *string  `json:"via,omitempty"`
	Status    *IntBool `json:"status,omitempty"`
}

type UpdateProfileRuleFolderBody struct {
	Groups []Group `json:"groups"`
}

type UpdateProfileRuleFolderResponse struct {
	Body CreateProfileRuleFolderBody `json:"body"`
	Response
}

type DeleteProfileRuleFolderParams struct {
	ProfileID string `json:"profile_id"`
	FolderID  string `json:"folder"`
}

type DeleteProfileRuleFolderResponse struct {
	Body any `json:"body"`
	Response
}

func (api *API) ListProfileRuleFolders(ctx context.Context, params ListProfileRuleFoldersParams) ([]Group, error) {
	if params.ProfileID == "" {
		return []Group{}, fmt.Errorf("list: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/groups", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileRuleFoldersResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Groups, nil
}

func (api *API) CreateProfileRuleFolder(ctx context.Context, params CreateProfileRuleFolderParams) ([]Group, error) {
	if params.ProfileID == "" {
		return []Group{}, fmt.Errorf("create: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/groups", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r CreateProfileRuleFolderResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Groups, nil
}

func (api *API) UpdateProfileRuleFolder(ctx context.Context, params UpdateProfileRuleFolderParams) ([]Group, error) {
	if params.ProfileID == "" {
		return []Group{}, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/groups/%s", params.ProfileID, params.FolderID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileRuleFolderResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Group{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Groups, nil
}

func (api *API) DeleteProfileRuleFolder(ctx context.Context, params DeleteProfileRuleFolderParams) (any, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("delete: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/groups/%s", params.ProfileID, params.FolderID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteProfileRuleFolderResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
