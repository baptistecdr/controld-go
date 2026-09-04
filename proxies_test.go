package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListProxies(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"proxies": [
				  {
					"city": "Tirana",
					"country": "AL",
					"country_name": "Albania",
					"PK": "TIA",
					"gps_lat": 41.3275,
					"gps_long": 19.8189,
					"uid": "Albania:Tirana"
				  },
				  {
					"city": "Res Toronto",
					"country": "CA",
					"country_name": "Canada",
					"PK": "RES_YYZ",
					"gps_lat": 43.5924,
					"gps_long": -79.7611,
					"uid": "Canada:Res Toronto",
					"hidden": true
				  }
				],
				"countries": [
				  {"country": "AL", "country_name": "Albania"},
				  {"country": "CA", "country_name": "Canada"}
				]
			  },
			  "success": true
			}
		`)
	}

	mux.HandleFunc("/proxies", handler)
	proxies, countries, err := client.ListProxies(context.Background())

	wantProxies := []Proxy{
		{
			PK:          "TIA",
			UID:         "Albania:Tirana",
			City:        "Tirana",
			Country:     "AL",
			CountryName: "Albania",
			Lat:         41.3275,
			Long:        19.8189,
		},
		{
			PK:          "RES_YYZ",
			UID:         "Canada:Res Toronto",
			City:        "Res Toronto",
			Country:     "CA",
			CountryName: "Canada",
			Lat:         43.5924,
			Long:        -79.7611,
			Hidden:      true,
		},
	}
	wantCountries := []ProxyCountry{
		{Country: "AL", CountryName: "Albania"},
		{Country: "CA", CountryName: "Canada"},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, wantProxies, proxies)
		assert.Equal(t, wantCountries, countries)
	}
}
