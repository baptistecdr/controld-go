package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Proxy struct {
	PK          string  `json:"PK"`
	UID         string  `json:"uid"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryName string  `json:"country_name"`
	Lat         float64 `json:"gps_lat"`
	Long        float64 `json:"gps_long"`
	Hidden      bool    `json:"hidden,omitempty"`
}

type ProxyCountry struct {
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
}

type ListProxiesBody struct {
	Proxies   []Proxy        `json:"proxies"`
	Countries []ProxyCountry `json:"countries"`
}

type ListProxiesResponse struct {
	Body ListProxiesBody `json:"body"`
	Response
}

// ListProxies returns the proxies that traffic can be redirected through
// (usable as the `via` target of a REDIRECT custom rule), along with the
// list of countries they belong to.
func (api *API) ListProxies(ctx context.Context) ([]Proxy, []ProxyCountry, error) {
	uri := buildURI("/proxies", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListProxiesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Proxies, r.Body.Countries, nil
}
