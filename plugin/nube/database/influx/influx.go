package main

import (
	"context"
	"encoding/json"
	"fmt"
	influxmodel "github.com/NubeIO/flow-framework/plugin/nube/database/influx/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

type InfluxSetting struct {
	Org       string
	Bucket    string
	ServerURL string
	AuthToken string
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
	return &InfluxSetting{
		Org:       s.Org,
		Bucket:    s.Bucket,
		ServerURL: s.ServerURL,
		AuthToken: s.AuthToken,
	}
}

// WriteHist function writes histories
func (i *InfluxSetting) WriteHist(t influxmodel.HistPayload) {
	client := influxdb2.NewClient(i.ServerURL, i.AuthToken)
	writeAPI := client.WriteAPI(i.Org, i.Bucket)
	p := influxdb2.NewPoint(
		measurementHist(),
		tagsHist(t),
		fieldsHist(t),
		t.Timestamp)
	writeAPI.WritePoint(p)
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

func tagsHist(t influxmodel.HistPayload) map[string]string {
	newMap := make(map[string]string)
	newMap["uuid"] = t.UUID
	return newMap
}

func fieldsHist(t influxmodel.HistPayload) map[string]interface{} {
	newMap := make(map[string]interface{})
	newMap["value"] = t.Value
	return newMap
}

func measurementHist() string {
	return "hist"
}
