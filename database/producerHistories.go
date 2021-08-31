package database

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"reflect"
	"strings"
	"time"
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

func CompareJSONToStruct(bytes []byte, empty interface{}) bool {
	var mapped map[string]interface{}
	if err := json.Unmarshal(bytes, &mapped); err != nil {
		return false
	}
	emptyValue := reflect.ValueOf(empty).Type()
	// check if number of fields is the same
	if len(mapped) != emptyValue.NumField() {
		return false
	}
	// check if field names are the same
	for key := range mapped {
		if field, found := emptyValue.FieldByName(key); found {
			if !strings.EqualFold(key, strings.Split(field.Tag.Get("json"), ",")[0]) {
				return false
			}
		}
	}
	return true
}

// GetProducerHistory returns the history for the given id or nil.
func (d *GormDatabase) GetProducerHistory(uuid string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("producer_uuid = ? ", uuid).Order("timestamp DESC").First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil

}


// HistoryByProducerUUID returns the history for the given id or nil.
func (d *GormDatabase) HistoryByProducerUUID(uuid string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("producer_uuid` = ? ", uuid).First(&historyModel)
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


func (d *GormDatabase) CreateBulkProducerHistory(history []*model.ProducerHistory) (bool, error) {
	for _, hist := range history {
		ph := new(model.ProducerHistory)
		ph.ProducerUUID = hist.ProducerUUID
		ph.DataStore = hist.DataStore
		ph.Timestamp = time.Now().UTC()
		_, err := d.CreateProducerHistory(ph)
		if err != nil {
			return true, err
		}
	}
	return false, nil
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



