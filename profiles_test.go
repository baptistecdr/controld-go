package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestListProfiles(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"profiles": [
				  {
					"PK": "PK",
					"updated": 1711223675,
					"name": "Default",
					"profile": {
					  "flt": {
						"count": 6
					  },
					  "cflt": {
						"count": 0
					  },
					  "ipflt": {
						"count": 1
					  },
					  "rule": {
						"count": 5
					  },
					  "svc": {
						"count": 5
					  },
					  "grp": {
						"count": 0
					  },
					  "opt": {
						"count": 0,
						"data": []
					  },
					  "da": {
						"do": 1,
						"status": 1
					  }
					}
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/profiles", handler)
	actual, err := client.ListProfiles(context.Background())

	want := []Profile{
		{
			PK:      "PK",
			Updated: UnixTime{time.Unix(1711223675, 0).UTC()},
			Name:    "Default",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestCreateProfile(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"profiles": [
				  {
					"PK": "PK",
					"updated": 1731247714,
					"name": "New Profile",
					"profile": {
					  "flt": {
						"count": 0
					  },
					  "cflt": {
						"count": 0
					  },
					  "ipflt": {
						"count": 0
					  },
					  "rule": {
						"count": 0
					  },
					  "svc": {
						"count": 0
					  },
					  "grp": {
						"count": 0
					  },
					  "opt": {
						"count": 0,
						"data": []
					  },
					  "da": []
					}
				  }
				]
			  },
			  "success": true,
			  "message": "Profile has been created"
			}
		`)
	}
	mux.HandleFunc("/profiles", handler)
	actual, err := client.CreateProfile(context.Background(), CreateProfileParams{
		Name: "New Profile",
	})

	want := []Profile{
		{
			PK:      "PK",
			Updated: UnixTime{time.Unix(1731247714, 0).UTC()},
			Name:    "New Profile",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestUpdateProfile(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"profiles": [
				  {
					"PK": "PK",
					"updated": 1731247714,
					"name": "New Profile Name",
					"profile": {
					  "flt": {
						"count": 0
					  },
					  "cflt": {
						"count": 0
					  },
					  "ipflt": {
						"count": 0
					  },
					  "rule": {
						"count": 0
					  },
					  "svc": {
						"count": 0
					  },
					  "grp": {
						"count": 0
					  },
					  "opt": {
						"count": 0,
						"data": []
					  },
					  "da": []
					}
				  }
				]
			  },
			  "success": true,
			  "message": "Profile has been created"
			}
		`)
	}
	newProfileName := "New Profile Name"
	params := UpdateProfileParams{
		ProfileID: "PK",
		Name:      &newProfileName,
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s", params.ProfileID), handler)

	actual, err := client.UpdateProfile(context.Background(), params)

	want := []Profile{
		{
			PK:      "PK",
			Updated: UnixTime{time.Unix(1731247714, 0).UTC()},
			Name:    "New Profile Name",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfile(context.Background(), UpdateProfileParams{
		ProfileID: "",
	})
	assert.Error(t, err, "Profile should not have been updated")
}

func TestDeleteProfile(t *testing.T) {
	setup()
	defer teardown()
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE", "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
		  "body": [],
		  "success": true,
		  "message": "Profile has been deleted"
	    }`)
	}
	profileID := "deviceID"
	mux.HandleFunc(fmt.Sprintf("/profiles/%s", profileID), handler)
	_, err := client.DeleteProfile(context.Background(), DeleteProfileParams{ProfileID: profileID})
	assert.NoError(t, err, "Profile should have been deleted")

	_, err = client.DeleteProfile(context.Background(), DeleteProfileParams{ProfileID: ""})
	assert.Error(t, err, "Profile should not have been deleted")
}

func TestListProfilesOptions(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"options": [
				  {
					"PK": "ai_malware",
					"title": "AI Malware Filter",
					"description": "EXPERIMENTAL: Blocks malicious domains using machine learning.",
					"type": "dropdown",
					"default_value": {
					  "0.9": "Relaxed Mode",
					  "0.7": "Balanced Mode",
					  "0.5": "Strict Mode"
					},
					"info_url": "https:\/\/docs.controld.com\/docs\/ai-malware-filter"
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/profiles/options", handler)
	actual, err := client.ListProfilesOptions(context.Background())

	var want = []ProfilesOption{
		{
			PK:          "ai_malware",
			Title:       "AI Malware Filter",
			Description: "EXPERIMENTAL: Blocks malicious domains using machine learning.",
			Type:        Dropdown,
			DefaultValue: map[string]interface{}{
				"0.9": "Relaxed Mode",
				"0.7": "Balanced Mode",
				"0.5": "Strict Mode",
			},
			InfoURL: "https://docs.controld.com/docs/ai-malware-filter",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestUpdateProfilesOption(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"options": [
				  {
					"PK": "ai_malware",
					"value": 0.9
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	value := "0.9"
	params := UpdateProfilesOption{
		ProfileID: "PK",
		Name:      "ai_malware",
		Status:    IntBool{true},
		Value:     &value,
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/options/%s", params.ProfileID, params.Name), handler)

	actual, err := client.UpdateProfilesOption(context.Background(), params)

	want := []interface{}{
		map[string]interface{}{
			"PK":    "ai_malware",
			"value": 0.9,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfilesOption(context.Background(), UpdateProfilesOption{
		ProfileID: "",
	})
	assert.Error(t, err, "Profile Option should not have been updated")

	_, err = client.UpdateProfilesOption(context.Background(), UpdateProfilesOption{
		ProfileID: "profileID",
		Name:      "",
	})
	assert.Error(t, err, "Profile Option should not have been updated")
}
