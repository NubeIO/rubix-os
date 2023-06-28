package bacnetschema

import "github.com/NubeIO/rubix-os/schema/schema"

type ObjectType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Object Type"`
	Options  []string `json:"enum" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\"]"`
	EnumName []string `json:"enumNames" default:"[\"analog input\",\"analog value\",\"analog output\",\"binary input\",\"binary value\",\"binary output\"]"`
	Default  string   `json:"default" default:"analog_input"`
}

type PointSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`

	ObjectId   schema.ObjectId `json:"object_id"`
	ObjectType ObjectType      `json:"object_type"`

	/*  TODO: THESE HAVE BEEN REMOVED UNTIL WE FIGURE OUT HOW BEST TO APPLY THEM TO BACNET POINTS
	MultiplicationFactor schema.MultiplicationFactor `json:"multiplication_factor"`
	ScaleEnable          schema.ScaleEnable          `json:"scale_enable"`
	ScaleInMin           schema.ScaleInMin           `json:"scale_in_min"`
	ScaleInMax           schema.ScaleInMax           `json:"scale_in_max"`
	ScaleOutMin          schema.ScaleOutMin          `json:"scale_out_min"`
	ScaleOutMax          schema.ScaleOutMax          `json:"scale_out_max"`
	Offset               schema.Offset               `json:"offset"`
	*/
	Decimal  schema.Decimal  `json:"decimal"`
	Fallback schema.Fallback `json:"fallback"`

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
