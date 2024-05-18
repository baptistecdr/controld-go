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

func TestListUser(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"resolver_ip": "76.76.2.162",
				"email_status": 1,
				"tutorials": 1,
				"v": 2,
				"resolver_status": 1,
				"rule_profile": "109129bma5vn",
				"date": "2023-01-11",
				"status": 1,
				"email": "email@icloud.com",
				"resolver_uid": "zk30ksddco",
				"proxy_access": 1,
				"stats_endpoint": "europe",
				"last_active": 1731223785,
				"PK": "119266bma1vr",
				"twofa": 1,
				"debug": []
			  },
			  "success": true
			}
		`)
	}
	mux.HandleFunc("/users", handler)
	actual, err := client.ListUser(context.Background())

	date, _ := time.Parse(time.DateOnly, "2023-01-11")

	want := User{
		PK:             "119266bma1vr",
		ResolverIP:     net.ParseIP("76.76.2.162"),
		EmailStatus:    IntBool{true},
		Tutorials:      IntBool{true},
		V:              2,
		ResolverStatus: IntBool{true},
		RuleProfile:    "109129bma5vn",
		Date:           Date{date},
		Status:         IntBool{true},
		Email:          "email@icloud.com",
		ResolverUid:    "zk30ksddco",
		ProxyAccess:    IntBool{true},
		StatsEndpoint:  "europe",
		LastActive:     UnixTime{time.Unix(1731223785, 0).UTC()},
		Twofa:          IntBool{true},
		Debug:          []any{},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
