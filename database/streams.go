package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) GetStreams(args api.Args) ([]*model.Stream, error) {
	var gatewaysModel []*model.Stream
	query := d.createStreamQuery(args)
	query.Find(&gatewaysModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return gatewaysModel, nil
}

func (d *GormDatabase) GetStream(uuid string, args api.Args) (*model.Stream, error) {
	var gatewayModel *model.Stream
	query := d.createStreamQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&gatewayModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return gatewayModel, nil
}

func (d *GormDatabase) CreateStream(body *model.Stream) (*model.Stream, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Stream)
	body.Name = nameIsNil(body.Name)
	err := d.DB.Create(&body).Error
	if err != nil {
		return nil, errorMsg("CreateStreamGateway", "error on trying to add a new stream gateway", nil)
	}
	return body, nil
}

func (d *GormDatabase) DeleteStream(uuid string) (bool, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("uuid = ? ", uuid).Delete(&gatewayModel)
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

func (d *GormDatabase) UpdateStream(uuid string, body *model.Stream) (*model.Stream, error) {
	var gatewayModel *model.Stream
	if err := d.DB.Preload("FlowNetworks").Where("uuid = ?", uuid).Find(&gatewayModel).Error; err != nil {
		return nil, err
	}
	if len(body.FlowNetworks) > 0 {
		if err := d.DB.Model(&gatewayModel).Association("FlowNetworks").Replace(body.FlowNetworks); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&gatewayModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return gatewayModel, nil
}

func (d *GormDatabase) DropStreams() (bool, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("1 = 1").Delete(&gatewayModel)
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

func (d *GormDatabase) GetFlowUUID(uuid string) (*model.Stream, *model.FlowNetwork, error) {
	var stream *model.Stream
	query := d.DB.Preload("FlowNetworks").Where("uuid = ? ", uuid).First(&stream)
	if query.Error != nil {
		return nil, nil, query.Error
	}
	flowUUID := ""
	for _, net := range stream.FlowNetworks {
		flowUUID = net.UUID
	}
	flow, err := d.GetFlowNetwork(flowUUID)
	if err != nil {
		return nil, nil, err
	}
	return stream, flow, nil
}
