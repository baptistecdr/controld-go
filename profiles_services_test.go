package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListProfileServices(t *testing.T) {
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
					"category": "social",
					"name": "4chan",
					"PK": "4chan",
					"action": {
					  "do": 0,
					  "status": 1
					}
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileServicesParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/services", params.ProfileID), handler)
	actual, err := client.ListProfileServices(context.Background(), params)

	warning := ""
	want := []ProfileService{
		{
			PK:             "4chan",
			Name:           "4chan",
			Category:       "social",
			UnlockLocation: "JFK",
			Locations:      nil,
			Action: Action{
				Do:     Block,
				Status: IntBool(true),
				Via:    nil,
				ViaV6:  nil,
			},
			Warning: &warning,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileServices(context.Background(), ListProfileServicesParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Services should not have been listed")
}

func TestUpdateProfileService(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"services": [
				  {
					"do": 2,
					"via": "google.com",
					"status": 1,
					"via_v6": "google.com"
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	via := "google.com"
	params := UpdateProfileServiceParams{
		ProfileID: "profileID",
		Service:   "4chan",
		Do:        Spoof,
		Status:    IntBool(true),
		Via:       &via,
		ViaV6:     &via,
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/services/%s", params.ProfileID, params.Service), handler)
	actual, err := client.UpdateProfileService(context.Background(), params)

	want := []Action{
		{
			Do:     Spoof,
			Status: IntBool(true),
			Via:    &via,
			ViaV6:  &via,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfileService(context.Background(), UpdateProfileServiceParams{
		ProfileID: "",
		Service:   "4chan",
	})
	require.Error(t, err, "Profile Service should not have been updated")

	_, err = client.UpdateProfileService(context.Background(), UpdateProfileServiceParams{
		ProfileID: "profileID",
		Service:   "",
	})
	require.Error(t, err, "Profile Service should not have been updated")
}
