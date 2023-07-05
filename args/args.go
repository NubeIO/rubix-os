package args

type Args struct {
	Enable                         *bool   `json:"enable,omitempty"`
	Sort                           string  `json:"sort,omitempty"`
	Order                          string  `json:"order,omitempty"`
	Offset                         string  `json:"offset,omitempty"`
	Limit                          string  `json:"limit,omitempty"`
	Search                         string  `json:"search,omitempty"`
	AskRefresh                     string  `json:"ask_refresh,omitempty"`
	AskResponse                    string  `json:"ask_response,omitempty"`
	Write                          string  `json:"write,omitempty"`
	ThingType                      string  `json:"thing_type,omitempty"`
	UUID                           *string `json:"uuid,omitempty"`
	WriteHistory                   string  `json:"write_history,omitempty"`
	Field                          string  `json:"field,omitempty"`
	Value                          string  `json:"value,omitempty"`
	CompactPayload                 string  `json:"compact_payload,omitempty"`
	CompactWithName                string  `json:"compact_with_name,omitempty"`
	Name                           *string `json:"name,omitempty"`
	AddToParent                    string  `json:"add_to_parent,omitempty"`
	GlobalUUID                     *string `json:"global_uuid,omitempty"`
	ClientId                       *string `json:"client_id,omitempty"`
	SiteId                         *string `json:"site_id,omitempty"`
	DeviceId                       *string `json:"device_id,omitempty"`
	SourceUUID                     *string `json:"source_uuid,omitempty"`
	ByPluginName                   bool    `json:"by_plugin_name,omitempty"`
	TimestampGt                    *string `json:"timestamp_gt,omitempty"`
	TimestampLt                    *string `json:"timestamp_lt,omitempty"`
	Networks                       bool    `json:"networks,omitempty"`
	WithDevices                    bool    `json:"with_devices,omitempty"`
	WithPoints                     bool    `json:"with_points,omitempty"`
	WithPriority                   bool    `json:"with_priority,omitempty"`
	WithTags                       bool    `json:"with_tags,omitempty"`
	PluginName                     string  `json:"plugin_name,omitempty"`
	NetworkName                    *string `json:"network_name,omitempty"`
	DeviceName                     *string `json:"device_name,omitempty"`
	PointName                      *string `json:"point_name,omitempty"`
	AddressUUID                    *string `json:"address_uuid,omitempty"`
	AddressID                      *string `json:"address_id,omitempty"`
	ObjectType                     *string `json:"object_type,omitempty"`
	IoNumber                       *string `json:"io_number,omitempty"`
	IdGt                           *string `json:"id_gt,omitempty"`
	IsRemote                       *bool   `json:"is_remote,omitempty"`
	IsMetadata                     bool    `json:"is_metadata,omitempty"`
	NetworkUUID                    *string `json:"network_uuid,omitempty"`
	DeviceUUID                     *string `json:"device_uuid,omitempty"`
	WithMetaTags                   bool    `json:"with_meta_tags,omitempty"`
	MetaTags                       *string `json:"meta_tags,omitempty"`
	MemberUUID                     *string `json:"member_uuid,omitempty"`
	ShowCloneNetworks              bool    `json:"show_clone_networks,omitempty"`
	PointSourceUUID                *string `json:"point_source_uuid,omitempty"`
	HostUUID                       *string `json:"host_uuid,omitempty"`
	WithMembers                    bool    `json:"with_members,omitempty"`
	WithMemberDevices              bool    `json:"with_member_devices,omitempty"`
	WithTeams                      bool    `json:"with_teams,omitempty"`
	WithViews                      bool    `json:"with_views,omitempty"`
	WithGroups                     bool    `json:"with_groups,omitempty"`
	WithHosts                      bool    `json:"with_hosts,omitempty"`
	WithComments                   bool    `json:"with_comments,omitempty"`
	WithWidgets                    bool    `json:"with_widgets,omitempty"`
	WithViewTemplateWidgets        bool    `json:"with_view_template_widgets,omitempty"`
	WithViewTemplateWidgetPointers bool    `json:"with_view_template_widget_pointers,omitempty"`
}

var ArgsType = struct {
	Sort                           string
	Order                          string
	Offset                         string
	Limit                          string
	Search                         string
	AskRefresh                     string
	AskResponse                    string
	Write                          string
	ThingType                      string
	UUID                           string
	WriteHistory                   string
	Field                          string
	Name                           string
	Value                          string
	CompactPayload                 string // for a point would be presentValue
	CompactWithName                string // for a point would be presentValue and pointName
	AddToParent                    string
	GlobalUUID                     string
	ClientId                       string
	SiteId                         string
	DeviceId                       string
	SourceUUID                     string
	ByPluginName                   string
	TimestampGt                    string
	TimestampLt                    string
	WithNetworks                   string
	WithDevices                    string
	WithPoints                     string
	WithPriority                   string
	WithTags                       string
	PluginName                     string
	NetworkName                    string
	DeviceName                     string
	PointName                      string
	AddressUUID                    string
	AddressID                      string
	ObjectType                     string
	IoNumber                       string
	IdGt                           string
	IsRemote                       string
	IsMetadata                     string
	NetworkUUID                    string
	DeviceUUID                     string
	WithMetaTags                   string
	MetaTags                       string
	ShowCloneNetworks              string
	PointSourceUUID                string
	HostUUID                       string
	WithMembers                    string
	WithMemberDevices              string
	WithTeams                      string
	WithViews                      string
	WithGroups                     string
	WithHosts                      string
	WithComments                   string
	WithWidgets                    string
	WithViewTemplateWidgets        string
	WithViewTemplateWidgetPointers string
}{
	Sort:                           "sort",
	Order:                          "order",
	Offset:                         "offset",
	Limit:                          "limit",
	Search:                         "search",
	AskRefresh:                     "ask_refresh",
	AskResponse:                    "ask_response",
	Write:                          "write",
	ThingType:                      "thing_type", // the type of thing like a point
	UUID:                           "uuid",
	WriteHistory:                   "write_history",
	Field:                          "field",
	Name:                           "name",
	Value:                          "value",
	CompactPayload:                 "compact_payload",
	CompactWithName:                "compact_with_name",
	AddToParent:                    "add_to_parent",
	GlobalUUID:                     "global_uuid",
	ClientId:                       "client_id",
	SiteId:                         "site_id",
	DeviceId:                       "device_id",
	SourceUUID:                     "source_uuid",
	ByPluginName:                   "by_plugin_name",
	TimestampGt:                    "timestamp_gt",
	TimestampLt:                    "timestamp_lt",
	WithNetworks:                   "with_networks",
	WithDevices:                    "with_devices",
	WithPoints:                     "with_points",
	WithPriority:                   "with_priority",
	WithTags:                       "with_tags",
	PluginName:                     "plugin_name",
	NetworkName:                    "network_name",
	DeviceName:                     "device_name",
	PointName:                      "point_name",
	AddressUUID:                    "address_uuid",
	AddressID:                      "address_id",
	ObjectType:                     "object_type",
	IoNumber:                       "io_number",
	IdGt:                           "id_gt",
	IsRemote:                       "is_remote",
	IsMetadata:                     "is_metadata",
	NetworkUUID:                    "network_uuid",
	DeviceUUID:                     "device_uuid",
	WithMetaTags:                   "with_meta_tags",
	MetaTags:                       "meta_tags",
	ShowCloneNetworks:              "show_clone_networks",
	PointSourceUUID:                "point_source_uuid",
	HostUUID:                       "host_uuid",
	WithMembers:                    "with_members",
	WithMemberDevices:              "with_member_devices",
	WithTeams:                      "with_teams",
	WithViews:                      "with_views",
	WithGroups:                     "with_groups",
	WithHosts:                      "with_hosts",
	WithComments:                   "with_comments",
	WithWidgets:                    "with_widgets",
	WithViewTemplateWidgets:        "with_view_template_widgets",
	WithViewTemplateWidgetPointers: "with_view_template_widget_pointers",
}

var ArgsDefault = struct {
	Sort                           string
	Order                          string
	Offset                         string
	Limit                          string
	Search                         string
	AskRefresh                     string
	AskResponse                    string
	Write                          string
	ThingType                      string
	Field                          string
	Value                          string
	CompactPayload                 string
	CompactWithName                string
	AddToParent                    string
	PluginName                     string
	WithNetworks                   string
	WithDevices                    string
	WithPoints                     string
	WithPriority                   string
	WithTags                       string
	NetworkName                    string
	DeviceName                     string
	PointName                      string
	IsMetadata                     string
	NetworkUUID                    string
	DeviceUUID                     string
	WithMetaTags                   string
	MetaTags                       string
	ShowCloneNetworks              string
	WithMembers                    string
	WithMemberDevices              string
	WithTeams                      string
	WithViews                      string
	WithGroups                     string
	WithHosts                      string
	WithComments                   string
	WithWidgets                    string
	WithViewTemplateWidgets        string
	WithViewTemplateWidgetPointers string
}{
	Sort:                           "ID",
	Order:                          "DESC", // ASC or DESC
	Offset:                         "0",
	Limit:                          "25",
	Search:                         "",
	AskRefresh:                     "false",
	AskResponse:                    "false",
	Write:                          "false",
	ThingType:                      "point",
	Field:                          "name",
	Value:                          "",
	CompactPayload:                 "false",
	CompactWithName:                "false",
	AddToParent:                    "",
	PluginName:                     "false",
	WithNetworks:                   "false",
	WithDevices:                    "false",
	WithPoints:                     "false",
	WithPriority:                   "false",
	WithTags:                       "false",
	NetworkName:                    "",
	DeviceName:                     "",
	PointName:                      "",
	IsMetadata:                     "false",
	NetworkUUID:                    "",
	DeviceUUID:                     "",
	WithMetaTags:                   "false",
	MetaTags:                       "",
	ShowCloneNetworks:              "false",
	WithMembers:                    "false",
	WithMemberDevices:              "false",
	WithTeams:                      "false",
	WithViews:                      "false",
	WithGroups:                     "false",
	WithHosts:                      "false",
	WithComments:                   "false",
	WithWidgets:                    "false",
	WithViewTemplateWidgets:        "false",
	WithViewTemplateWidgetPointers: "false",
}
