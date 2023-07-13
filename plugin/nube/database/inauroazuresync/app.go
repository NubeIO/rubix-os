package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"time"
)

func (inst *Instance) syncAzureSensorHistories() (bool, error) {
	inst.inauroazuresyncDebugMsg("azure sensor history sync has been called...")

	hosts, err := inst.db.GetHosts(argspkg.Args{})
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
		inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() getPluginConfStorage() err:", err)
	}
	inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() pluginStorage: %+v", pluginStorage))
	if pluginStorage == nil {
		newPluginStorage := PluginConfStorage{}
		newPluginStorage.LastSyncByGateway = make(map[string]time.Time)
		pluginStorage = &newPluginStorage
	}

	now := time.Now()
	var histories []*model.History
	for _, host := range hosts {
		// get the lastSyncTime from module storage
		// lastSyncTime, _ := time.Parse(time.RFC3339, "2023-06-25T00:00:00Z")
		lastSyncTime, ok := pluginStorage.LastSyncByGateway[host.UUID]
		if !ok {
			inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() host last sync time not found.  Host: ", host.UUID)
			lastSyncTime = time.Now().Add(-1 * emptyModuleStorageResyncPeriod)
			pluginStorage.LastSyncByGateway[host.UUID] = lastSyncTime
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() lastSyncTime: %v", lastSyncTime))

		// get the gateway histories that still need to be sync'd to Azure
		histories, err = inst.db.GetHistoriesByHostUUID(host.UUID, lastSyncTime, now) // fetches histories that have been added since the last sync
		if err != nil {
			inst.inauroazuresyncErrorMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID() error: %v", err))
			inst.inauroazuresyncErrorMsg(err)
			return false, err
		}
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID(): %v", len(histories)))
		/*
			for _, history := range histories {
				inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() GetHistoriesByHostUUID() history: %+v", history))
			}
		*/

		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() history count: %+v", len(histories)))

		if len(histories) > 0 {
			bulkInauroHistoryPayloadsArray, latestHistoryTime, _ := inst.packageHistoriesToInauroPayloads(host.UUID, histories)
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() packageHistoriesToInauroPayloads() bulkInauroHistoryPayloadsArray: %+v", bulkInauroHistoryPayloadsArray))
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureSensorHistories() packageHistoriesToInauroPayloads() latestHistoryTime: %+v", latestHistoryTime))

			if len(bulkInauroHistoryPayloadsArray) < 0 {
				inst.inauroazuresyncErrorMsg(fmt.Sprintf("syncAzureSensorHistories() no histories to store on hostUUID: %v. There were histories, so check for errors", host.UUID))
				continue
			}

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
			} else {
				// if the Azure push was successful save the latest history sync time to the host in JSON storage.
				pluginStorage.LastSyncByGateway[host.UUID] = latestHistoryTime
			}
			azureClient.Close()
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureSensorHistories() azureClient.Close() error:", err)
			} else {
				inst.inauroazuresyncDebugMsg("syncAzureSensorHistories()  azureClient.Close() CLOSED")
			}
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("azure iot hub: Stored %v new sensor records", len(histories)))

		} else {
			inst.inauroazuresyncErrorMsg(fmt.Sprintf("syncAzureSensorHistories() no histories to store on hostUUID: %v", host.UUID))
		}
	}
	// save the updated lastSyncTime to module storage
	inst.setPluginConfStorage(pluginStorage)
	return true, nil
}

func (inst *Instance) syncAzureGatewayPayloads() (bool, error) {
	inst.inauroazuresyncDebugMsg("azure gateway payload sync has been called...")

	hosts, err := inst.db.GetHosts(argspkg.Args{})
	if err != nil {
		return false, err
	}
	inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() hosts: %v", len(hosts)))

	gatewayDetailsMap, err := inst.getGatewayDetailsFromConfig()
	if err != nil {
		inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() getGatewayDetailsFromConfig() error:", err)
		return false, err
	}

	for _, host := range hosts {
		inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads()  host: %+v", host))

		timestamp := time.Now().Truncate(time.Second)

		foundGatewayDetails := false
		for hostKey, gatewayDetails := range gatewayDetailsMap {
			if hostKey != host.UUID {
				continue
			}
			foundGatewayDetails = true

			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() gatewayDetails: %+v", gatewayDetails))

			gatewayPayload := InauroGatewayPayload{
				TimestampUTC: timestamp.UTC().Format(time.RFC3339),
				GatewayID:    gatewayDetails.AzureDeviceId,
				GatewayIMEI:  gatewayDetails.RouterIMEI,
				GatewayICCID: gatewayDetails.SIMICCID,
				Latitude:     gatewayDetails.Latitude,
				Longitude:    gatewayDetails.Longitude,
				Network:      gatewayDetails.NetworkType,
			}

			// ping check on the gateway/host
			cli := cligetter.GetEdgeClientFastTimeout(host)
			globalUUID, pingable, isValidToken := cli.Ping()
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() cli.Ping() globalUUID %v, pingable %v, isValidToken %v", globalUUID, pingable, isValidToken))
			if pingable == false {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() cli.Ping() error: make it accessible at first")
			}
			if isValidToken == false || globalUUID == nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() cli.Ping() error: configure valid token at first")
			}
			gatewayPayload.Ping = pingable

			networkingInfo, err := cli.GetNetworking()
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() cli.GetNetworking() gateway: ", gatewayDetails.AzureDeviceId, "  error:", err)
			}

			networkingPayloadInfo := make([]InauroGatewayNetworking, 0)
			for i, netInfo := range networkingInfo {
				inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() i: %v, netInfo: %+v", i, netInfo))
				// if strings.HasPrefix(netInfo.Interface, "eth") {
				if netInfo.Interface == "eth0" || netInfo.Interface == "eth1" {
					networkingPayloadInfo = append(networkingPayloadInfo, InauroGatewayNetworking{
						Interface: netInfo.Interface,
						IP:        netInfo.IP,
						MAC:       netInfo.MacAddress,
					})
				}
			}
			gatewayPayload.Networking = networkingPayloadInfo
			inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() networkingPayloadInfo: %+v", networkingPayloadInfo))

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
			} else {
				inst.inauroazuresyncDebugMsg(fmt.Sprintf("syncAzureGatewayPayloads() SendEvent() Success Host/Gateway: %v, %v", host.UUID, gatewayDetails.AzureDeviceId))
			}
			azureClient.Close()
			if err != nil {
				inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() azureClient.Close() error:", err)
			}
			inst.inauroazuresyncDebugMsg("syncAzureGatewayPayloads()  azureClient.Close() gateway: ", gatewayDetails.AzureDeviceId, "  CLOSED")
			break
		}
		if !foundGatewayDetails {
			inst.inauroazuresyncErrorMsg("syncAzureGatewayPayloads() error: missing gateway details for host: ", host.UUID)
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
