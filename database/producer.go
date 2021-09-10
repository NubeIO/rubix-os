package database

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/datatypes"
	"time"
)

type Producer struct {
	*model.Producer
}

// GetProducers get all of them
func (d *GormDatabase) GetProducers() ([]*model.Producer, error) {
	var producersModel []*model.Producer
	query := d.DB.Preload("WriterClone").Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateProducer make it
func (d *GormDatabase) CreateProducer(body *model.Producer) (*model.Producer, error) {
	//call points and make it exists
	_, err := d.GetStream(body.StreamUUID)
	if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the gateway UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	err = d.DB.Create(&body).Error
	if err != nil {
		return nil, errorMsg("CreateProducer", "error on trying to add a new Producer", nil)
	}
	return body, nil
}

// GetProducer get it
func (d *GormDatabase) GetProducer(uuid string) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}

// UpdateProducer  update it
func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&producerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}

type Point struct {
	Priority model.Priority `json:"priority"`
}

// ProducerCOV  update it
func (d *GormDatabase) ProducerCOV(uuid string, writeData datatypes.JSON) error {
	ph := new(model.ProducerHistory)
	ph.UUID = utils.MakeTopicUUID("")
	ph.ProducerUUID = uuid
	ph.DataStore = writeData
	ph.Timestamp = time.Now().UTC()
	_, err := d.CreateProducerHistory(ph)
	if err != nil {
		return err
	}
	pnt := Point{}
	err = json.Unmarshal(writeData, &pnt)
	if err != nil {
		return err
	}
	return err
}

// GetProducerByField returns the point for the given field ie name or nil.
//for example get a producer by its producer_thing_uuid
func (d *GormDatabase) GetProducerByField(field string, value string) (*model.Producer, error) {
	var producerModel *model.Producer
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).First(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil

}

// UpdateProducerByField get by field and update.
//for example update a producer by its producer_thing_uuid
func (d *GormDatabase) UpdateProducerByField(field string, value string, body *model.Producer) (*model.Producer, error) {
	var producerModel *model.Producer
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Find(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&producerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}

// DeleteProducer deletes it
func (d *GormDatabase) DeleteProducer(uuid string) (bool, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&producerModel)
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

// DropProducers delete all.
func (d *GormDatabase) DropProducers() (bool, error) {
	var producerModel *model.Producer
	query := d.DB.Where("1 = 1").Delete(&producerModel)
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
