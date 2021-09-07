package model

// PluginConf holds information about the plugin.
type PluginConf struct {
	UUID          string `json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	UserID        uint
	Name          string `json:"name"`
	ModulePath    string `json:"module_path" gorm:"type:text"`
	Token         string `gorm:"type:varchar(180);unique_index"`
	ApplicationID uint
	Enabled       bool `json:"enabled"`
	Config        []byte
	Storage       []byte
	Network       Network     `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Integration   Integration `json:"integration" gorm:"constraint:OnDelete:CASCADE"`
	Job           []Job       `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}

// PluginConfExternal Model
// Holds information about a plugin instance for one user.
type PluginConfExternal struct {
	UUID         string        `json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	Name         string        `json:"name"`
	Token        string        `binding:"required" json:"token" query:"token" form:"token"`
	ModulePath   string        `json:"module_path" form:"modulePath" query:"modulePath"`
	Author       string        `json:"author,omitempty" form:"author" query:"author"`
	Website      string        `json:"website,omitempty" form:"website" query:"website"`
	License      string        `json:"license,omitempty" form:"license" query:"license"`
	Enabled      bool          `json:"enabled"`
	Capabilities []string      `json:"capabilities"`
	Network      Network       `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Integration  []Integration `json:"integration" gorm:"constraint:OnDelete:CASCADE"`
	Job          []Job         `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}
