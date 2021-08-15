package model

import "time"

type CommonName struct {
	Name        			string `json:"name" validate:"min=1,max=255"  gorm:"type:varchar(255);not null"`
}

type CommonNameUnique struct {
	Name        			string `json:"name" validate:"min=1,max=255"  gorm:"type:varchar(255);unique;not null"`

}

type CommonUUID struct {
	Uuid			string 		`json:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
}

type Created struct {
	CreatedAt 				time.Time `json:"created_on"`
	UpdatedAt 				time.Time  `json:"updated_on"`
}


type Common struct {
	Description 			string `json:"description"`
	Id						string `json:"id"`
	Enable 					bool `json:"enable"`
	Fault					bool  `json:"fault"`
	FaultMessage 			string  `json:"fault_message"`
	EnableHistory 			bool   `json:"history_enable"`
}

