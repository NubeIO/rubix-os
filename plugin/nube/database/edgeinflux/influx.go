package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var influxConnectionInstances []*InfluxConnection

type InfluxSetting struct {
	ServerURL                string
	AuthToken                string
	Org                      string
	Bucket                   string
	Measurement              string
	once                     sync.Once
	influxConnectionInstance *InfluxConnection
}

type InfluxConnection struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func New(s *InfluxSetting) *InfluxSetting {
	return &InfluxSetting{
		ServerURL:                s.ServerURL,
		AuthToken:                s.AuthToken,
		Org:                      s.Org,
		Bucket:                   s.Bucket,
		Measurement:              s.Measurement,
		once:                     sync.Once{},
		influxConnectionInstance: nil,
	}
}

func (i *InfluxSetting) getInfluxConnectionInstance() *InfluxConnection {
	i.once.Do(func() {
		client := influxdb2.NewClient(i.ServerURL, i.AuthToken)
		i.influxConnectionInstance = &InfluxConnection{
			client:   client,
			writeAPI: client.WriteAPI(i.Org, i.Bucket),
		}
		influxConnectionInstances = append(influxConnectionInstances, i.influxConnectionInstance)
	})
	return i.influxConnectionInstance
}

func (i *InfluxSetting) WriteHistories(tags map[string]string, fields map[string]interface{}, ts time.Time) {
	influxConnectionInstance := i.getInfluxConnectionInstance()
	point := influxdb2.NewPoint(i.Measurement, tags, fields, ts)
	influxConnectionInstance.writeAPI.WritePoint(point)
}

func (i *InfluxSetting) GetLastSyncId() (value int, isError bool) {
	client := i.getInfluxConnectionInstance().client
	queryAPI := client.QueryAPI(i.Org)
	fluxQuery := fmt.Sprintf(
		`from(bucket: "%v")
				  |> range(start:-1)
				  |> filter(fn: (r) => r["_measurement"] == "%v")
				  |> filter(fn: (r) => r["_field"] == "id")
				  |> aggregateWindow(every: 1y, fn: max, createEmpty: false)
				  |> yield(name: "max")`, i.Bucket, i.Measurement)
	log.Debugf("Flux Query: %s", fluxQuery)
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Errorf("Error :%v", err)
		return 0, true
	}
	value = 0
	for result.Next() {
		values := result.Record().Values()
		value = int(values["_value"].(int64))
	}
	return value, false
}

func fieldsHistory(t *History) map[string]interface{} {
	newMap := make(map[string]interface{})
	newMap["value"] = t.Value
	return newMap
}
