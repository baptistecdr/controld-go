package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type KnownIP struct {
	IP      net.IP   `json:"ip"`
	Ts      UnixTime `json:"ts"`
	Country string   `json:"country"`
	City    string   `json:"city"`
	ISP     string   `json:"isp"`
	Asn     int      `json:"asn"`
	AsName  string   `json:"as_name"`
}

type ListKnownIPsBody struct {
	IPs []KnownIP `json:"ips"`
}

type ListKnownIPsResponse struct {
	Body ListKnownIPsBody `json:"body"`
	Response
}

type ListKnownIPsParams struct {
	DeviceID string `json:"device_id"`
}

type LearnNewIPsParams struct {
	DeviceID string   `json:"device_id"`
	IPs      []net.IP `json:"ips"`
}

type LearnNewIPsResponse struct {
	Body []any `json:"body"`
	Response
	Message string `json:"message"`
}

type DeleteLearnedIPsParams struct {
	DeviceID string   `json:"device_id"`
	IPs      []net.IP `json:"ips"`
}

type DeleteLearnedIPsResponse struct {
	Body []any `json:"body"`
	Response
	Message string `json:"message"`
}

func (api *API) ListKnownIPs(ctx context.Context, params ListKnownIPsParams) ([]KnownIP, error) {
	uri := buildURI("/access", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, params)
	if err != nil {
		return []KnownIP{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r ListKnownIPsResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return []KnownIP{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.IPs, nil
}

func (api *API) LearnNewIPs(ctx context.Context, params LearnNewIPsParams) ([]any, error) {
	uri := buildURI("/access", nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r LearnNewIPsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) DeleteLearnedIPs(ctx context.Context, params DeleteLearnedIPsParams) ([]any, error) {
	uri := buildURI("/access", nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, params)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteLearnedIPsResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
