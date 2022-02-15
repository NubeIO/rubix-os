package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	writerClone := model.WriterClone{}
	producer, err := d.GetProducer(body.ProducerUUID, api.Args{})
	if err != nil {
		return nil, errors.New("producer does not exist")
	}
	writerClone.WriterThingName = producer.ProducerThingName
	writerClone.WriterThingUUID = producer.ProducerThingUUID
	writerClone.WriterThingClass = producer.ProducerThingClass
	writerClone.WriterThingType = producer.ProducerThingType
	writerClone.ProducerUUID = body.ProducerUUID
	writerClone.SourceUUID = body.WriterUUID
	writerClone.FlowFrameworkUUID = body.FlowFrameworkUUID
	var writerCloneModel []*model.WriterClone
	if err = d.DB.Where("source_uuid = ? ", writerClone.SourceUUID).Find(&writerCloneModel).Error; err != nil {
		return nil, err
	}
	if len(writerCloneModel) == 0 {
		writerClone.UUID = utils.MakeTopicUUID(model.CommonNaming.WriterClone)
		if err = d.DB.Create(&writerClone).Error; err != nil {
			return nil, err
		}
	} else {
		writerClone.UUID = writerCloneModel[0].UUID
		if err = d.DB.Model(&writerCloneModel[0]).Updates(writerClone).Error; err != nil {
			return nil, err
		}
	}
	return &writerClone, nil
}

func (d *GormDatabase) SyncCOV(body *model.SyncCOV) error {
	writer, err := d.GetWriter(body.WriterUUID)
	if err != nil {
		return err
	}
	if writer.WriterThingClass == model.ThingClass.Point {
		err = d.updatePointFromCOV(writer.WriterThingUUID, body)
	}
	return err
}

func (d *GormDatabase) SyncWriterAction(body *model.SyncWriterAction) error {
	writerClone, err := d.GetOneWriterCloneByArgs(api.Args{SourceUUID: &body.WriterUUID})
	if err != nil {
		return err
	}
	if writerClone.WriterThingClass == model.ThingClass.Point {
		data, _ := json.Marshal(body.Priority)
		writerCloneBody := model.WriterClone{CommonWriter: model.CommonWriter{DataStore: data}}
		err = d.UpdateWriterClone(writerClone, &writerCloneBody)
		if err != nil {
			return nil
		}
		point := model.Point{Priority: body.Priority}
		_, _ = d.PointWrite(writerClone.WriterThingUUID, &point, true)
	} else if writerClone.WriterThingClass == model.ThingClass.Schedule {
		data, _ := json.Marshal(body.Schedule)
		writerCloneBody := model.WriterClone{CommonWriter: model.CommonWriter{DataStore: data}}
		err = d.UpdateWriterClone(writerClone, &writerCloneBody)
		if err != nil {
			return nil
		}
		_ = d.ScheduleWrite(writerClone.WriterThingUUID, body.Schedule)
	} else {
		return errors.New("no match writer thing class")
	}
	return err
}

func (d *GormDatabase) updatePointFromCOV(pointUUID string, body *model.SyncCOV) error {
	pointModel := model.Point{
		CommonUUID: model.CommonUUID{UUID: pointUUID},
		Priority:   body.Priority,
	}
	_, err := d.PointWrite(pointUUID, &pointModel, false)
	return err
}
