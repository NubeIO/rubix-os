package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type InfluxSetting struct {
	Org         string
	Bucket      string
	ServerURL   string
	AuthToken   string
	Measurement string
}

func New(s *InfluxSetting) *InfluxSetting {
	if s.Org == "" {
		s.Org = "mydb"
	}
	if s.Bucket == "" {
		s.Bucket = "mydb"
	}
	if s.ServerURL == "" {
		s.ServerURL = "http://localhost:8086"
	}
	if s.Measurement == "" {
		s.Measurement = "points"
	}
	return &InfluxSetting{
		Org:         s.Org,
		Bucket:      s.Bucket,
		ServerURL:   s.ServerURL,
		AuthToken:   s.AuthToken,
		Measurement: s.Measurement,
	}
}

// WriteHistories function writes histories
func (i *InfluxSetting) WriteHistories(tags map[string]string, fields map[string]interface{}, ts time.Time) {
	client := influxdb2.NewClient(i.ServerURL, i.AuthToken)
	writeAPI := client.WriteAPI(i.Org, i.Bucket)
	point := influxdb2.NewPoint(
		i.Measurement,
		tags,
		fields,
		ts)
	writeAPI.WritePoint(point)
	writeAPI.Flush()
	client.Close()
}

// Read functions reads all the histories saved inside of InfluxDB and returns them as array
func (i *InfluxSetting) Read(measurement string) [][]byte {
	client := influxdb2.NewClient(i.ServerURL, i.AuthToken)
	queryAPI := client.QueryAPI(i.Org)
	fluxQuery := fmt.Sprintf(`from(bucket:"%v") |> range(start:-5) |> filter(fn:(r) => r._measurement == "%v")`, i.Bucket, measurement)
	logrus.Infof("FLUX QUERY: %v", fluxQuery)
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		logrus.Errorf("Error :%v", err)
		panic(err)
	}
	var temperaturesArray [][]byte
	for result.Next() {
		j, err := json.Marshal(result.Record().Values())
		if err != nil {
			panic(err)
		}
		temperaturesArray = append(temperaturesArray, j)
	}
	client.Close()
	return temperaturesArray
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
