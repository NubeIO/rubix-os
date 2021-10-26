package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	var flowNetworkClonesModel []*model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkClonesModel, nil
}

func (d *GormDatabase) GetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkCloneByArgs(args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) RefreshFlowNetworkClonesConnections() (*bool, error) {
	var flowNetworkClones []*model.FlowNetworkClone
	d.DB.Where("is_master_slave = ?", false).Find(&flowNetworkClones)
	for _, fnc := range flowNetworkClones {
		cli := client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
		token, err := cli.Login(&model.LoginBody{Username: *fnc.FlowUsername, Password: *fnc.FlowPassword})
		if err != nil {
			fnc.IsError = utils.NewTrue()
			fnc.ErrorMsg = utils.NewStringAddress(err.Error())
		} else {
			fnc.IsError = utils.NewFalse()
			fnc.ErrorMsg = nil
			fnc.FlowToken = utils.NewStringAddress(token.AccessToken)
		}
		d.DB.Model(&fnc).Updates(fnc)
	}
	return utils.NewTrue(), nil
}
