package model

// The User holds information about the credentials of a user and its application and client tokens.
type User struct {
	ID           uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name         string `gorm:"type:varchar(180);unique_index"`
	Pass         []byte
	Admin        bool
	Applications []Application
	Clients      []Client
	Plugins      []PluginConf
}

type UserExternal struct {
	ID    uint   `json:"id"`
	Name  string `binding:"required" json:"name" query:"name" form:"name"`
	Admin bool   `json:"admin" form:"admin" query:"admin"`
}

// UserExternalWithPass Model
type UserExternalWithPass struct {
	UserExternal
	UserExternalPass
}

// UserExternalPass Model
type UserExternalPass struct {
	Pass string `json:"pass,omitempty" form:"pass" query:"pass" binding:"required"`
}
