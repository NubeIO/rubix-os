package pgmodel

import (
	"time"
)

type History struct {
	ID        int       `json:"id" gorm:"primary_key"`
	UUID      string    `json:"uuid" gorm:"primary_key"`
	Value     float64   `json:"value" gorm:"primary_key"`
	Timestamp time.Time `json:"timestamp" gorm:"primary_key"`
}

type Point struct {
	UUID         string `json:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
	Name         string `json:"name"`
	DeviceUUID   string `json:"device_uuid,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	NetworkUUID  string `json:"network_uuid"`
	NetworkName  string `json:"network_name"`
	GlobalUUID   string `json:"global_uuid"`
	HostUUID     string `json:"host_uuid"`
	HostName     string `json:"host_name"`
	GroupUUID    string `json:"group_uuid"`
	GroupName    string `json:"group_name"`
	LocationUUID string `json:"location_uuid"`
	LocationName string `json:"location_name"`
}

type NetworkTag struct {
	NetworkUUID string `json:"network_uuid,omitempty" gorm:"primaryKey"`
	Tag         string `json:"tag" gorm:"primaryKey"`
}

type DeviceTag struct {
	DeviceUUID string `json:"device_uuid,omitempty" gorm:"primaryKey"`
	Tag        string `json:"tag" gorm:"primaryKey"`
}

type PointTag struct {
	PointUUID string `json:"point_uuid,omitempty" gorm:"primaryKey"`
	Tag       string `json:"tag" gorm:"primaryKey"`
}

type NetworkMetaTag struct {
	NetworkUUID string `json:"network_uuid,omitempty" gorm:"primaryKey"`
	Key         string `json:"key,omitempty" gorm:"primaryKey"`
	Value       string `json:"value,omitempty"`
}

type DeviceMetaTag struct {
	DeviceUUID string `json:"device_uuid,omitempty" gorm:"primaryKey"`
	Key        string `json:"key,omitempty" gorm:"primaryKey"`
	Value      string `json:"value,omitempty"`
}

type PointMetaTag struct {
	PointUUID string `json:"point_uuid,omitempty" gorm:"primaryKey"`
	Key       string `json:"key,omitempty" gorm:"primaryKey"`
	Value     string `json:"value,omitempty"`
}

type HistoryData struct {
	HostData
	NetworkUUID string    `json:"network_uuid"`
	NetworkName string    `json:"network_name"`
	DeviceUUID  string    `json:"device_uuid"`
	DeviceName  string    `json:"device_name"`
	PointUUID   string    `json:"point_uuid"`
	PointName   string    `json:"point_name"`
	Value       float64   `json:"value"`
	Timestamp   time.Time `json:"timestamp"`
}

type HostData struct {
	LocationUUID string `json:"location_uuid,omitempty"`
	LocationName string `json:"location_name,omitempty"`
	GroupUUID    string `json:"group_uuid,omitempty"`
	GroupName    string `json:"group_name,omitempty"`
	HostUUID     string `json:"host_uuid,omitempty"`
	HostName     string `json:"host_name,omitempty"`
	GlobalUUID   string `json:"global_uuid,omitempty"`
}

type HistoryDataResponse struct {
	HostData
	NetworkUUID string             `json:"network_uuid"`
	NetworkName string             `json:"network_name"`
	DeviceUUID  string             `json:"device_uuid"`
	DeviceName  string             `json:"device_name"`
	PointUUID   string             `json:"point_uuid"`
	PointName   string             `json:"point_name"`
	Histories   []*HistoryResponse `json:"histories"`
}

type HistoryResponse struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
