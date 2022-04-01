package main

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) PointWriteModeTest() {

	net := model.Network{}
	net.Name = "TEST_NET_1"
	net.TransportType = "ip"
	net.PluginPath = "modbus"
	createdNetwork, err := inst.db.CreateNetwork(&net, true)
	log.Infof("network err: %+v\n", err)
	log.Infof("network: %+v\n", createdNetwork)

	dev := model.Device{}
	dev.NetworkUUID = createdNetwork.UUID
	dev.Name = "TEST_DEV_1"

	createdDevice, err := inst.db.CreateDevice(&dev)

	log.Infof("device err: %+v\n", err)
	log.Infof("device: %+v\n", createdDevice)

	//Create Point with priority array mode of PriorityArrayToWriteValue.
	//This type uses the priority array to get a write value, then polls the protocol (eg modbus) to update the point presentValue.
	var pnt model.Point
	pnt.Name = "polledPointTest"
	pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
	pnt.DeviceUUID = createdDevice.UUID

	createdPoint, _ := inst.db.CreatePoint(&pnt, true, false)

	log.Infof("createdPoint err: %+v\n", err)
	log.Infof("createdPoint: %+v\n", createdPoint)

	//Now update point write value.  Write value should be 10 @ priority 10.
	//At this point the presentValue should still be null as there has not been a poll/write operation done.

	var pri model.Priority
	pri.P16 = utils.NewFloat64(16)
	pri.P10 = utils.NewFloat64(10)
	createdPoint.Priority = &pri

	updatedPoint, _ := inst.db.UpdatePoint(createdPoint.UUID, &pnt, false)

	log.Infof("updatedPoint err: %+v\n", err)
	log.Infof("updatedPoint: %+v\n", updatedPoint)

	updatedPoint.PrintPointValues()

	//THIS SECTION IS IN PLACE OF MODBUS (or other protocol) PLUGIN WHICH DOES A WRITE AND THEN READ TO GET THE PRESENT VALUE.
	updatedPoint.ValueUpdatedFlag = utils.NewTrue()
	updatedPoint.PresentValue = utils.NewFloat64(10)
	//updatedPoint.InSync = utils.NewTrue()

	polledPoint, _ := inst.db.UpdatePoint(updatedPoint.UUID, updatedPoint, true)
	//polledPoint, _ := inst.db.UpdatePointPresentValue(updatedPoint, true)

	log.Infof("polledPoint err: %+v\n", err)
	log.Infof("polledPoint: %+v\n", polledPoint)

	polledPoint.PrintPointValues()

}
