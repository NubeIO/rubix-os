package main

import (
	"context"
	"encoding/json"
	"fmt"
	influxmodel "github.com/NubeDev/flow-framework/plugin/nube/database/influx/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
	"time"
)

const token = "biwutj0neNcj5oI_SbvIKMh7jFNf7Y_hRfApSmqpCKH9bpuHBRXLoJtyEZua2LmeYJWvDKxCisy0Kzc4_qYX2A=="
const bucket = "mydb"
const org = "mydb"

// Write function writes
func Write(t influxmodel.Temperature) {
	client := influxdb2.NewClient("http://localhost:8086", token)
	writeAPI := client.WriteAPI(org, bucket)
	p := influxdb2.NewPoint(Measurement(t), Tags(t), Fields(t), time.Now())
	writeAPI.WritePoint(p)
}

// Read functions reads all the temperatures saved inside of InfluxDB and returns them as array
func Read(measurement string) [][]byte {
	client := influxdb2.NewClient("http://localhost:8086", token)
	queryAPI := client.QueryAPI(org)
	fluxQuery := fmt.Sprintf(`from(bucket:"%v") |> range(start:-5) |> filter(fn:(r) => r._measurement == "%v")`, bucket, measurement)
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

	return temperaturesArray
}

func Tags(t influxmodel.Temperature) map[string]string {
	newMap := make(map[string]string)
	newMap["city"] = t.City
	newMap["country"] = t.Country
	newMap["temperature_scale"] = t.TemperatureScale

	return newMap
}

func Fields(t influxmodel.Temperature) map[string]interface{} {
	newMap := make(map[string]interface{})
	newMap["temperature_value"] = t.TemperatureValue

	return newMap
}

func Measurement(t influxmodel.Temperature) string {
	return "temperatures"
}
