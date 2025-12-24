package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListProfileDefaultRule(t *testing.T) {
	setup()
	defer teardown()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"default": []
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileDefaultRuleParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/default", params.ProfileID), handler)
	var actual, err = client.ListProfileDefaultRule(context.Background(), params)

	want := DefaultRule{
		Do:     Bypass,
		Status: IntBool(true),
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	setup()

	handler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"default": {
				  "do": 1,
				  "status": 1
				}
			  },
			  "success": true
			}
		`)
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/default", params.ProfileID), handler)
	actual, err = client.ListProfileDefaultRule(context.Background(), params)
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileDefaultRule(context.Background(), ListProfileDefaultRuleParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Default Rule should not have been listed")
}

func TestUpdateProfileDefaultRule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"default": {
				  "do": 0,
				  "status": 1
				}
			  },
			  "success": true
			}
		`)
	}

	params := UpdateProfileDefaultRuleParams{
		ProfileID: "profileID",
		Do:        Block,
		Status:    IntBool(true),
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/default", params.ProfileID), handler)
	actual, err := client.UpdateProfileDefaultRule(context.Background(), params)

	want := DefaultRule{
		Do:     Block,
		Status: IntBool(true),
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfileDefaultRule(context.Background(), UpdateProfileDefaultRuleParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Default Rule should not have been updated")
}
