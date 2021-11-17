package database

import (
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) GetLocalStorageFlowNetwork() (*model.LocalStorageFlowNetwork, error) {
	var localStorageFlowNetwork *model.LocalStorageFlowNetwork
	if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
		return nil, err
	}
	return localStorageFlowNetwork, nil
}

func (d *GormDatabase) UpdateLocalStorageFlowNetwork(body *model.LocalStorageFlowNetwork) (*model.LocalStorageFlowNetwork, error) {
	var lsfn *model.LocalStorageFlowNetwork
	if err := d.DB.First(&lsfn).Error; err != nil {
		return nil, err
	}
	conf := config.Get()
	token, err := client.GetFlowToken(conf.Server.ListenAddr, conf.Server.RSPort, body.FlowUsername, body.FlowPassword)
	if err != nil {
		return nil, err
	}
	body.FlowToken = *token
	if err := d.DB.Model(&lsfn).Updates(&body).Error; err != nil {
		return nil, err
	}
	return lsfn, nil
}

func (d *GormDatabase) RefreshLocalStorageFlowToken() (*bool, error) {
	var lsfn *model.LocalStorageFlowNetwork
	if err := d.DB.First(&lsfn).Error; err != nil {
		return nil, err
	}
	conf := config.Get()
	token, err := client.GetFlowToken(conf.Server.ListenAddr, conf.Server.RSPort, lsfn.FlowUsername, lsfn.FlowPassword)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Model(&lsfn).Updates(model.LocalStorageFlowNetwork{FlowToken: *token}).Error; err != nil {
		return nil, err
	}
	return utils.NewTrue(), nil
}
