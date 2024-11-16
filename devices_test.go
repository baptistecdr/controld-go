package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestListDevices(t *testing.T) {
	setup()
	defer teardown()
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
		  "body": {
			"devices": [
			  {
				"PK": "PK",
				"ts": 1674972821,
				"name": "deviceName",
				"user": "user",
				"stats": 2,
				"device_id": "deviceID",
				"status": 1,
				"learn_ip": 0,
				"resolvers": {
				  "uid": "resolverUID",
				  "doh": "https:\/\/dns.controld.com\/deviceID",
				  "dot": "deviceID.dns.controld.com",
				  "v6": [
					"ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a",
					"526a:2dc5:a3fe:0732:b240:08ef:ece4:c828"
				  ]
				},
				"icon": "desktop-mac",
				"profile": {
				  "PK": "ProfileID",
				  "updated": 1714923836,
				  "name": "tvOS"
				},
				"last_activity": 1716043811
			  }
			]
		  },
		  "success": true
		}`)
	}
	stats := Full
	ipv6 := []net.IP{net.ParseIP("ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a"), net.ParseIP("526a:2dc5:a3fe:0732:b240:08ef:ece4:c828")}
	icon := DesktopMac

	want := []Device{
		{
			PK:       "PK",
			Ts:       UnixTime{time.Unix(1674972821, 0).UTC()},
			Name:     "deviceName",
			User:     "user",
			Stats:    &stats,
			DeviceID: "deviceID",
			Status:   Active,
			LearnIP:  IntBool{false},
			Resolvers: Resolvers{
				Uid: "resolverUID",
				DoH: "https://dns.controld.com/deviceID",
				DoT: "deviceID.dns.controld.com",
				V6:  &ipv6,
			},
			Icon: &icon,
			Profile: Profile{
				PK:      "ProfileID",
				Updated: UnixTime{time.Unix(1714923836, 0).UTC()},
				Name:    "tvOS",
			},
		},
	}
	mux.HandleFunc("/devices", handler)
	actual, err := client.ListDevices(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestCreateDevice(t *testing.T) {
	setup()
	defer teardown()
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
		  "body": {
			"PK": "devicePK",
			"ts": 1716053105,
			"name": "deviceName",
			"user": "user",
			"device_id": "deviceID",
			"status": 0,
			"learn_ip": 0,
			"resolvers": {
			  "uid": "resolverUID",
			  "doh": "https:\/\/dns.controld.com\/deviceID",
			  "dot": "deviceID.dns.controld.com",
			  "v6": [
				"ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a",
				"526a:2dc5:a3fe:0732:b240:08ef:ece4:c828"
			  ]
			},
			"icon": "desktop-mac",
			"profile": {
			  "PK": "profilePK",
			  "updated": 1709483590,
			  "name": "profileName"
			}
		  },
		  "success": true,
		  "message": "Device has been added"
		}`)
	}

	mux.HandleFunc("/devices", handler)
	actual, err := client.CreateDevice(context.Background(), CreateDeviceParams{
		Name:      "deviceName",
		ProfileID: "ProfileID",
		Icon:      DesktopMac,
	})

	v6 := []net.IP{net.ParseIP("ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a"), net.ParseIP("526a:2dc5:a3fe:0732:b240:08ef:ece4:c828")}
	icon := DesktopMac
	want := Device{
		PK:       "devicePK",
		Ts:       UnixTime{time.Unix(1716053105, 0).UTC()},
		Name:     "deviceName",
		User:     "user",
		DeviceID: "deviceID",
		Status:   Pending,
		LearnIP:  IntBool{false},
		Resolvers: Resolvers{
			Uid: "resolverUID",
			DoH: "https://dns.controld.com/deviceID",
			DoT: "deviceID.dns.controld.com",
			V6:  &v6,
		},
		Icon: &icon,
		Profile: Profile{
			PK:      "profilePK",
			Updated: UnixTime{time.Unix(1709483590, 0).UTC()},
			Name:    "profileName",
		},
	}

	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestListDeviceTypes(t *testing.T) {
	setup()
	defer teardown()
	deviceTypes := loadFixture("devices", "ListDeviceTypes")
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", deviceTypes)
	}
	mux.HandleFunc("/devices/types", handler)
	_, err := client.ListDeviceType(context.Background())
	assert.NoError(t, err)
}

func TestUpdateDevice(t *testing.T) {
	setup()
	defer teardown()

	var deviceID = "deviceID"
	var newDeviceName = "newDeviceName"

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
		  "body": {
			"PK": "PK",
			"ts": 1716053105,
			"name": "newDeviceName",
			"user": "user",
			"device_id": "deviceID",
			"status": 0,
			"learn_ip": 0,
			"resolvers": {
			  "uid": "resolverUID",
			  "doh": "https:\/\/dns.controld.com\/deviceID",
			  "dot": "deviceID.dns.controld.com",
			  "v6": [
				"ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a",
				"526a:2dc5:a3fe:0732:b240:08ef:ece4:c828"
			  ]
			},
			"icon": "desktop-mac",
			"profile": {
			  "PK": "PK",
			  "updated": 1709483590,
			  "name": "profileName"
			}
		  },
		  "success": true,
		  "message": "Device has been updated"
		}`)
	}

	mux.HandleFunc(fmt.Sprintf("/devices/%s", deviceID), handler)
	actual, err := client.UpdateDevice(context.Background(), UpdateDeviceParams{
		DeviceID: deviceID,
		Name:     &newDeviceName,
	})

	v6 := []net.IP{net.ParseIP("ef4f:81ab:0618:4663:d938:4cff:8bb2:0d2a"), net.ParseIP("526a:2dc5:a3fe:0732:b240:08ef:ece4:c828")}
	icon := DesktopMac
	want := Device{
		PK:       "PK",
		Ts:       UnixTime{time.Unix(1716053105, 0).UTC()},
		Name:     "newDeviceName",
		User:     "user",
		DeviceID: "deviceID",
		Status:   Pending,
		LearnIP:  IntBool{false},
		Resolvers: Resolvers{
			Uid: "resolverUID",
			DoH: "https://dns.controld.com/deviceID",
			DoT: "deviceID.dns.controld.com",
			V6:  &v6,
		},
		Icon: &icon,
		Profile: Profile{
			PK:      "PK",
			Updated: UnixTime{time.Unix(1709483590, 0).UTC()},
			Name:    "profileName",
		},
	}
	if err == nil {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateDevice(context.Background(), UpdateDeviceParams{DeviceID: ""})
	require.Error(t, err)
}

func TestDeleteDevice(t *testing.T) {
	setup()
	defer teardown()
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
		  "body": [],
		  "success": true,
		  "message": "Device has been deleted"
	    }`)
	}
	deviceID := "deviceID"
	mux.HandleFunc(fmt.Sprintf("/devices/%s", deviceID), handler)
	_, err := client.DeleteDevice(context.Background(), DeleteDeviceParams{DeviceID: deviceID})
	require.NoError(t, err, "Device should have been deleted")

	_, err = client.DeleteDevice(context.Background(), DeleteDeviceParams{DeviceID: ""})
	require.Error(t, err, "Device should not have been deleted")
}
