package model

import "time"

type CommonPoint struct {
	Writeable 		bool   `json:"writeable"`
	Cov  			float64 `json:"cov"`
	ObjectType		string `json:"object_type"`
}

//type EquipRef struct {
//	PointUuid     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
//	Ref  			string `json:"equip_ref"`
//}
//
//type Association struct {
//	PointUuid     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
//	Association  	string `json:"equip_ref"`
//}


type Point struct {
	Uuid						string `json:"uuid" gorm:"type:varchar(255);unique;not null;default:null;primaryKey"`
	CommonName
	Common
	Created
	DeviceUuid     				string `json:"device_uuid" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	CommonPoint
	//EquipRef 					[]EquipRef `json:"equip_ref" gorm:"default:null"`
	//PriorityArrayModel 			PriorityArrayModel `json:"priority_array" gorm:"constraint:OnDelete:CASCADE"`
	//PointStore 					PointStore `json:"point_store" gorm:"constraint:OnDelete:CASCADE"`
}


type PointStore struct {
	PointUuid     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	Created
	Value  			string `json:"value" sql:"DEFAULT:NULL"`

}

type PriorityArrayModel struct {
	PointUuid     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	Created
	P1  			string `json:"_1"`
	P1Ts			time.Time  `json:"ts_1_updated_on"`
	P2  			string `json:"_2"`
	P2Ts			time.Time  `json:"ts_2_updated_on"`
	P3  			string `json:"_3"`
	P4  			string `json:"_4"`
	P5  			string `json:"_5"`
	P6  			string `json:"_6"`
	P7  			string `json:"_7"`
	P8  			string `json:"_8"`
	P9  			string `json:"_9"`
	P10  			string `json:"_10"`
	P11  			string `json:"_11"`
	P12  			string `json:"_12"`
	P13  			string `json:"_13"`
	P14  			string `json:"_14"`
	P15  			string `json:"_15"`
	P16  			string `json:"_16"`
	P17  			string `json:"_17"`
}
