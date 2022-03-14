package main

import (
	"github.com/NubeIO/flow-framework/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"sync"
	"time"
)

var influxConnectionInstance *InfluxConnection
var influxConnectionInstances []*InfluxConnection
var once sync.Once

type InfluxSetting struct {
	ServerURL   string
	AuthToken   string
	Org         string
	Bucket      string
	Measurement string
}

type InfluxConnection struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func New(s *InfluxSetting) *InfluxSetting {
	return &InfluxSetting{
		ServerURL:   s.ServerURL,
		AuthToken:   s.AuthToken,
		Org:         s.Org,
		Bucket:      s.Bucket,
		Measurement: s.Measurement,
	}
}

func (i *InfluxSetting) getInfluxConnectionInstance() *InfluxConnection {
	once.Do(func() {
		client := influxdb2.NewClient(i.ServerURL, i.AuthToken)
		influxConnectionInstance = &InfluxConnection{
			client:   client,
			writeAPI: client.WriteAPI(i.Org, i.Bucket),
		}
		influxConnectionInstances = append(influxConnectionInstances, influxConnectionInstance)
	})
	return influxConnectionInstance
}

func (i *InfluxSetting) WriteHistories(tags map[string]string, fields map[string]interface{}, ts time.Time) {
	influxConnectionInstance := i.getInfluxConnectionInstance()
	point := influxdb2.NewPoint(i.Measurement, tags, fields, ts)
	influxConnectionInstance.writeAPI.WritePoint(point)
}

func tagsHistory(ht *model.HistoryInfluxTag) map[string]string {
	newMap := make(map[string]string)
	newMap["clint_id"] = ht.ClientId
	newMap["client_name"] = ht.ClientName
	newMap["site_id"] = ht.SiteId
	newMap["site_name"] = ht.SiteName
	newMap["device_id"] = ht.DeviceId
	newMap["device_name"] = ht.DeviceName
	newMap["rubix_point_uuid"] = ht.RubixPointUUID
	newMap["rubix_point_name"] = ht.RubixPointName
	newMap["rubix_device_uuid"] = ht.RubixDeviceUUID
	newMap["rubix_device_name"] = ht.RubixDeviceName
	newMap["rubix_network_uuid"] = ht.RubixNetworkUUID
	newMap["rubix_network_name"] = ht.RubixNetworkName
	return newMap
}

func fieldsHistory(t *model.History) map[string]interface{} {
	newMap := make(map[string]interface{})
	newMap["id"] = t.ID
	newMap["value"] = t.Value
	return newMap
}
