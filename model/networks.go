package model



//https://project-haystack.org/doc/appendix/protocol

type Network struct {
	CommonUUID
	CommonNameUnique
	Common
	Created
	Manufacture 	string `json:"manufacture"`
	Model 			string `json:"model"`
	NetworkType		string `json:"network_type"`
	Device 			[]Device `json:"devices" gorm:"constraint:OnDelete:CASCADE;"`
}

