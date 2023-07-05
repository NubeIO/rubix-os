package database

import (
	"bytes"
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/pkg/errors"
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

func (d *GormDatabase) GetSchedulesByArgs(args argspkg.Args) ([]*model.Schedule, error) {
	var scheduleModel []*model.Schedule
	query := d.buildScheduleQuery(args)
	if err := query.Find(&scheduleModel).Error; err != nil {
		return nil, err
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetOneScheduleByArgs(args argspkg.Args) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.buildScheduleQuery(args)
	if err := query.First(&scheduleModel).Error; err != nil {
		return nil, err
	}
	return scheduleModel, nil
}

func (d *GormDatabase) CreateSchedule(body *model.Schedule) (*model.Schedule, error) {
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
		globalUUID, err := d.getGlobalUUID()
		if err != nil {
			return nil, err
		}
		body.GlobalUUID = globalUUID
	}
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
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

	query = d.DB.Model(&scheduleModel).Select("*").Omit("IsActive", "ActiveWeekly", "ActiveException", "ActiveEvent", "Payload", "PeriodStart", "PeriodStop", "NextStart", "NextStop", "PeriodStartString", "PeriodStopString", "NextStartString", "NextStopString", "CreatedAt").Updates(&body)
	// query = d.DB.Model(&scheduleModel).Updates(body)  // This line doesn't update properties to 0 (zero values).  Example is NextStart and NextStop
	if query.Error != nil {
		return nil, query.Error
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

	// don't restrict to update it, coz this doesn't get called from the API
	query = d.DB.Model(&scheduleModel).Select("*").Updates(&body)
	// query = d.DB.Model(&scheduleModel).Updates(body)  // This line doesn't update properties to 0 (zero values).  Example is NextStart and NextStop
	if query.Error != nil {
		return nil, query.Error
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
	}
	return nil
}

func (d *GormDatabase) DeleteSchedule(uuid string) (bool, error) {
	schedule, err := d.GetSchedule(uuid)
	if err != nil {
		return false, err
	}
	query := d.DB.Delete(&schedule)
	return d.deleteResponseBuilder(query)
}
