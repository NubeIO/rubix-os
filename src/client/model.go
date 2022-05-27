package client

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type Token struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

type ResponsePlugins struct {
	Response Plugins `json:"response"`
	Status   string  `json:"status"`
	Count    int     `json:"count"`
}

type Plugins struct {
	Items []model.PluginConf
}

type ResponseBody struct {
	Response ResponseCommon `json:"response"`
	Status   string         `json:"status"`
	Count    string         `json:"count"`
}

type ResponseCommon struct {
	UUID string `json:"uuid"`
	// Name        string `json:"name"`
}

type Stream struct {
	Name     string `json:"name"`
	IsRemote bool   `json:"is_remote"`
}

type Consumer struct {
	Name                string `json:"name"`
	Enable              bool   `json:"enable"`
	ProducerType        string `json:"producer_type"`
	ProducerApplication string `json:"producer_application"`
	StreamUUID          string `json:"stream_uuid"`
	ToUUID              string `json:"to_uuid"`
	IsRemote            bool   `json:"is_remote"`
	RemoteRubixUUID     string `json:"remote_rubix_uuid"`
}
