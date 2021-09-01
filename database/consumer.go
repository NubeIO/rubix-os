package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Consumers struct {
	*model.Consumer
}

// GetConsumers get all of them
func (d *GormDatabase) GetConsumers() ([]*model.Consumer, error) {
	var consumersModel []*model.Consumer
	query := d.DB.Preload("Writer").Find(&consumersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumersModel, nil
}

// CreateConsumer make it
func (d *GormDatabase) CreateConsumer(body *model.Consumer) (*model.Consumer, error) {
	_, err := d.GetStream(body.StreamUUID);if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the stream UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
	body.Name = nameIsNil(body.Name)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetConsumer get it
func (d *GormDatabase) GetConsumer(uuid string) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ? ", uuid).First(&consumerModel); if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil
}

// DeleteConsumer deletes it
func (d *GormDatabase) DeleteConsumer(uuid string) (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&consumerModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateConsumer  update it
func (d *GormDatabase) UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ?", uuid).Find(&consumerModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&consumerModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil

}

// DropConsumers delete all.
func (d *GormDatabase) DropConsumers() (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("1 = 1").Delete(&consumerModel)
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

