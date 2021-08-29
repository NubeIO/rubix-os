package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)



// GetProducerHistories returns all histories.
func (d *GormDatabase) GetProducerHistories() ([]*model.ProducerHistory, error) {
	var historiesModel []*model.ProducerHistory
	query := d.DB.Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil


}

// GetProducerHistory returns the history for the given id or nil.
func (d *GormDatabase) GetProducerHistory(uuid string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("uuid = ? ", uuid).First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil

}


// CreateProducerHistory creates a thing.
func (d *GormDatabase) CreateProducerHistory(body *model.ProducerHistory) (*model.ProducerHistory, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
	//todo make sure thing_uuid is provided
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}


// DeleteProducerHistory delete a history. TODO //add in by thing_uuid
func (d *GormDatabase) DeleteProducerHistory(uuid string) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("uuid = ? ", uuid).Delete(&historyModel)
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

// DropProducerHistories delete all. TODO //add in by thing_uuid
func (d *GormDatabase) DropProducerHistories() (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("1 = 1").Delete(&historyModel)
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



