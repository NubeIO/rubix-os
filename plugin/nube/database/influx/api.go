package main

import (
	"encoding/json"
	influxmodel "github.com/NubeDev/flow-framework/plugin/nube/database/influx/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyDevice(ctx *gin.Context) (dto influxmodel.Temperature, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET("/influx/temperatures", func(ctx *gin.Context) {
		records := Read("temperatures")
		var temperatures []influxmodel.FluxTemperature
		for _, data := range records {
			var temp influxmodel.FluxTemperature
			if err := json.Unmarshal(data, &temp); err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			}
			temperatures = append(temperatures, temp)
		}
		ctx.JSON(http.StatusOK, temperatures)

	})
	mux.POST("/influx/temperature", func(ctx *gin.Context) {
		body, err := bodyDevice(ctx)
		if err != nil {
			log.Info(err, "ERROR ON influx write")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			Write(body)
			ctx.JSON(http.StatusOK, "ok")
		}
	})

}
