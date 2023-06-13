package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func (d *GormDatabase) GetPointHistories(args api.Args) ([]*model.PointHistory, error) {
	var historiesModel []*model.PointHistory
	query := d.buildPointHistoryQuery(args)
	query.Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

func (d *GormDatabase) GetPointHistoriesByPointUUID(pUuid string, args api.Args) ([]*model.PointHistory, int64, error) {
	var count int64
	var historiesModel []*model.PointHistory
	query := d.buildPointHistoryQuery(args)
	query = query.Where("point_uuid = ?", pUuid)
	query.Find(&historiesModel)
	query.Count(&count)
	return historiesModel, count, nil
}

func (d *GormDatabase) GetLatestPointHistoryByPointUUID(pUuid string) (*model.PointHistory, error) {
	var historyModel *model.PointHistory
	query := d.DB.Where("point_uuid = ? ", pUuid).Order("timestamp desc").First(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil
}

func (d *GormDatabase) GetPointHistoriesByPointUUIDs(pointUUIDs []string, args api.Args) ([]*model.PointHistory, error) {
	var historiesModel []*model.PointHistory
	query := d.buildPointHistoryQuery(args)
	if err := query.Where("point_uuid IN ?", pointUUIDs).Order("point_uuid").
		Find(&historiesModel).Error; err != nil {
		return nil, err
	}
	return historiesModel, nil
}

func (d *GormDatabase) GetPointHistoriesForSync(id string, timeStamp string) ([]*model.PointHistory, error) {
	var pointHistoriesModel []*model.PointHistory
	query := d.DB.Where("id = ?", id).Where("datetime(timestamp) = datetime(?)", timeStamp).
		Find(&pointHistoriesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(pointHistoriesModel) == 0 {
		id = "0"
	}
	pointHistories, err := d.GetPointHistories(api.Args{IdGt: nstring.New(id)})
	if err != nil {
		return nil, err
	}
	return pointHistories, nil
}

func (d *GormDatabase) CreatePointHistory(body *model.PointHistory) (*model.PointHistory, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateBulkPointHistory(histories []*model.PointHistory) (bool, error) {
	if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(histories, 1000).Error; err != nil {
		log.Error("Issue on creating bulk point histories")
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) DeletePointHistoriesByPointUUID(pUuid string, args api.Args) (bool, error) {
	var historyModel *model.PointHistory
	query := d.buildPointHistoryQuery(args)
	query = query.Where("point_uuid = ? ", pUuid).Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeletePointHistoriesBeforeTimestamp(ts string) (bool, error) {
	var historyModel *model.PointHistory
	query := d.DB.Where("timestamp < datetime(?)", ts)
	query.Delete(&historyModel)
	return d.deleteResponseBuilder(query)
}
