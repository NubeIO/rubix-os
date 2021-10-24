package rubixmodel

type WiresPlat struct {
	GlobalUuid  string      `json:"global_uuid"`
	CreatedOn   interface{} `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	ClientId    string      `json:"client_id"`
	ClientName  string      `json:"client_name"`
	SiteId      string      `json:"site_id"`
	SiteName    string      `json:"site_name"`
	DeviceId    string      `json:"device_id"`
	DeviceName  string      `json:"device_name"`
	SiteAddress string      `json:"site_address"`
	SiteCity    string      `json:"site_city"`
	SiteState   string      `json:"site_state"`
	SiteZip     string      `json:"site_zip"`
	SiteCountry string      `json:"site_country"`
	SiteLat     string      `json:"site_lat"`
	SiteLon     string      `json:"site_lon"`
	TimeZone    string      `json:"time_zone"`
}
