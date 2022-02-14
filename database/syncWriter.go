package database

import (
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

func (d *GormDatabase) SyncWriterCOV(body *model.SyncWriterCOV) error {
	writer, err := d.GetWriter(body.WriterUUID)
	if err != nil {
		return err
	}
	if writer.WriterThingClass == model.ThingClass.Point {
		err = d.updatePointFromCOV(writer.WriterThingUUID, body)
	}
	return err
}

func (d *GormDatabase) updatePointFromCOV(writerThingUUID string, body *model.SyncWriterCOV) error {
	var pointModel *model.Point
	priorityMap, _, _, _ := d.parsePriority(body.Priority)
	query := d.DB.Where("uuid = ?", writerThingUUID).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return query.Error
	}
	d.DB.Model(&pointModel.Priority).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMap)
	pointModel.OriginalValue = body.OriginalValue
	pointModel.PresentValue = body.PresentValue
	pointModel.CurrentPriority = body.CurrentPriority
	d.DB.Model(&pointModel).Updates(pointModel)
	return nil
}
