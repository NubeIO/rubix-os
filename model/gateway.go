package model




type Gateway struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	IsRemote 			bool `json:"is_remote"`
	Subscriber			[]Subscriber `json:"subscribers" gorm:"constraint:OnDelete:CASCADE;"`
	Subscriptions		[]Subscriptions `json:"subscriptions" gorm:"constraint:OnDelete:CASCADE;"`
}