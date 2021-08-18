package model

/*
Example tag usage
pointUUID:123
pointName: test
tagRef:  {
	tagRef uuid-123 //user selected these tags
	tagRef uuid-456
}
TaggedRef-Table
uuid-123 [Tags.temp, Tags.room]
uuid-456 [Tags.other]
 */



//TaggedRef is an item that has a tags added to and for example a point would use this tagRef fo when a new tag is added it can be used
type TaggedRef struct {
	Uuid		string `json:"uuid" gorm:"type:varchar(255);unique;not null;default:null;primaryKey"`
	Key 		string //roomTemp
	Tags  		[]Tags  //temp, room
}


type Tags struct {
	Tag string //temp, room

}


