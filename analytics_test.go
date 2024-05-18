package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListLogLevels(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"levels": [
				  {
					"PK": 0,
					"title": "No Analytics"
				  },
				  {
					"PK": 1,
					"title": "Some Analytics"
				  },
				  {
					"PK": 2,
					"title": "Full Analytics"
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/analytics/levels", handler)
	actual, err := client.ListLogLevels(context.Background())

	want := []LogLevel{
		{
			PK:    Off,
			Title: "No Analytics",
		},
		{
			PK:    Basic,
			Title: "Some Analytics",
		},
		{
			PK:    Full,
			Title: "Full Analytics",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListStorageRegions(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"endpoints": [
				  {
					"PK": "america",
					"title": "New York, US",
					"country_code": "US"
				  },
				  {
					"PK": "europe",
					"title": "Amsterdam, NL",
					"country_code": "NL"
				  },
				  {
					"PK": "asia",
					"title": "Sydney, AU",
					"country_code": "AU"
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/analytics/endpoints", handler)
	actual, err := client.ListStorageRegions(context.Background())

	want := []Endpoint{
		{
			PK:          "america",
			Title:       "New York, US",
			CountryCode: "US",
		},
		{
			PK:          "europe",
			Title:       "Amsterdam, NL",
			CountryCode: "NL",
		},
		{
			PK:          "asia",
			Title:       "Sydney, AU",
			CountryCode: "AU",
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
