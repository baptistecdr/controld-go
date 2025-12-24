package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type UnixTime struct {
	time.Time
}

func (s UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Unix())
}

func (s *UnixTime) UnmarshalJSON(data []byte) error {
	seconds, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	s.Time = time.Unix(seconds, 0).UTC()
	return nil
}

type IntBool bool

func (s IntBool) MarshalJSON() ([]byte, error) {
	if s {
		return json.Marshal(1)
	} else {
		return json.Marshal(0)
	}
}

func (s *IntBool) UnmarshalJSON(data []byte) error {
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*s = IntBool(value == 1)
	return nil
}

type AnalyticsLevel int

const (
	Off   AnalyticsLevel = 0
	Basic AnalyticsLevel = 1
	Full  AnalyticsLevel = 2
)

type IconName string

const (
	DesktopWindows    IconName = "desktop-windows"
	DesktopMac        IconName = "desktop-mac"
	DesktopLinux      IconName = "desktop-linux"
	MobileIOS         IconName = "mobile-ios"
	MobileAndroid     IconName = "mobile-android"
	BrowserChrome     IconName = "browser-chrome"
	BrowserFirefox    IconName = "browser-firefox"
	BrowserEdge       IconName = "browser-edge"
	BrowserBrave      IconName = "browser-brave"
	BrowserOther      IconName = "browser-other"
	TVApple           IconName = "tv-apple"
	TVAndroid         IconName = "tv-android"
	TVFireTV          IconName = "tv-firetv"
	TVSamsung         IconName = "tv-samsung"
	TVOther           IconName = "tv"
	RouterAsus        IconName = "router-asus"
	RouterDDWRT       IconName = "router-ddwrt"
	RouterFirewalla   IconName = "router-firewalla"
	RouterFreshTomato IconName = "router-freshtomato"
	RouterGLiNET      IconName = "router-glinet"
	RouterOpenWRT     IconName = "router-openwrt"
	RouterOPNsense    IconName = "router-opnsense"
	RouterPfSense     IconName = "router-pfsense"
	RouterSynology    IconName = "router-synology"
	RouterUbiquiti    IconName = "router-ubiquiti"
	RouterWindows     IconName = "router-windows"
	RouterLinux       IconName = "router-linux"
	RouterOther       IconName = "router"
)

type DDNS struct {
	Status    int    `json:"status"`
	Subdomain string `json:"subdomain"`
	Hostname  string `json:"hostname"`
	Record    string `json:"record"`
}

type DDNSExt struct {
	Status int    `json:"status"`
	Host   string `json:"host"`
}

type Resolvers struct {
	Uid string    `json:"uid"`
	DoH string    `json:"doh"`
	DoT string    `json:"dot"`
	V4  *[]net.IP `json:"v4,omitempty"`
	V6  *[]net.IP `json:"v6,omitempty"`
}

type LegacyIPv4 struct {
	Resolver string  `json:"resolver"`
	Status   IntBool `json:"status"`
}

type DeviceStatus int

const (
	Pending      = 0
	Active       = 1
	SoftDisabled = 2
	HardDisabled = 3
)

type Device struct {
	PK         string          `json:"PK"`
	Ts         UnixTime        `json:"ts"`
	Name       string          `json:"name"`
	User       string          `json:"user"`
	Stats      *AnalyticsLevel `json:"stats,omitempty"`
	DeviceID   string          `json:"device_id"`
	Status     DeviceStatus    `json:"status"`
	Restricted *IntBool        `json:"restricted,omitempty"`
	LearnIP    IntBool         `json:"learn_ip"`
	Desc       string          `json:"desc"`
	DDNS       *DDNS           `json:"ddns,omitempty"`
	Resolvers  Resolvers       `json:"resolvers"`
	LegacyIPv4 LegacyIPv4      `json:"legacy_ipv4"`
	Profile    Profile         `json:"profile"`
	Icon       *IconName       `json:"icon"`
}

type ListDevicesBody struct {
	Devices []Device `json:"devices"`
}

type ListDevicesResponse struct {
	Body ListDevicesBody `json:"body"`
	Response
}

type CreateDeviceParams struct {
	Name             string          `json:"name"`
	ProfileID        string          `json:"profile_id"`
	ProfileID2       *string         `json:"profile_id2,omitempty"`
	Icon             IconName        `json:"icon"`
	Stats            *AnalyticsLevel `json:"stats,omitempty"`
	LegacyIPv4Status *IntBool        `json:"legacy_ipv4_status,omitempty"`
	LearnIP          *IntBool        `json:"learn_ip,omitempty"`
	Restricted       *IntBool        `json:"restricted,omitempty"`
	Desc             *string         `json:"desc,omitempty"`
	DDNSStatus       *IntBool        `json:"ddns_status,omitempty"`
	DDNSSubdomain    *string         `json:"ddns_subdomain,omitempty"`
	DDNSExtStatus    *IntBool        `json:"ddns_ext_host,omitempty"`
	DDNSExtHost      *string         `json:"ddns_ext_status,omitempty"`
	RemapDeviceID    *string         `json:"remap_device_id,omitempty"`
	RemapClientID    *string         `json:"remap_client_id,omitempty"`
}

type CreateDeviceResponse struct {
	Body Device `json:"body"`
	Response
}
type Icon struct {
	Name string `json:"name"`
}

type OSIcons struct {
	MobileIos      Icon `json:"mobile-ios"`
	MobileAndroid  Icon `json:"mobile-android"`
	DesktopWindows Icon `json:"desktop-windows"`
	DesktopMac     Icon `json:"desktop-mac"`
	DesktopLinux   Icon `json:"desktop-linux"`
}

type BrowserIcons struct {
	BrowserChrome  Icon `json:"browser-chrome"`
	BrowserFirefox Icon `json:"browser-firefox"`
	BrowserEdge    Icon `json:"browser-edge"`
	BrowserBrave   Icon `json:"browser-brave"`
	BrowserOther   Icon `json:"browser-other"`
}

type Browser struct {
	Name  string  `json:"name"`
	Icons OSIcons `json:"icons"`
}

type OS struct {
	Name  string  `json:"name"`
	Icons OSIcons `json:"icons"`
}

type TVIcons struct {
	TV        Icon `json:"tv"`
	TVApple   Icon `json:"tv-apple"`
	TVAndroid Icon `json:"tv-android"`
	TVFireTV  Icon `json:"tv-firetv"`
	TVSamsung Icon `json:"tv-samsung"`
}

type TV struct {
	Name  string  `json:"name"`
	Icons TVIcons `json:"icons"`
}

type RouterIcons struct {
	Router         Icon `json:"router"`
	RouterOpenWRT  Icon `json:"router-openwrt"`
	RouterUbiquiti Icon `json:"router-ubiquiti"`
	RouterAsus     Icon `json:"router-asus"`
	RouterDDWRT    Icon `json:"router-ddwrt"`
}

type Router struct {
	Name     string      `json:"name"`
	Icons    RouterIcons `json:"icons"`
	SetupURL string      `json:"setup_url"`
}

type DeviceTypes struct {
	OS      OS      `json:"os"`
	Browser Browser `json:"browser"`
	TV      TV      `json:"tv"`
	Router  Router  `json:"router"`
}

type ListDeviceTypesBody struct {
	Types DeviceTypes `json:"types"`
}

type ListDeviceTypesResponse struct {
	Body ListDeviceTypesBody `json:"body"`
	Response
}

type UpdateDeviceParams struct {
	DeviceID          string          `json:"id"`
	Name              *string         `json:"name,omitempty"`
	ProfileID         *string         `json:"profile_id,omitempty"`
	ProfileID2        *string         `json:"profile_id2,omitempty"`
	Stats             *AnalyticsLevel `json:"stats,omitempty"`
	LegacyIPv4Status  *IntBool        `json:"legacy_ipv4_status,omitempty"`
	LearnIP           *IntBool        `json:"learn_ip,omitempty"`
	Restricted        *IntBool        `json:"restricted,omitempty"`
	BumpTLS           *IntBool        `json:"bump_tls,omitempty"`
	Desc              *string         `json:"desc,omitempty"`
	DDNSStatus        *IntBool        `json:"ddns_status,omitempty"`
	DDNSSubdomain     *string         `json:"ddns_subdomain,omitempty"`
	DDNSExtHost       *string         `json:"ddns_ext_host,omitempty"`
	DDNSExtStatus     *IntBool        `json:"ddns_ext_status,omitempty"`
	Status            *DeviceStatus   `json:"status,omitempty"`
	CtrldCustomConfig *string         `json:"ctrld_custom_config,omitempty"`
}

type UpdateDeviceResponse struct {
	Body    Device `json:"body"`
	Message string `json:"message"`
	Response
}

type DeleteDeviceParams struct {
	DeviceID string `json:"device-id"`
}

type DeleteDeviceResponse struct {
	Body    []any  `json:"body"`
	Message string `json:"message"`
	Response
}

func (api *API) ListDevices(ctx context.Context) ([]Device, error) {
	uri := buildURI("/devices", nil)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Device{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r ListDevicesResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Device{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Devices, nil
}

func (api *API) CreateDevice(ctx context.Context, params CreateDeviceParams) (Device, error) {
	uri := buildURI("/devices", nil)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return Device{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r CreateDeviceResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return Device{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) ListDeviceType(ctx context.Context) (DeviceTypes, error) {
	uri := buildURI("/devices/types", nil)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return DeviceTypes{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r ListDeviceTypesResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return DeviceTypes{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body.Types, nil
}

func (api *API) UpdateDevice(ctx context.Context, params UpdateDeviceParams) (Device, error) {
	if params.DeviceID == "" {
		return Device{}, fmt.Errorf("update: no device ID provided")
	}
	baseURL := fmt.Sprintf("/devices/%s", params.DeviceID)
	uri := buildURI(baseURL, nil)
	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return Device{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}
	var r UpdateDeviceResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Device{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}

func (api *API) DeleteDevice(ctx context.Context, params DeleteDeviceParams) ([]any, error) {
	if params.DeviceID == "" {
		return []any{}, fmt.Errorf("delete: no device ID provided")
	}

	baseURL := fmt.Sprintf("/devices/%s", params.DeviceID)
	uri := buildURI(baseURL, nil)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, params)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var r DeleteDeviceResponse

	err = json.Unmarshal(res, &r)
	if err != nil {
		return []any{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Body, nil
}
