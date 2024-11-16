package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomRule Action

type Rule struct {
	PK     string `json:"PK"`
	Order  int    `json:"order"`
	Group  int    `json:"group"`
	Action Action `json:"action"`
}

type ListProfileCustomRulesParams struct {
	ProfileID string `json:"profile_id"`
	FolderID  string `json:"folder_id"`
}

type ListProfileCustomRulesBody struct {
	Rules []Rule `json:"rules"`
}

type ListProfileCustomRulesResponse struct {
	Body ListProfileCustomRulesBody `json:"body"`
	Response
}

type CreateProfileCustomRuleParams struct {
	ProfileID string   `json:"profile_id"`
	Do        DoType   `json:"do"`
	Status    IntBool  `json:"status"`
	Via       *string  `json:"via,omitempty"`
	ViaV6     *string  `json:"via_v6,omitempty"`
	Group     *int     `json:"group,omitempty"`
	Hostnames []string `json:"hostnames"`
}

type CreateProfileCustomRuleBody struct {
	Rules []CustomRule `json:"rules"`
}

type CreateProfileCustomRuleResponse struct {
	Body CreateProfileCustomRuleBody `json:"body"`
	Response
}

type UpdateProfileCustomRuleParams struct {
	ProfileID string   `json:"profile_id"`
	Do        DoType   `json:"do"`
	Status    IntBool  `json:"status"`
	Via       *string  `json:"via,omitempty"`
	ViaV6     *string  `json:"via_v6,omitempty"`
	Group     *int     `json:"group,omitempty"`
	Hostnames []string `json:"hostnames"`
}

type UpdateProfileCustomRuleBody struct {
	Rules []CustomRule `json:"rules"`
}

type UpdateProfileCustomRuleResponse struct {
	Body UpdateProfileCustomRuleBody `json:"body"`
	Response
}

type DeleteProfileCustomRuleParams struct {
	ProfileID string `json:"profile_id"`
	Hostname  string `json:"hostname"`
}

type DeleteProfileCustomRuleResponse struct {
	Body any `json:"body"`
	Response
}

func (api *API) ListProfileCustomRules(ctx context.Context, params ListProfileCustomRulesParams) ([]Rule, error) {
	if params.ProfileID == "" {
		return []Rule{}, fmt.Errorf("list: no profile ID provided")
	}
	if params.FolderID == "" {
		return []Rule{}, fmt.Errorf("list: no folder ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/rules/%s", params.ProfileID, params.FolderID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Rule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileCustomRulesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Rule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Rules, nil
}

func (api *API) CreateProfileCustomRule(ctx context.Context, params CreateProfileCustomRuleParams) ([]CustomRule, error) {
	if params.ProfileID == "" {
		return []CustomRule{}, fmt.Errorf("create: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/rules", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return []CustomRule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r CreateProfileCustomRuleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []CustomRule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Rules, nil
}

func (api *API) UpdateProfileCustomRule(ctx context.Context, params UpdateProfileCustomRuleParams) ([]CustomRule, error) {
	if params.ProfileID == "" {
		return []CustomRule{}, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/rules", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return []CustomRule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileCustomRuleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []CustomRule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Rules, nil
}

func (api *API) DeleteProfileCustomRule(ctx context.Context, params DeleteProfileCustomRuleParams) (any, error) {
	if params.ProfileID == "" {
		return nil, fmt.Errorf("delete: no profile ID provided")
	}
	if params.Hostname == "" {
		return nil, fmt.Errorf("delete: no hostname provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/rules/%s", params.ProfileID, params.Hostname)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteProfileCustomRuleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
