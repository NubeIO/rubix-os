package main

import (
	lwmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/model"
	"github.com/gin-gonic/gin"
	"github.com/simonvetter/modbus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func bodyDevice(ctx *gin.Context) (dto lwmodel.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

var restMB  *modbus.ModbusClient
var connected bool

func setClient(url string) error {
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     url,
		Timeout:  1 * time.Second,
	})
	if err != nil {
		connected = false
		return err
	}
	connected = true
	err = c.Open()
	restMB = c
	if err != nil {
		connected = false
		return err
	}
	return nil
}

func getClient() *modbus.ModbusClient {
	return restMB
}

func isConnected() bool {
 return connected
}


// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.POST("/client", func(ctx *gin.Context) {
		err := setClient("tcp://192.168.15.202:502")
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, "ok")
		}
	})

	mux.POST("/read/coils", func(ctx *gin.Context) {
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		}
		var o operation
		o.op = readBools
		o.isCoil = true
		o.addr = uint16(1)
		o.quantity = uint16(1)
		r, err := operations(cli,  o)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, r)


	})

}
