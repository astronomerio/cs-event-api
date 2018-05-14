package v1

import (
	"net"
	"github.com/arizz96/event-api/types"
)

// Context provides the representation of the `context` object
type Context struct {
	App       *AppInfo      `json:"app,omitempty"`
	Campaign  *CampaignInfo `json:"campaign,omitempty"`
	Device    *DeviceInfo   `json:"device,omitempty"`
	Library   *LibraryInfo  `json:"library,omitempty"`
	Location  *LocationInfo `json:"location,omitempty"`
	Network   *NetworkInfo  `json:"network,omitempty"`
	OS        *OSInfo       `json:"os,omitempty"`
	Page      *PageInfo     `json:"page,omitempty"`
	Referrer  *ReferrerInfo `json:"referrer,omitempty"`
	Screen    *ScreenInfo   `json:"screen,omitempty"`
	IP        net.IP        `json:"ip,omitempty"`
	Locale    string        `json:"locale,omitempty"`
	Timezone  string        `json:"timezone,omitempty"`
	UserAgent string        `json:"userAgent,omitempty"`
	Traits    *Traits       `json:"traits,omitempty"`
}

// AppInfo provides the representation of the `context.app` object
type AppInfo struct {
	Name      string                  `json:"name,omitempty"`
	Version   string                  `json:"version,omitempty"`
	Build     types.ConvertibleString `json:"build,omitempty"`
	Namespace string                  `json:"namespace,omitempty"`
}

// CampaignInfo provides the representation of the `context.campaign` object
type CampaignInfo struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Term    string `json:"term,omitempty"`
	Content string `json:"content,omitempty"`
}

// DeviceInfo provides the representation of the `context.device` object
type DeviceInfo struct {
	ID            string `json:"id,omitempty"`
	Manufacturer  string `json:"manufacturer,omitempty"`
	Model         string `json:"model,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Version       string `json:"version,omitempty"`
	AdvertisingID string `json:"advertisingId,omitempty"`
}

// LibraryInfo provides the representation of the `context.library` object
type LibraryInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// LocationInfo provides the representation of the `context.location` object
type LocationInfo struct {
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

// NetworkInfo provides the representation of the `context.network` object
type NetworkInfo struct {
	Bluetooth types.ConvertibleBoolean `json:"bluetooth,omitempty"`
	Cellular  types.ConvertibleBoolean `json:"cellular,omitempty"`
	WIFI      types.ConvertibleBoolean `json:"wifi,omitempty"`
	Carrier   string                   `json:"carrier,omitempty"`
}

// OSInfo provides the representation of the `context.os` object
type OSInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// PageInfo provides the representation of the `context.page` object
type PageInfo struct {
	Hash     string `json:"hash,omitempty"`
	Path     string `json:"path,omitempty"`
	Referrer string `json:"referrer,omitempty"`
	Search   string `json:"search,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
}

// ReferrerInfo provides the representation of the `context.referrer` object
type ReferrerInfo struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Link string `json:"link,omitempty"`
}

// ScreenInfo provides the representation of the `context.screen` object
type ScreenInfo struct {
	Density int `json:"density,omitempty"`
	Width   int `json:"width,omitempty"`
	Height  int `json:"height,omitempty"`
}
