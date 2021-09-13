package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/gorm"
)

// GetStreams get all of them
func (d *GormDatabase) GetStreams(withChildren bool) ([]*model.Stream, error) {
	var gatewaysModel []*model.Stream
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("FlowNetworks").Preload("Producer").Preload("Producer.WriterClone").Preload("Consumer").Preload("Consumer.Writer").Find(&gatewaysModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return gatewaysModel, nil
	} else {
		query := d.DB.Find(&gatewaysModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return gatewaysModel, nil
	}

}

// CreateStream make it
func (d *GormDatabase) CreateStream(body *model.Stream) (*model.Stream, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Stream)
	body.Name = nameIsNil(body.Name)
	err := d.DB.Create(&body).Error
	if err != nil {
		return nil, errorMsg("CreateStreamGateway", "error on trying to add a new stream gateway", nil)
	}
	return body, nil
}

// GetStream get it
func (d *GormDatabase) GetStream(uuid string, withChildren bool) (*model.Stream, error) {
	var gatewayModel *model.Stream
	var query *gorm.DB = nil
	if withChildren {
		query = d.DB.Preload("FlowNetworks").Preload("Producer").Preload("Producer.WriterClone").Preload("Consumer").Preload("Consumer.Writer").Find(&gatewayModel).Where("uuid = ? ", uuid).First(&gatewayModel)
	} else {
		query = d.DB.Find(&gatewayModel)
	}
	if query.Error != nil {
		return nil, query.Error
	}
	return gatewayModel, nil
}

// GetFlowUUID get flowNetwork stream and flow network
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

// DeleteStream deletes it
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

// UpdateStream  update it
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

// DropStreams delete all.
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
