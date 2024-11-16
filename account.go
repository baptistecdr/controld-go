package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Date struct {
	time.Time
}

func (s Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Format(time.DateOnly))
}

func (s *Date) UnmarshalJSON(data []byte) error {
	dateStr := string(data[1 : len(data)-1])
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return err
	}
	s.Time = date
	return nil
}

type User struct {
	PK             string   `json:"PK"`
	ResolverIP     net.IP   `json:"resolver_ip"`
	EmailStatus    IntBool  `json:"email_status"`
	Tutorials      IntBool  `json:"tutorials"`
	V              int      `json:"v"`
	ResolverStatus IntBool  `json:"resolver_status"`
	RuleProfile    string   `json:"rule_profile"`
	Date           Date     `json:"date"`
	Status         IntBool  `json:"status"`
	Email          string   `json:"email"`
	ResolverUid    string   `json:"resolver_uid"`
	ProxyAccess    IntBool  `json:"proxy_access"`
	StatsEndpoint  string   `json:"stats_endpoint"`
	LastActive     UnixTime `json:"last_active"`
	Twofa          IntBool  `json:"twofa"`
	Debug          []any    `json:"debug"`
}

type ListUserResponse struct {
	Body User `json:"body"`
	Response
}

func (api *API) ListUser(ctx context.Context) (User, error) {
	uri := buildURI("/users", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return User{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r ListUserResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return User{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
