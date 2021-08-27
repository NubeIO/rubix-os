package eventbus

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/go-resty/resty/v2"
)

func publishHTTP(sensorStruct *model.Point) {
	client := resty.New()
	resp, err := client.R().SetPathParams(map[string]string{
		"name": sensorStruct.Name,
	}).Post("http://0.0.0.0:8080/stream/{name}")
	fmt.Println(sensorStruct.Name, resp.String())
	fmt.Println(sensorStruct.Name, err)
}

