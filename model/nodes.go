package model

import "gorm.io/datatypes"

// Node table
type Node struct {
	CommonUUID
	CommonName
	CommonType
	NodeSettings datatypes.JSON  `json:"node_settings"`
	In1 string `json:"in_1"`
	In1FromUUID string `json:"in_1_from_uuid"`
	In2 string `json:"in_2"`
	Out1 string `json:"out_1"`
	Out2 string `json:"out_2"`

}


// NodeAdd table
type NodeAdd struct {
	CommonName
	CommonType
	Inputs 				datatypes.JSON  `json:"inputs"`
	Output 				datatypes.JSON  `json:"output"`
	NodeSettings 		datatypes.JSON  `json:"node_settings"`
	DataStore 			datatypes.JSON  `json:"data_store"`
}

type NodePayload struct {
	Payload string `json:"payload"`
}


type NodeIn1 struct {
	CommonUUID
	NodeListUUID     		string  `json:"node_list_uuid" gorm:"TYPE:varchar(255) REFERENCES node_lists;not null;default:null"`
	FromUUID 				string  `json:"from_uuid"`
	Connection 				string  `json:"connection"`
}

type NodeOut1 struct {
	CommonUUID
	NodeListUUID     		string  `json:"node_list_uuid" gorm:"TYPE:varchar(255) REFERENCES node_lists;not null;default:null"`
	ToUUID 					string  `json:"to_uuid"`
	Connection 				string  `json:"connection"`
}

//NodeList table
type NodeList struct {
	CommonUUID
	CommonName
	CommonNodeType
	CommonAction
	In1					string `json:"in_1"`
	In1FromUUID         string `json:"in_1_from_uuid"`
	In2					string `json:"in_2"`
	In2FromUUID 		string `json:"in_2_from_uuid"`
	NodeIn1 			[]NodeIn1 	`json:"node_in_1" gorm:"constraint:OnDelete:CASCADE;"`
	NodeOut1 			[]NodeOut1 `json:"node_out_1" gorm:"constraint:OnDelete:CASCADE;"`
	NodeSettings 		datatypes.JSON  `json:"node_settings"`
	Out1Value string `json:"out_1_value"`
}


//NodeBody could be a local network, job or alarm and so on
type NodeBody struct {
	Action 		string `json:"action"`  //read, write and so on
	AskRefresh 	bool `json:"ask_refresh"`
	CommonValue CommonValue `json:"common_value"`
	Priority 	Priority `json:"priority"`
	NodeAdd     NodeAdd `json:"node_add"`
	DataStore 			datatypes.JSON  `json:"data_store"`
}
