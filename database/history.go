package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"time"
)

// GetHistoriesForSync returns all histories after id.
// We order by `uuid` i.e. `producer_uuid`, so all similar data comes on same block which helps to reduce data query
// for fetching the data from points, devices, networks etc.
// TODO: only used in influx db, remove this when influx plugin gets removed
func (d *GormDatabase) GetHistoriesForSync(lastSyncId int) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("id > (?)", lastSyncId)
	query.Order("uuid ASC").Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesForPostgresSync returns all histories after history_id ordered by history_id.
func (d *GormDatabase) GetHistoriesForPostgresSync(lastSyncId int) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("history_id > (?)", lastSyncId)
	query.Order("history_id ASC").Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesByHostUUID fetches histories based on a given HostUUID and startTime
func (d *GormDatabase) GetHistoriesByHostUUID(hostUUID string, startTime, endTime time.Time) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("host_uuid = ? AND timestamp > ? AND timestamp <= ?", hostUUID, startTime, endTime).Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// CreateBulkHistory bulk creates a thing.
func (d *GormDatabase) CreateBulkHistory(histories []*model.History) (bool, error) {
	if err := d.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(histories, 1000).Error; err != nil {
		log.Error("Issue on creating bulk histories")
		return false, err
	}
	return true, nil
}
