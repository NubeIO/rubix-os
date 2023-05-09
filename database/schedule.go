package database

import (
	"bytes"
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) GetSchedules() ([]*model.Schedule, error) {
	var scheduleModel []*model.Schedule
	query := d.DB.Find(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetSchedule(uuid string) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).First(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetSchedulesResult() ([]*model.Schedule, error) {
	var scheduleModel []*model.Schedule
	query := d.DB.Find(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, schedule := range scheduleModel {
		schedule.Schedule = nil
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetScheduleResult(uuid string) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).First(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	scheduleModel.Schedule = nil
	return scheduleModel, nil
}

func (d *GormDatabase) GetSchedulesByArgsTransaction(db *gorm.DB, args api.Args) ([]*model.Schedule, error) {
	var scheduleModel []*model.Schedule
	query := d.buildScheduleQueryTransaction(db, args)
	if err := query.Find(&scheduleModel).Error; err != nil {
		return nil, err
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetSchedulesByArgs(args api.Args) ([]*model.Schedule, error) {
	return d.GetSchedulesByArgsTransaction(d.DB, args)
}

func (d *GormDatabase) GetOneScheduleByArgsTransaction(db *gorm.DB, args api.Args) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.buildScheduleQueryTransaction(db, args)
	if err := query.First(&scheduleModel).Error; err != nil {
		return nil, err
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetOneScheduleByArgs(args api.Args) (*model.Schedule, error) {
	return d.GetOneScheduleByArgsTransaction(d.DB, args)
}

func (d *GormDatabase) CreateScheduleTransaction(db *gorm.DB, body *model.Schedule) (*model.Schedule, error) {
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Schedule)
	body.Name = nameIsNil(body.Name)
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	body.Schedule = validSchedule
	if body.GlobalUUID == "" {
		deviceInfo, err := deviceinfo.GetDeviceInfo()
		if err != nil {
			return nil, err
		}
		body.GlobalUUID = deviceInfo.GlobalUUID
	}
	if err = db.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateSchedule(body *model.Schedule) (*model.Schedule, error) {
	return d.CreateScheduleTransaction(d.DB, body)
}

func (d *GormDatabase) validateSchedule(schedule *model.Schedule) ([]byte, error) {
	if schedule.Schedule == nil {
		return nil, nil
	}
	scheduleDataModel := new(model.ScheduleData)
	err := json.Unmarshal(schedule.Schedule, &scheduleDataModel)
	if err != nil {
		return nil, err
	}
	validSchedule, err := json.Marshal(scheduleDataModel)
	if err != nil {
		return nil, err
	}
	return validSchedule, nil
}

func (d *GormDatabase) UpdateScheduleTransactionForAutoMapping(db *gorm.DB, uuid string, body *model.Schedule) (*model.Schedule, error) {
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	scheduleModel := model.Schedule{CommonUUID: model.CommonUUID{UUID: uuid}}
	body.Name = strings.TrimSpace(body.Name)
	body.Schedule = validSchedule
	if err := db.Model(&scheduleModel).Select("*").Updates(&body).Error; err != nil {
		return nil, err
	}
	return &scheduleModel, nil
}

func (d *GormDatabase) UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	if uuid == "" {
		return nil, errors.New("UpdateSchedule() requires a valid schedule UUID.")
	}
	body.UUID = uuid
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	body.Schedule = validSchedule
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	query := d.DB.Where("uuid = ?", uuid).First(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Name == "" {
		body.Name = scheduleModel.Name
	}
	if body.Enable == nil {
		if scheduleModel.Enable == nil {
			body.Enable = boolean.NewFalse()
		} else {
			body.Enable = scheduleModel.Enable
		}
	}
	if body.TimeZone == "" {
		body.TimeZone = scheduleModel.TimeZone
	}

	scheduleData := new(model.ScheduleData)
	_ = json.Unmarshal(body.Schedule, &scheduleData)
	_ = d.ScheduleWrite(uuid, scheduleData, false)

	if boolean.IsFalse(scheduleModel.CreatedFromAutoMapping) {
		query = d.DB.Model(&scheduleModel).Select("*").Omit("IsActive", "ActiveWeekly", "ActiveException", "ActiveEvent", "Payload", "PeriodStart", "PeriodStop", "NextStart", "NextStop", "PeriodStartString", "PeriodStopString", "NextStartString", "NextStopString", "CreatedAt").Updates(&body)
		// query = d.DB.Model(&scheduleModel).Updates(body)  // This line doesn't update properties to 0 (zero values).  Example is NextStart and NextStop
		if query.Error != nil {
			return nil, query.Error
		}
		d.UpdateProducerByProducerThingUUID(scheduleModel.UUID, scheduleModel.Name, nil, "", nil)
	}
	return scheduleModel, nil
}

func (d *GormDatabase) UpdateScheduleAllProps(uuid string, body *model.Schedule) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	if uuid == "" {
		return nil, errors.New("UpdateScheduleAllProps() requires a valid schedule UUID.")
	}
	body.UUID = uuid
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	body.Schedule = validSchedule
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	query := d.DB.Where("uuid = ?", uuid).First(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Name == "" {
		body.Name = scheduleModel.Name
	}

	scheduleData := new(model.ScheduleData)
	_ = json.Unmarshal(body.Schedule, &scheduleData)
	_ = d.ScheduleWrite(uuid, scheduleData, false)

	if boolean.IsFalse(scheduleModel.CreatedFromAutoMapping) {
		query = d.DB.Model(&scheduleModel).Select("*").Updates(&body)
		// query = d.DB.Model(&scheduleModel).Updates(body)  // This line doesn't update properties to 0 (zero values).  Example is NextStart and NextStop
		if query.Error != nil {
			return nil, query.Error
		}
		d.UpdateProducerByProducerThingUUID(scheduleModel.UUID, scheduleModel.Name, nil, "", nil)
	}
	return scheduleModel, nil
}

func (d *GormDatabase) ScheduleWrite(uuid string, body *model.ScheduleData, forceWrite bool) error {
	var scheduleModel *model.Schedule
	query := d.DB.Where("uuid = ?", uuid).First(&scheduleModel)
	if query.Error != nil {
		return query.Error
	}
	scheduleModuleScheduleData, err := json.Marshal(scheduleModel.Schedule)
	if err != nil {
		return err
	}

	scheduleData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	schedule := map[string]interface{}{}
	schedule["schedule"] = scheduleData

	isScheduleDataChange := !bytes.Equal(scheduleModuleScheduleData, scheduleData)
	if forceWrite || isScheduleDataChange {
		err = d.DB.Model(&scheduleModel).Updates(schedule).Error
		if err != nil {
			return err
		}
		d.ConsumersScheduleWrite(uuid, body)
		err = d.ProducersScheduleWrite(uuid, body)
		if err != nil {
			return err
		}
	}
	return d.ProducersScheduleWrite(uuid, body)
}

func (d *GormDatabase) DeleteSchedule(uuid string) (bool, error) {
	schedule, err := d.GetSchedule(uuid)
	if err != nil {
		return false, err
	}
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &schedule.UUID})
	if producers != nil {
		for _, producer := range producers {
			_, _ = d.DeleteProducer(producer.UUID)
		}
	}
	writers, _ := d.GetWriters(api.Args{WriterThingUUID: &schedule.UUID})
	if writers != nil {
		for _, writer := range writers {
			_, _ = d.DeleteWriter(writer.UUID)
		}
	}
	query := d.DB.Delete(&schedule)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncSchedules() error {
	schedules, err := d.GetSchedules()
	var firstErr error
	if err != nil {
		return err
	}
	uniqueAutoMappingFlowNetworkNames := GetUniqueAutoMappingScheduleFlowNetworkNames(schedules)
	for _, fnName := range uniqueAutoMappingFlowNetworkNames {
		err = d.CreateAutoMappingsSchedules(fnName, schedules)
		if err != nil {
			log.Error("Auto mapping error: ", err)
		}
	}

	if err != nil {
		return err
	}
	return firstErr
}

func GetUniqueAutoMappingScheduleFlowNetworkNames(schedules []*model.Schedule) []string {
	uniqueAutoMappingFlowNetworkNamesMap := make(map[string]struct{})
	var uniqueAutoMappingFlowNetworkNames []string

	for _, schedule := range schedules {
		if _, ok := uniqueAutoMappingFlowNetworkNamesMap[schedule.AutoMappingFlowNetworkName]; !ok {
			uniqueAutoMappingFlowNetworkNamesMap[schedule.AutoMappingFlowNetworkName] = struct{}{}
			uniqueAutoMappingFlowNetworkNames = append(uniqueAutoMappingFlowNetworkNames, schedule.AutoMappingFlowNetworkName)
		}
	}

	return uniqueAutoMappingFlowNetworkNames
}

func (d *GormDatabase) UpdateScheduleConnectionErrors(uuid string, schedule *model.Schedule) error {
	return UpdateScheduleConnectionErrorsTransaction(d.DB, uuid, schedule)
}

func UpdateScheduleConnectionErrorsTransaction(db *gorm.DB, uuid string, schedule *model.Schedule) error {
	return db.Model(&model.Schedule{}).
		Where("uuid = ?", uuid).
		Select("Connection", "ConnectionMessage").
		Updates(&schedule).
		Error
}
