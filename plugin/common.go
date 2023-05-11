package plugin

import (
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

func GetFlowNetworkNames(fns []*model.FlowNetwork) []string {
	fnsNames := make([]string, 0)
	for _, fn := range fns {
		fnsNames = append(fnsNames, fn.Name)
	}
	return fnsNames
}

const (
	NetworksURL    = "/networks"
	DevicesURL     = "/devices"
	PointsURL      = "/points"
	PointsWriteURL = "/points/write/:uuid"

	JsonSchemaNetwork = "/schema/json/network"
	JsonSchemaDevice  = "/schema/json/device"
	JsonSchemaPoint   = "/schema/json/point"

	SchemaLegacyNetwork = "/schema/network"
	SchemaLegacyDevice  = "/schema/device"
	SchemaLegacyPoint   = "/schema/point"
)

func GetBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBodyPointWriter(ctx *gin.Context) (dto *model.PointWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func ResolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func PointWrite(pnt *model.Point) (out float64) {
	out = float.NonNil(pnt.WriteValue)
	return
}

func SetStatusCode(code, defaultCode int) int {
	if code == 0 {
		return defaultCode
	} else {
		return code
	}
}
