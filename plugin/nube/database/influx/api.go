package main

import (
	"encoding/json"
	influxmodel "github.com/NubeIO/flow-framework/plugin/nube/database/influx/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyDevice(ctx *gin.Context) (dto influxmodel.HistPayload, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	influxSetting := new(InfluxSetting)
	isc := New(influxSetting)
	i.basePath = basePath
	mux.GET("/influx/histories", func(ctx *gin.Context) {
		records := isc.Read("hist")
		var histories []influxmodel.HistPayload
		for _, data := range records {
			var temp influxmodel.HistPayload
			if err := json.Unmarshal(data, &temp); err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			}
			histories = append(histories, temp)
		}
		ctx.JSON(http.StatusOK, histories)
	})
	mux.POST("/influx/histories", func(ctx *gin.Context) {
		body, err := bodyDevice(ctx)
		if err != nil {
			log.Info(err, "ERROR ON influx write")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			isc.WriteHist(body)
			ctx.JSON(http.StatusOK, "ok")
		}
	})
}
