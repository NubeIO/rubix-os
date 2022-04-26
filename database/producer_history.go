package database

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
		priority := new(model.Priority)
		_ = json.Unmarshal(pHis.DataStore, &priority)
		highestPriorityValue := priority.GetHighestPriorityValue()
		value := 0.0
		if highestPriorityValue != nil {
			value = *highestPriorityValue
		}
		historiesModel = append(historiesModel,
			&model.History{ID: pHis.ID, UUID: pHis.ProducerUUID, Value: value, Timestamp: pHis.Timestamp})
	}
	return historiesModel, nil
}

func (d *GormDatabase) CreateProducerHistory(history *model.ProducerHistory) (*model.ProducerHistory, error) {
	return d.AppendProducerHistory(history)
}

// DeleteProducerHistoriesByProducerUUID delete all history for given producer_uuid.
func (d *GormDatabase) DeleteProducerHistoriesByProducerUUID(pUuid string, args api.Args) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("producer_uuid = ? ", pUuid)
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

// DeleteAllProducerHistories delete all histories.
func (d *GormDatabase) DeleteAllProducerHistories(args api.Args) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.buildProducerHistoryQuery(args)
	query = query.Where("1 = 1")
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

func (d *GormDatabase) AppendProducerHistory(body *model.ProducerHistory) (*model.ProducerHistory, error) {
	var limit = 100 // TODO add in the limit as a field in the producer
	var count int64
	ids := d.DB.Model(&model.ProducerHistory{}).
		Select("id").
		Where("producer_uuid = ?", body.ProducerUUID).
		Order("timestamp desc").
		Limit(limit - 1)
	ids.Count(&count)
	if count >= int64(limit) {
		query := d.DB.
			Where("producer_uuid = ?", body.ProducerUUID).
			Where("id not in (?)", ids).
			Delete(&model.ProducerHistory{})
		if query.Error != nil {
			return nil, query.Error
		}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// DeleteProducerHistoriesBeforeTimestamp delete producer histories before timestamp
func (d *GormDatabase) DeleteProducerHistoriesBeforeTimestamp(ts string) (bool, error) {
	var historyModel *model.ProducerHistory
	query := d.DB.Where("timestamp < datetime(?)", ts)
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