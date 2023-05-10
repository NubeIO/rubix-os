package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) CreateAutoMappingSchedule(autoMapping *interfaces.AutoMapping) *interfaces.AutoMappingScheduleResponse {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	d.cleanScheduleAutoMappedModels(tx, autoMapping)
	d.clearScheduleStreamClonesAndConsumers(tx)

	var syncWriters []*interfaces.SyncWriter
	for _, amSchedule := range autoMapping.Schedules {
		if amSchedule.CreateSchedule {
			amRes := d.createScheduleAutoMapping(tx, amSchedule, autoMapping.FlowNetworkUUID, autoMapping.GlobalUUID)
			if amRes.HasError {
				tx.Rollback()
				return amRes
			}
			syncWriters = append(syncWriters, amRes.SyncWriters...)
		}
	}

	tx.Commit()

	return &interfaces.AutoMappingScheduleResponse{
		HasError:    false,
		SyncWriters: syncWriters,
	}
}

func (d *GormDatabase) cleanScheduleAutoMappedModels(tx *gorm.DB, autoMapping *interfaces.AutoMapping) {
	// delete those which is not deleted when we delete edge
	var edgeSchedules []string
	for _, schedule := range autoMapping.Schedules {
		edgeSchedules = append(edgeSchedules, schedule.UUID)
	}

	schedules, _ := d.GetSchedulesByArgsTransaction(tx, api.Args{GlobalUUID: &autoMapping.GlobalUUID})
	for _, schedule := range schedules {
		if boolean.IsTrue(schedule.CreatedFromAutoMapping) &&
			schedule.AutoMappingUUID != nil && !nstring.ContainsString(edgeSchedules, *schedule.AutoMappingUUID) {
			tx.Delete(&schedule)
		}
	}
}

func (d *GormDatabase) clearScheduleStreamClonesAndConsumers(tx *gorm.DB) {
	// delete those which is not deleted when we delete schedules
	tx.Where("created_from_auto_mapping IS TRUE AND IFNULL(auto_mapping_network_uuid,'') = '' AND "+
		"IFNULL(auto_mapping_device_uuid,'') = '' AND auto_mapping_schedule_uuid NOT IN (?)",
		tx.Where("created_from_auto_mapping IS TRUE").Model(&model.Schedule{}).Select("uuid")).
		Delete(&model.StreamClone{})
	tx.Where("created_from_auto_mapping IS TRUE AND producer_thing_class = ? AND producer_thing_uuid NOT IN (?)",
		model.ThingClass.Schedule, tx.Where("created_from_auto_mapping IS TRUE").Model(&model.Schedule{}).
			Select("auto_mapping_uuid")).
		Delete(&model.Consumer{})
}

func (d *GormDatabase) createScheduleAutoMapping(tx *gorm.DB, amSchedule *interfaces.AutoMappingSchedule, fnUUID, globalUUID string) *interfaces.AutoMappingScheduleResponse {
	amRes := &interfaces.AutoMappingScheduleResponse{
		ScheduleUUID: amSchedule.UUID,
		HasError:     true,
	}
	fnc, err := d.GetOneFlowNetworkCloneByArgsTransaction(tx, api.Args{SourceUUID: nstring.New(fnUUID)})
	if err != nil {
		amRes.Error = err.Error()
		return amRes
	}

	scheduleName := getAutoMappedScheduleName(fnc.Name, amSchedule.Name)

	schedule, err := d.GetOneScheduleByArgsTransaction(tx, api.Args{Name: nstring.New(scheduleName)})
	if schedule != nil {
		if schedule.GlobalUUID != globalUUID {
			amRes.Error = fmt.Sprintf("schedule.name %s already exists in fnc side with different global_uuid", schedule.Name)
			return amRes
		} else if boolean.IsFalse(schedule.CreatedFromAutoMapping) {
			amRes.Error = fmt.Sprintf("manually created network.schedule %s already exists in fnc side", schedule.Name)
			return amRes
		}
	}

	schedule, _ = d.GetOneScheduleByArgsTransaction(tx, api.Args{AutoMappingUUID: nstring.New(amSchedule.UUID), GlobalUUID: nstring.New(globalUUID)})
	if schedule == nil {
		if amSchedule.AutoMappingEnable {
			schedule = &model.Schedule{}
			schedule.Name = getTempAutoMappedName(scheduleName)
			d.setScheduleModel(fnc, amSchedule, schedule, globalUUID)
			schedule, err = d.CreateScheduleTransaction(tx, schedule)
			if err != nil {
				amRes.Error = err.Error()
				return amRes
			}
		} else {
			return &interfaces.AutoMappingScheduleResponse{
				HasError:    false,
				SyncWriters: nil,
			}
		}
	} else {
		schedule.Name = getTempAutoMappedName(scheduleName)
		d.setScheduleModel(fnc, amSchedule, schedule, globalUUID)
		schedule, err = d.UpdateScheduleTransactionForAutoMapping(tx, schedule.UUID, schedule)
		if err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	}

	streamClone, _ := GetOneStreamCloneByArgTransaction(tx, api.Args{SourceUUID: nstring.New(amSchedule.StreamUUID)})
	if streamClone == nil && amSchedule.AutoMappingEnable {
		streamClone = &model.StreamClone{}
		streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
		d.setScheduleStreamCloneModel(fnc, schedule, amSchedule, streamClone)
		if err = tx.Create(&streamClone).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	} else {
		d.setScheduleStreamCloneModel(fnc, schedule, amSchedule, streamClone)
		if err = tx.Model(&streamClone).Where("uuid = ?", streamClone.UUID).Updates(streamClone).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	}

	consumer, _ := GetOneConsumerByArgsTransaction(tx, api.Args{ProducerThingUUID: nstring.New(amSchedule.UUID)})
	if consumer == nil {
		consumer = &model.Consumer{}
		consumer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Consumer)
		d.setScheduleConsumerModel(amSchedule, streamClone.UUID, amSchedule.Name, consumer)
		if err = tx.Create(&consumer).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	} else {
		d.setScheduleConsumerModel(amSchedule, streamClone.UUID, amSchedule.Name, consumer)
		if err = tx.Model(&consumer).Where("uuid = ?", consumer.UUID).Updates(consumer).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	}

	writer, _ := GetOneWriterByArgsTransaction(tx, api.Args{WriterThingUUID: nstring.New(schedule.UUID)})
	if writer == nil {
		writer = &model.Writer{}
		writer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Writer)
		d.setScheduleWriterModel(amSchedule.Name, schedule.UUID, consumer.UUID, writer)
		if err = tx.Create(&writer).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	} else {
		d.setScheduleWriterModel(amSchedule.Name, schedule.UUID, consumer.UUID, writer)
		if err = tx.Model(&writer).Where("uuid = ?", writer.UUID).Updates(writer).Error; err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	}
	var syncWriters []*interfaces.SyncWriter
	syncWriters = append(syncWriters, &interfaces.SyncWriter{
		ProducerUUID:      amSchedule.ProducerUUID,
		WriterUUID:        writer.UUID,
		FlowFrameworkUUID: fnc.SourceUUID,
		UUID:              amSchedule.UUID,
		Name:              amSchedule.Name,
	})

	amRes_ := d.swapScheduleMapperNames(tx, amSchedule, fnc.Name, scheduleName)
	if amRes_ != nil {
		return amRes_
	}

	return &interfaces.AutoMappingScheduleResponse{
		HasError:    false,
		SyncWriters: syncWriters,
	}
}

func (d *GormDatabase) swapScheduleMapperNames(db *gorm.DB, amSchedule *interfaces.AutoMappingSchedule, fncName, scheduleName string) *interfaces.AutoMappingScheduleResponse {
	if err := db.Model(&model.StreamClone{}).
		Where("source_uuid = ?", amSchedule.StreamUUID).
		Update("name", getScheduleAutoMappedStreamName(fncName, amSchedule.Name)).
		Error; err != nil {
		return &interfaces.AutoMappingScheduleResponse{
			ScheduleUUID: amSchedule.UUID,
			HasError:     true,
			Error:        err.Error(),
		}
	}

	if err := db.Model(&model.Schedule{}).
		Where("auto_mapping_uuid = ?", amSchedule.UUID).
		Update("name", scheduleName).
		Error; err != nil {
		return &interfaces.AutoMappingScheduleResponse{
			ScheduleUUID: amSchedule.UUID,
			HasError:     true,
			Error:        err.Error(),
		}
	}

	if err := db.Model(&model.Consumer{}).
		Where("producer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", amSchedule.UUID).
		Update("name", amSchedule.Name).
		Error; err != nil {
		return &interfaces.AutoMappingScheduleResponse{
			ScheduleUUID: amSchedule.UUID,
			HasError:     true,
			Error:        err.Error(),
		}
	}

	writer := model.Writer{}
	if err := db.Model(&writer).
		Where("writer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", amSchedule.UUID).
		Update("writer_thing_name", amSchedule.Name).
		Error; err != nil {
		return &interfaces.AutoMappingScheduleResponse{
			ScheduleUUID: amSchedule.UUID,
			HasError:     true,
			Error:        err.Error(),
		}
	}
	return nil
}

func (d *GormDatabase) setScheduleModel(fnc *model.FlowNetworkClone, amSchedule *interfaces.AutoMappingSchedule, scheduleModel *model.Schedule, globalUUID string) {
	scheduleModel.Enable = &amSchedule.Enable
	scheduleModel.AutoMappingEnable = &amSchedule.AutoMappingEnable
	scheduleModel.ThingClass = "schedule"
	scheduleModel.ThingType = "schedule"
	scheduleModel.GlobalUUID = globalUUID
	scheduleModel.AutoMappingFlowNetworkName = fnc.Name
	scheduleModel.CreatedFromAutoMapping = boolean.NewTrue()
	scheduleModel.AutoMappingUUID = &amSchedule.UUID
	scheduleModel.TimeZone = amSchedule.TimeZone
	scheduleModel.EnablePayload = amSchedule.EnablePayload
	scheduleModel.MinPayload = amSchedule.MinPayload
	scheduleModel.MaxPayload = amSchedule.MaxPayload
	scheduleModel.DefaultPayload = amSchedule.DefaultPayload
	scheduleModel.Schedule = amSchedule.Schedule
}

func (d *GormDatabase) setScheduleStreamCloneModel(fnc *model.FlowNetworkClone, schedule *model.Schedule,
	amSchedule *interfaces.AutoMappingSchedule, streamClone *model.StreamClone) {
	streamClone.Name = getTempAutoMappedName(getScheduleAutoMappedStreamName(fnc.Name, amSchedule.Name))
	streamClone.Enable = boolean.New(amSchedule.Enable && amSchedule.AutoMappingEnable)
	streamClone.SourceUUID = amSchedule.StreamUUID
	streamClone.FlowNetworkCloneUUID = fnc.UUID
	streamClone.CreatedFromAutoMapping = boolean.NewTrue()
	streamClone.AutoMappingScheduleUUID = nstring.New(schedule.UUID)
}

func (d *GormDatabase) setScheduleConsumerModel(amSchedule *interfaces.AutoMappingSchedule, stcUUID, scheduleName string, consumerModel *model.Consumer) {
	consumerModel.Name = getTempAutoMappedName(scheduleName)
	consumerModel.Enable = boolean.New(amSchedule.Enable && amSchedule.AutoMappingEnable)
	consumerModel.StreamCloneUUID = stcUUID
	consumerModel.ProducerUUID = amSchedule.ProducerUUID
	consumerModel.ProducerThingName = scheduleName
	consumerModel.ProducerThingUUID = amSchedule.UUID
	consumerModel.ProducerThingClass = "schedule"
	consumerModel.CreatedFromAutoMapping = boolean.NewTrue()
}

func (d *GormDatabase) setScheduleWriterModel(scheduleName, scheduleUUID, consumerUUID string, writer *model.Writer) {
	writer.WriterThingName = scheduleName
	writer.WriterThingClass = "schedule"
	writer.WriterThingUUID = scheduleUUID
	writer.ConsumerUUID = consumerUUID
	writer.CreatedFromAutoMapping = boolean.NewTrue()
}
