package model

// PluginConf holds information about the plugin.
type PluginConf struct {
	UUID	string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	UserID        uint
	ModulePath    string `gorm:"type:text"`
	Token         string `gorm:"type:varchar(180);unique_index"`
	ApplicationID uint
	Enabled       bool
	Config        []byte
	Storage       []byte
	Network 		Network `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Job 			[]Job `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}

// PluginConfExternal Model
//
// Holds information about a plugin instance for one user.
//
// swagger:model PluginConf
type PluginConfExternal struct {
	// The plugin id.
	//
	// read only: true
	// required: true
	// example: 25
	UUID	string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	// The plugin name.
	//
	// read only: true
	// required: true
	// example: RSS poller
	Name string `json:"name"`
	// The user name. For login.
	//
	// required: true
	// example: P1234
	Token string `binding:"required" json:"token" query:"token" form:"token"`
	// The module path of the plugin.
	//
	// example: github.com/flow/server/plugin/example/echo
	// read only: true
	// required: true
	ModulePath string `json:"modulePath" form:"modulePath" query:"modulePath"`
	// The author of the plugin.
	//
	// example: jmattheis
	// read only: true
	Author string `json:"author,omitempty" form:"author" query:"author"`
	// The website of the plugin.
	//
	// example: flow.net
	// read only: true
	Website string `json:"website,omitempty" form:"website" query:"website"`
	// The license of the plugin.
	//
	// example: MIT
	// read only: true
	License string `json:"license,omitempty" form:"license" query:"license"`
	// Whether the plugin instance is enabled.
	//
	// example: true
	// required: true
	Enabled bool `json:"enabled"`
	// Capabilities the plugin provides
	//
	// example: ["webhook","display"]
	// required: true
	Capabilities []string `json:"capabilities"`
	Network 		Network `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Job 			[]Job `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}
