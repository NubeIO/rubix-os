package main

type Point struct {
	UUID         string `json:"uuid" gorm:"type:varchar(255);primaryKey"`
	Name         string `json:"name"`
	DeviceUUID   string `json:"device_uuid,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	NetworkUUID  string `json:"network_uuid"`
	NetworkName  string `json:"network_name"`
	GlobalUUID   string `json:"global_uuid"`
	HostUUID     string `json:"host_uuid" gorm:"type:varchar(255);primaryKey"`
	HostName     string `json:"host_name"`
	GroupUUID    string `json:"group_uuid"`
	GroupName    string `json:"group_name"`
	LocationUUID string `json:"location_uuid"`
	LocationName string `json:"location_name"`
}

type DeviceMetaTag struct {
	DeviceUUID string `json:"device_uuid,omitempty" gorm:"primaryKey"`
	Key        string `json:"key,omitempty" gorm:"primaryKey"`
	Value      string `json:"value,omitempty"`
}

type DoorProcessingPoint struct {
	UUID                   string `json:"uuid" dataframe:"point_uuid" gorm:"column:point_uuid"`
	Name                   string `json:"name" dataframe:"name" gorm:"column:name"`
	HostUUID               string `json:"host_uuid" dataframe:"host_uuid" gorm:"column:host_uuid"`
	SiteRef                string `json:"site_ref" dataframe:"site_ref" gorm:"column:site_ref"`
	AssetRef               string `json:"assetRef" dataframe:"assetRef" gorm:"column:asset_ref"`
	AssetFunc              string `json:"assetFunc" dataframe:"assetFunc" gorm:"column:asset_func"`
	FloorRef               string `json:"floorRef" dataframe:"floorRef" gorm:"column:floor_ref"`
	GenderRef              string `json:"genderRef" dataframe:"genderRef" gorm:"column:gender_ref"`
	LocationRef            string `json:"locationRef" dataframe:"locationRef" gorm:"column:location_ref"`
	PointFunction          string `json:"pointFunction" dataframe:"pointFunction" gorm:"column:point_function"`
	MeasurementRef         string `json:"measurementRef" dataframe:"measurementRef" gorm:"column:measurement_ref"`
	DoorType               string `json:"doorType" dataframe:"doorType" gorm:"column:door_type"`
	NormalPosition         string `json:"normalPosition" dataframe:"normalPosition" gorm:"column:normal_position"`
	EnableCleaningTracking string `json:"enableCleaningTracking" dataframe:"enableCleaningTracking" gorm:"column:enable_cleaning_tracking"`
	EnableUseCounting      string `json:"enableUseCounting" dataframe:"enableUseCounting" gorm:"column:enable_use_counting"`
	IsEOT                  bool   `json:"isEOT" dataframe:"isEOT" gorm:"column:is_eot"`
	AvailabilityID         string `json:"availabilityID" dataframe:"availabilityID" gorm:"column:availability_id"`
	ResetID                string `json:"resetID" dataframe:"resetID" gorm:"column:reset_id"`
}

type DoorResetPoint struct {
	UUID           string `json:"uuid" dataframe:"point_uuid" gorm:"column:point_uuid"`
	Name           string `json:"name" dataframe:"name" gorm:"column:name"`
	HostUUID       string `json:"host_uuid" dataframe:"host_uuid" gorm:"column:host_uuid"`
	SiteRef        string `json:"site_ref" dataframe:"site_ref" gorm:"column:site_ref"`
	PointFunction  string `json:"pointFunction" dataframe:"pointFunction" gorm:"column:point_function"`
	MeasurementRef string `json:"measurementRef" dataframe:"measurementRef" gorm:"column:measurement_ref"`
	IsEOT          bool   `json:"isEOT" dataframe:"isEOT" gorm:"column:is_eot"`
	ResetID        string `json:"resetID" dataframe:"resetID" gorm:"column:reset_id"`
}

type History struct {
	ID        int     `json:"id" dataframe:"id" gorm:"primary_key"`
	PointUUID string  `json:"point_uuid" dataframe:"point_uuid" gorm:"primary_key"`
	HostUUID  string  `json:"host_uuid" dataframe:"host_uuid" gorm:"primary_key"`
	Value     float64 `json:"value" dataframe:"value"`
	Timestamp string  `json:"timestamp" dataframe:"timestamp" gorm:"column:timestamp"`
}

type LastProcessedData struct {
	DoorPosition           int `json:"door_position" dataframe:"door_position"`
	CubicleOccupancy       int `json:"cubicleOccupancy" dataframe:"cubicleOccupancy"`
	TotalUses              int `json:"totalUses" dataframe:"totalUses"`
	CurrentUses            int `json:"currentUses" dataframe:"currentUses"`
	PendingStatus          int `json:"pendingStatus" dataframe:"pendingStatus"`
	OverdueStatus          int `json:"overdueStatus" dataframe:"overdueStatus"`
	LastToPendingTimestamp string
	LastToCleanTimestamp   string
}

type DoorInfo struct {
	IsEOT                  bool   `json:"is_eot" dataframe:"is_eot"`
	DoorTypeTag            string `json:"doorType" dataframe:"doorType"`
	NormalPosition         DoorNormalPosition
	DoorTypeID             DoorType
	EnableCleaningTracking bool
	EnableUseCounting      bool
	AssetFunction          string
	AvailabilityID         string
	ResetID                string
}
