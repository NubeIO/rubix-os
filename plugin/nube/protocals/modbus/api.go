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
	Host    string        `json:"ip"`
	Port    string        `json:"port"`
	Timeout time.Duration `json:"device_timeout_in_ms"`
}

type Scan struct {
	Start  uint32 `json:"start"`
	Count  uint32 `json:"count"`
	IsCoil bool   `json:"is_coil"`
}

type Body struct {
	Client    `json:"client"`
	Operation `json:"request_body"`
	Scan      `json:"scan"`
}

func bodyClient(ctx *gin.Context) (dto Body, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

var restMB *modbus.ModbusClient
var connected bool

func setClient(client Client) error {
	var cli utils.URLParts
	cli.Transport = "tcp"
	cli.Host = client.Host
	cli.Port = client.Port
	url := utils.JoinURL(cli)

	if client.Timeout < 10 {
		client.Timeout = 500
	}
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     url,
		Timeout: client.Timeout * time.Millisecond,
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
	mux.POST("/point/tcp/operation", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := setClient(body.Client)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			request, err := parseRequest(body.Operation)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, "request was invalid, try readCoil or writeCoil")
				return
			}
			r, err := DoOperations(cli, request)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, r)
			}

		}
	})
	mux.POST("/scan/bool", func(ctx *gin.Context) {
		body, _ := bodyClient(ctx)
		err := setClient(body.Client)
		if err != nil {
			log.Info(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err)
		}
		cli := getClient()
		if !isConnected() {
			ctx.JSON(http.StatusBadRequest, "modbus not enabled")
		} else {
			found, _ := performBoolScan(cli, body.Scan.IsCoil, body.Scan.Start, body.Scan.Count)
			ctx.JSON(http.StatusOK, found)
		}

	})

}
