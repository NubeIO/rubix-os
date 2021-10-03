package database

import (
	"github.com/NubeDev/flow-framework/model"
)

func (d *GormDatabase) GetLocalStorageFlowNetwork() (*model.LocalStorageFlowNetwork, error) {
	var localStorageFlowNetwork *model.LocalStorageFlowNetwork
	if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
		return nil, err
	}
	return localStorageFlowNetwork, nil
}

func (d *GormDatabase) UpdateLocalStorageFlowNetwork(body *model.LocalStorageFlowNetwork) (*model.LocalStorageFlowNetwork, error) {
	var localStorageFlowNetwork *model.LocalStorageFlowNetwork
	if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
		return nil, err
	}
	if err := d.DB.Model(&localStorageFlowNetwork).Updates(&body).Error; err != nil {
		return nil, err
	}
	return localStorageFlowNetwork, nil
}
