package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
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

func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var fn *model.FlowNetwork
	if err := d.DB.Where("uuid = ?", uuid).First(&fn).Error; err != nil {
		return nil, err
	}
	if len(body.Streams) > 0 {
		return d.updateStreamsOnFlowNetwork(fn, body.Streams) // normally we just either edit flow_network or assign stream
	}
	if err := d.DB.Model(&fn).Updates(body).Error; err != nil {
		return nil, err
	}
	isMasterSlave, cli, isRemote, tx, err := d.editFlowNetworkBody(body)
	if err != nil {
		return nil, err
	}
	if err := tx.Model(&fn).Updates(&body).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	return d.afterCreateUpdateFlowNetwork(body, isMasterSlave, cli, isRemote, tx)
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
		accessToken, err := client.GetFlowToken(*fn.FlowIP, *fn.FlowPort, *fn.FlowUsername, *fn.FlowPassword)
		fnModel := model.FlowNetworkClone{}
		if err != nil {
			fnModel.IsError = utils.NewTrue()
			fnModel.ErrorMsg = utils.NewStringAddress(err.Error())
			fnModel.FlowToken = fn.FlowToken
		} else {
			fnModel.IsError = utils.NewFalse()
			fnModel.ErrorMsg = nil
			fnModel.FlowToken = accessToken
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fn).Select("IsError", "ErrorMsg", "FlowToken").Updates(&fnModel).Error; err != nil {
			log.Error(err)
		}
	}
	return utils.NewTrue(), nil
}

func (d *GormDatabase) updateStreamsOnFlowNetwork(fn *model.FlowNetwork, streams []*model.Stream) (*model.FlowNetwork, error) {
	if err := d.DB.Model(&fn).Association("Streams").Replace(streams); err != nil {
		return nil, err
	}
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	for _, stream := range fn.Streams {
		_ = d.SyncStreamFunction(fn, stream, deviceInfo)
	}
	return fn, nil
}
