package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

func TestListIP(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"ip": "23.251.148.254",
				"type": "v4",
				"org": "Org SAS",
				"country": "FR",
				"handler": "cdg-h01",
				"pop": "CDG"
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/ip", handler)
	actual, err := client.ListIP(context.Background())

	want := IP{
		IP:      net.ParseIP("23.251.148.254"),
		Type:    "v4",
		Org:     "Org SAS",
		Country: "FR",
		Handler: "cdg-h01",
		Pop:     "CDG",
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListNetwork(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"network": [
				  {
					"iata_code": "AMS",
					"city_name": "Amsterdam",
					"country_name": "NL",
					"location": {
					  "lat": 52.377956,
					  "long": 4.89707
					},
					"status": {
					  "api": 1,
					  "dns": 1,
					  "pxy": 1
					}
				  }
				],
				"time": 1716095598,
				"current_pop": "CDG"
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/network", handler)
	actual, err := client.ListNetwork(context.Background())

	want := []Network{
		{
			IataCode:    "AMS",
			CityName:    "Amsterdam",
			CountryName: "NL",
			Location: Location{
				Lat:  52.377956,
				Long: 4.89707,
			},
			Status: Status{
				API: 1,
				DNS: 1,
				Pxy: 1,
			},
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
