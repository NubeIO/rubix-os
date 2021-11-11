package model

import (
	"time"
)

// Message holds information about a message.
type Message struct {
	ID       uint   `gorm:"AUTO_INCREMENT;primary_key;index"`
	Message  string `gorm:"type:text"`
	Title    string `gorm:"type:text"`
	Priority int
	Extras   []byte
	Date     time.Time
}

// MessageExternal Model
type MessageExternal struct {
	ID       uint                   `json:"id"`
	Message  string                 `form:"message" query:"message" json:"message" binding:"required"`
	Title    string                 `form:"title" query:"title" json:"title"`
	Priority int                    `form:"priority" query:"priority" json:"priority"`
	Extras   map[string]interface{} `form:"-" query:"-" json:"extras,omitempty" gorm:"type:text[]"`
	Date     time.Time              `json:"date"`
}
