package model

type CommonDevice struct {
	Manufacture string `json:"manufacture"` // nube
	Model       string `json:"model"`       // thml
	AddressId   int    `json:"address_id"`  // for example a modbus address or bacnet address
	ZeroMode    *bool  `json:"zero_mode"`
	AddressUUID string `json:"address_uuid"` // AAB1213
	CommonIP
	TransportType string `json:"transport_type"  gorm:"type:varchar(255);not null"` //serial, ip

}

type Device struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	CommonThingClass //point, job
	CommonThingUse
	CommonThingType //for example temp, rssi, voltage
	CommonDevice
	NetworkUUID string   `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Points      []*Point `json:"points" gorm:"constraint:OnDelete:CASCADE"`
	Tags        []*Tag   `json:"tags" gorm:"many2many:devices_tags;constraint:OnDelete:CASCADE"`
}
