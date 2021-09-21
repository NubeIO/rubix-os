package model

type Tag struct {
	Tag     string    `json:"tag" gorm:"type:varchar(255);unique;not null;default:null;primaryKey"`
	Streams []*Stream `json:"streams,omitempty" gorm:"many2many:streams_tags;constraint:OnDelete:CASCADE"`
}
