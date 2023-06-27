package main

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/float"
	"strconv"
	"time"
)

type InauroGatewaySensorTrackingList []string // all sensor ids that have pushed since the last gateway payload

type InauroGatewayPayload struct {
	TimestampUTC     string `json:"timestampUTC"`
	GatewayID        string `json:"gatewayID"`
	GatewayICCID     string `json:"gatewayICCID"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	Network          string `json:"network"`
	ConnectedSensors int    `json:"connectedSensors"`
	Ping             bool   `json:"ping"`
}

type MeasurementPayloadMap map[string]float64

type InauroSensorPayload struct { // stores sensor payloads with multiple measurements (grouped by individual sensor push)
	TimestampUTC   string                `json:"timestampUTC"`
	GatewayID      string                `json:"gatewayID"`
	DataRate       string                `json:"dataRate"`
	NubeSensorUUID string                `json:"nubeSensorUUID"`
	Points         MeasurementPayloadMap `json:"points"`
}

type InauroHistoryArrayPayload []InauroSensorPayload // bulk histories in azure payload format

type InauroPackagedSensorHistoriesByTimestamp map[time.Time]InauroSensorPayload // histories for one sensor stored by timestamp

type InauroPackagedSensorHistoriesByDevice map[string]InauroPackagedSensorHistoriesByTimestamp // all sensor histories stored by device UUID

type InauroPackagedSensorHistoriesByGateway map[string]InauroPackagedSensorHistoriesByDevice // all sensor histories stored by host/gateway UUID

type InauroMultipleGatewayHistoryPayloads map[string]InauroHistoryArrayPayload // grouped histories by host/gateway

// packageHistoriesToInauroPayloadsByGateway this function takes histories, packages them by similar timestamps and splits them into grouped inauro history payload arrays (mapped by gateway/host).
func (inst *Instance) packageHistoriesToInauroPayloadsByGateway(histories []*model.History) (bulkHistoryPayloadsByGateway InauroMultipleGatewayHistoryPayloads, err error) {
	if len(histories) <= 0 {
		return nil, errors.New("histories are empty")
	}

	historiesByGateway := InauroPackagedSensorHistoriesByGateway{} // TODO: histories need to be grouped by host so that they can be sent out together
	deviceData := make(map[string]*model.Device)
	pointData := make(map[string]*model.Point)

	for _, history := range histories {
		// TODO: point shouldn't be required, only device (for sensorID)
		pnt, pntExists := pointData[history.PointUUID]
		if !pntExists {
			pnt, err = inst.db.GetPoint(history.PointUUID, api.Args{}) // needed to get device UUID
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloadsByGateway() GetPoint() uuid: ", history.PointUUID, "    err: ", err)
				continue
			}
			pointData[pnt.UUID] = pnt // store for later to save DB calls
		}

		// TODO: it would be much better to get the device by point UUID + host UUID (request added function/api)
		// TODO: HostUUID will be available on the new history model.
		dev, devExists := deviceData[pnt.DeviceUUID]
		if !devExists {
			dev, err = inst.db.GetDevice(pnt.DeviceUUID, api.Args{WithPoints: true}) // needed for azure values stored on Device Description
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloadsByGateway() GetDevice() uuid: ", pnt.DeviceUUID, "    err: ", err)
				continue
			}
			deviceData[dev.UUID] = dev // store for later to save DB calls
		}

		// TODO: HostUUID will be available on the new history model.
		hostUUID := history.HostUUID
		_, hostExists := historiesByGateway[hostUUID]
		if !hostExists {
			historiesByGateway[hostUUID] = InauroPackagedSensorHistoriesByDevice{} // store for later to save DB calls
		}

		timestamp := history.Timestamp.Truncate(time.Second)

		_, sensorExists := historiesByGateway[hostUUID][pnt.DeviceUUID]
		if !sensorExists {
			historiesByGateway[hostUUID][pnt.DeviceUUID] = InauroPackagedSensorHistoriesByTimestamp{}
		}

		timestampExists, mapTimestamp := SimilarTimestampExistsInSensorHistoryMap(timestamp, historiesByGateway[hostUUID][pnt.DeviceUUID])
		if !timestampExists {
			sensorID, err := inst.GetSensorIDFromDeviceDescription(dev) // TODO: update to get SensorID from meta tags
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloadsByGateway() GetSensorIDFromDeviceDescription() uuid: ", dev.UUID, "    err: ", err)
				continue
			}
			dataRate, err := inst.GetDataRateFromDevice(dev)
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloadsByGateway() GetDataRateFromDevice() err: ", err)
			}

			_, _, azureDeviceID, _, err := inst.getAzureDetailsFromConfigByHostUUID(hostUUID)
			historiesByGateway[hostUUID][pnt.DeviceUUID][mapTimestamp] = InauroSensorPayload{
				TimestampUTC:   mapTimestamp.UTC().Format(time.RFC3339), // .Format(time.RFC3339)
				GatewayID:      azureDeviceID,
				DataRate:       strconv.Itoa(int(dataRate)),
				NubeSensorUUID: sensorID,
				Points:         make(MeasurementPayloadMap),
			}
		}

		// add this measurement to the sensor payload
		sensorPayload := historiesByGateway[hostUUID][pnt.DeviceUUID][mapTimestamp]
		sensorPayload.Points[pnt.Name] = float.NonNil(history.Value)
		historiesByGateway[hostUUID][pnt.DeviceUUID][mapTimestamp] = sensorPayload
	}

	// now reformat the history data to be an array of inaruo formatted histories
	bulkInauroHistoryPayloadByGateway := InauroMultipleGatewayHistoryPayloads{}
	for gateway, _ := range historiesByGateway {
		_, exists := bulkInauroHistoryPayloadByGateway[gateway]
		if !exists {
			bulkInauroHistoryPayloadByGateway[gateway] = InauroHistoryArrayPayload{}
		}
		for _, device := range historiesByGateway[gateway] {
			for _, history := range device {
				bulkInauroHistoryPayloadByGateway[gateway] = append(bulkInauroHistoryPayloadByGateway[gateway], history)
			}
		}
	}

	return bulkInauroHistoryPayloadByGateway, nil
}

// packageHistoriesToInauroPayloads this function takes histories and packages them by sensor and similar timestamps, then formats them for export to Azure as an array of azure histories.
func (inst *Instance) packageHistoriesToInauroPayloads(hostUUID string, histories []*model.History) (bulkHistoryPayloadsArray InauroHistoryArrayPayload, err error) {
	if len(histories) <= 0 {
		return nil, errors.New("histories are empty")
	}

	historiesByDevice := InauroPackagedSensorHistoriesByDevice{}
	pointData := make(map[string]*model.Point)

	for _, history := range histories {
		if history.HostUUID != hostUUID {
			inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() history.HostUUID != hostUUID")
			continue
		}

		pnt, pntExists := pointData[history.PointUUID]
		if !pntExists {
			pnt, err = inst.db.GetPoint(history.PointUUID, api.Args{}) // TODO: Replace with get point by pointUUID and hostUUID (request added function/api)
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() GetPoint() uuid: ", history.PointUUID, "    err: ", err)
				continue
			}
			pointData[pnt.UUID] = pnt // store for later to save DB calls
		}

		// TODO: get the device by point UUID + host UUID (request added function/api)
		dev := GetDeviceByHostUUIDPointUUID(hostUUID, history.PointUUID)

		timestamp := history.Timestamp.Truncate(time.Second)

		_, sensorExists := historiesByDevice[hostUUID][dev.UUID]
		if !sensorExists {
			historiesByDevice[dev.UUID] = InauroPackagedSensorHistoriesByTimestamp{}
		}

		timestampExists, mapTimestamp := SimilarTimestampExistsInSensorHistoryMap(timestamp, historiesByDevice[dev.UUID])
		if !timestampExists {
			sensorID, err := inst.GetSensorIDFromDeviceDescription(dev) // TODO: update to get SensorID from meta tags
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() GetSensorIDFromDeviceDescription() uuid: ", dev.UUID, "    err: ", err)
				continue
			}
			dataRate, err := inst.GetDataRateFromDevice(dev)
			if err != nil {
				inst.inauroazuresyncErrorMsg("packageHistoriesToInauroPayloads() GetDataRateFromDevice() err: ", err)
			}

			_, _, azureDeviceID, _, err := inst.getAzureDetailsFromConfigByHostUUID(hostUUID)
			historiesByDevice[dev.UUID][mapTimestamp] = InauroSensorPayload{
				TimestampUTC:   mapTimestamp.UTC().Format(time.RFC3339), // .Format(time.RFC3339)
				GatewayID:      azureDeviceID,
				DataRate:       strconv.Itoa(int(dataRate)),
				NubeSensorUUID: sensorID,
				Points:         make(MeasurementPayloadMap),
			}
		}

		// add this measurement to the sensor payload
		sensorPayload := historiesByDevice[dev.UUID][mapTimestamp]
		sensorPayload.Points[pnt.Name] = float.NonNil(history.Value)
		historiesByDevice[pnt.DeviceUUID][mapTimestamp] = sensorPayload
	}

	// now reformat the history data to be an array of inaruo formatted histories
	bulkHistoryPayloadsArray = InauroHistoryArrayPayload{}
	for _, device := range historiesByDevice {
		for _, inauroHistory := range device {
			bulkHistoryPayloadsArray = append(bulkHistoryPayloadsArray, inauroHistory)
		}
	}

	return bulkHistoryPayloadsArray, nil
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
