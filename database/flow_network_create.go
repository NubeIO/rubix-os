package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
)

/*
CreateFlowNetwork
- Create UUID
- Create Name if doesn't exist
- Create SyncUUID
- If it's pointing local device make that is_remote=false forcefully (straightforward: 0.0.0.0 can be remote)
- If local then don't apply rollback feature, it leads deadlock
- Create flow_network
- Edit body with localstorage FlowNetwork details
- Edit body with device_info
- Create FlowNetworkClone
- Update FlowNetwork with FlowNetworkClone details
- Update sync_uuid with FlowNetworkClone's sync_uuid
*/
func (d *GormDatabase) CreateFlowNetwork(body *model.FlowNetwork) (*model.FlowNetwork, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetwork)
	if body.FlowNetworkParentUUID != "" { //is a sub-flow-network
		flowNetwork, err := d.GetFlowNetwork(body.FlowNetworkParentUUID, api.Args{})
		if err != nil {
			return nil, err
		}
		flowNetwork.Name = body.Name
		flowNetwork.FlowNetworkParentUUID = body.FlowNetworkParentUUID
		isMasterSlave, cli, isRemote, tx, err := d.editFlowNetworkBody2(flowNetwork)
		if err != nil {
			return nil, err
		}
		if err := tx.Create(&body).Error; err != nil {
			if isRemote {
				tx.Rollback()
			}
			return nil, err
		}
		network, err := d.afterCreateUpdateFlowNetwork(flowNetwork, isMasterSlave, cli, isRemote, tx)
		if err != nil {
			return nil, err
		}
		return network, err
	} else {
		//decide what type of network it is, as in master, remote and so on and update the body to suit
		isMasterSlave, cli, isRemote, tx, err := d.editFlowNetworkBody(body)
		if err != nil {
			return nil, err
		}
		if err := tx.Create(&body).Error; err != nil {
			if isRemote {
				tx.Rollback()
			}
			return nil, err
		}
		network, err := d.afterCreateUpdateFlowNetwork(body, isMasterSlave, cli, isRemote, tx)
		if err != nil {
			return nil, err
		}
		return network, err
	}

}

func (d *GormDatabase) editFlowNetworkBody2(body *model.FlowNetwork) (bool, *client.FlowClient, bool, *gorm.DB, error) {
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = utils.MakeUUID()
	body.IsRemote = utils.NewTrue()
	isMasterSlave := utils.IsTrue(body.IsMasterSlave)
	deviceInfo, err := d.GetDeviceInfo()
	cli := client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, model.IsFNCreator(body))
	remoteDeviceInfo, err := cli.DeviceInfo()
	if err != nil {
		return false, nil, false, nil, err
	} else {
		if deviceInfo.GlobalUUID == remoteDeviceInfo.GlobalUUID {
			//body.IsRemote = utils.NewFalse()
			//if !isMasterSlave {
			//	body.FlowHTTPS = utils.NewFalse()
			//	body.FlowIP = utils.NewStringAddress("0.0.0.0")
			//	body.IsRemote = utils.NewFalse()
			//}
		}
	}
	isRemote := utils.IsTrue(body.IsRemote)
	// rollback is needed only when flow-network is remote,
	// if we make it true in local it blocks the next transaction of clone creation which leads deadlock
	var tx *gorm.DB
	if tx = d.DB; isRemote {
		tx = d.DB.Begin()
	}
	return isMasterSlave, cli, isRemote, tx, nil
}

func (d *GormDatabase) editFlowNetworkBody(body *model.FlowNetwork) (bool, *client.FlowClient, bool, *gorm.DB, error) {
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = utils.MakeUUID()
	body.IsRemote = utils.NewTrue()
	isMasterSlave := utils.IsTrue(body.IsMasterSlave)
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return false, nil, false, nil, err
	}
	if isMasterSlave { //if it's a master
		body.FlowHTTPS = nil
		body.FlowIP = nil
		body.FlowPort = nil
		body.FlowToken = nil
	} else {
		conf := config.Get()
		if body.FlowIP == nil {
			body.FlowIP = &conf.Server.ListenAddr
		}
		if body.FlowPort == nil {
			body.FlowPort = &conf.Server.Port
		}
		if body.FlowUsername == nil {
			return false, nil, false, nil, errors.New("FlowUsername can't be null when we it's not master/slave flow network")
		}
		if body.FlowPassword == nil {
			return false, nil, false, nil, errors.New("FlowPassword can't be null when we it's not master/slave flow network")
		}
		accessToken, err := client.GetFlowToken(*body.FlowIP, *body.FlowPort, *body.FlowUsername, *body.FlowPassword)
		if err != nil {
			return false, nil, false, nil, err
		}
		body.FlowToken = accessToken
	}

	cli := client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, model.IsFNCreator(body))
	remoteDeviceInfo, err := cli.DeviceInfo()
	if err != nil {
		return false, nil, false, nil, err
	} else {
		if deviceInfo.GlobalUUID == remoteDeviceInfo.GlobalUUID {
			body.IsRemote = utils.NewFalse()
			if !isMasterSlave {
				body.FlowHTTPS = utils.NewFalse()
				body.FlowIP = utils.NewStringAddress("0.0.0.0")
				body.IsRemote = utils.NewFalse()
			}
		}
	}
	isRemote := utils.IsTrue(body.IsRemote)
	// rollback is needed only when flow-network is remote,
	// if we make it true in local it blocks the next transaction of clone creation which leads deadlock
	var tx *gorm.DB
	if tx = d.DB; isRemote {
		tx = d.DB.Begin()
	}
	return isMasterSlave, cli, isRemote, tx, nil
}

func (d *GormDatabase) afterCreateUpdateFlowNetwork(body *model.FlowNetwork, isMasterSlave bool, cli *client.FlowClient, isRemote bool, tx *gorm.DB) (*model.FlowNetwork, error) {
	bodyToSync := *body
	if body.FlowNetworkParentUUID != "" {

	} else if !isMasterSlave {
		var localStorageFlowNetwork *model.LocalStorageFlowNetwork
		if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
			if isRemote {
				tx.Rollback()
			}
			return nil, err
		}
		bodyToSync.FlowHTTPS = localStorageFlowNetwork.FlowHTTPS
		bodyToSync.FlowIP = utils.NewStringAddress(localStorageFlowNetwork.FlowIP)
		bodyToSync.FlowPort = utils.NewInt(localStorageFlowNetwork.FlowPort)
		bodyToSync.FlowUsername = utils.NewStringAddress(localStorageFlowNetwork.FlowUsername)
		bodyToSync.FlowPassword = utils.NewStringAddress(localStorageFlowNetwork.FlowPassword)
		bodyToSync.FlowToken = utils.NewStringAddress(localStorageFlowNetwork.FlowToken)
	}
	err := d.syncAndEditFlowNetwork(cli, body, &bodyToSync)
	if err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	if isRemote {
		tx.Commit()
	}
	d.DB.Model(&body).Updates(body)
	return body, nil
}

func (d *GormDatabase) syncAndEditFlowNetwork(cli *client.FlowClient, body *model.FlowNetwork, bodyToSync *model.FlowNetwork) error { //

	if body.FlowNetworkParentUUID != "" { //is a sub-flow-network
		_, err := cli.SyncFlowNetwork(body)
		if err != nil {
			return err
		}
		return nil
	} else {
		deviceInfo, err := d.GetDeviceInfo()
		if err != nil {
			return err
		}
		bodyToSync.GlobalUUID = deviceInfo.GlobalUUID
		bodyToSync.ClientId = deviceInfo.ClientId
		bodyToSync.ClientName = deviceInfo.ClientName
		bodyToSync.SiteId = deviceInfo.SiteId
		bodyToSync.SiteName = deviceInfo.SiteName
		bodyToSync.DeviceId = deviceInfo.DeviceId
		bodyToSync.DeviceName = deviceInfo.DeviceName
		res, err := cli.SyncFlowNetwork(bodyToSync)
		if err != nil {
			return err
		}
		body.SyncUUID = res.SyncUUID
		body.GlobalUUID = res.GlobalUUID
		body.ClientId = res.ClientId
		body.ClientName = res.ClientName
		body.SiteId = res.SiteId
		body.SiteName = res.SiteName
		body.DeviceId = res.DeviceId
		body.DeviceName = res.DeviceName
		return nil
	}

}
