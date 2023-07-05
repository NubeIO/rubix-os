package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/go-gota/gota/dataframe"
	"os"
	"strings"
	"time"
)

func (inst *Instance) CPSProcessing() {
	inst.cpsDebugMsg("CPSProcessing()")

	periodStart, _ := time.Parse(time.RFC3339, "2023-06-19T07:55:00Z")
	periodEnd, _ := time.Parse(time.RFC3339, "2023-06-19T10:00:00Z")

	// TODO: get site thresholds from thresholds table in pg database
	siteRef := "cps_b49e0c73919c47ef"
	var thresholds Threshold
	err := postgresSetting.postgresConnectionInstance.db.Last(&thresholds, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("CPSProcessing() db.Last(&thresholds, \"site_ref = ?\", siteRef) error: ", err)
	}
	// fmt.Println("thresholds")
	// fmt.Println(thresholds)
	var thresholdsSlice []Threshold
	thresholdsSlice = append(thresholdsSlice, thresholds)
	dfSiteThresholds := dataframe.LoadStructs(thresholdsSlice)
	// dfSiteThresholds := dataframe.ReadCSV(strings.NewReader(csvSiteThresholds))
	fmt.Println("dfSiteThresholds")
	fmt.Println(dfSiteThresholds)

	// TODO: pull data for each sensor for the given time range.  Do the below functions for each sensor

	// get raw sensor data for time range
	dfRawDoor := dataframe.ReadCSV(strings.NewReader(csvRawDoor))
	// fmt.Println("dfRawDoor")
	// fmt.Println(dfRawDoor)

	// get last stored processed data values
	dfLastProcessedDoor := dataframe.ReadCSV(strings.NewReader(csvLastProNODoor))
	fmt.Println("dfLastProcessedDoor")
	fmt.Println(dfLastProcessedDoor)

	// TODO: need to get last totalUses at the time when the last 15MinRollup value was stored (prior to the processing time range)
	// get last totalUses at the time when the last 15MinRollup value was stored (prior to the processing time range)
	dfLastTotalUsesAt15Min := dataframe.ReadCSV(strings.NewReader(csvLastTotalUsesAt15Min))
	fmt.Println("dfLastTotalUsesAt15Min")
	fmt.Println(dfLastTotalUsesAt15Min)

	dfResets := dataframe.ReadCSV(strings.NewReader(csvRawResets))
	// fmt.Println("dfResets")
	// fmt.Println(dfResets)

	dfDailyResets, err := inst.MakeDailyResetsDF(periodStart, periodEnd, dfSiteThresholds)
	if err != nil {
		inst.cpsErrorMsg("MakeDailyResetsDF() error: ", err)
		return
	}
	fmt.Println("dfDailyResets")
	fmt.Println(dfDailyResets)

	// join daily reset timestamps with the manual resets
	dfAllResets := dfResets.Concat(*dfDailyResets)
	dfAllResets = dfAllResets.Arrange(dataframe.Sort(string(timestampColName)))
	fmt.Println("dfAllResets")
	fmt.Println(dfAllResets)

	lastToPending := "2023-06-19T07:42:00Z"
	lastToClean := "2023-06-17T07:15:00Z"

	// Normally OPEN Door Usage Count and Occupancy DF
	DoorResultDF, err := inst.CalculateDoorUses(facilityToilet, normallyOpen, dfRawDoor, dfLastProcessedDoor, dfAllResets, dfSiteThresholds, lastToPending)
	if err != nil {
		inst.cpsErrorMsg("CalculateDoorUses() error: ", err)
		return
	}
	fmt.Println("DoorResultDF")
	fmt.Println(DoorResultDF)

	RollupResultDF, err := inst.Calculate15MinUsageRollup(periodStart, periodEnd, *DoorResultDF, dfLastTotalUsesAt15Min)
	if err != nil {
		inst.cpsErrorMsg("Calculate15MinUsageRollup() error: ", err)
		return
	}
	fmt.Println("RollupResultDF")
	fmt.Println(RollupResultDF)

	OverdueResultDF, err := inst.CalculateOverdueCubicles(facilityToilet, periodStart, periodEnd, *DoorResultDF, dfLastProcessedDoor, dfSiteThresholds, lastToPending, lastToClean)
	if err != nil {
		inst.cpsErrorMsg("CalculateOverdueCubicles() error: ", err)
		return
	}
	fmt.Println("OverdueResultDF")
	fmt.Println(OverdueResultDF)

	ResultFile, err := os.Create("/home/marc/Documents/Nube/CPS/Development/Data_Processing/1_Results.csv")
	if err != nil {
		inst.cpsErrorMsg(err)
	}
	defer ResultFile.Close()
	OverdueResultDF.WriteCSV(ResultFile)
}

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.cpsDebugMsg("addNetwork(): ", body.Name)

	body.HistoryEnable = boolean.NewTrue()

	network, err = inst.db.CreateNetwork(body)
	if network == nil || err != nil {
		inst.cpsErrorMsg("addNetwork(): failed to create cps network: ", body.Name)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("failed to create cps network")
	}
	return network, err
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.cpsDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.cpsDebugMsg("addDevice(): ", body.Name)

	if body.Name == "Cleaning Resets" {
		body.HistoryEnable = boolean.NewTrue()
	} else {
		body.HistoryEnable = boolean.NewFalse()
	}

	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.cpsDebugMsg("addDevice(): failed to create cps device: ", body.Name)
		return nil, errors.New("failed to create cps device")
	}

	if device.Name == "Cleaning Resets" || device.Name == "Availability" {
		return device, nil
	}

	// Automatically create door sensor processed data points
	createThesePoints := make([]model.Point, 0)
	occupancyPoint := model.Point{
		Name: string(cubicleOccupancyColName),
	}
	createThesePoints = append(createThesePoints, occupancyPoint)

	totalUsesPoint := model.Point{
		Name: string(totalUsesColName),
	}
	createThesePoints = append(createThesePoints, totalUsesPoint)

	currentUsesPoint := model.Point{
		Name: string(currentUsesColName),
	}
	createThesePoints = append(createThesePoints, currentUsesPoint)

	fifteenMinUsesPoint := model.Point{
		Name: string(fifteenMinRollupUsesColName),
	}
	createThesePoints = append(createThesePoints, fifteenMinUsesPoint)

	pendingStatusPoint := model.Point{
		Name: string(pendingStatusColName),
	}
	createThesePoints = append(createThesePoints, pendingStatusPoint)

	overdueStatusPoint := model.Point{
		Name: string(overdueStatusColName),
	}
	createThesePoints = append(createThesePoints, overdueStatusPoint)

	toPendingPoint := model.Point{
		Name: string(toPendingColName),
	}
	createThesePoints = append(createThesePoints, toPendingPoint)

	toCleanPoint := model.Point{
		Name: string(toCleanColName),
	}
	createThesePoints = append(createThesePoints, toCleanPoint)

	toOverduePoint := model.Point{
		Name: string(toOverdueColName),
	}
	createThesePoints = append(createThesePoints, toOverduePoint)

	cleaningTimePoint := model.Point{
		Name: string(cleaningTimeColName),
	}
	createThesePoints = append(createThesePoints, cleaningTimePoint)

	devNameSplit := strings.Split(device.Name, "-")
	if len(devNameSplit) < 3 {
		inst.cpsErrorMsg("addDevice(): device name should be of the form Level-Gender-Location")
		return device, nil
	}
	level := devNameSplit[0]
	gender := strings.ToLower(devNameSplit[1])
	location := devNameSplit[2]

	for _, point := range createThesePoints {
		point.DeviceUUID = device.UUID
		point.MetaTags = make([]*model.PointMetaTag, 0)

		metaTag1 := model.PointMetaTag{Key: "assetFunc", Value: "managedCubicle"}
		point.MetaTags = append(point.MetaTags, &metaTag1)
		metaTag2 := model.PointMetaTag{Key: "doorType", Value: "toilet"}
		point.MetaTags = append(point.MetaTags, &metaTag2)
		metaTag3 := model.PointMetaTag{Key: "enableCleaningTracking", Value: "true"}
		point.MetaTags = append(point.MetaTags, &metaTag3)
		metaTag4 := model.PointMetaTag{Key: "enableUseCounting", Value: "true"}
		point.MetaTags = append(point.MetaTags, &metaTag4)
		metaTag5 := model.PointMetaTag{Key: "floorRef", Value: level}
		point.MetaTags = append(point.MetaTags, &metaTag5)
		metaTag6 := model.PointMetaTag{Key: "genderRef", Value: gender}
		point.MetaTags = append(point.MetaTags, &metaTag6)
		metaTag7 := model.PointMetaTag{Key: "isEOT", Value: "false"}
		point.MetaTags = append(point.MetaTags, &metaTag7)
		metaTag8 := model.PointMetaTag{Key: "locationRef", Value: location}
		point.MetaTags = append(point.MetaTags, &metaTag8)
		metaTag9 := model.PointMetaTag{Key: "measurementRefs", Value: "door_position"}
		point.MetaTags = append(point.MetaTags, &metaTag9)
		metaTag10 := model.PointMetaTag{Key: "normalPosition", Value: "NO"}
		point.MetaTags = append(point.MetaTags, &metaTag10)
		metaTag11 := model.PointMetaTag{Key: "siteRef", Value: "cps_b49e0c73919c47ef"}
		point.MetaTags = append(point.MetaTags, &metaTag11)
		metaTag12 := model.PointMetaTag{Key: "assetRef", Value: device.Description}
		point.MetaTags = append(point.MetaTags, &metaTag12)
		metaTag13 := model.PointMetaTag{Key: "resetID", Value: "rst_1"}
		point.MetaTags = append(point.MetaTags, &metaTag13)
		metaTag14 := model.PointMetaTag{Key: "availabilityID", Value: "avl_1"}
		point.MetaTags = append(point.MetaTags, &metaTag14)

		inst.addPoint(&point)
	}

	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil || body.DeviceUUID == "" {
		inst.cpsDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.cpsDebugMsg("addPoint(): ", body.Name)

	device, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if device == nil || err != nil {
		inst.cpsDebugMsg("addPoint(): failed to find device", body.DeviceUUID)
		return nil, err
	}

	if device.Name == "Cleaning Resets" {
		body.HistoryEnable = boolean.NewTrue()
		body.HistoryType = model.HistoryTypeCov
		body.HistoryCOVThreshold = float.New(0.1)
	} else {
		body.HistoryEnable = boolean.NewFalse()
	}

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.cpsDebugMsg("addPoint(): failed to create cps point: ", body.Name)
		return nil, err
	}
	// inst.cpsDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))
	return point, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.cpsDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updateNetwork():  nil network object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "network disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}

	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.cpsDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
		body.CommonFault.Message = "device disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(body.UUID, body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.cpsDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updatePoint(): nil point object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}
	body.CommonFault.InFault = false
	body.CommonFault.MessageLevel = model.MessageLevel.Info
	body.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	body.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	body.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.cpsErrorMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.cpsDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deleteNetwork(): nil network object")
		return
	}

	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.cpsDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deleteDevice(): nil device object")
		return
	}
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.cpsDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deletePoint(): nil point object")
		return
	}
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}
