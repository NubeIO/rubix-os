package masterschema

import "github.com/NubeIO/rubix-os/schema/schema"

type PointSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`

	ObjectId   schema.ObjectId      `json:"object_id"`
	ObjectType schema.ObjectType    `json:"object_type"`
	WriteMode  BACnetPointWriteMode `json:"write_mode"`

	PollPriority schema.PollPriority `json:"poll_priority"`
	PollRate     schema.PollRate     `json:"poll_rate"`

	MultiplicationFactor schema.MultiplicationFactor `json:"multiplication_factor"`
	ScaleEnable          schema.ScaleEnable          `json:"scale_enable"`
	ScaleInMin           schema.ScaleInMin           `json:"scale_in_min"`
	ScaleInMax           schema.ScaleInMax           `json:"scale_in_max"`
	ScaleOutMin          schema.ScaleOutMin          `json:"scale_out_min"`
	ScaleOutMax          schema.ScaleOutMax          `json:"scale_out_max"`
	Offset               schema.Offset               `json:"offset"`
	Decimal              schema.Decimal              `json:"decimal"`
	Fallback             schema.Fallback             `json:"fallback"`

	HistoryEnable       schema.HistoryEnableDefaultTrue `json:"history_enable"`
	HistoryType         schema.HistoryType              `json:"history_type"`
	HistoryInterval     schema.HistoryInterval          `json:"history_interval"`
	HistoryCOVThreshold schema.HistoryCOVThreshold      `json:"history_cov_threshold"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}

type BACnetPointWriteMode struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Write Mode"`
	Options  []string `json:"enum" default:"[\"read_only\",\"write_once_then_read\",\"write_always\"]"`
	EnumName []string `json:"enumNames" default:"[\"read only\",\"write once then read\",\"write always\"]"`
	Default  string   `json:"default" default:"read_only"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}
