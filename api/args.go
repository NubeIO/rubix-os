package api

type Args struct {
	Sort                 string
	Order                string
	Offset               string
	Limit                string
	Search               string
	AskRefresh           string
	AskResponse          string
	Write                string
	ThingType            string
	UUID                 *string
	FlowNetworkUUID      string
	FlowNetworkCloneUUID *string
	WriteHistory         string
	WriteConsumer        string
	Field                string
	Value                string
	UpdateProducer       string
	CompactPayload       string
	CompactWithName      string
	FlowUUID             string
	Name                 *string
	StreamUUID           *string
	ProducerUUID         *string
	ConsumerUUID         *string
	WriterUUID           string
	AddToParent          string
	GlobalUUID           *string
	ClientId             *string
	SiteId               *string
	DeviceId             *string
	WriterThingClass     *string
	SourceUUID           *string
	ProducerThingUUID    *string
	WriterThingUUID      *string
	ByPluginName         bool
	TimestampGt          *string
	TimestampLt          *string
	WithFlowNetworks     bool
	WithStreams          bool
	WithStreamClones     bool
	WithProducers        bool
	WithConsumers        bool
	WithCommandGroups    bool
	WithWriters          bool
	WithWriterClones     bool
	Networks             bool
	WithDevices          bool
	WithPoints           bool
	WithPriority         bool
	WithTags             bool
	PluginName           string
	NetworkName          *string
	DeviceName           *string
	PointName            *string
	AddressUUID          *string
	AddressID            *string
	ObjectType           *string
	IoNumber             *string
	IdGt                 *string
	IsRemote             *bool
	IsMetadata           bool
	NetworkUUID          *string
	DeviceUUID           *string
	WithMetaTags         bool
	AutoMappingUUID      *string
	AutoMappingEnable    *bool
	MetaTags             *string
}

var ArgsType = struct {
	Sort                 string
	Order                string
	Offset               string
	Limit                string
	Search               string
	AskRefresh           string
	AskResponse          string
	Write                string
	ThingType            string
	UUID                 string
	FlowNetworkUUID      string
	FlowNetworkCloneUUID string
	WriteHistory         string
	WriteConsumer        string
	Field                string
	Name                 string
	Value                string
	UpdateProducer       string
	CompactPayload       string // for a point would be presentValue
	CompactWithName      string // for a point would be presentValue and pointName
	FlowUUID             string
	StreamUUID           string
	ProducerUUID         string
	ConsumerUUID         string
	WriterUUID           string
	AddToParent          string
	GlobalUUID           string
	ClientId             string
	SiteId               string
	DeviceId             string
	WriterThingClass     string
	SourceUUID           string
	ProducerThingUUID    string
	WriterThingUUID      string
	ByPluginName         string
	TimestampGt          string
	TimestampLt          string
	WithFlowNetworks     string
	WithStreams          string
	WithStreamClones     string
	WithProducers        string
	WithConsumers        string
	WithCommandGroups    string
	WithWriters          string
	WithWriterClones     string
	WithNetworks         string
	WithDevices          string
	WithPoints           string
	WithPriority         string
	WithTags             string
	PluginName           string
	NetworkName          string
	DeviceName           string
	PointName            string
	AddressUUID          string
	AddressID            string
	ObjectType           string
	IoNumber             string
	IdGt                 string
	IsRemote             string
	IsMetadata           string
	NetworkUUID          string
	DeviceUUID           string
	WithMetaTags         string
	AutoMappingUUID      string
	AutoMappingEnable    string
	MetaTags             string
}{
	Sort:                 "sort",
	Order:                "order",
	Offset:               "offset",
	Limit:                "limit",
	Search:               "search",
	AskRefresh:           "ask_refresh",
	AskResponse:          "ask_response",
	Write:                "write",      // consumer to write a value
	ThingType:            "thing_type", // the type of thing like a point
	UUID:                 "uuid",
	FlowNetworkUUID:      "flow_network_uuid",
	FlowNetworkCloneUUID: "flow_network_clone_uuid",
	WriteHistory:         "write_history",
	WriteConsumer:        "write_consumer",
	Field:                "field",
	Name:                 "name",
	Value:                "value",
	UpdateProducer:       "update_producer",
	CompactPayload:       "compact_payload",
	CompactWithName:      "compact_with_name",
	FlowUUID:             "flow_uuid",
	StreamUUID:           "stream_uuid",
	ProducerUUID:         "producer_uuid",
	ConsumerUUID:         "consumer_uuid",
	WriterUUID:           "writer_uuid",
	AddToParent:          "add_to_parent",
	GlobalUUID:           "global_uuid",
	ClientId:             "client_id",
	SiteId:               "site_id",
	DeviceId:             "device_id",
	WriterThingClass:     "writer_thing_class",
	SourceUUID:           "source_uuid",
	ProducerThingUUID:    "producer_thing_uuid",
	WriterThingUUID:      "writer_thing_uuid",
	ByPluginName:         "by_plugin_name",
	TimestampGt:          "timestamp_gt",
	TimestampLt:          "timestamp_lt",
	WithFlowNetworks:     "with_flow_networks",
	WithStreams:          "with_streams",
	WithStreamClones:     "with_stream_clones",
	WithProducers:        "with_producers",
	WithConsumers:        "with_consumers",
	WithCommandGroups:    "with_command_groups",
	WithWriters:          "with_writers",
	WithWriterClones:     "with_writer_clones",
	WithNetworks:         "with_networks",
	WithDevices:          "with_devices",
	WithPoints:           "with_points",
	WithPriority:         "with_priority",
	WithTags:             "with_tags",
	PluginName:           "plugin_name",
	NetworkName:          "network_name",
	DeviceName:           "device_name",
	PointName:            "point_name",
	AddressUUID:          "address_uuid",
	AddressID:            "address_id",
	ObjectType:           "object_type",
	IoNumber:             "io_number",
	IdGt:                 "id_gt",
	IsRemote:             "is_remote",
	IsMetadata:           "is_metadata",
	NetworkUUID:          "network_uuid",
	DeviceUUID:           "device_uuid",
	WithMetaTags:         "with_meta_tags",
	AutoMappingUUID:      "auto_mapping_uuid",
	AutoMappingEnable:    "auto_mapping_enable",
	MetaTags:             "meta_tags",
}

var ArgsDefault = struct {
	Sort              string
	Order             string
	Offset            string
	Limit             string
	Search            string
	AskRefresh        string
	AskResponse       string
	Write             string
	ThingType         string
	FlowNetworkUUID   string
	Field             string
	Value             string
	UpdateProducer    string
	CompactPayload    string
	CompactWithName   string
	FlowUUID          string
	StreamUUID        string
	ProducerUUID      string
	ConsumerUUID      string
	WriterUUID        string
	AddToParent       string
	PluginName        string
	WithFlowNetworks  string
	WithStreams       string
	WithStreamClones  string
	WithProducers     string
	WithConsumers     string
	WithCommandGroups string
	WithWriters       string
	WithWriterClones  string
	WithNetworks      string
	WithDevices       string
	WithPoints        string
	WithPriority      string
	WithTags          string
	NetworkName       string
	DeviceName        string
	PointName         string
	IsMetadata        string
	NetworkUUID       string
	DeviceUUID        string
	WithMetaTags      string
	AutoMappingUUID   string
	AutoMappingEnable string
	MetaTags          string
}{
	Sort:              "ID",
	Order:             "DESC", // ASC or DESC
	Offset:            "0",
	Limit:             "25",
	Search:            "",
	AskRefresh:        "false",
	AskResponse:       "false",
	Write:             "false",
	ThingType:         "point",
	FlowNetworkUUID:   "",
	Field:             "name",
	Value:             "",
	UpdateProducer:    "false",
	CompactPayload:    "false",
	CompactWithName:   "false",
	FlowUUID:          "",
	ProducerUUID:      "",
	ConsumerUUID:      "",
	WriterUUID:        "",
	AddToParent:       "",
	PluginName:        "false",
	WithFlowNetworks:  "false",
	WithStreams:       "false",
	WithStreamClones:  "false",
	WithProducers:     "false",
	WithConsumers:     "false",
	WithCommandGroups: "false",
	WithWriters:       "false",
	WithWriterClones:  "false",
	WithNetworks:      "false",
	WithDevices:       "false",
	WithPoints:        "false",
	WithPriority:      "false",
	WithTags:          "false",
	NetworkName:       "",
	DeviceName:        "",
	PointName:         "",
	IsMetadata:        "false",
	NetworkUUID:       "",
	DeviceUUID:        "",
	WithMetaTags:      "false",
	AutoMappingUUID:   "",
	AutoMappingEnable: "false",
	MetaTags:          "",
}
