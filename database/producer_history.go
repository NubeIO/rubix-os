package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
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

// GetProducerHistoriesByProducerUUID returns the history for the given producer_uuid or nil.
func (d *GormDatabase) GetProducerHistoriesByProducerUUID(pUuid string, args api.Args) ([]*model.ProducerHistory, int64, error) {
	var count int64
	var historiesModel []*model.ProducerHistory
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("producer_uuid = ?", pUuid)
	query.Find(&historiesModel)
	query.Count(&count)
	return historiesModel, count, nil
}

// GetProducerHistoriesByProducerName returns the history for the given producer name.
func (d *GormDatabase) GetProducerHistoriesByProducerName(name string) ([]*model.ProducerHistory, int64, error) {
	var count int64
	var historiesModel []*model.ProducerHistory
	p, err := d.GetOneProducerByArgs(api.Args{Name: &name})
	if err != nil {
		return nil, 0, err
	}
	query := d.DB.Where("producer_uuid = ? ", p.UUID).Order("timestamp desc")
	if query.Error != nil {
		return nil, 0, query.Error
	}
	query.Find(&historiesModel)
	query.Count(&count)
	return historiesModel, count, nil
}

// GetLatestProducerHistoryByProducerName returns the latest history for the given producer_name or nil.
func (d *GormDatabase) GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	p, err := d.GetOneProducerByArgs(api.Args{Name: &name})
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("producer_uuid = ? ", p.UUID).Order("timestamp desc").First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil
}

// GetLatestProducerHistoryByProducerUUID returns the latest history for the given producer_uuid or nil.
func (d *GormDatabase) GetLatestProducerHistoryByProducerUUID(pUuid string) (*model.ProducerHistory, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("producer_uuid = ? ", pUuid).Order("timestamp desc").First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil
}

// GetProducerHistoriesByPointUUIDs returns producer histories of PointUUIDs
func (d *GormDatabase) GetProducerHistoriesByPointUUIDs(pointUUIDs []string, args api.Args) ([]*model.ProducerHistory, error) {
	var proHistoriesModel []*model.ProducerHistory
	subQuery := d.DB.Model(&model.Producer{}).Select("uuid").Where("producer_thing_uuid IN ?", pointUUIDs)
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("producer_uuid in (?)", subQuery).Find(&proHistoriesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return proHistoriesModel, nil
}

// GetProducerHistoriesPoints returns point histories of producer histories
func (d *GormDatabase) GetProducerHistoriesPoints(args api.Args) ([]*model.History, error) {
	var historiesModel []*model.History
	var proHistoriesModel []*model.ProducerHistory
	subQuery := d.DB.Model(&model.Producer{}).Select("uuid").Where("producer_thing_class = ?", "point")
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("producer_uuid in (?)", subQuery).Find(&proHistoriesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, pHis := range proHistoriesModel {
		if pHis.PresentValue == nil {
			continue
		}
		historiesModel = append(historiesModel,
			&model.History{ID: pHis.ID, UUID: pHis.ProducerUUID, Value: *pHis.PresentValue, Timestamp: pHis.Timestamp})
	}
	return historiesModel, nil
}

// GetProducerHistoriesPointsForSync returns point histories of producer histories
func (d *GormDatabase) GetProducerHistoriesPointsForSync(id string, timeStamp string) ([]*model.History, error) {
	var historiesModel []*model.History
	var proHistoriesModel []*model.ProducerHistory
	query := d.DB.Where("id = ?", id).Where("datetime(timestamp) = datetime(?)", timeStamp).
		Find(&proHistoriesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(proHistoriesModel) == 0 {
		id = "0"
	}
	subQuery := d.DB.Model(&model.Producer{}).Select("uuid").Where("producer_thing_class = ?", "point")
	query = d.DB.Where("id > ?", id).Where("producer_uuid in (?)", subQuery).Find(&proHistoriesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, pHis := range proHistoriesModel {
		if pHis.PresentValue == nil {
			continue
		}
		historiesModel = append(historiesModel,
			&model.History{ID: pHis.ID, UUID: pHis.ProducerUUID, Value: *pHis.PresentValue, Timestamp: pHis.Timestamp})
	}
	return historiesModel, nil
}

func (d *GormDatabase) CreateProducerHistory(body *model.ProducerHistory) (*model.ProducerHistory, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateBulkProducerHistory(histories []*model.ProducerHistory) (bool, error) {
	if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(histories, 1000).Error; err != nil {
		log.Error("Issue on creating bulk producer histories")
		return false, err
	}
	return true, nil
}

// DeleteProducerHistoriesByProducerUUID delete all history for given producer_uuid.
func (d *GormDatabase) DeleteProducerHistoriesByProducerUUID(pUuid string, args api.Args) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("producer_uuid = ? ", pUuid)
	query = query.Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}

// DeleteProducerHistoriesBeforeTimestamp delete producer histories before timestamp
func (d *GormDatabase) DeleteProducerHistoriesBeforeTimestamp(ts string) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("timestamp < datetime(?)", ts)
	query.Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}
