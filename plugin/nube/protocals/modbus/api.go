package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

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

type Enable struct {
	Type        string
	ReadOnly    url.URL
	WriteValues interface{}
}

type T struct {
	Networks struct {
		Enable Urls        `json:"capabilities"`
		Body   interface{} `json:"body"`
	} `json:"networks"`
}

type Urls struct {
	Help       []string `json:"help"`
	ListSerial string   `json:"list_serial"`
}

type T2 struct {
	Capabilities struct {
		Help []string `json:"help,omitempty"`
	} `json:"capabilities,omitempty"`
	ObjectType struct {
		Options  interface{} `json:"options"`
		Type     string      `json:"type"`
		Required bool        `json:"required"`
	} `json:"object_type"`
}

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	objs := utils.ArrayValues(model.ObjectTypes)
	for _, obj := range objs {
		switch obj {
		case model.ObjectTypes.AnalogInput:
			out.Add(obj)
		case model.ObjectTypes.AnalogOutput:
			out.Add(obj)
		case model.ObjectTypes.AnalogValue:
			out.Add(obj)
		case model.ObjectTypes.BinaryInput:
			out.Add(obj)
		case model.ObjectTypes.BinaryOutput:
			out.Add(obj)
		case model.ObjectTypes.BinaryValue:
			out.Add(obj)
		default:
		}
	}
	return out
}

const (
	help       = "/modbus/help"
	listSerial = "/modbus/list/serial"
)

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(help, func(ctx *gin.Context) {

		var h T2
		//h.Enabled = model.CommonNaming
		//a1 := []string{fmt.Sprintf("http://0.0.0.0:1660/api/plugins/api/%s%s" ,name, help), "GET", "POST", "PATCH"}
		//a1 := []string{fmt.Sprintf("http://0.0.0.0:1660/api/plugins/api/%s%s" ,name, help), "GET", "POST", "PATCH"}
		//h.Capabilities.Help = a1
		h.ObjectType.Options = supportedObjects()
		h.ObjectType.Type = "array"
		h.ObjectType.Required = true

		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, h)
		}
	})
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/modbus/point/tcp/operation", func(ctx *gin.Context) {
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
	mux.POST("/modbus/wizard/tcp", func(ctx *gin.Context) {
		serial, err := i.wizardTCP()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST("/modbus/wizard/serial", func(ctx *gin.Context) {
		serial, err := i.wizardSerial()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
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
