package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListServiceCategories(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"categories": [
				  {
					"PK": "audio",
					"name": "Audio",
					"description": "Audio streaming services and radio stations",
					"count": 16
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/services/categories", handler)
	actual, err := client.ListServiceCategories(context.Background())

	want := []Category{
		{
			PK:          "audio",
			Name:        "Audio",
			Description: "Audio streaming services and radio stations",
			Count:       16,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListServices(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"services": [
				  {
					"warning": "",
					"unlock_location": "JFK",
					"category": "tools",
					"name": "1Password",
					"PK": "1password"
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	params := ListServicesParams{
		Category: "tools",
	}
	mux.HandleFunc(fmt.Sprintf("/services/categories/%s", params.Category), handler)
	actual, err := client.ListServices(context.Background(), params)

	warning := ""
	want := []Service{
		{
			PK:             "1password",
			Name:           "1Password",
			Category:       "tools",
			UnlockLocation: "JFK",
			Warning:        &warning,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListServices(context.Background(), ListServicesParams{
		Category: "",
	})
	require.Error(t, err, "Services should not have been listed")
}
