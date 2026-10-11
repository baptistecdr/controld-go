package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

const endpointScheduleFixture = `{
	"PK": "schedulePK",
	"name": "iPad schedule",
	"time_zone": "Europe/Paris",
	"status": 1,
	"created_at": "2026-10-11 12:52:56",
	"updated_at": "2026-10-11 12:52:57",
	"windows": [
		{
			"weekday": "mon",
			"time_start": "00:00",
			"time_end": "08:00",
			"profile_id": "profileID",
			"profile_name": "profileName",
			"start_minute": 0,
			"end_minute": 480
		},
		{
			"weekday": "mon",
			"time_start": "21:00",
			"time_end": "23:59",
			"profile_id": "profileID",
			"profile_name": "profileName",
			"start_minute": 1260,
			"end_minute": 1439
		}
	],
	"weekdays": {"mon": 1, "tue": 0, "wed": 0, "thu": 0, "fri": 0, "sat": 0, "sun": 0},
	"enforcing": 0,
	"device": {
		"PK": "deviceID",
		"device_id": "deviceID",
		"name": "iPad",
		"resolvers": {
			"doh": "https:\/\/dns.controld.com\/deviceID",
			"dot": "deviceID.dns.controld.com"
		}
	},
	"profile_id": "profileID",
	"profile_name": "profileName",
	"time_start": "00:00",
	"time_end": "08:00",
	"active_window": null
}`

func wantEndpointSchedule() EndpointSchedule {
	return EndpointSchedule{
		PK:        "schedulePK",
		Name:      "iPad schedule",
		TimeZone:  "Europe/Paris",
		Status:    IntBool(true),
		CreatedAt: "2026-10-11 12:52:56",
		UpdatedAt: "2026-10-11 12:52:57",
		Windows: []EndpointScheduleWindow{
			{
				Weekday:     ScheduleMonday,
				TimeStart:   "00:00",
				TimeEnd:     "08:00",
				ProfileID:   "profileID",
				ProfileName: "profileName",
				StartMinute: 0,
				EndMinute:   480,
			},
			{
				Weekday:     ScheduleMonday,
				TimeStart:   "21:00",
				TimeEnd:     "23:59",
				ProfileID:   "profileID",
				ProfileName: "profileName",
				StartMinute: 1260,
				EndMinute:   1439,
			},
		},
		Weekdays: ScheduleWeekdays{Mon: IntBool(true)},
		Device: Device{
			PK:       "deviceID",
			DeviceID: "deviceID",
			Name:     "iPad",
			Resolvers: Resolvers{
				DoH: "https://dns.controld.com/deviceID",
				DoT: "deviceID.dns.controld.com",
			},
		},
		ProfileID:   "profileID",
		ProfileName: "profileName",
		TimeStart:   "00:00",
		TimeEnd:     "08:00",
	}
}

func TestListEndpointSchedules(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"body": {"schedules": [%s]}, "success": true}`, endpointScheduleFixture)
	}

	mux.HandleFunc("/endpointschedules", handler)
	actual, err := client.ListEndpointSchedules(context.Background())

	if assert.NoError(t, err) {
		assert.Equal(t, []EndpointSchedule{wantEndpointSchedule()}, actual)
	}
}

func TestCreateEndpointSchedule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"name": "iPad schedule",
			"device_id": "deviceID",
			"time_zone": "Europe/Paris",
			"status": 1,
			"windows": [
				{"weekday": "mon", "time_start": "00:00", "time_end": "08:00", "profile_id": "profileID", "start_minute": 0, "end_minute": 480, "sort_order": 0},
				{"weekday": "mon", "time_start": "21:00", "time_end": "23:59", "profile_id": "profileID"}
			]
		}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"body": %s, "success": true, "message": "iPad schedule schedule has been added."}`, endpointScheduleFixture)
	}

	mux.HandleFunc("/endpointschedules", handler)

	status := IntBool(true)
	zero, endMinute := 0, 480
	actual, err := client.CreateEndpointSchedule(context.Background(), CreateEndpointScheduleParams{
		Name:     "iPad schedule",
		DeviceID: "deviceID",
		TimeZone: "Europe/Paris",
		Status:   &status,
		Windows: []EndpointScheduleWindowParams{
			{
				Weekday:     ScheduleMonday,
				TimeStart:   "00:00",
				TimeEnd:     "08:00",
				ProfileID:   "profileID",
				StartMinute: &zero,
				EndMinute:   &endMinute,
				SortOrder:   &zero,
			},
			{
				Weekday:   ScheduleMonday,
				TimeStart: "21:00",
				TimeEnd:   "23:59",
				ProfileID: "profileID",
			},
		},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, wantEndpointSchedule(), actual)
	}
}

func TestUpdateEndpointSchedule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		// Only the fields that were set are sent, and the schedule ID is only in the path.
		assert.JSONEq(t, `{"name": "iPad schedule"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"body": %s, "success": true, "message": "iPad schedule schedule has been updated."}`, endpointScheduleFixture)
	}

	scheduleID := "schedulePK"
	mux.HandleFunc(fmt.Sprintf("/endpointschedules/%s", scheduleID), handler)

	name := "iPad schedule"
	actual, err := client.UpdateEndpointSchedule(context.Background(), UpdateEndpointScheduleParams{
		ScheduleID: scheduleID,
		Name:       &name,
	})
	if assert.NoError(t, err) {
		assert.Equal(t, wantEndpointSchedule(), actual)
	}

	_, err = client.UpdateEndpointSchedule(context.Background(), UpdateEndpointScheduleParams{ScheduleID: ""})
	require.Error(t, err, "Endpoint schedule should not have been updated")
}

func TestDeleteEndpointSchedule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"body": [],
			"success": true,
			"message": "iPad schedule schedule has been removed."
		}`)
	}

	scheduleID := "schedulePK"
	mux.HandleFunc(fmt.Sprintf("/endpointschedules/%s", scheduleID), handler)
	_, err := client.DeleteEndpointSchedule(context.Background(), DeleteEndpointScheduleParams{ScheduleID: scheduleID})
	require.NoError(t, err, "Endpoint schedule should have been deleted")

	_, err = client.DeleteEndpointSchedule(context.Background(), DeleteEndpointScheduleParams{ScheduleID: ""})
	require.Error(t, err, "Endpoint schedule should not have been deleted")
}
