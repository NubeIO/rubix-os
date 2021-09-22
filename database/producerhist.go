package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"strings"
	"time"
)

// GetProducerHistories returns all histories.
func (d *GormDatabase) GetProducerHistories(args api.Args) ([]*model.ProducerHistory, error) {
	var historiesModel []*model.ProducerHistory
	query := d.buildProducerHistoryQuery(args)
	query.Find(&historiesModel)
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

// HistoryLatestByProducerUUID returns the history for the given id or nil.
func (d *GormDatabase) HistoryLatestByProducerUUID(uuid string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	t := fmt.Sprintf("timestamp %s", "DESC")
	query := d.DB.Where("producer_uuid = ? ", uuid).Order(t).First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil
}

// HistoriesAllByProducerUUID returns the history for the given id or nil.
func (d *GormDatabase) HistoriesAllByProducerUUID(uuid string, order string) ([]*model.ProducerHistory, int64, error) {
	var count int64
	var historiesModel []*model.ProducerHistory
	order = strings.ToUpper(strings.TrimSpace(order))
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}
	t := fmt.Sprintf("timestamp %s", order)
	q := d.DB.Where("producer_uuid = ? ", uuid).Order(t).Find(&historiesModel) //ASC or DESC
	q.Count(&count)
	return historiesModel, count, nil
}

// CreateProducerHistory creates a thing.
func (d *GormDatabase) CreateProducerHistory(body *model.ProducerHistory) (*model.ProducerHistory, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
	hist, count, err := d.HistoriesAllByProducerUUID(body.ProducerUUID, "DESC")
	if err != nil {
		return nil, err
	}
	var limit int64 = 10
	//TODO add in the limit as a field in the producer
	if count >= limit {
		for i, e := range hist {
			if i >= int(limit) {
				_, err := d.DeleteProducerHistory(e.UUID)
				if err != nil {
					return nil, err
				}
			}
		}
	}
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

// DeleteProducerHistory delete a history.
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

// DropProducerHistories delete all.
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
