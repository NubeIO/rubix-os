package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
)

const token = "biwutj0neNcj5oI_SbvIKMh7jFNf7Y_hRfApSmqpCKH9bpuHBRXLoJtyEZua2LmeYJWvDKxCisy0Kzc4_qYX2A=="
const bucket = "mydb"
const org = "mydb"

func main() {

	min := 10
	max := 30
	ran := rand.Intn(max-min) + min
	ran2 := rand.Intn(max-min) + 15
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://localhost:8086", token)
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": ran, "max": ran2},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)
	// Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", ran).
		AddField("max", ran2).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)

	// Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 2.55, 6.5)
	writeAPI.WriteRecord(context.Background(), line)

	// Get query client
	queryAPI := client.QueryAPI(org)
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"my-bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	// Ensures background processes finishes
	client.Close()
}
