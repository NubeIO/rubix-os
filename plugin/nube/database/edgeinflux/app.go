package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	log "github.com/sirupsen/logrus"
)

type InfluxDetail struct {
	InfluxSetting *InfluxSetting
	MaxId         int
	Records       int
	IsError       bool
}

func (inst *Instance) initializeInfluxSettings() []*InfluxSetting {
	var influxSettings []*InfluxSetting
	influxConnections := inst.config.Influx
	for _, influx := range influxConnections {
		influxSetting := new(InfluxSetting)
		schema := "http"
		if influx.Port == 443 {
			schema = "https"
		}
		influxSetting.ServerURL = fmt.Sprintf("%s://%s:%d", schema, influx.Host, influx.Port)
		if influx.Token == nil {
			log.Warn("influx: Token is null, please update it")
			continue
		}
		influxSetting.AuthToken = *influx.Token
		influxSetting.Org = influx.Org
		influxSetting.Bucket = influx.Bucket
		influxSetting.Measurement = influx.Measurement
		influxSettings = append(influxSettings, influxSetting)
	}
	return influxSettings
}

func (inst *Instance) setupInfluxInstances(influxSettings []*InfluxSetting) ([]InfluxDetail, error) {

	var influxDetails []InfluxDetail
	allError := true
	for _, influxSetting := range influxSettings {
		// lastSyncId, isError := influxSetting.GetLastSyncId()
		lastSyncId := 0
		isError := false
		influxDetail := InfluxDetail{
			InfluxSetting: influxSetting,
			MaxId:         lastSyncId,
			Records:       0,
			IsError:       isError,
		}
		if !isError {
			allError = false
		}
		influxDetails = append(influxDetails, influxDetail)
	}

	if allError {
		err := "influx: no connections are valid"
		log.Warn(err)
		return nil, errors.New(err)
	}
	return influxDetails, nil
}

func (inst *Instance) sendHistoriesToInflux(influxDetails []InfluxDetail, histories []*History) (bool, error) {

	for _, history := range histories {
		tags := history.Tags
		fields := fieldsHistory(history)
		for i, influxDetail := range influxDetails {
			influxDetail.InfluxSetting.WriteHistories(tags, fields, history.Timestamp)
			influxDetails[i].Records += 1 // directly updating to reflect value
		}
	}

	// forcing to push bulk writes
	for _, influxDetail := range influxDetails {
		influxDetail.InfluxSetting.getInfluxConnectionInstance().writeAPI.Flush()
		if influxDetail.Records > 0 {
			inst.edgeinfluxDebugMsg(fmt.Sprintf("Stored %v rows on %v", influxDetail.Records, path))
		} else {
			inst.edgeinfluxDebugMsg("influx: Nothing to store, no new records")
		}
	}
	return true, nil
}

func (inst *Instance) syncInflux(influxSettings []*InfluxSetting) (bool, error) {
	inst.edgeinfluxDebugMsg("InfluxDB sync has is been called...")
	if len(influxSettings) == 0 {
		err := "influx: InfluxDB sync failure: no any valid InfluxDB connection with not NULL token"
		log.Warn(err)
		return false, errors.New(err)
	}

	influxDetails, err := inst.setupInfluxInstances(influxSettings)
	if err != nil {
		log.Warn(err)
	}

	histories, err := inst.GetHistoryValues(inst.config.Job.Networks)
	if err != nil {
		log.Warn(err)
		return false, err
	}

	_, err = inst.sendHistoriesToInflux(influxDetails, histories)
	if err != nil {
		log.Warn(err)
		return false, err
	}
	return true, nil
}

func (inst *Instance) GetHistoryValues(requiredNetworksArray []string) ([]*History, error) {
	inst.edgeinfluxDebugMsg("GetHistoryValues()")
	if requiredNetworksArray == nil || len(requiredNetworksArray) == 0 {
		requiredNetworksArray = []string{"system"}
	}
	nowTimestamp := time.Now()
	var historyArray []*History
	for _, reqNet := range requiredNetworksArray {
		net, err := inst.db.GetNetworkByName(reqNet, api.Args{WithDevices: true, WithPoints: true})
		if err != nil || net == nil || net.Devices == nil {
			inst.edgeinfluxErrorMsg("GetHistoryValues() issue getting network: ", reqNet)
			continue
		}
		inst.edgeinfluxDebugMsg("GetHistoryValues() Net: ", net.Name)
		for _, dev := range net.Devices {
			if dev == nil || dev.Points == nil {
				inst.edgeinfluxErrorMsg("GetHistoryValues() issue getting device: ", reqNet)
				continue
			}
			for _, pnt := range dev.Points {
				point, err := inst.db.GetPoint(pnt.UUID, api.Args{WithTags: true})
				if point == nil || err != nil {
					inst.edgeinfluxErrorMsg("GetHistoryValues() Point is nil: ", pnt.Name, pnt.UUID)
					continue
				}
				if (inst.config.Job.RequireHistoryEnable && !boolean.NonNil(point.HistoryEnable)) || (point.HistoryType != model.HistoryTypeInterval && point.HistoryType != model.HistoryTypeCovAndInterval && point.HistoryType != "") {
					inst.edgeinfluxDebugMsg("GetHistoryValues() History save not required.")
					continue
				}
				// inst.edgeinfluxDebugMsg(fmt.Sprintf("GetHistoryValues() point: %+v", point))
				if point.PresentValue != nil {
					tagMap := make(map[string]string)
					tagMap["plugin_name"] = "lorawan"
					tagMap["rubix_network_name"] = net.Name
					tagMap["rubix_network_uuid"] = net.UUID
					tagMap["rubix_device_name"] = dev.Name
					tagMap["rubix_device_uuid"] = dev.UUID
					tagMap["rubix_point_name"] = point.Name
					tagMap["rubix_point_uuid"] = point.UUID

					pointHistory := History{
						UUID:      point.UUID,
						Value:     float.NonNil(point.PresentValue),
						Timestamp: nowTimestamp,
						Tags:      tagMap,
					}
					inst.edgeinfluxDebugMsg(fmt.Sprintf("GetHistoryValues() history: %+v", pointHistory))
					historyArray = append(historyArray, &pointHistory)
				}
			}
		}
	}
	return historyArray, nil
}

func (inst *Instance) SendPointWriteHistory(pntUUID string) error {
	inst.edgeinfluxDebugMsg("InfluxDB COV sync has is been called...")
	if pntUUID == "" {
		inst.edgeinfluxDebugMsg("Invalid Point UUID (empty) to SendPointWriteHistory()")
		return errors.New("invalid Point UUID (empty) to SendPointWriteHistory()")
	}
	if len(inst.influxDetails) == 0 {
		err := "influx: InfluxDB sync failure: no any valid InfluxDB connection with not NULL token"
		log.Warn(err)
		return errors.New(err)
	}

	influxDetails, err := inst.setupInfluxInstances(inst.influxDetails)
	if err != nil {
		log.Warn(err)
	}

	inst.edgeinfluxDebugMsg("SendPointWriteHistory()")

	point, err := inst.db.GetPoint(pntUUID, api.Args{WithTags: true})
	if err != nil || point == nil {
		inst.edgeinfluxErrorMsg("SendPointWriteHistory() GetPoint() err: ", err)
		return errors.New("SendPointWriteHistory() GetPoint() error")
	}

	/*(
	if (inst.config.Job.RequireHistoryEnable && !boolean.NonNil(point.HistoryEnable)) || (point.HistoryType != model.HistoryTypeCov && point.HistoryType != model.HistoryTypeCovAndInterval) {
		return nil
	}

	*/
	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.edgeinfluxErrorMsg("SendPointWriteHistory() GetDevice() err: ", err)
		return errors.New("SendPointWriteHistory() GetDevice() error")
	}
	net, _ := inst.db.GetNetwork(dev.NetworkUUID, api.Args{})
	if err != nil || dev == nil {
		inst.edgeinfluxErrorMsg("SendPointWriteHistory() GetNetwork() err: ", err)
		return errors.New("SendPointWriteHistory() GetNetwork() error")
	}

	// Check that the point is from a plugin network with history enabled (from config file)
	networkIsHistoryEnabled := false
	for _, reqNet := range inst.config.Job.Networks {
		if reqNet == net.Name {
			networkIsHistoryEnabled = true
			break
		}
	}
	if !networkIsHistoryEnabled {
		return nil
	}
	inst.edgeinfluxDebugMsg(fmt.Sprintf("GetHistoryValues() net: %+v", net))
	// inst.edgeinfluxDebugMsg(fmt.Sprintf("GetHistoryValues() point: %+v", point))
	if point.PresentValue != nil {
		tagMap := make(map[string]string)
		tagMap["plugin_name"] = "lorawan"
		tagMap["rubix_network_name"] = net.Name
		tagMap["rubix_network_uuid"] = net.UUID
		tagMap["rubix_device_name"] = dev.Name
		tagMap["rubix_device_uuid"] = dev.UUID
		tagMap["rubix_point_name"] = point.Name
		tagMap["rubix_point_uuid"] = point.UUID

		pointHistory := History{
			UUID:      point.UUID,
			Value:     float.NonNil(point.PresentValue),
			Timestamp: time.Now(),
			Tags:      tagMap,
		}
		inst.edgeinfluxDebugMsg(fmt.Sprintf("GetHistoryValues() history: %+v", pointHistory))

		var historyArray []*History
		historyArray = append(historyArray, &pointHistory)
		_, err = inst.sendHistoriesToInflux(influxDetails, historyArray)
		if err != nil {
			log.Warn(err)
			return err
		}
		return nil
	}
	return errors.New("no point present value found")
}
