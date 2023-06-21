package lorawanschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	UUID          schema.UUID        `json:"uuid"`
	Name          schema.Name        `json:"name"`
	Description   schema.Description `json:"description"`
	Enable        schema.Enable      `json:"enable"`
	SkipFCntCheck struct {
		Type     string `json:"type" default:"boolean"`
		Title    string `json:"title" default:"Disable Frame-Counter Validation"`
		Default  bool   `json:"default" default:"true"`
		ReadOnly bool   `json:"readOnly" default:"false"`
	} `json:"zero_mode"`
	DevEUI struct {
		Type     string `json:"type" default:"string"`
		Title    string `json:"title" default:"Device EUI"`
		Min      int    `json:"minLength" default:"16"`
		Max      int    `json:"maxLength" default:"16"`
		Default  string `json:"default" default:""`
		ReadOnly bool   `json:"readOnly" default:"false"`
	} `json:"address_uuid"`
	ApplicationID struct {
		Type     string                  `json:"type" default:"number"`
		Title    string                  `json:"title" default:"Application"`
		Options  []schema.OptionOneOfInt `json:"oneOf"`
		Default  int                     `json:"default" default:"0"`
		ReadOnly bool                    `json:"readOnly" default:"false"`
	} `json:"address_id"`
	DeviceProfileID struct {
		Type     string               `json:"type" default:"string"`
		Title    string               `json:"title" default:"Device profile"`
		Options  []schema.OptionOneOf `json:"oneOf"`
		Default  string               `json:"default" default:""`
		ReadOnly bool                 `json:"readOnly" default:"false"`
	} `json:"model"`
	AppKey struct {
		Type     string `json:"type" default:"string"`
		Title    string `json:"title" default:"Application Key"`
		Min      int    `json:"minLength" default:"32"`
		Max      int    `json:"maxLength" default:"32"`
		Default  string `json:"default" default:""`
		ReadOnly bool   `json:"readOnly" default:"false"`
	} `json:"manufacture"`
	HistoryEnable schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
