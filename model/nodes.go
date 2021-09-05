package model

import "gorm.io/datatypes"


type NodeAdd struct {
	NodeSettings 		datatypes.JSON  `json:"node_settings"`
}

type In1Connections struct {
	CommonUUID
	NodeUUID     		string  `json:"node_uuid" gorm:"TYPE:varchar(255) REFERENCES node_lists;not null;default:null"`
	FromUUID 			string  `json:"from_uuid"`
	Connection 			string  `json:"connection"`
}

type In2Connections struct {
	CommonUUID
	NodeUUID     		string  `json:"node_uuid" gorm:"TYPE:varchar(255) REFERENCES node_lists;not null;default:null"`
	FromUUID 			string  `json:"from_uuid"`
	Connection 			string  `json:"connection"`
}

type Out1Connections struct {
	CommonUUID
	NodeUUID     		string  `json:"node_uuid" gorm:"TYPE:varchar(255) REFERENCES nodes;not null;default:null"`
	ToUUID 				string  `json:"to_uuid"`
	Connection 			string  `json:"connection"`
}

type Out2Connections struct {
	CommonUUID
	NodeUUID     		string  `json:"node_uuid" gorm:"TYPE:varchar(255) REFERENCES nodes;not null;default:null"`
	ToUUID 				string  `json:"to_uuid"`
	Connection 			string  `json:"connection"`
}

//Node table
type Node struct {
	CommonUUID
	CommonName
	CommonNodeType
	CommonHelp
	In1					string `json:"in_1"`
	In2					string `json:"in_2"`
	In1Connections 		[]In1Connections 	`json:"in_1_connections" gorm:"constraint:OnDelete:CASCADE;"`
	In2Connections 		[]In2Connections 	`json:"in_2_connections" gorm:"constraint:OnDelete:CASCADE;"`
	Out1Connections 	[]Out1Connections 	`json:"out_1_connections" gorm:"constraint:OnDelete:CASCADE;"`
	Out2Connections 	[]Out2Connections 	`json:"out_2_connections" gorm:"constraint:OnDelete:CASCADE;"`
	Out1Value           string `json:"out_1_value"`
	Out2Value           string `json:"out_2_value"`
	NodeSettings 		datatypes.JSON  `json:"node_settings"`
}
