package model




type Gateway struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	IsRemote 			bool `json:"is_remote"`
	CommonRubixUUID
	Subscriber			[]Subscriber `json:"subscribers" gorm:"constraint:OnDelete:CASCADE;"`
	Subscription		[]Subscription `json:"subscriptions" gorm:"constraint:OnDelete:CASCADE;"`
}