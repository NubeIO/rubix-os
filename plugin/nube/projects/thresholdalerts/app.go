package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.thresholdalertsDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body)
	if err != nil {
		inst.thresholdalertsErrorMsg("addNetwork(): failed to create network: ", body.Name)
		return nil, errors.New("failed to create network")
	}
	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.thresholdalertsDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.thresholdalertsDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.thresholdalertsDebugMsg("addDevice(): failed to create tmv device: ", body.Name)
		return nil, errors.New("failed to create tmv device")
	}
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.thresholdalertsDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.thresholdalertsDebugMsg("addPoint(): ", body.Name)

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.thresholdalertsDebugMsg("addPoint(): failed to create tmv point: ", body.Name)
		return nil, errors.New("failed to create tmv point")
	}
	inst.thresholdalertsDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	return point, nil

}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.thresholdalertsDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("updateNetwork():  nil network object")
		return
	}
	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		// DO POLLING DISABLE ACTIONS
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.thresholdalertsDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("updateDevice(): nil device object")
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
	if err != nil || device == nil {
		return nil, err
	}

	if boolean.IsFalse(device.Enable) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	} else {
		// TODO: Currently on every device update, all device points are removed, and re-added.
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Info
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(device.UUID, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.thresholdalertsDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("updatePoint(): nil point object")
		return
	}

	inst.thresholdalertsDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.thresholdalertsDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

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
		inst.thresholdalertsDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	point = nil
	inst.thresholdalertsDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("writePoint(): nil point object")
		return
	}

	inst.thresholdalertsDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.thresholdalertsDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	point, _, _, _, err = inst.db.PointWrite(pntUUID, body)
	if err != nil {
		inst.thresholdalertsDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.thresholdalertsDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("deleteNetwork(): nil network object")
		return
	}

	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.thresholdalertsDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("deleteDevice(): nil device object")
		return
	}

	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.thresholdalertsDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.thresholdalertsDebugMsg("deletePoint(): nil point object")
		return
	}

	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// THE FOLLOWING FUNCTIONS ARE CALLED FROM WITHIN THE PLUGIN

func (inst *Instance) ProcessThresholdAlerts() error {
	// Check that there is a High or Low Threshold configured
	if inst.config.Job.HighLimitThreshold == nil && inst.config.Job.LowLimitThreshold == nil {
		inst.thresholdalertsErrorMsg("ProcessThresholdAlerts()  Error: no high or low limit thresholds were set")
		return errors.New("no high or low limit thresholds were set")
	}

	if !inst.config.Job.HighLimitEnable && !inst.config.Job.LowLimitEnable {
		inst.thresholdalertsErrorMsg("ProcessThresholdAlerts()  Error: no high or low limit thresholds alerts enabled")
		return errors.New("no high or low limit thresholds alerts enabled")
	}

	// Get FF Token
	ffTokenResp, err := inst.GetFFToken("admin", "N00BWires")
	if ffTokenResp == nil || err != nil {
		inst.thresholdalertsErrorMsg("ProcessThresholdAlerts()  GetFFToken() err: ", err)
	}

	// Assemble Query Parameters (From Config File)
	queryParams, err := inst.ProcessQueryParams(inst.config.Job)
	inst.thresholdalertsDebugMsg("ProcessQueryParams() queryParams ", queryParams)

	ffHistoryArray, err := inst.GetFFHistories(*ffTokenResp, queryParams)
	if ffHistoryArray == nil || len(ffHistoryArray) <= 0 || err != nil {
		inst.thresholdalertsErrorMsg("ProcessThresholdAlerts() GetFFHistories(): ", err)
		return errors.New("ProcessThresholdAlerts() GetFFHistories(): error getting histories")
	}
	inst.thresholdalertsDebugMsg(fmt.Sprintf("ProcessQueryParams() ffHistoryArray: %+v", ffHistoryArray))
	highThresholdAlerts, lowThresholdAlerts, err := inst.ThresholdAnalysis(ffHistoryArray, inst.config.Job)
	if err != nil {
		inst.thresholdalertsErrorMsg("ProcessThresholdAlerts() ThresholdAnalysis(): ", err)
		return errors.New("ProcessThresholdAlerts() ThresholdAnalysis(): error converting histories to dataframe")
	}
	inst.thresholdalertsDebugMsg("ProcessQueryParams() highThresholdAlerts ", highThresholdAlerts)
	inst.thresholdalertsDebugMsg("ProcessQueryParams() lowThresholdAlerts ", lowThresholdAlerts)

	return err
}

func (inst *Instance) ProcessQueryParams(jobConfig Job) (string, error) {
	inst.thresholdalertsDebugMsg("ProcessQueryParams()")

	// Filter Params
	filterParams := "filter="
	filterParamsExist := false

	// Site Includes
	if jobConfig.SiteNamesInclude != nil && len(jobConfig.SiteNamesInclude) > 0 {
		filterParams += "("
		for ind, siteIncludeEntry := range jobConfig.SiteNamesInclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%7C%7C"
			}
			filterParams += fmt.Sprintf("site_name=%s", siteIncludeEntry)
		}
		filterParams += ")"
	}

	// Site Excludes
	if jobConfig.SiteNamesExclude != nil && len(jobConfig.SiteNamesExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, siteExcludeEntry := range jobConfig.SiteNamesExclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%26%26"
			}
			filterParams += fmt.Sprintf("site_name!=%s", siteExcludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Network Includes
	if jobConfig.RubixNetworkNamesInclude != nil && len(jobConfig.RubixNetworkNamesInclude) > 0 {
		for ind, netIncludeEntry := range jobConfig.RubixNetworkNamesInclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%7C%7C"
			}
			filterParams += fmt.Sprintf("rubix_network_name=%s", netIncludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Network Excludes
	if jobConfig.RubixNetworkNamesExclude != nil && len(jobConfig.RubixNetworkNamesExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, netExcludeEntry := range jobConfig.RubixNetworkNamesExclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%26%26"
			}
			filterParams += fmt.Sprintf("rubix_network_name!=%s", netExcludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Device Includes
	if jobConfig.RubixDeviceNamesInclude != nil && len(jobConfig.RubixDeviceNamesInclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, devIncludeEntry := range jobConfig.RubixDeviceNamesInclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%7C%7C"
			}
			filterParams += fmt.Sprintf("rubix_device_name=%s", devIncludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Device Excludes
	if jobConfig.RubixDeviceNamesExclude != nil && len(jobConfig.RubixDeviceNamesExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, devExcludeEntry := range jobConfig.RubixDeviceNamesExclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%26%26"
			}
			filterParams += fmt.Sprintf("rubix_device_name!=%s", devExcludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Point Includes
	if jobConfig.RubixPointNamesInclude != nil && len(jobConfig.RubixPointNamesInclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, pntIncludeEntry := range jobConfig.RubixPointNamesInclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%7C%7C"
			}
			filterParams += fmt.Sprintf("rubix_point_name=%s", pntIncludeEntry)
		}
		filterParams += ")"
	}

	// Rubix Point Excludes
	if jobConfig.RubixPointNamesExclude != nil && len(jobConfig.RubixPointNamesExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, pntExcludeEntry := range jobConfig.RubixPointNamesExclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%26%26"
			}
			filterParams += fmt.Sprintf("rubix_point_name!=%s", pntExcludeEntry)
		}
		filterParams += ")"
	}

	// Tags Includes
	if jobConfig.TagsInclude != nil && len(jobConfig.TagsInclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, tagIncludeEntry := range jobConfig.TagsInclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%7C%7C"
			}
			filterParams += fmt.Sprintf("tag=%s", tagIncludeEntry)
		}
		filterParams += ")"
	}

	// Tags Excludes
	if jobConfig.TagsExclude != nil && len(jobConfig.TagsExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		for ind, tagExcludeEntry := range jobConfig.TagsExclude {
			filterParamsExist = true
			if ind > 0 {
				filterParams += "%26%26"
			}
			filterParams += fmt.Sprintf("tag!=%s", tagExcludeEntry)
		}
		filterParams += ")"
	}

	// Meta Tags Includes
	if jobConfig.MetaTagsInclude != nil && len(jobConfig.MetaTagsInclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		first := true
		for metatagIncludeKey, metatagIncludeValue := range jobConfig.MetaTagsInclude {
			filterParamsExist = true
			if first {
				filterParams += "%7C%7C"
				first = false
			}
			filterParams += fmt.Sprintf("meta_tag_key=%s", metatagIncludeKey)
			filterParams += "%26%26"
			filterParams += fmt.Sprintf("meta_tag_value=%s", metatagIncludeValue)
		}
		filterParams += ")"
	}

	// Tags Excludes
	if jobConfig.MetaTagsExclude != nil && len(jobConfig.MetaTagsExclude) > 0 {
		if filterParamsExist {
			filterParams += "%26%26"
		}
		filterParams += "("
		first := true
		for metatagExcludeKey, metatagExcludeValue := range jobConfig.MetaTagsExclude {
			filterParamsExist = true
			if first {
				filterParams += "%26%26"
				first = false
			}
			filterParams += fmt.Sprintf("meta_tag_key=%s", metatagExcludeKey)
			filterParams += "%26%26"
			filterParams += fmt.Sprintf("meta_tag_value=%s", metatagExcludeValue)
		}
		filterParams += ")"
	}

	inst.thresholdalertsDebugMsg("ProcessQueryParams() filterParams: ", filterParams)

	if filterParamsExist {
		paramString := fmt.Sprintf("?%s", filterParams)
		// TODO: We may need to extend the history period if there are no values within the `alertDelay` period.  If we get no values in the history, then we should get the last history value.
		alertDelay, err := time.ParseDuration(fmt.Sprintf("%fm", float.NonNil(jobConfig.AlertDelayMins)))
		if err != nil {
			alertDelay = time.Minute * 60
		}
		// anHourAgo := time.Now().Add(time.Hour * -1)
		// anHourAgoString := anHourAgo.UTC().Format("2006-01-02%2015:04:05")
		// periodStartDateString := anHourAgo.UTC().Format("2006-01-02")
		// periodStartTimeString := anHourAgo.UTC().Format("15:04:05")
		periodStartTime := time.Now().Add(alertDelay * -1)
		periodStartDateString := periodStartTime.UTC().Format("2006-01-02")
		periodStartTimeString := periodStartTime.UTC().Format("15:04:05")
		inst.thresholdalertsDebugMsg("ProcessQueryParams() periodStartDateString: ", periodStartDateString)
		inst.thresholdalertsDebugMsg("ProcessQueryParams() periodStartTimeString: ", periodStartTimeString)
		paramString += "%26%26(timestamp%3E"
		paramString += periodStartDateString
		paramString += "%20"
		paramString += periodStartTimeString
		paramString += ")"
		return paramString, nil
	} else {
		return "", nil
	}
}
