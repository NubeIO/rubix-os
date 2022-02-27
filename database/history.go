package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"time"
)

// GetHistories returns all histories.
func (d *GormDatabase) GetHistories(args api.Args) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.buildHistoryQuery(args)
	query.Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesForSync returns all histories after id.
func (d *GormDatabase) GetHistoriesForSync(args api.Args) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("id > (?)", d.DB.Table("history_influx_logs").
		Select("IFNULL(MAX(last_sync_id),0)"))
	query.Order("uuid ASC").Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesByUUID returns the history for the given uuid or nil.
func (d *GormDatabase) GetHistoriesByUUID(uuid string, args api.Args) ([]*model.History, int64, error) {
	var count int64
	var historiesModel []*model.History
	query := d.buildHistoryQuery(args)
	query = query.Where("uuid = ?", uuid)
	query.Find(&historiesModel)
	query.Count(&count)
	return historiesModel, count, nil
}

// CreateHistory creates a thing.
func (d *GormDatabase) CreateHistory(body *model.History) (*model.History, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// CreateBulkHistory bulk creates a thing.
func (d *GormDatabase) CreateBulkHistory(history []*model.History) (bool, error) {
	for _, hist := range history {
		ph := new(model.History)
		ph.UUID = hist.UUID
		ph.Value = hist.Value
		ph.Timestamp = time.Now().UTC()
		_, err := d.CreateHistory(ph)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// DeleteHistoriesByUUID delete all history for given uuid.
func (d *GormDatabase) DeleteHistoriesByUUID(uuid string, args api.Args) (bool, error) {
	var historyModel *model.History
	query := d.buildHistoryQuery(args)
	query = query.Where("uuid = ? ", uuid)
	query.Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteHistory delete a history for given id
func (d *GormDatabase) DeleteHistory(id int) (bool, error) {
	var historyModel *model.History
	query := d.DB.Where("id = ? ", id).Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteAllHistories delete all histories.
func (d *GormDatabase) DeleteAllHistories() (bool, error) {
	var historyModel *model.History
	query := d.DB.Where("1 = 1")
	query.Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}
