package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/datatypes"
)

type Writer struct {
	*model.Writer
}

func (d *GormDatabase) GetWriters(args api.Args) ([]*model.Writer, error) {
	var writers []*model.Writer
	query := d.buildWriterQuery(args)
	if err := query.Find(&writers).Error; err != nil {
		return nil, err
	}
	return writers, nil
}

func (d *GormDatabase) CreateWriter(body *model.Writer) (*model.Writer, error) {
	name := ""
	switch body.WriterThingClass {
	case model.ThingClass.Point:
		point, err := d.GetPoint(body.WriterThingUUID, api.Args{})
		if err != nil {
			return nil, errors.New("point not found, please supply a valid point writer_thing_uuid")
		}
		name = point.Name
	case model.ThingClass.Schedule:
		schedule, err := d.GetSchedule(body.WriterThingUUID)
		if err != nil {
			return nil, errors.New("schedule not found, please supply a valid point writer_thing_uuid")
		}
		name = schedule.Name
	default:
		return nil, errors.New("we are not supporting writer_thing_uuid other than point for now")
	}

	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Writer)
	body.WriterThingName = name
	body.DataStore = nil
	body.SyncUUID, _ = nuuid.MakeUUID()
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	err := d.syncAfterCreateUpdateWriter(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) GetWriter(uuid string) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

func (d *GormDatabase) GetWriterByThing(producerThingUUID string) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

func (d *GormDatabase) GetOneWriterByArgs(args api.Args) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.buildWriterQuery(args)
	if err := query.First(&writerModel).Error; err != nil {
		return nil, err
	}
	return writerModel, nil
}

func (d *GormDatabase) DeleteWriter(uuid string) (bool, error) {
	aType := api.ArgsType
	writer, err := d.GetWriter(uuid)
	if err != nil {
		return false, err
	}
	consumer, _ := d.GetConsumer(writer.ConsumerUUID, api.Args{})
	streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
	fnc, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCliFromFNC(fnc)
	url := urls.SingularUrlByArg(urls.WriterCloneUrl, aType.SourceUUID, writer.UUID)
	_ = cli.DeleteQuery(url)
	query := d.DB.Delete(&writer)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error) {
	writerModel, err := d.updateWriterWithoutSync(uuid, body)
	if err != nil {
		return nil, err
	}
	err = d.syncAfterCreateUpdateWriter(writerModel)
	if err != nil {
		return nil, err
	}
	return writerModel, nil
}

func (d *GormDatabase) updateWriterWithoutSync(uuid string, body *model.Writer) (*model.Writer, error) {
	var writerModel *model.Writer
	body.DataStore = nil
	query := d.DB.Where("uuid = ?", uuid).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&writerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

func (d *GormDatabase) WriterAction(uuid string, body *model.WriterBody) *model.WriterActionOutput {
	writer, err := d.GetWriter(uuid)
	output := &model.WriterActionOutput{IsError: false}
	output.UUID = uuid
	output.Action = body.Action
	if err != nil {
		output.IsError = true
		output.Message = nstring.NewStringAddress(err.Error())
		return output
	}
	consumer, _ := d.GetConsumer(writer.ConsumerUUID, api.Args{})
	streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
	fnc, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCliFromFNC(fnc)
	if body.Action == model.CommonNaming.Sync {
		err = cli.SyncWriterReadAction(uuid)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
		}
		dataStore, presentValue, err := d.getDataStoreAndPresentValues(writer)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
		}
		output.DataStore = dataStore
		output.PresentValue = presentValue
		return output
	} else if body.Action == model.CommonNaming.Read {
		dataStore, presentValue, err := d.getDataStoreAndPresentValues(writer)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
		}
		output.DataStore = dataStore
		output.PresentValue = presentValue
		return output
	} else {
		bytes, err := d.validateWriterWriteBody(writer.WriterThingClass, body)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
			return output
		}
		writer.DataStore = bytes
		d.DB.Model(&writer).Updates(writer)
		syncWriterAction := model.SyncWriterAction{
			Priority: body.Priority,
			Schedule: body.Schedule,
		}
		err = cli.SyncWriterWriteAction(uuid, &syncWriterAction)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
			return output
		}
		dataStore, presentValue, err := d.getDataStoreAndPresentValues(writer)
		if err != nil {
			output.IsError = true
			output.Message = nstring.NewStringAddress(err.Error())
		}
		output.DataStore = dataStore
		output.PresentValue = presentValue
		return output
	}
}

func (d *GormDatabase) WriterBulkAction(body []*model.WriterBulkBody) []*model.WriterActionOutput {
	arr := make([]*model.WriterActionOutput, len(body))
	for index, singleWriterBulkBody := range body {
		writerBody := &model.WriterBody{}
		writerBodyBytes, _ := json.Marshal(singleWriterBulkBody)
		_ = json.Unmarshal(writerBodyBytes, &writerBody)
		out := d.WriterAction(singleWriterBulkBody.WriterUUID, writerBody)
		arr[index] = out
	}
	return arr
}

func (d *GormDatabase) validateWriterWriteBody(thingClass string, body *model.WriterBody) ([]byte, error) {
	if thingClass == model.ThingClass.Point {
		if body.Priority == nil {
			return nil, errors.New("error: invalid json on writer body")
		}
		b, err := json.Marshal(body.Priority)
		if err != nil {
			return nil, errors.New("error: failed to marshal priority on write body")
		}
		return b, err
	} else {
		if body.Schedule == nil {
			return nil, errors.New("error: invalid json on writer body")
		}
		b, err := json.Marshal(body.Schedule)
		if err != nil {
			return nil, errors.New("error: failed to marshal schedule on write body")
		}
		return b, err
	}
}

func (d *GormDatabase) syncAfterCreateUpdateWriter(body *model.Writer) error {
	consumer, _ := d.GetConsumer(body.ConsumerUUID, api.Args{})
	streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
	fnc, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCliFromFNC(fnc)
	syncWriterBody := model.SyncWriter{
		ProducerUUID:      consumer.ProducerUUID,
		WriterUUID:        body.UUID,
		FlowFrameworkUUID: fnc.SourceUUID,
	}
	_, err := cli.SyncWriter(&syncWriterBody)
	return err
}

func (d *GormDatabase) getDataStoreAndPresentValues(writer *model.Writer) (*datatypes.JSON, *float64, error) {
	if writer.WriterThingClass == model.ThingClass.Point {
		point, err := d.GetPoint(writer.WriterThingUUID, api.Args{WithPriority: true})
		if err != nil {
			return nil, nil, err
		}
		priorityBytes, _ := json.Marshal(point.Priority)
		priorities := &model.Priorities{}
		_ = json.Unmarshal(priorityBytes, priorities)
		prioritiesBytes, _ := json.Marshal(priorities)
		prioritiesJSON := (datatypes.JSON)(prioritiesBytes)
		return &prioritiesJSON, point.PresentValue, nil
	} else {
		schedule, err := d.GetSchedule(writer.WriterThingUUID)
		if err != nil {
			return nil, nil, err
		}
		scheduleBytes, _ := json.Marshal(schedule.Schedule)
		scheduleJSON := (datatypes.JSON)(scheduleBytes)
		return &scheduleJSON, nil, nil
	}
}
