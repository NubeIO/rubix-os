package rubixmodel

type DiscoveredSlaves struct {
	Slaves []Slaves `json:"slaves"`
}

type Slaves struct {
	GlobalUUID   string `json:"global_uuid"`
	CreatedOn    string `json:"created_on,omitempty"`
	UpdatedOn    string `json:"updated_on,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientName   string `json:"client_name,omitempty"`
	SiteID       string `json:"site_id,omitempty"`
	SiteName     string `json:"site_name,omitempty"`
	DeviceID     string `json:"device_id,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	SiteAddress  string `json:"site_address,omitempty"`
	SiteCity     string `json:"site_city,omitempty"`
	SiteState    string `json:"site_state,omitempty"`
	SiteZip      string `json:"site_zip,omitempty"`
	SiteCountry  string `json:"site_country,omitempty"`
	SiteLat      string `json:"site_lat,omitempty"`
	SiteLon      string `json:"site_lon,omitempty"`
	TimeZone     string `json:"time_zone,omitempty"`
	IsMaster     bool   `json:"is_master,omitempty"`
	Count        int    `json:"count,omitempty"`
	IsOnline     bool   `json:"is_online,omitempty"`
	TotalCount   int    `json:"total_count,omitempty"`
	FailureCount int    `json:"failure_count,omitempty"`
}
