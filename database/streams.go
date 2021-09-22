package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) GetStreams(args api.Args) ([]*model.Stream, error) {
	var streamsModel []*model.Stream
	query := d.buildStreamQuery(args)
	query.Find(&streamsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamsModel, nil
}

func (d *GormDatabase) GetStream(uuid string, args api.Args) (*model.Stream, error) {
	var streamModel *model.Stream
	query := d.buildStreamQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamModel, nil
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
	var streamModel *model.Stream
	query := d.DB.Where("uuid = ? ", uuid).Delete(&streamModel)
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
	var streamModel *model.Stream
	if err := d.DB.Preload("FlowNetworks").Where("uuid = ?", uuid).Find(&streamModel).Error; err != nil {
		return nil, err
	}
	if len(body.FlowNetworks) > 0 {
		if err := d.DB.Model(&streamModel).Association("FlowNetworks").Replace(body.FlowNetworks); err != nil {
			return nil, err
		}
	}
	if len(body.Tags) > 0 {
		if err := d.DB.Model(&streamModel).Association("Tags").Replace(body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&streamModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return streamModel, nil
}

func (d *GormDatabase) DropStreams() (bool, error) {
	var streamModel *model.Stream
	query := d.DB.Where("1 = 1").Delete(&streamModel)
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
	flow, err := d.GetFlowNetwork(flowUUID, api.Args{})
	if err != nil {
		return nil, nil, err
	}
	return stream, flow, nil
}

// GetStreamByField ie: get stream by its name as an example
func (d *GormDatabase) GetStreamByField(field string, value string, args api.Args) (*model.Stream, error) {
	var streamModel *model.Stream
	f := fmt.Sprintf("%s = ? ", field)
	query := d.buildStreamQuery(args)
	query = query.Where(f, value).First(&streamModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamModel, nil
}
