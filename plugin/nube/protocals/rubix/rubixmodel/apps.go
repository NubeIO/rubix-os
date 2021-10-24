package rubixmodel

/*
USER
*/

type TokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	Message     *string `json:"message,omitempty"`
}

type UserResponse struct {
	Username string `json:"username"`
}

/*
APPs
*/

type AppControl []AppAction

type AppAction struct {
	Service string  `json:"service"`
	Action  string  `json:"action"`
	Error   *string `json:"error,omitempty"`
}

type AppsInstall []AppsInstallElement

type AppsInstallElement struct {
	Version            string `json:"version"`
	AppType            string `json:"app_type"`
	GatewayAccess      bool   `json:"gateway_access"`
	MinSupportVersion  string `json:"min_support_version"`
	Port               int    `json:"port"`
	DisplayName        string `json:"display_name"`
	Service            string `json:"service"`
	IsInstalled        bool   `json:"is_installed"`
	State              string `json:"state"`
	Status             bool   `json:"status"`
	DateSince          string `json:"date_since"`
	TimeSince          string `json:"time_since"`
	IsEnabled          bool   `json:"is_enabled"`
	BrowserDownloadURL string `json:"browser_download_url"`
	LatestVersion      string `json:"latest_version"`
}

//AppsDownload install an app
type AppsDownload []AppsDownloadElement

type AppsDownloadElement struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

//AppsLatestVersions get version
type AppsLatestVersions struct {
	WIRES             string `json:"WIRES"`
	RUBIXPLAT         string `json:"RUBIX_PLAT"`
	USERMANAGEMENT    string `json:"USER_MANAGEMENT"`
	POINTSERVER       string `json:"POINT_SERVER"`
	MODBUS            string `json:"MODBUS"`
	BACNETSERVER      string `json:"BACNET_SERVER"`
	LORARAW           string `json:"LORA_RAW"`
	DATAPUSH          string `json:"DATA_PUSH"`
	RUBIXBACNETMASTER string `json:"RUBIX_BACNET_MASTER"`
	FLOWFRAMEWORK     string `json:"FLOW_FRAMEWORK"`
}

//DownloadState state
type DownloadState struct {
	State    string `json:"state"`
	Services []struct {
		Service  string `json:"service"`
		Version  string `json:"version"`
		Download bool   `json:"download"`
		Error    string `json:"error"`
	} `json:"services"`
}

type GeneralResponse struct {
	GlobalUUID string `json:"global_uuid,omitempty"`
	Message    string `json:"message,omitempty"`
}
