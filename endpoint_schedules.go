package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ScheduleWeekday is a day of the week as used by endpoint schedules.
type ScheduleWeekday string

const (
	ScheduleMonday    ScheduleWeekday = "mon"
	ScheduleTuesday   ScheduleWeekday = "tue"
	ScheduleWednesday ScheduleWeekday = "wed"
	ScheduleThursday  ScheduleWeekday = "thu"
	ScheduleFriday    ScheduleWeekday = "fri"
	ScheduleSaturday  ScheduleWeekday = "sat"
	ScheduleSunday    ScheduleWeekday = "sun"
)

// ScheduleWeekdays tells on which days of the week a schedule has windows.
type ScheduleWeekdays struct {
	Mon IntBool `json:"mon"`
	Tue IntBool `json:"tue"`
	Wed IntBool `json:"wed"`
	Thu IntBool `json:"thu"`
	Fri IntBool `json:"fri"`
	Sat IntBool `json:"sat"`
	Sun IntBool `json:"sun"`
}

type EndpointScheduleWindow struct {
	Weekday     ScheduleWeekday `json:"weekday"`
	TimeStart   string          `json:"time_start"`
	TimeEnd     string          `json:"time_end"`
	ProfileID   string          `json:"profile_id"`
	ProfileName string          `json:"profile_name"`
	StartMinute int             `json:"start_minute"`
	EndMinute   int             `json:"end_minute"`
}

// EndpointSchedule is a weekly set of time windows attached to an endpoint
// (device), during which another profile is enforced. This endpoint is
// undocumented; its shape was determined by observing the live API.
type EndpointSchedule struct {
	PK       string  `json:"PK"`
	Name     string  `json:"name"`
	TimeZone string  `json:"time_zone"`
	Status   IntBool `json:"status"`
	// CreatedAt and UpdatedAt are "2006-01-02 15:04:05" timestamps without an
	// explicit time zone, so they are kept as returned by the API.
	CreatedAt string                   `json:"created_at"`
	UpdatedAt string                   `json:"updated_at"`
	Windows   []EndpointScheduleWindow `json:"windows"`
	Weekdays  ScheduleWeekdays         `json:"weekdays"`
	Enforcing IntBool                  `json:"enforcing"`
	// Device is only partially populated by ListEndpointSchedules.
	Device Device `json:"device"`
	// ProfileID, ProfileName, TimeStart and TimeEnd mirror the first window.
	ProfileID    string `json:"profile_id"`
	ProfileName  string `json:"profile_name"`
	TimeStart    string `json:"time_start"`
	TimeEnd      string `json:"time_end"`
	ActiveWindow any    `json:"active_window"`
}

type ListEndpointSchedulesBody struct {
	Schedules []EndpointSchedule `json:"schedules"`
}

type ListEndpointSchedulesResponse struct {
	Body ListEndpointSchedulesBody `json:"body"`
	Response
}

// EndpointScheduleWindowParams describes a window of a schedule. Windows on
// the same day cannot overlap.
type EndpointScheduleWindowParams struct {
	Weekday ScheduleWeekday `json:"weekday"`
	// TimeStart and TimeEnd use the 24h HH:MM format.
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	ProfileID string `json:"profile_id"`
	// StartMinute, EndMinute and SortOrder are derived by ControlD when omitted.
	StartMinute *int `json:"start_minute,omitempty"`
	EndMinute   *int `json:"end_minute,omitempty"`
	SortOrder   *int `json:"sort_order,omitempty"`
}

type CreateEndpointScheduleParams struct {
	Name     string `json:"name"`
	DeviceID string `json:"device_id"`
	// TimeZone is an IANA time zone name, e.g. "Europe/Paris".
	TimeZone string `json:"time_zone"`
	// Status defaults to enabled when omitted.
	Status  *IntBool                       `json:"status,omitempty"`
	Windows []EndpointScheduleWindowParams `json:"windows"`
}

type CreateEndpointScheduleResponse struct {
	Body    EndpointSchedule `json:"body"`
	Message string           `json:"message"`
	Response
}

// UpdateEndpointScheduleParams only modifies the fields that are set. When
// Windows is set, it replaces all the windows of the schedule.
type UpdateEndpointScheduleParams struct {
	ScheduleID string                         `json:"-"`
	Name       *string                        `json:"name,omitempty"`
	TimeZone   *string                        `json:"time_zone,omitempty"`
	Status     *IntBool                       `json:"status,omitempty"`
	Windows    []EndpointScheduleWindowParams `json:"windows,omitempty"`
}

type UpdateEndpointScheduleResponse struct {
	Body    EndpointSchedule `json:"body"`
	Message string           `json:"message"`
	Response
}

type DeleteEndpointScheduleParams struct {
	ScheduleID string `json:"-"`
}

type DeleteEndpointScheduleResponse struct {
	Body    []any  `json:"body"`
	Message string `json:"message"`
	Response
}

func (api *API) ListEndpointSchedules(ctx context.Context) ([]EndpointSchedule, error) {
	uri := buildURI("/endpointschedules", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []EndpointSchedule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListEndpointSchedulesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []EndpointSchedule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Schedules, nil
}

func (api *API) CreateEndpointSchedule(ctx context.Context, params CreateEndpointScheduleParams) (EndpointSchedule, error) {
	uri := buildURI("/endpointschedules", nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return EndpointSchedule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r CreateEndpointScheduleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return EndpointSchedule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) UpdateEndpointSchedule(ctx context.Context, params UpdateEndpointScheduleParams) (EndpointSchedule, error) {
	if params.ScheduleID == "" {
		return EndpointSchedule{}, fmt.Errorf("update: no schedule ID provided")
	}
	baseURL := fmt.Sprintf("/endpointschedules/%s", params.ScheduleID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return EndpointSchedule{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r UpdateEndpointScheduleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return EndpointSchedule{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) DeleteEndpointSchedule(ctx context.Context, params DeleteEndpointScheduleParams) ([]any, error) {
	if params.ScheduleID == "" {
		return []any{}, fmt.Errorf("delete: no schedule ID provided")
	}
	baseURL := fmt.Sprintf("/endpointschedules/%s", params.ScheduleID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteEndpointScheduleResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
