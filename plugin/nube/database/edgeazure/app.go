package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/amenzhinsky/iothub/iotdevice"
	"time"

	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
)

func (inst *Instance) syncAzure() error {
	inst.edgeazureDebugMsg("syncAzure()")

	azureDetails, err := inst.checkDeviceConnectionDetails()
	if err != nil {
		inst.edgeazureErrorMsg("makeDeviceConnectionDetails() err:", err)
		return err
	}
	inst.AzureDetails = azureDetails

	/*
		azureClient, err := inst.getAzureClient(inst.AzureDetails)
		if err != nil {
			inst.edgeazureErrorMsg("getAzureClient() err:", err)
			return err
		}

		inst.edgeazureDebugMsg("syncAzure() DeviceConnectionString:", inst.AzureDetails.DeviceConnectionString)

		// connect to the iothub
		err = azureClient.Connect(context.Background())
		if err != nil {
			inst.edgeazureErrorMsg("azureClient.Connect() err:", err)
			return err
		}
		// defer azureClient.Close()

	*/

	histories, err := inst.GetHistoryValues(inst.config.Job.Networks)
	if err != nil {
		inst.edgeazureErrorMsg("GetHistoryValues() err:", err)
		return err
	}

	err = inst.sendHistoriesToAzureWithHttp(azureDetails, histories)
	if err != nil {
		inst.edgeazureErrorMsg("sendHistoriesToAzureWithHttp() err:", err)
		return err
	}

	/*
		inst.edgeazureDebugMsg(fmt.Sprintf("syncAzure() azureClient: %+v", azureClient))
		err = inst.sendHistoriesToAzure(azureClient, histories)
		if err != nil {
			inst.edgeazureErrorMsg("sendHistoriesToAzure() err:", err)
			return err
		}
	*/

	return nil
}

func (inst *Instance) sendHistoriesToAzureWithHttp(azureDetails *AzureDeviceConnectionDetails, histories []*History) error {
	inst.edgeazureDebugMsg("sendHistoriesToAzureWithHttp()")
	var err error
	rest := inst.NewAzureRestClient(azureDetails)
	if rest == nil {
		return errors.New("failed to create rest client")
	}
	for _, history := range histories {
		resp, err := rest.sendAzureDeviceEventHttp(history)
		inst.edgeazureDebugMsg("sendHistoriesToAzureWithHttp() response Status: ", resp.Status())
		inst.edgeazureDebugMsg("sendHistoriesToAzureWithHttp() request URL: ", resp.Request.URL)
		inst.edgeazureDebugMsg("sendHistoriesToAzureWithHttp() request Token: ", resp.Request.Token)
		inst.edgeazureDebugMsg("sendHistoriesToAzureWithHttp() request Authorization: ", resp.Request.Header.Get("Authorization"))
		if err != nil {
			inst.edgeazureErrorMsg("sendAzureDeviceEventHttp() err:", err)
			continue
		}
		/*
			var historyAsByteArray []byte
			historyAsByteArray, err = inst.EncodeToBytes(history)
			if err != nil || len(historyAsByteArray) <= 0 {
				inst.edgeazureErrorMsg("sendHistoriesToAzureWithHttp() couldn't encode history to bytes. err:", err)
				continue
			}

			err = azureClient.SendEvent(context.Background(), historyAsByteArray)
			if err != nil {
				inst.edgeazureErrorMsg("sendHistoriesToAzureWithHttp() SendEvent() err:", err)
				continue
			}

		*/
	}
	return err
}

func (inst *Instance) sendHistoriesToAzure(azureClient *iotdevice.Client, histories []*History) error {
	inst.edgeazureDebugMsg("sendHistoriesToAzure()")
	var err error
	for _, history := range histories {
		inst.edgeazureDebugMsg(fmt.Sprintf("sendHistoriesToAzure() history: %+v", history))
		var historyAsByteArray []byte
		historyAsByteArray, err = inst.EncodeToBytes(history)
		if err != nil || len(historyAsByteArray) <= 0 {
			inst.edgeazureErrorMsg("sendHistoriesToAzure() couldn't encode history to bytes. err:", err)
			continue
		}

		err = azureClient.SendEvent(context.Background(), historyAsByteArray)
		if err != nil {
			inst.edgeazureErrorMsg("sendHistoriesToAzure() SendEvent() err:", err)
			continue
		}
	}
	return err
}

func (inst *Instance) GetHistoryValues(requiredNetworksArray []string) ([]*History, error) {
	inst.edgeazureDebugMsg("GetHistoryValues()")
	if requiredNetworksArray == nil || len(requiredNetworksArray) == 0 {
		requiredNetworksArray = []string{"system"}
	}
	nowTimestamp := time.Now()
	var historyArray []*History
	for _, reqNet := range requiredNetworksArray {
		net, err := inst.db.GetNetworkByName(reqNet, argspkg.Args{WithDevices: true, WithPoints: true})
		if err != nil || net == nil || net.Devices == nil {
			inst.edgeazureErrorMsg("GetHistoryValues() issue getting network: ", reqNet)
			continue
		}
		inst.edgeazureDebugMsg("GetHistoryValues() Net: ", net.Name)
		for _, dev := range net.Devices {
			if dev == nil || dev.Points == nil {
				inst.edgeazureErrorMsg("GetHistoryValues() issue getting device: ", reqNet)
				continue
			}
			for _, pnt := range dev.Points {
				point, err := inst.db.GetPoint(pnt.UUID, argspkg.Args{WithTags: true})
				if point == nil || err != nil {
					inst.edgeazureErrorMsg("GetHistoryValues() Point is nil: ", pnt.Name, pnt.UUID)
					continue
				}
				// if (inst.config.Job.RequireHistoryEnable && !boolean.NonNil(point.HistoryEnable)) || (point.HistoryType != model.HistoryTypeInterval && point.HistoryType != model.HistoryTypeCovAndInterval) {
				if inst.config.Job.RequireHistoryEnable && !boolean.NonNil(point.HistoryEnable) {
					continue
				}
				// inst.edgeazureDebugMsg(fmt.Sprintf("GetHistoryValues() point: %+v", point))
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
					inst.edgeazureDebugMsg(fmt.Sprintf("GetHistoryValues() history: %+v", pointHistory))
					historyArray = append(historyArray, &pointHistory)
				}
			}
		}
	}

	return historyArray, nil
}

func (inst *Instance) EncodeToBytes(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		inst.edgeazureErrorMsg("EncodeToBytes() err: ", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
func (inst *Instance) SendPointWriteHistory(pntUUID string) error {
	inst.edgeazureDebugMsg("SendPointWriteHistory()")
	if pntUUID == "" {
		return errors.New("invalid Point UUID (empty) to SendPointWriteHistory()")
	}

	azureDetails, err := inst.checkDeviceConnectionDetails()
	if err != nil {
		inst.edgeazureErrorMsg("makeDeviceConnectionDetails() err:", err)
		return err
	}
	inst.AzureDetails = azureDetails

	rest := inst.NewAzureRestClient(azureDetails)
	if rest == nil {
		return errors.New("failed to create rest client")
	}


	point, err := inst.db.GetPoint(pntUUID, api.Args{WithTags: true})
	if err != nil || point == nil {
		inst.edgeazureErrorMsg("SendPointWriteHistory() GetPoint() err: ", err)
		return errors.New("SendPointWriteHistory() GetPoint() error")
	}


	if (inst.config.Job.RequireHistoryEnable && !boolean.NonNil(point.HistoryEnable)) || (point.HistoryType != model.HistoryTypeCov && point.HistoryType != model.HistoryTypeCovAndInterval) {
		return nil
	}


	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.edgeazureErrorMsg("SendPointWriteHistory() GetDevice() err: ", err)
		return errors.New("SendPointWriteHistory() GetDevice() error")
	}
	net, _ := inst.db.GetNetwork(dev.NetworkUUID, api.Args{})
	if err != nil || dev == nil {
		inst.edgeazureErrorMsg("SendPointWriteHistory() GetNetwork() err: ", err)
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
	inst.edgeazureDebugMsg(fmt.Sprintf("SendPointWriteHistory() net: %+v", net))
	// inst.edgeinfluxDebugMsg(fmt.Sprintf("SendPointWriteHistory() point: %+v", point))
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
		inst.edgeazureDebugMsg(fmt.Sprintf("GetHistoryValues() history: %+v", pointHistory))

		var historyArray []*History
		historyArray = append(historyArray, &pointHistory)

		err = inst.sendHistoriesToAzure(rest, historyArray)
		if err != nil {
			inst.edgeazureErrorMsg("sendHistoriesToAzure() err:", err)
			return err
		}
		return nil
	}
	return errors.New("no point present value found")
}

*/
