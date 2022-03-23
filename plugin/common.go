package plugin

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
)

const (
	NetworksURL = "/networks"
	DevicesURL  = "/devices"
	PointsURL   = "/points"
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

func PointWrite(pnt *model.Point) (out float64) {
	out = utils.Float64IsNil(pnt.WriteValue)
	//log.Infof("modbus-write: pointWrite() ObjectType: %s  Addr: %d WriteValue: %v\n", pnt.ObjectType, utils.IntIsNil(pnt.AddressID), out)
	//if pnt.Priority != nil {
	//	if (*pnt.Priority).P16 != nil {
	//		out = *pnt.Priority.P16
	//		//log.Infof("modbus-write: pointWrite() ObjectType: %s  Addr: %d WriteValue: %v\n", pnt.ObjectType, utils.IntIsNil(pnt.AddressID), out)
	//	}
	//}
	return
}
