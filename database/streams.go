package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetStreams(args api.Args) ([]*model.Stream, error) {
	var gatewaysModel []*model.Stream
	query := d.createStreamQueryWithArgs(args)
	query.Find(&gatewaysModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return gatewaysModel, nil
}

func (d *GormDatabase) GetStream(uuid string, args api.Args) (*model.Stream, error) {
	var gatewayModel *model.Stream
	query := d.createStreamQueryWithArgs(args)
	query = query.Where("uuid = ? ", uuid).First(&gatewayModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return gatewayModel, nil
}

func (d *GormDatabase) CreateStream(body *model.Stream, AddToParent string) (*model.Stream, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Stream)
	body.Name = nameIsNil(body.Name)
	var flowNetwork model.FlowNetwork
	flowNetwork.UUID = AddToParent
	body.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
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
	query := d.DB.Where("uuid = ?", uuid).Find(&gatewayModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&gatewayModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
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

func (d *GormDatabase) createStreamQueryWithArgs(args api.Args) *gorm.DB {
	query := d.DB
	if args.FlowNetworks {
		query = query.Preload("FlowNetworks")
	}
	if args.Producers {
		query = query.Preload("Producers")
		if args.Writers {
			query = query.Preload("Producers.WriterClone")
		}
	}
	if args.Consumers {
		query = query.Preload("Consumers")
		if args.Writers {
			query = query.Preload("Consumers.Writer")
		}
	}
	if args.CommandGroups {
		query = query.Preload("CommandGroups")
	}
	return query
}
