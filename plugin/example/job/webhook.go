package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/mustafaturan/bus/v3"
	"github.com/tidwall/gjson"
)

type messageBody struct {
	JobUUID  	string   `json:"job_uuid"`
	Enable 		bool   `json:"enable"`

}

type payloadBody struct {
	UUID  			string   		`json:"uuid"`
	Delete  		bool   			`json:"delete"`
	MessageString  	string   		`json:"message_string"`

}


var handler = bus.Handler {
	Handle: func(ctx context.Context, e bus.Event) {
		fmt.Printf(e.Topic)
		fmt.Printf("inside plugin handler")
		fmt.Println(e.Data)
	},
	Matcher: ".*", // matches all topics
}

func getBODY(ctx *gin.Context) (dto *messageBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}


// RegisterWebhook implements plugin.Webhooker
func (c *PluginTest) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	mux.PATCH("/job/subscribe/:uuid", func(ctx *gin.Context) {
		body, _ := getBODY(ctx)
		uuid := ctx.Param("uuid")
		client := resty.New()
		address := "http://0.0.0.0"
		port := "1660"
		url := fmt.Sprintf("%s:%s", address, port)
		urlClient := fmt.Sprintf("%s/%s", url, "client")
		//get token
		getToken, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"name":"admin"}`).
			SetBasicAuth("admin", "admin").
			Post(urlClient)
		r := gjson.Get(string(getToken.Body()), "token")
		if err != nil {
			fmt.Println(err)
		}
		token := r.Str
		urlJobs := fmt.Sprintf("%s/%s", url, "api/jobs/producer/{uuid}")
		//edit job
		subJob, _ := client.NewRequest().
			SetHeader("Authorization", token).
			SetBody(map[string]interface{}{"enable": body.Enable}).
			SetPathParams(map[string]string{
				"uuid": uuid,
			}).
			Patch(urlJobs)
		if subJob.Status() == "200 OK" {
			topic := fmt.Sprintf("%s:%s", "job",body.JobUUID)
			if body.Enable{
				// subscribe
				eventbus.BUS.RegisterHandler(topic, handler)
				ctx.JSON(200, "subscribe")
			} else {
				// unsubscribe
				eventbus.BUS.DeregisterHandler(topic)
				ctx.JSON(200, "unsubscribe")
			}
		} else {
			fmt.Println("subJob", "FAIL", subJob.Status())
			ctx.JSON(404, subJob.String())
		}
	})
}
