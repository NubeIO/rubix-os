package model

//Application holds information about an app which can send notifications.
type Application struct {
	ID          uint              `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"id"`
	Token       string            `gorm:"type:varchar(180);unique_index" json:"token"`
	UserID      uint              `gorm:"index" json:"-"`
	Name        string            `gorm:"type:text" form:"name" query:"name" json:"name" binding:"required"`
	Description string            `gorm:"type:text" form:"description" query:"description" json:"description"`
	Internal    bool              `form:"internal" query:"internal" json:"internal"`
	Image       string            `gorm:"type:text" json:"image"`
	Messages    []MessageExternal `json:"-"`
}
