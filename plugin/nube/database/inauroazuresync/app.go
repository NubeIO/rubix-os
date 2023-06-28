package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/float"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Instance) syncAzureSensorHistories() (bool, error) {
	log.Info("azure sensor history sync has been called...")

	hosts, err := inst.db.GetHosts()
	if err != nil {
		return false, err
	}
	inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() hosts: %v", len(hosts)))
	for i, host := range hosts {
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() i: %v, host: %+v", i, host))
	}

	// Get the plugin storage that holds the last sync times for each host/gateway
	pluginStorage, err := inst.getPluginConfStorage()
	if err != nil {
		inst.inauroazuresyncErrorMsg("Enable() getPluginConfStorage() err:", err)
	}
	inst.inauroazuresyncDebugMsg(fmt.Sprintf("Enable() pluginStorage: %+v", pluginStorage))
	if pluginStorage == nil {
		newPluginStorage := PluginConfStorage{}
		newPluginStorage.LastSyncByGateway = make(map[string]time.Time)
		/*
			for _, host := range hosts {
				newPluginStorage.LastSyncByGateway[host.UUID] = time.Now().Add(-time.Hour)
			}
		*/
		pluginStorage = &newPluginStorage
	}

	now := time.Now()
	var histories []*model.History
	for _, host := range hosts {
		// TODO: add startTime from module storage
		// lastSyncTime, _ := time.Parse(time.RFC3339, "2023-06-25T00:00:00Z")
		lastSyncTime, ok := pluginStorage.LastSyncByGateway[host.UUID]
		if !ok {
			inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() host last sync time not found.  Host: ", host.UUID)
			lastSyncTime = time.Now().Add(-12 * time.Hour)
			pluginStorage.LastSyncByGateway[host.UUID] = lastSyncTime
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() lastSyncTime: %v", lastSyncTime))

		// TODO: fetch each gateway histories since the last syncTime
		// next get the gateway histories that still need to be sync'd to Azure
		histories, err = inst.db.GetHistoriesByHostUUID(host.UUID, lastSyncTime, now) // fetches histories that have been added since the last sync
		if err != nil {
			inst.inauroazuresyncErrorMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID() error: %v", err))
			inst.inauroazuresyncErrorMsg(err)
			return false, err
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID(): %v", len(histories)))
		for _, history := range histories {
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID() history: %+v", history))
		}

		// TODO: Delete these sample test histories
		sampleHistory1 := model.History{
			HistoryID: 1,
			ID:        1,
			PointUUID: "pnt_592b0dc2ba8d4434",
			Value:     float.New(111),
			Timestamp: time.Now().Add(-1 * time.Hour),
		}
		histories = append(histories, &sampleHistory1)

		sampleHistory2 := model.History{
			HistoryID: 2,
			ID:        2,
			PointUUID: "pnt_76003bbae99846a3",
			Value:     float.New(222),
			Timestamp: time.Now().Add(-1 * time.Hour).Add(1 * time.Second),
		}
		histories = append(histories, &sampleHistory2)

		sampleHistory3 := model.History{
			HistoryID: 3,
			ID:        3,
			PointUUID: "pnt_592b0dc2ba8d4434",
			Value:     float.New(333),
			Timestamp: time.Now().Add(-30 * time.Minute),
		}
		histories = append(histories, &sampleHistory3)

		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() history count: %+v", len(histories)))

		if len(histories) > 0 {
			bulkInauroHistoryPayloadsArray, latestHistoryTime, _ := inst.packageHistoriesToInauroPayloads(host.UUID, histories)
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() packageHistoriesToInauroPayloads() bulkInauroHistoryPayloadsArray: %+v", bulkInauroHistoryPayloadsArray))
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() packageHistoriesToInauroPayloads() latestHistoryTime: %+v", latestHistoryTime))

			byteData, err := json.Marshal(bulkInauroHistoryPayloadsArray)
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() json.Marshal(bulkInauroHistoryPayloadsArray) gateway: ", host.UUID, "  error:", err)
				continue
			}

			// azure open mqtt client connection and checks.
			azureClient, err := inst.newAzureMQTTClientByHostUUID(host.UUID)
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() inst.newAzureMQTTClientByHostUUID() error:", err)
				continue
			}
			// at this point we have a connected Azure MQTT Client
			// now we push the histories to Azure
			// send a device-to-cloud message
			if err = azureClient.SendEvent(context.Background(), byteData); err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() SendEvent() error:", err)
			}
			azureClient.Close()
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() azureClient.Close() error:", err)
			} else {
				inst.inauroazuresyncDebugMsg("syncAzureSensorHistories()  azureClient.Close() CLOSED")
				pluginStorage.LastSyncByGateway[host.UUID] = latestHistoryTime
			}

			inst.inauroazuresyncDebugMsg(fmt.Sprintf("azure iot hub: Stored %v new sensor records", len(histories)))

			// TODO: figure out how to implement this on a per gateway basis.  Need to re-sync old histories if they fail to send to a specific gateway.

			// If the Azure push was successful save the latest history time to the host in JSON storage.
			inst.setPluginConfStorage(pluginStorage)

		} else {
			inst.inauroazuresyncDebugMsg("azure iot hub: Nothing to store, no new sensor records")
		}
	}
	return true, nil
}

func (inst *Instance) syncAzureGatewayPayloads() (bool, error) {
	log.Info("azure gateway payload sync has been called...")

	// TODO: I don't think this is neccesary.  Should only be using the local db, not the timescale pg db
	/*
		_, err := inst.initializePostgresDBConnection()
		if err != nil {
			inst.inauroazuresyncErrorMsg(err)
			return false, err
		}
	*/

	hosts, err := inst.db.GetHosts()
	if err != nil {
		return false, err
	}
	inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() hosts: %v", len(hosts)))
	for i, host := range hosts {
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() i: %v, host: %+v", i, host))
	}

	for _, host := range hosts {
		fmt.Println(host)

		// TODO: replace this history fetching logic to get the last history for each sensor during the last 15 minutes, grouped by gateway(host).
		// next get the sensor histories for the current gateway payload period (last 15 mins, or as set by config)
		lastSyncId, err := inst.db.GetHistoryPostgresLogLastSyncHistoryId() // fetches the ID of the last history that was sync'd
		if err != nil {
			inst.inauroazuresyncErrorMsg(err)
			return false, err
		}
		// TODO: only get the last history for each sensor within the last gateway payload period (15mins)
		histories, err := inst.db.GetHistoriesForPostgresSync(lastSyncId) // fetches histories that have been added since the last sync
		if err != nil {
			inst.inauroazuresyncErrorMsg(err)
			return false, err
		}

		gatewayDetailsMap, err := inst.getGatewayDetailsFromConfig()
		if err != nil {
			inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() getGatewayDetailsFromConfig() error:", err)
			return false, err
		}

		timestamp := time.Now().Truncate(time.Second)

		for _, gatewayDetails := range gatewayDetailsMap {

			gatewayPayload := InauroGatewayPayload{
				TimestampUTC: timestamp.UTC().Format(time.RFC3339),
				GatewayID:    gatewayDetails.AzureDeviceId,
				GatewayICCID: gatewayDetails.SIMICCID,
				Latitude:     gatewayDetails.Latitude,
				Longitude:    gatewayDetails.Longitude,
				Network:      gatewayDetails.NetworkType,
			}

			// TODO: add ping check.
			// cli := cligetter.GetEdgeClientFastTimeout(host)
			// _, pingable, _ := cli.Ping()

			// TODO: add connected sensors
			gatewayPayload.ConnectedSensors = len(histories)

			byteData, err := json.Marshal(gatewayPayload)
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() json.Marshal(gatewayPayload) gateway: ", gatewayDetails.AzureDeviceId, "  error:", err)
				continue
			}

			// azure open mqtt client connection and checks.
			azureClient, err := inst.newAzureMQTTClientByGatewayDetails(gatewayDetails)
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() inst.newAzureMQTTClientByGatewayDetails() error:", err)
				continue
			}
			// at this point we have a connected Azure MQTT Client

			// now we push the gateway payload to Azure
			// send a device-to-cloud message
			if err = azureClient.SendEvent(context.Background(), byteData); err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() SendEvent() error:", err)
			}
			azureClient.Close()
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() azureClient.Close() error:", err)
			}
			inst.inauroazuresyncDebugMsg("syncAzureGatewayPayloads()  azureClient.Close() gateway: ", gatewayDetails.AzureDeviceId, "  CLOSED")

			inst.inauroazuresyncDebugMsg(fmt.Sprintf("azure iot hub: Stored %v new gateway records", len(histories)))
		}
	}
	return true, nil
}

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.inauroazuresyncDebugMsg("addNetwork(): ", body.Name)

	network, err = inst.db.CreateNetwork(body)
	if err != nil {
		inst.inauroazuresyncErrorMsg("addNetwork(): failed to create inaruoazuresync network: ", body.Name)
		return nil, errors.New("failed to create inaruoazuresync network")
	}
	return network, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.inauroazuresyncDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.inauroazuresyncErrorMsg("updateNetwork():  nil network object")
		return
	}
	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.inauroazuresyncDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.inauroazuresyncErrorMsg("deleteNetwork(): nil network object")
		return
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}
