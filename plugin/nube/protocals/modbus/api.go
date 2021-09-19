package main

import (
	"github.com/NubeDev/flow-framework/utils"
	"github.com/gin-gonic/gin"
	"github.com/simonvetter/modbus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Client struct {
	Host string `json:"ip"`
	Port string `json:"port"`
}

type Bool struct {
	Client    `json:"client"`
	Operation `json:"request_body"`
}

func bodyClient(ctx *gin.Context) (dto Bool, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

var restMB *modbus.ModbusClient
var connected bool

func setClient(u utils.URLParts) error {
	url := utils.JoinURL(u)
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     url,
		Timeout: 1 * time.Second,
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

	mux.POST("/bool", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		var u utils.URLParts
		u.Transport = "tcp"
		u.Host = body.Client.Host
		u.Port = body.Client.Port
		err := setClient(u)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			var o Operation
			request, err := parseRequest(body.Operation)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "request was invalid, try readCoil or writeCoil")
				return
			}
			o.Op = request.Op
			o.IsCoil = request.IsCoil
			o.Addr = request.Addr
			o.Length = request.Length
			r, err := Operations(cli, o)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, r)
			}

		}

	})

}
