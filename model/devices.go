package model


type CommonDevice struct {
	Manufacture 	string `json:"manufacture"` // nube
	DeviceType		string `json:"device_type"` // droplet
	Model 			string `json:"model"` // thml
	AddressId     	int     `json:"address_id"`     // for example a modbus address or bacnet address
	AddressCode   	string  `json:"address_code"`  // thml

}

type Device struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	CommonDevice
	NetworkUUID     	string  `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Point 				[]Point `json:"points" gorm:"constraint:OnDelete:CASCADE"`

}

