package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type IP struct {
	IP      net.IP `json:"ip"`
	Type    string `json:"type"`
	Org     string `json:"org"`
	Country string `json:"country"`
	Handler string `json:"handler"`
	Pop     string `json:"pop"`
}

type ListIPResponse struct {
	Body IP `json:"body"`
	Response
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Status struct {
	API int64 `json:"api"`
	DNS int64 `json:"dns"`
	Pxy int64 `json:"pxy"`
}

type Network struct {
	IataCode    string   `json:"iata_code"`
	CityName    string   `json:"city_name"`
	CountryName string   `json:"country_name"`
	Location    Location `json:"location"`
	Status      Status   `json:"status"`
}

type ListNetworkBody struct {
	Network    []Network `json:"network"`
	Time       UnixTime  `json:"time"`
	CurrentPop string    `json:"current_pop"`
}

type ListNetworkResponse struct {
	Body ListNetworkBody `json:"body"`
	Response
}

func (api *API) ListIP(ctx context.Context) (IP, error) {
	baseURL := fmt.Sprintf("/ip")
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return IP{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r ListIPResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return IP{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) ListNetwork(ctx context.Context) ([]Network, error) {
	baseURL := fmt.Sprintf("/network")
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Network{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r ListNetworkResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return []Network{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Network, nil
}
