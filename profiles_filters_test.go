package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
)

func TestListProfilesNativeFilters(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"filters": [
				  {
					"PK": "ads",
					"name": "Ads & Trackers",
					"description": "One of the most comprehensive ad and tracker blocking lists out there derived from over 2 dozen block lists, hand curated to remove thousands of false positives that plague most community maintained lists.",
					"additional": "<p><strong>Relaxed Mode<\/strong><\/p>Smallest block list, will only block very common ads and trackers. Least prone to false positives. <p><strong>Balanced Mode<\/strong><\/p>Medium sized block list, will block a lot more than the Relaxed mode, but still allow common affiliate and email tracking links. <p><strong>Strict Mode<\/strong><\/p>Will block the most ads and trackers, but also has the highest chance of false positives.",
					"sources": [],
					"levels": [
					  {
						"title": "Relaxed",
						"type": "filter",
						"name": "ads_small",
						"status": 0
					  }
					],
					"status": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileFiltersParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/filters", params.ProfileID), handler)
	actual, err := client.ListProfileNativeFilters(context.Background(), params)

	additional := "<p><strong>Relaxed Mode</strong></p>Smallest block list, will only block very common ads and trackers. Least prone to false positives. <p><strong>Balanced Mode</strong></p>Medium sized block list, will block a lot more than the Relaxed mode, but still allow common affiliate and email tracking links. <p><strong>Strict Mode</strong></p>Will block the most ads and trackers, but also has the highest chance of false positives."
	want := []Filter{
		{
			PK:          "ads",
			Name:        "Ads & Trackers",
			Description: "One of the most comprehensive ad and tracker blocking lists out there derived from over 2 dozen block lists, hand curated to remove thousands of false positives that plague most community maintained lists.",
			Additional:  &additional,
			Sources:     []string{},
			Levels: []FilterLevel{
				{
					Title:  "Relaxed",
					Type:   "filter",
					Name:   "ads_small",
					Status: IntBool{false},
					Opt:    nil,
				},
			},
			Status:    IntBool{false},
			Resolvers: nil,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileNativeFilters(context.Background(), ListProfileFiltersParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Filters should not have been listed")
}

func TestListProfilesExternalFilters(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"filters": [
				  {
					"PK": "x-1hosts-lite",
					"name": "1Hosts - Lite",
					"description": "World's most advanced general-purpose DNS filter\/blocklists!",
					"additional": "<p><strong>Warning<\/strong><\/p>This is a community list and is not maintained by Control D. Any false positives\/negatives should be reported to the list maintainer.",
					"sources": ["https:\/\/github.com\/badmojr\/1Hosts"],
					"resolvers": {
					  "v4": ["244.161.205.7"],
					  "v6": ["5e47:583e:8421:ff1c:d410:abaf:3787:a6e2"]
					},
					"status": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileFiltersParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/filters/external", params.ProfileID), handler)
	actual, err := client.ListProfileExternalFilters(context.Background(), params)

	additional := "<p><strong>Warning</strong></p>This is a community list and is not maintained by Control D. Any false positives/negatives should be reported to the list maintainer."
	resolvers := FilterResolvers{
		V4: []net.IP{net.ParseIP("244.161.205.7")},
		V6: []net.IP{net.ParseIP("5e47:583e:8421:ff1c:d410:abaf:3787:a6e2")},
	}
	want := []Filter{
		{
			PK:          "x-1hosts-lite",
			Name:        "1Hosts - Lite",
			Description: "World's most advanced general-purpose DNS filter/blocklists!",
			Additional:  &additional,
			Sources:     []string{"https://github.com/badmojr/1Hosts"},
			Levels:      nil,
			Status:      IntBool{false},
			Resolvers:   &resolvers,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileExternalFilters(context.Background(), ListProfileFiltersParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Filters should not have been listed")
}

func TestUpdateProfileFilter(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"filters": {
				  "x-1hosts-lite": {
					"do": 0,
					"status": 1
				  }
				}
			  },
			  "success": true
			}
		`)
	}

	params := UpdateProfileFilterParams{
		ProfileID: "profileID",
		Filter:    "x-1hosts-lite",
		Status:    IntBool{true},
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/filters/filter/%s", params.ProfileID, params.Filter), handler)
	actual, err := client.UpdateProfileFilter(context.Background(), params)

	want := map[string]interface{}{
		"x-1hosts-lite": map[string]interface{}{
			"do":     float64(0),
			"status": float64(1),
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfileFilter(context.Background(), UpdateProfileFilterParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Filter should not have been updated")

	_, err = client.UpdateProfileFilter(context.Background(), UpdateProfileFilterParams{
		ProfileID: "profileID",
		Filter:    "",
	})
	require.Error(t, err, "Profile Filter should not have been updated")
}
