package model

//Client holds information about a device which can receive notifications (and other stuff).
type Client struct {
	ID     uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"id"`
	Token  string `gorm:"type:varchar(180);unique_index" json:"token"`
	UserID uint   `gorm:"index" json:"-"`
	Name   string `gorm:"type:text" form:"name" query:"name" json:"name" binding:"required"`
}
