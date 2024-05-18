package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestListKnownIPs(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"ips": [
				  {
					"ip": "114.157.113.63",
					"ts": 1731235525,
					"country": "FR",
					"city": "Saint-Quentin-Fallavier",
					"isp": "ISP SAS",
					"asn": 14310,
					"as_name": "ISP SAS"
				  }
				]
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/access", handler)
	actual, err := client.ListKnownIPs(context.Background(), ListKnownIPsParams{DeviceID: "deviceID"})

	want := []KnownIP{{
		IP:      net.ParseIP("114.157.113.63"),
		Ts:      UnixTime{time.Unix(1731235525, 0).UTC()},
		Country: "FR",
		City:    "Saint-Quentin-Fallavier",
		ISP:     "ISP SAS",
		Asn:     14310,
		AsName:  "ISP SAS",
	}}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestLearnNewIPs(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": [],
			  "success": true,
			  "message": "1 IPs added"
			}
		`)
	}
	mux.HandleFunc("/access", handler)
	actual, err := client.LearnNewIPs(context.Background(), LearnNewIPsParams{
		DeviceID: "deviceID",
		IPs:      []net.IP{net.ParseIP("114.157.113.63")},
	})

	want := []any{}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestDeleteLearnedIPs(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": [],
			  "success": true,
			  "message": "1 IPs deleted"
			}
		`)
	}
	mux.HandleFunc("/access", handler)
	actual, err := client.DeleteLearnedIPs(context.Background(), DeleteLearnedIPsParams{
		DeviceID: "deviceID",
		IPs:      []net.IP{net.ParseIP("114.157.113.63")},
	})

	want := []any{}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
