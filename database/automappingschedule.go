package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) CreateAutoMappingsSchedules(fnName string, schedules []*model.Schedule) error {
	if fnName == "" {
		return nil
	}

	d.clearSchedulesStreamsAndProducers()

	err := d.createSchedulesAutoMappings(fnName, schedules)
	if err != nil {
		return err
	}

	err = d.updateSchedulesConnectionInCloneSide(fnName)
	if err != nil {
		return err
	}

	return nil
}

func (d *GormDatabase) createSchedulesAutoMappings(fnName string, schedules []*model.Schedule) error {
	doAutoMapping := false
	var amSchedules []*interfaces.AutoMappingSchedule
	fn, fnError := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(fnName)})

	for _, schedule := range schedules {
		if boolean.IsTrue(schedule.CreatedFromAutoMapping) {
			continue
		}

		// we are sending extra schedules to make sure whether it's available or not in fn side
		if schedule.AutoMappingFlowNetworkName != fnName || boolean.IsFalse(schedule.AutoMappingEnable) {
			amSchedule := &interfaces.AutoMappingSchedule{
				Enable:            boolean.IsTrue(schedule.Enable),
				AutoMappingEnable: boolean.IsTrue(schedule.AutoMappingEnable),
				UUID:              schedule.UUID,
				Name:              schedule.Name,
				CreateSchedule:    false,
			}
			amSchedules = append(amSchedules, amSchedule)

			if boolean.IsFalse(schedule.AutoMappingEnable) {
				schedule.Connection = connection.Connected.String()
				schedule.ConnectionMessage = nstring.New(nstring.NotAvailable)
				_ = d.UpdateScheduleConnectionErrors(schedule.UUID, schedule)
			}
			continue
		}

		// if fnError has issue then return that just right away
		if fnError != nil {
			var msg string
			if schedule.AutoMappingFlowNetworkName == "" {
				msg = fmt.Sprintf("No flow-network has been selected for the enabled auto-mapping schedule.")
			} else {
				msg = fmt.Sprintf("The flow network with the name '%s' could not be found.", schedule.AutoMappingFlowNetworkName)
			}
			schedule.Connection = connection.Broken.String()
			schedule.ConnectionMessage = nstring.New(msg)
			_ = d.UpdateScheduleConnectionErrors(schedule.UUID, schedule)
			return errors.New(msg)
		} else {
			schedule.Connection = connection.Broken.String()
			schedule.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = d.UpdateScheduleConnectionErrors(schedule.UUID, schedule)
		}

		doAutoMapping = true // this is the case where it has auto_mapping creator with valid flow_network

		tx := d.DB.Begin()
		streamUUIDMap, err := createScheduleAutoMappingStreamsTransaction(tx, fn, schedule)
		if err != nil {
			tx.Rollback()
			log.Error(err)
			return err
		}
		tx.Commit()

		producerUUIDMap, err := d.createScheduleAutoMappingProducers(streamUUIDMap, schedule)
		if err != nil {
			amRes := interfaces.AutoMappingScheduleResponse{
				HasError:     true,
				ScheduleUUID: schedule.UUID,
			}
			updateScheduleCascadeConnectionError(d.DB, amRes)
			log.Error(err)
			return err
		}
		amSchedule := &interfaces.AutoMappingSchedule{
			Enable:            boolean.IsTrue(schedule.Enable),
			AutoMappingEnable: boolean.IsTrue(schedule.AutoMappingEnable),
			UUID:              schedule.UUID,
			Name:              schedule.Name,
			StreamUUID:        streamUUIDMap,
			ProducerUUID:      producerUUIDMap,
			CreateSchedule:    true,
		}
		amSchedules = append(amSchedules, amSchedule)
	}

	if !doAutoMapping {
		return nil
	}

	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	autoMapping := &interfaces.AutoMapping{
		GlobalUUID:      deviceInfo.GlobalUUID,
		FlowNetworkUUID: fn.UUID,
		Schedules:       amSchedules,
	}
	cli := client.NewFlowClientCliFromFN(fn)
	amRes := cli.CreateAutoMappingSchedule(autoMapping)
	if amRes.HasError {
		errMsg := fmt.Sprintf("Flow network clone side: %s", amRes.Error)
		log.Error(errMsg)
		updateScheduleCascadeConnectionError(d.DB, amRes)
	} else {
		for _, amSchedule := range amSchedules {
			if amSchedule.CreateSchedule { // just update its own schedule
				err := d.clearScheduleConnectionError(amSchedule)
				if err != nil {
					return err
				}
			}
		}

		scheduleUUID, err := d.createScheduleWriterClones(amRes.SyncWriters)
		if scheduleUUID != nil && err != nil {
			amRes.ScheduleUUID = *scheduleUUID
			updateScheduleCascadeConnectionError(d.DB, amRes)
		}
	}
	return nil
}

func createScheduleAutoMappingStreamsTransaction(tx *gorm.DB, flowNetwork *model.FlowNetwork, schedule *model.Schedule) (
	string, error) {
	streamName := getScheduleAutoMappedStreamName(flowNetwork.Name, schedule.Name)
	stream, _ := GetOneStreamByArgsTransaction(tx, api.Args{Name: nils.NewString(streamName)})
	if stream != nil {
		if boolean.IsFalse(stream.CreatedFromAutoMapping) {
			errMsg := fmt.Sprintf("manually created stream_name %s already exists", streamName)
			amRes := interfaces.AutoMappingScheduleResponse{
				ScheduleUUID: schedule.UUID,
				HasError:     true,
				Error:        errMsg,
			}
			updateScheduleCascadeConnectionError(tx, amRes)
			return "", errors.New(errMsg)
		}
	}
	stream, _ = GetOneStreamByArgsTransaction(tx, api.Args{AutoMappingScheduleUUID: nstring.New(schedule.UUID), WithFlowNetworks: true})
	if stream == nil {
		if boolean.IsTrue(schedule.AutoMappingEnable) { // create stream only when auto_mapping is enabled
			stream = &model.Stream{}
			stream.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Stream)
			setScheduleStreamModel(flowNetwork, schedule, stream)
			if err := tx.Create(&stream).Error; err != nil {
				schedule.Connection = connection.Broken.String()
				errMsg := fmt.Sprintf("create stream: %s", err.Error())
				amRes := interfaces.AutoMappingScheduleResponse{
					ScheduleUUID: schedule.UUID,
					HasError:     true,
					Error:        errMsg,
				}
				updateScheduleCascadeConnectionError(tx, amRes)
				return "", errors.New(errMsg)
			}
		}
	} else {
		schedule.Connection = connection.Connected.String()
		schedule.ConnectionMessage = nstring.New(nstring.NotAvailable)
		_ = UpdateScheduleConnectionErrorsTransaction(tx, schedule.UUID, schedule)
		if err := tx.Model(&stream).Association("FlowNetworks").Replace([]*model.FlowNetwork{flowNetwork}); err != nil {
			errMsg := fmt.Sprintf("update flow_networks on stream: %s", err.Error())
			amRes := interfaces.AutoMappingScheduleResponse{
				ScheduleUUID: schedule.UUID,
				HasError:     true,
				Error:        errMsg,
			}
			updateScheduleCascadeConnectionError(tx, amRes)
			return "", err
		}
		setScheduleStreamModel(flowNetwork, schedule, stream)
		if err := tx.Model(&stream).Updates(&stream).Error; err != nil {
			errMsg := fmt.Sprintf("update stream: %s", err.Error())
			amRes := interfaces.AutoMappingScheduleResponse{
				ScheduleUUID: schedule.UUID,
				HasError:     true,
				Error:        errMsg,
			}
			updateScheduleCascadeConnectionError(tx, amRes)
			return "", err
		}
	}

	// swap back the names

	if err := tx.Model(&model.Stream{}).
		Where("auto_mapping_schedule_uuid = ? AND created_from_auto_mapping IS TRUE", schedule.UUID).
		Update("name", streamName).
		Error; err != nil {
		errMsg := fmt.Sprintf("update stream: %s", err.Error())
		amRes := interfaces.AutoMappingScheduleResponse{
			ScheduleUUID: schedule.UUID,
			HasError:     true,
			Error:        errMsg,
		}
		updateScheduleCascadeConnectionError(tx, amRes)
		return "", err
	}

	return stream.UUID, nil
}

func (d *GormDatabase) createScheduleAutoMappingProducers(streamUUID string, schedule *model.Schedule) (string, error) {
	tx := d.DB.Begin()

	producer, _ := d.GetOneProducerByArgs(api.Args{StreamUUID: nils.NewString(streamUUID), ProducerThingUUID: nils.NewString(schedule.UUID)})
	if producer == nil {
		if boolean.IsTrue(schedule.AutoMappingEnable) { // create stream only when auto_mapping is enabled
			producer = &model.Producer{}
			producer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
			d.setScheduleProducerModel(streamUUID, schedule, producer)
			if err := tx.Create(&producer).Error; err != nil {
				tx.Rollback()
				return "", err
			}
		}
	} else {
		d.setScheduleProducerModel(streamUUID, schedule, producer)
		if err := tx.Save(producer).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	// swap back the names
	if err := tx.Model(&model.Producer{}).
		Where("producer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", schedule.UUID).
		Update("name", schedule.Name).
		Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return producer.UUID, nil
}

func (d *GormDatabase) clearSchedulesStreamsAndProducers() {
	// delete those which is not deleted when we delete stream
	d.DB.Where("created_from_auto_mapping IS TRUE AND IFNULL(auto_mapping_network_uuid,'') = '' AND "+
		"IFNULL(auto_mapping_device_uuid,'') = '' AND auto_mapping_schedule_uuid NOT IN (?)",
		d.DB.Model(&model.Schedule{}).Select("uuid")).Delete(&model.Stream{})
	d.DB.Where("created_from_auto_mapping IS TRUE AND producer_thing_class = ? AND producer_thing_uuid NOT IN (?)",
		model.ThingClass.Schedule, d.DB.Model(&model.Schedule{}).Select("uuid")).Delete(&model.Producer{})
}

func (d *GormDatabase) updateSchedulesConnectionInCloneSide(fnName string) error {
	schedules, err := d.GetSchedules()
	if err != nil {
		return err
	}
	for _, schedule := range schedules {
		if boolean.IsTrue(schedule.CreatedFromAutoMapping) && schedule.AutoMappingFlowNetworkName == fnName {
			err = d.updateScheduleConnectionInCloneSide(schedule)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *GormDatabase) updateScheduleConnectionInCloneSide(schedule *model.Schedule) error {
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: &schedule.AutoMappingFlowNetworkName})
	if err != nil {
		schedule.Connection = connection.Connected.String()
		schedule.ConnectionMessage = nstring.New(err.Error())
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
		return err
	}

	cli := client.NewFlowClientCliFromFNC(fnc)

	sch, connectionErr, _ := cli.GetScheduleV2(*schedule.AutoMappingUUID)
	if connectionErr != nil {
		schedule.Connection = connection.Broken.String()
		schedule.ConnectionMessage = nstring.New("connection error: " + connectionErr.Error())
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
	} else if sch == nil {
		schedule.Connection = connection.Broken.String()
		schedule.ConnectionMessage = nstring.New("The schedule creator has already been deleted. Delete manually if needed.")
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
	} else if boolean.IsFalse(sch.AutoMappingEnable) && boolean.IsTrue(sch.CreatedFromAutoMapping) {
		// here we use 'sch' instead 'schedule' because schedule wouldn't get updated in auto-mapped disabled case
		schedule.Connection = connection.Broken.String()
		msg := fmt.Sprintf("The auto-mapping feature for the schedule creator '%s' is currently disabled. Delete manually if needed.", sch.Name)
		schedule.ConnectionMessage = nstring.New(msg)
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
	} else if sch.AutoMappingFlowNetworkName != sch.AutoMappingFlowNetworkName {
		schedule.Connection = connection.Broken.String()
		var msg string
		if sch.AutoMappingFlowNetworkName != "" {
			msg = fmt.Sprintf("The schedule creator '%s' is attached to a different flow network named '%s'. Delete manually if needed.", sch.Name, sch.AutoMappingFlowNetworkName)
		} else {
			msg = fmt.Sprintf("The schedule creator '%s' isn't attached with any flow network. Delete manually if needed.", sch.Name)
		}
		schedule.ConnectionMessage = nstring.New(msg)
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
	} else {
		schedule.Connection = connection.Connected.String()
		schedule.ConnectionMessage = nstring.New(nstring.NotAvailable)
		_ = UpdateScheduleConnectionErrorsTransaction(d.DB, schedule.UUID, schedule)
	}
	tx := d.DB.Begin()
	tx.Commit()
	return nil

}

func (d *GormDatabase) createScheduleWriterClones(syncWriters []*interfaces.SyncWriter) (*string, error) {
	tx := d.DB.Begin()
	for _, syncWriter := range syncWriters {
		// it will restrict duplicate creation of writer_clone
		wc, _ := d.GetOneWriterCloneByArgs(api.Args{ProducerUUID: &syncWriter.ProducerUUID, CreatedFromAutoMapping: boolean.NewTrue()})
		if wc == nil {
			wc = &model.WriterClone{}
			wc.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setScheduleWriterCloneModel(syncWriter, wc)
			if err := tx.Create(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.UUID, err
			}
		} else {
			d.setScheduleWriterCloneModel(syncWriter, wc)
			if err := tx.Model(&wc).Where("uuid = ?", wc.UUID).Updates(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.UUID, err
			}
		}
	}
	tx.Commit()
	return nil, nil
}

func (d *GormDatabase) clearScheduleConnectionError(amSchedule *interfaces.AutoMappingSchedule) error {
	tx := d.DB.Begin()
	scheduleModel := model.Schedule{
		Connection:        connection.Connected.String(),
		ConnectionMessage: nstring.New(nstring.NotAvailable),
	}

	err := UpdateScheduleConnectionErrorsTransaction(tx, amSchedule.UUID, &scheduleModel)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func updateScheduleCascadeConnectionError(tx *gorm.DB, amsRes interfaces.AutoMappingScheduleResponse) {
	scheduleModel := model.Schedule{}

	connection_ := connection.Connected.String()
	if amsRes.HasError {
		connection_ = connection.Broken.String()
	}

	scheduleModel.Connection = connection_
	scheduleModel.ConnectionMessage = &amsRes.Error
	_ = UpdateScheduleConnectionErrorsTransaction(tx, amsRes.ScheduleUUID, &scheduleModel)
}

func setScheduleStreamModel(flowNetwork *model.FlowNetwork, schedule *model.Schedule, streamModel *model.Stream) {
	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	streamModel.Name = getTempAutoMappedName(getScheduleAutoMappedStreamName(flowNetwork.Name, schedule.Name))
	streamModel.Enable = boolean.New(boolean.IsTrue(schedule.Enable) && boolean.IsTrue(schedule.AutoMappingEnable))
	streamModel.CreatedFromAutoMapping = boolean.NewTrue()
	streamModel.AutoMappingScheduleUUID = nstring.New(schedule.UUID)
}

func (d *GormDatabase) setScheduleProducerModel(streamUUID string, schedule *model.Schedule, producerModel *model.Producer) {
	producerModel.Name = getTempAutoMappedName(schedule.Name)
	producerModel.Enable = boolean.New(boolean.IsTrue(schedule.Enable) && boolean.IsTrue(schedule.AutoMappingEnable))
	producerModel.StreamUUID = streamUUID
	producerModel.ProducerThingUUID = schedule.UUID
	producerModel.ProducerThingName = schedule.Name
	producerModel.ProducerThingClass = "schedule"
	producerModel.ProducerApplication = "mapping"
	producerModel.CreatedFromAutoMapping = boolean.NewTrue()
}

func (d *GormDatabase) setScheduleWriterCloneModel(syncWriter *interfaces.SyncWriter, writerClone *model.WriterClone) {
	writerClone.WriterThingName = syncWriter.Name
	writerClone.WriterThingClass = "schedule"
	writerClone.FlowFrameworkUUID = syncWriter.FlowFrameworkUUID
	writerClone.WriterThingUUID = syncWriter.UUID
	writerClone.ProducerUUID = syncWriter.ProducerUUID
	writerClone.SourceUUID = syncWriter.WriterUUID
	writerClone.CreatedFromAutoMapping = boolean.NewTrue()
}
