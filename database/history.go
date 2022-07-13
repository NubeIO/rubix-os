package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
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
// We order by `uuid` i.e. `producer_uuid`, so all similar data comes on same block which helps to reduce data query
// for fetching the data from points, devices, networks etc.
func (d *GormDatabase) GetHistoriesForSync(lastSyncId int) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("id > (?)", lastSyncId)
	query.Order("uuid ASC").Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesForPostgresSync returns all histories after id ordered by id.
func (d *GormDatabase) GetHistoriesForPostgresSync(lastSyncId int) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("id > (?)", lastSyncId)
	query.Order("id ASC").Find(&historiesModel)
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
func (d *GormDatabase) CreateBulkHistory(histories []*model.History) (bool, error) {
	if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(histories, 1000).Error; err != nil {
		log.Error("Issue on creating bulk histories")
		return false, err
	}
	return true, nil
}

// DeleteHistoriesByUUID delete all history for given uuid.
func (d *GormDatabase) DeleteHistoriesByUUID(uuid string, args api.Args) (bool, error) {
	var historyModel *model.History
	query := d.buildHistoryQuery(args)
	query = query.Where("uuid = ? ", uuid)
	query.Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}

// DeleteHistory delete a history for given id
func (d *GormDatabase) DeleteHistory(id int) (bool, error) {
	var historyModel *model.History
	query := d.DB.Where("id = ? ", id).Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}
