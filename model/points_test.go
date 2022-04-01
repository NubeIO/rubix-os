package model

import (
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"testing"
)

func TestPrintPointValues(*testing.T) {

	// Instance is plugin instance
	type Instance struct {
		db dbhandler.Handler
	}
	inst := &Instance{}

	//Create Point with priority array mode of PriorityArrayToWriteValue.
	//This type uses the priority array to get a write value, then polls the protocol (eg modbus) to update the point presentValue.
	var pnt Point
	pnt.Name = "polledPointTest"
	pnt.PointPriorityArrayMode = PriorityArrayToWriteValue

	createdPoint, _ := inst.db.DB.CreatePoint(&pnt, true)

	//Now update point write value.  Write value should be 10 @ priority 10.
	//At this point the presentValue should still be null as there has not been a poll/write operation done.
	value16 := 16.0
	value10 := 10.0
	var pri Priority
	pri.P16 = &value16
	pri.P10 = &value10
	createdPoint.Priority = &pri

	updatedPoint, _ := inst.db.DB.UpdatePoint(createdPoint.UUID, &pnt, false)

	//THIS SECTION IS IN PLACE OF MODBUS (or other protocol) PLUGIN WHICH DOES A WRITE AND THEN READ TO GET THE PRESENT VALUE.
	trueVar1 := true
	updatedPoint.ValueUpdatedFlag = &trueVar1
	floatVar1 := 10.0
	updatedPoint.PresentValue = &floatVar1
	trueVar2 := true
	updatedPoint.InSync = &trueVar2

	polledPoint, _ := inst.db.DB.UpdatePoint(updatedPoint.UUID, &pnt, true)

	polledPoint.PrintPointValues()

}
