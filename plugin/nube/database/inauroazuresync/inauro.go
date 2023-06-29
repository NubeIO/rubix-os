package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/float"
	"time"
)

type InauroGatewaySensorTrackingList []string // all sensor ids that have pushed since the last gateway payload

// TODO: add in network adapter properties that we can get from the edge/host via API call.

type InauroGatewayPayload struct {
	TimestampUTC string `json:"timestampUTC"`
	GatewayID    string `json:"gatewayID"`
	GatewayICCID string `json:"gatewayICCID"`
	GatewayIMEI  string `json:"gatewayIMEI"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Network      string `json:"network"`
	Ping         bool   `json:"ping"`
}

type MeasurementPayloadMap map[string]float64

type InauroSensorPayload struct { // stores sensor payloads with multiple measurements (grouped by individual sensor push)
	TimestampUTC string                `json:"timestampUTC"`
	GatewayID    string                `json:"gatewayID"`
	SensorMake   string                `json:"sensorMake"`
	SensorModel  string                `json:"sensorModel"`
	SensorID     string                `json:"sensorID"`
	Points       MeasurementPayloadMap `json:"points"`
}

type InauroHistoryArrayPayload []InauroSensorPayload // bulk histories in azure payload format

type InauroPackagedSensorHistoriesByTimestamp map[time.Time]InauroSensorPayload // histories for one sensor stored by timestamp

type InauroPackagedSensorHistoriesByDevice map[string]InauroPackagedSensorHistoriesByTimestamp // all sensor histories stored by device UUID

type InauroPackagedSensorHistoriesByGateway map[string]InauroPackagedSensorHistoriesByDevice // all sensor histories stored by host/gateway UUID

type InauroMultipleGatewayHistoryPayloads map[string]InauroHistoryArrayPayload // grouped histories by host/gateway

type InauroReqPayloadInfo struct { // stores the required information from device and point to be added to sensor history payloads
	PointName   string `json:"pointName"`   // point name that is sent with the `points` measurements
	SensorMake  string `json:"sensorMake"`  // this is the sensor Manufacturer (eg. Milesight, Nube, Elsys)
	SensorModel string `json:"sensorModel"` // this is the sensor Model (eg. TH301)
	SensorID    string `json:"sensorID"`    // this is the sensor ID that is sent to Azure (usually the last 8 digits of the EUI)
	DeviceID    string `json:"deviceID"`    // reference to get the `dataRate` value for each device(sensor)
}

// packageHistoriesToInauroPayloads this function takes histories and packages them by sensor and similar timestamps, then formats them for export to Azure as an array of azure histories.
func (inst *Instance) packageHistoriesToInauroPayloads(hostUUID string, histories []*model.History) (bulkHistoryPayloadsArray InauroHistoryArrayPayload, latestHistoryTime time.Time, err error) {
	if len(histories) <= 0 {
		return nil, time.Time{}, errors.New("histories are empty")
	}

	historiesByDevice := InauroPackagedSensorHistoriesByDevice{}
	pointPayloadDataMap := make(map[string]InauroReqPayloadInfo) // save sensorID (from device meta-tags) and point name (from point).  reduces DB calls if there are multiple histories for the same point.

	for _, history := range histories {
		if history.HostUUID != hostUUID {
			inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() history.HostUUID != hostUUID")
			continue
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("packageHistoriesToInauroPayloads() history: %+v", history))

		// save the latest history timestamp (per gateway) to be saved to plugin storage
		if history.Timestamp.After(latestHistoryTime) {
			latestHistoryTime = history.Timestamp
		}

		_, pntExists := pointPayloadDataMap[history.PointUUID]

		if !pntExists {
			// sensorIDMetaTag := sensorIDMetaTagKey
			var dev *model.Device
			dev, err = inst.db.GetOneDeviceByArgs(api.Args{PointSourceUUID: &history.PointUUID, HostUUID: &history.HostUUID, WithPoints: true, WithMetaTags: true})
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() GetOneDeviceByArgs() uuid: ", history.PointUUID, "  err: ", err)
				continue
			}
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("packageHistoriesToInauroPayloads() dev: %+v", dev))
			sensorID, pointName, sensorMake, sensorModel, err := inst.GetSensorPayloadInfoFromDevice(dev, history.PointUUID)
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() GetSensorPayloadInfoFromDevice() uuid: ", history.PointUUID, "  err: ", err)
				continue
			}
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("GetSensorPayloadInfoFromDevice() sensorID: %v, pointName: %v, sensorMake: %v, sensorModel: %v", sensorID, pointName, sensorMake, sensorModel))
			// store for later to save DB calls
			pointPayloadDataMap[history.PointUUID] = InauroReqPayloadInfo{
				PointName:   pointName,
				SensorID:    sensorID,
				SensorMake:  sensorMake,
				SensorModel: sensorModel,
				DeviceID:    dev.UUID,
			}
		}
		devUUID := pointPayloadDataMap[history.PointUUID].DeviceID
		pointName := pointPayloadDataMap[history.PointUUID].PointName

		_, sensorExists := historiesByDevice[devUUID]
		if !sensorExists {
			historiesByDevice[devUUID] = InauroPackagedSensorHistoriesByTimestamp{}
		}

		timestamp := history.Timestamp.Truncate(time.Second)
		timestampExists, mapTimestamp := SimilarTimestampExistsInSensorHistoryMap(timestamp, historiesByDevice[devUUID])
		if !timestampExists {
			_, _, azureDeviceID, _, err := inst.getAzureDetailsFromConfigByHostUUID(hostUUID)
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() getAzureDetailsFromConfigByHostUUID() hostUUID: ", hostUUID, "  err: ", err)
				continue
			}

			historiesByDevice[devUUID][mapTimestamp] = InauroSensorPayload{
				TimestampUTC: mapTimestamp.UTC().Format(time.RFC3339), // .Format(time.RFC3339)
				GatewayID:    azureDeviceID,
				SensorMake:   pointPayloadDataMap[history.PointUUID].SensorMake,
				SensorModel:  pointPayloadDataMap[history.PointUUID].SensorModel,
				SensorID:     pointPayloadDataMap[history.PointUUID].SensorID,
				Points:       make(MeasurementPayloadMap),
			}
		}

		// add this measurement to the sensor payload
		sensorPayload := historiesByDevice[devUUID][mapTimestamp]
		sensorPayload.Points[pointName] = float.NonNil(history.Value)
		historiesByDevice[devUUID][mapTimestamp] = sensorPayload
	}

	printHistoriesByDevice(historiesByDevice)

	// now reformat the history data to be an array of inaruo formatted histories
	bulkHistoryPayloadsArray = InauroHistoryArrayPayload{}
	for _, device := range historiesByDevice {
		for _, inauroHistory := range device {
			bulkHistoryPayloadsArray = append(bulkHistoryPayloadsArray, inauroHistory)
		}
	}

	return bulkHistoryPayloadsArray, latestHistoryTime, nil
}

func HostExistsOnHistoriesByGatewayMap(gatewayHostUUID string, historiesByGatewayHost InauroPackagedSensorHistoriesByGateway) bool {
	_, exists := historiesByGatewayHost[gatewayHostUUID]
	if !exists {
		return false
	}
	return true
}

func DeviceExistsOnHistoriesByGatewayMap(gatewayHostUUID string, historiesByGatewayHost InauroPackagedSensorHistoriesByGateway) bool {
	_, exists := historiesByGatewayHost[gatewayHostUUID]
	if !exists {
		return false
	}
	return true
}

func SensorExistsOnHistoriesByGatewayMap(gatewayHostUUID, sensorDeviceUUID string, historiesByGatewayHost InauroPackagedSensorHistoriesByGateway) bool {
	historiesBySensorDevice, exists := historiesByGatewayHost[gatewayHostUUID]
	if !exists {
		return false
	}
	_, exists = historiesBySensorDevice[gatewayHostUUID]
	if !exists {
		return false
	}
	return true
}

// SimilarTimestampExistsInSensorHistoryMap Checks to see whether the new timestamp is close enough to be considered the same sensor push.  range defined by timestampRangeToCombine constant.
func SimilarTimestampExistsInSensorHistoryMap(newTime time.Time, myMap InauroPackagedSensorHistoriesByTimestamp) (exists bool, existingKey time.Time) {
	for key, _ := range myMap {
		diff := key.Sub(newTime)
		absDiff := diff
		if diff < 0 {
			absDiff = -diff
		}
		if absDiff <= timestampRangeToCombine {
			return true, key
		}
	}
	return false, newTime
}

// GetSensorPayloadInfoFromDevice gets sensorID, pointName, and data-rate from device. Device must have Meta-Tags and Points included (use api.Args)
func (inst *Instance) GetSensorPayloadInfoFromDevice(dev *model.Device, pointUUID string) (sensorID, pointName, sensorMake, sensorModel string, err error) {
	// get sensorID from sensorRef meta-tag
	if dev.MetaTags == nil || len(dev.MetaTags) == 0 {
		err = errors.New(fmt.Sprintf("no meta tags on device, %v", dev.UUID))
		return
	}
	foundSensorID := false
	foundSensorManu := false
	foundSensorModel := false
	for _, tag := range dev.MetaTags {
		if foundSensorID && foundSensorManu && foundSensorModel {
			break
		}
		if !foundSensorID && tag.Key == sensorIDMetaTagKey {
			sensorID = tag.Value
			foundSensorID = true
			continue
		} else if !foundSensorManu && tag.Key == sensorMakeMetaTagKey {
			sensorMake = tag.Value
			foundSensorManu = true
			continue
		} else if !foundSensorModel && tag.Key == sensorModelMetaTagKey {
			sensorModel = tag.Value
			foundSensorModel = true
			continue
		}
	}
	if !foundSensorID {
		err = errors.New(fmt.Sprintf("'%s' meta tag not found on device, %v", sensorIDMetaTagKey, dev.UUID))
		return
	}
	if !foundSensorManu {
		inst.inauroazuresyncErrorMsg(fmt.Sprintf("GetSensorPayloadInfoFromDevice() error: '%s' meta tag not found on device, %v", sensorMakeMetaTagKey, dev.UUID))
	}
	if !foundSensorModel {
		inst.inauroazuresyncErrorMsg(fmt.Sprintf("GetSensorPayloadInfoFromDevice() error: '%s' meta tag not found on device, %v", sensorModelMetaTagKey, dev.UUID))
	}

	// get pointName
	if dev.Points == nil || len(dev.Points) == 0 {
		err = errors.New(fmt.Sprintf("no points found on device, %v", dev.UUID))
		return
	}
	pointNameFound := false
	for _, point := range dev.Points {
		// inst.inauroazuresyncDebugMsg(fmt.Sprintf("GetSensorPayloadInfoFromDevice() point: %+v", point))
		// inst.inauroazuresyncDebugMsg(fmt.Sprintf("GetSensorPayloadInfoFromDevice() point sourceUUID: %+v", *point.SourceUUID))
		if point.SourceUUID != nil && *point.SourceUUID == pointUUID {
			pointName = point.Name
			pointNameFound = true
			break
		}
	}
	if !pointNameFound {
		err = errors.New(fmt.Sprintf("point (%v) not found on device (%v)", pointUUID, dev.UUID))
		return
	}
	return
}
