package controld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type DefaultRule Action

type ListProfileDefaultRuleParams struct {
	ProfileID string `json:"profile_id"`
}

type ListProfileDefaultRuleBody struct {
	Default any `json:"default"`
}

type ListProfileDefaultRuleResponse struct {
	Body ListProfileDefaultRuleBody `json:"body"`
	Response
}

type UpdateProfileDefaultRuleParams struct {
	ProfileID string  `json:"profile_id"`
	Do        DoType  `json:"do"`
	Status    IntBool `json:"status"`
	Via       *string `json:"via,omitempty"`
}

type UpdateProfileDefaultRuleBody struct {
	Default DefaultRule `json:"default"`
}

type UpdateProfileDefaultRuleResponse struct {
	Body UpdateProfileDefaultRuleBody `json:"body"`
	Response
}

func (api *API) ListProfileDefaultRule(ctx context.Context, params ListProfileDefaultRuleParams) (DefaultRule, error) {
	if params.ProfileID == "" {
		return DefaultRule{}, fmt.Errorf("list: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/default", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return DefaultRule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProfileDefaultRuleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return DefaultRule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	switch rule := r.Body.Default.(type) {
	case map[string]any:
		jsonDefaultRule, err := json.Marshal(rule)
		if err != nil {
			return DefaultRule{}, fmt.Errorf("%s: %w", errMarshalError, err)
		}
		var defaultRule DefaultRule
		err = json.Unmarshal(jsonDefaultRule, &defaultRule)
		if err != nil {
			return DefaultRule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
		}
		return defaultRule, nil
	// ControlD return an empty array when the Default Rule has never been modified
	case []any:
		return DefaultRule{
			Do:     Bypass,
			Status: IntBool(true),
		}, nil
	default:
		return DefaultRule{}, fmt.Errorf("%s: %w", errTypeError, errors.New("type: unknown field type"))
	}
}

func (api *API) UpdateProfileDefaultRule(ctx context.Context, params UpdateProfileDefaultRuleParams) (DefaultRule, error) {
	if params.ProfileID == "" {
		return DefaultRule{}, fmt.Errorf("update: no profile ID provided")
	}
	baseURL := fmt.Sprintf("/profiles/%s/default", params.ProfileID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return DefaultRule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateProfileDefaultRuleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return DefaultRule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Default, nil
}
