package model

// PluginConf holds information about the plugin.
type PluginConf struct {
	UUID        string `json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	Name        string `json:"name" gorm:"name:varchar(255);not null"`
	ModulePath  string `json:"module_path"  gorm:"type:varchar(255);unique;not null"`
	Enabled     bool   `json:"enabled"`
	HasNetwork  bool   `json:"has_network"`
	Config      []byte
	Storage     []byte
	Network     Network     `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Integration Integration `json:"integration" gorm:"constraint:OnDelete:CASCADE"`
	Job         []Job       `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}

// PluginConfExternal Model
// Holds information about a plugin instance for one user.
type PluginConfExternal struct {
	UUID         string        `json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	Name         string        `json:"name" gorm:"type:varchar(255);not null"`
	ModulePath   string        `json:"module_path"  gorm:"type:varchar(255);unique;not null"`
	Author       string        `json:"author,omitempty" form:"author" query:"author"`
	Website      string        `json:"website,omitempty" form:"website" query:"website"`
	License      string        `json:"license,omitempty" form:"license" query:"license"`
	Enabled      bool          `json:"enabled"`
	HasNetwork   bool          `json:"has_network"`
	ProtocolType string        `json:"protocol_type" gorm:"type:varchar(255);not null"`
	Capabilities []string      `json:"capabilities"`
	Network      Network       `json:"networks" gorm:"constraint:OnDelete:CASCADE"`
	Integration  []Integration `json:"integration" gorm:"constraint:OnDelete:CASCADE"`
	Job          []Job         `json:"jobs" gorm:"constraint:OnDelete:CASCADE"`
}
