package pluginapi

import "net/url"

type Help struct {
	Name               string `json:"name"`
	PluginType         string `json:"plugin_type"` // protocol
	AllowConfigWrite   bool   `json:"allow_config_write"`
	IsNetwork          bool   `json:"is_network"`
	MaxAllowedNetworks int    `json:"max_allowed_networks"` // as an example for lora only 1 network is allowed to be added
	NetworkType        string `json:"network_type"`         // lora
	TransportType      string `json:"transport_type"`       // serial
}

// Response ...
type Response struct {
	Details Help `json:"details"`
}

// Displayer is the interface plugins should implement to show a text to the user.
// The text will appear on the plugin details page and can be multi-line.
// Markdown syntax is allowed. Good for providing dynamically generated instructions to the user.
// Location is the current location the user is accessing the API from, nil if not recoverable.
// Location contains the path to the display api endpoint, you may only need the base url.
type Displayer interface {
	Plugin
	GetDisplay(location *url.URL) Response
}
