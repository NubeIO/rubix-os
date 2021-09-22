package model

type Tag struct {
	Tag       string      `json:"tag" gorm:"type:varchar(255);unique;not null;default:null;primaryKey"`
	Streams   []*Stream   `json:"streams,omitempty" gorm:"many2many:streams_tags;constraint:OnDelete:CASCADE"`
	Producers []*Producer `json:"producers,omitempty" gorm:"many2many:producers_tags;constraint:OnDelete:CASCADE"`
	Consumers []*Consumer `json:"consumers,omitempty" gorm:"many2many:consumers_tags;constraint:OnDelete:CASCADE"`
}
