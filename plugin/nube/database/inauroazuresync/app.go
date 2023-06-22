package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Instance) syncAzureSensorHistories() (bool, error) {
	log.Info("azure sensor history sync has been called...")

	// TODO: replace this history fetching logic to get histories on a per-gateway basis.
	// TODO: Last sync will need to be stored on a per-gateway basis. This will likely be done with the module JSON storage.
	// next get the histories that need to be sync'd to Azure
	lastSyncId, err := inst.db.GetHistoryPostgresLogLastSyncHistoryId() // fetches the ID of the last history that was sync'd
	if err != nil {
		inst.inauroazuresyncErrorMsg(err)
		return false, err
	}
	histories, err := inst.db.GetHistoriesForPostgresSync(lastSyncId) // fetches histories that have been added since the last sync
	if err != nil {
		inst.inauroazuresyncErrorMsg(err)
		return false, err
	}

	sampleHistory1 := model.History{
		HistoryID: 1,
		ID:        1,
		UUID:      "pnt_592b0dc2ba8d4434",
		Value:     111,
		Timestamp: time.Now().Add(-1 * time.Hour),
	}
	histories = append(histories, &sampleHistory1)

	sampleHistory2 := model.History{
		HistoryID: 2,
		ID:        2,
		UUID:      "pnt_76003bbae99846a3",
		Value:     222,
		Timestamp: time.Now().Add(-1 * time.Hour).Add(1 * time.Second),
	}
	histories = append(histories, &sampleHistory2)

	sampleHistory3 := model.History{
		HistoryID: 3,
		ID:        3,
		UUID:      "pnt_592b0dc2ba8d4434",
		Value:     333,
		Timestamp: time.Now().Add(-30 * time.Minute),
	}
	histories = append(histories, &sampleHistory3)

	inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() history count: %+v", len(histories)))

	if len(histories) > 0 {
		bulkInauroHistoryPayloadsByGateway, _ := inst.packageHistoriesToInauroPayloads(histories)
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() packageHistoriesToInauroPayloads() bulkInauroPayload: %+v", bulkInauroHistoryPayloadsByGateway))

		for gateway, gatewayHistories := range bulkInauroHistoryPayloadsByGateway {
			gateway = "host_xxxxxxxxxx" // TODO: DELETE ME
			byteData, err := json.Marshal(gatewayHistories)
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() json.Marshal(gatewayHistories) gateway: ", gateway, "  error:", err)
				continue
			}

			// azure open mqtt client connection and checks.
			// TODO: replace with: azureClient, err := inst.newAzureMQTTClientByHostUUID(gateway)
			azureClient, err := inst.newAzureMQTTClientByHostUUID(gateway)
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
			}
			inst.inauroazuresyncDebugMsg("syncAzureSensorHistories()  azureClient.Close() CLOSED")
		}

		// TODO: figure out how to implement this on a per gateway basis.  Need to re-sync old histories if they fail to send to a specific gateway.
		// TODO: probably use the Module JSON DB Storage to log the last sync'd history for each gateway
		lastHistory := histories[len(histories)-1]
		lastDeliveredHistoryLog := &model.HistoryPostgresLog{
			ID:        lastHistory.ID,
			UUID:      lastHistory.UUID,
			Value:     lastHistory.Value,
			Timestamp: lastHistory.Timestamp,
		}

		// TODO: Last sync will need to be stored on a per-gateway basis. This will likely be done with the module JSON storage.
		_, _ = inst.db.UpdateHistoryPostgresLog(lastDeliveredHistoryLog)
		if err != nil {
			inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() UpdateHistoryPostgresLog err:", err)
			return false, err
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("azure iot hub: Stored %v new sensor records", len(histories)))
	} else {
		inst.inauroazuresyncDebugMsg("azure iot hub: Nothing to store, no new sensor records")
	}
	return true, nil
}

func (inst *Instance) syncAzureGatewayPayloads() (bool, error) {
	log.Info("azure gateway payload sync has been called...")

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
