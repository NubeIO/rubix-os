package database

import (
	"errors"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetFlowNetworks(args api.Args) ([]*model.FlowNetwork, error) {
	var flowNetworksModel []*model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Find(&flowNetworksModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworksModel, nil
}

func (d *GormDatabase) GetFlowNetwork(uuid string, args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkByArgs(args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkModel, nil
}

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
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = utils.MakeUUID()
	isMasterSlave := utils.IsTrue(body.IsMasterSlave)
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	if isMasterSlave {
		body.FlowHTTPS = nil
		body.FlowIP = nil
		body.FlowPort = nil
		body.IsRemote = nil
		body.FlowToken = nil
		cli := client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, model.IsFNCreator(body))
		remoteDeviceInfo, err := cli.DeviceInfo()
		if err != nil {
			return nil, err
		} else {
			body.IsRemote = utils.NewTrue()
			if deviceInfo.GlobalUUID == remoteDeviceInfo.GlobalUUID {
				body.IsRemote = utils.NewFalse()
			}
		}
	} else {
		if body.FlowIP == nil {
			body.FlowIP = utils.NewStringAddress("0.0.0.0")
		}
		if body.FlowPort == nil {
			body.FlowPort = utils.NewInt(1660)
		}
		if body.FlowToken == nil {
			body.FlowToken = utils.NewStringAddress("token")
		}
		if body.FlowUsername == nil {
			return nil, errors.New("FlowUsername can't be null when we it's not master/slave flow network")
		}
		if body.FlowPassword == nil {
			return nil, errors.New("FlowPassword can't be null when we it's not master/slave flow network")
		}
		if *body.FlowIP == "0.0.0.0" || *body.FlowIP == "127.0.0.0" || *body.FlowIP == "localhost" {
			body.FlowHTTPS = utils.NewFalse()
			body.FlowIP = utils.NewStringAddress("0.0.0.0")
			body.FlowPort = utils.NewInt(1660)
			body.IsRemote = utils.NewFalse()
		} else {
			body.IsRemote = utils.NewTrue()
		}
	}

	isRemote := utils.IsTrue(body.IsRemote)
	//rollback is needed only when flow-network is remote,
	//if we make it true in local it blocks the next transaction of clone creation which leads deadlock
	var tx *gorm.DB
	if tx = d.DB; isRemote {
		tx = d.DB.Begin()
	}
	cli := client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, model.IsFNCreator(body))
	token, err := cli.Login(&model.LoginBody{
		Username: *body.FlowUsername,
		Password: *body.FlowPassword,
	})
	if err != nil {
		return nil, err
	}
	body.FlowToken = utils.NewStringAddress(token.AccessToken)
	cli = client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, model.IsFNCreator(body))
	if err := tx.Create(&body).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	fnb := *body
	if !isMasterSlave {
		var localStorageFlowNetwork *model.LocalStorageFlowNetwork
		if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
			if isRemote {
				tx.Rollback()
			}
			return nil, err
		}
		fnb.FlowHTTPS = localStorageFlowNetwork.FlowHTTPS
		fnb.FlowIP = utils.NewStringAddress(localStorageFlowNetwork.FlowIP)
		fnb.FlowPort = utils.NewInt(localStorageFlowNetwork.FlowPort)
		fnb.FlowUsername = utils.NewStringAddress(localStorageFlowNetwork.FlowUsername)
		fnb.FlowPassword = utils.NewStringAddress(localStorageFlowNetwork.FlowPassword)
		fnb.FlowToken = utils.NewStringAddress(localStorageFlowNetwork.FlowToken)
		token, err := GetFlowToken(*body.FlowIP, *body.FlowPort, *body.FlowUsername, *body.FlowPassword)
		if err != nil {
			return nil, err
		}
		body.FlowToken = token
	}
	fnb.GlobalUUID = deviceInfo.GlobalUUID
	fnb.ClientId = deviceInfo.ClientId
	fnb.ClientName = deviceInfo.ClientName
	fnb.SiteId = deviceInfo.SiteId
	fnb.SiteName = deviceInfo.SiteName
	fnb.DeviceId = deviceInfo.DeviceId
	fnb.DeviceName = deviceInfo.DeviceName
	res, err := cli.SyncFlowNetwork(&fnb)
	if err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	if isRemote {
		tx.Commit()
	}
	body.SyncUUID = res.SyncUUID
	body.GlobalUUID = res.GlobalUUID
	body.ClientId = res.ClientId
	body.ClientName = res.ClientName
	body.SiteId = res.SiteId
	body.SiteName = res.SiteName
	body.DeviceId = res.DeviceId
	body.DeviceName = res.DeviceName
	d.DB.Model(&body).Updates(body)
	return body, nil
}

func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	if err := d.DB.Where("uuid = ?", uuid).Find(&flowNetworkModel).Error; err != nil {
		return nil, err
	}
	if len(body.Streams) > 0 {
		if err := d.DB.Model(&flowNetworkModel).Association("Streams").Replace(body.Streams); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&flowNetworkModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) DeleteFlowNetwork(uuid string) (bool, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.DB.Where("uuid = ? ", uuid).Delete(&flowNetworkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) DropFlowNetworks() (bool, error) {
	var networkModel *model.FlowNetwork
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) RefreshFlowNetworksConnections() (*bool, error) {
	var flowNetworks []*model.FlowNetwork
	d.DB.Where("is_master_slave IS NOT TRUE").Find(&flowNetworks)
	for _, fn := range flowNetworks {
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))
		token, err := cli.Login(&model.LoginBody{Username: *fn.FlowUsername, Password: *fn.FlowPassword})
		fnModel := model.FlowNetworkClone{}
		if err != nil {
			fnModel.IsError = utils.NewTrue()
			fnModel.ErrorMsg = utils.NewStringAddress(err.Error())
			fnModel.FlowToken = fn.FlowToken
		} else {
			fnModel.IsError = utils.NewFalse()
			fnModel.ErrorMsg = nil
			fnModel.FlowToken = utils.NewStringAddress(token.AccessToken)
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fn).Select("IsError", "ErrorMsg", "FlowToken").Updates(&fnModel).Error; err != nil {
			log.Error(err)
		}
	}
	return utils.NewTrue(), nil
}
