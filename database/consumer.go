package database

import (
	"errors"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"

	"github.com/NubeDev/flow-framework/utils"
)

type Consumers struct {
	*model.Consumer
}

// GetConsumers get all of them
func (d *GormDatabase) GetConsumers() ([]*model.Consumer, error) {
	var consumersModel []*model.Consumer
	query := d.DB.Preload("Writer").Find(&consumersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumersModel, nil
}

// CreateConsumer make it
func (d *GormDatabase) CreateConsumer(body *model.Consumer) (*model.Consumer, error) {
	_, err := d.GetStream(body.StreamUUID, false)
	if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the stream UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
	body.Name = nameIsNil(body.Name)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetConsumer get it
func (d *GormDatabase) GetConsumer(uuid string, withChildren bool) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	if withChildren {
		query := d.DB.Preload("Writer").Where("uuid = ? ", uuid).First(&consumerModel)
		if query.Error != nil {
			return nil, query.Error
		}
	} else {
		query := d.DB.Where("uuid = ? ", uuid).First(&consumerModel)
		if query.Error != nil {
			return nil, query.Error
		}
	}
	return consumerModel, nil
}

// DeleteConsumer deletes it
func (d *GormDatabase) DeleteConsumer(uuid string) (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&consumerModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// UpdateConsumer  update it
func (d *GormDatabase) UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ?", uuid).Find(&consumerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&consumerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil

}

// DropConsumers delete all.
func (d *GormDatabase) DropConsumers() (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("1 = 1").Delete(&consumerModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

/*
add new consumer auto add, writer and writer clone for a remote network and local
body
-- consumerStreamUUID
-- producerUUID

- get the streamUUID from the producerUUID
- from the consumerStreamUUID get the flowUUID
- first make sure that the producer device is online and the streamUUID is valid
- work out if the producer is local or remote
- add the new consumer, writer and writerClone
*/

func (d *GormDatabase) AddConsumerWizard(consumerStreamUUID, producerUUID string, consumerModel *model.Consumer) (*model.Consumer, error) {
	streamUUID := consumerStreamUUID
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	if producerUUID == "" {
		return nil, errors.New("error: no producer uuid provided")
	}
	if consumerStreamUUID == "" {
		return nil, errors.New("error: no stream uuid provided")
	}
	stream, flow, err := d.GetFlowUUID(streamUUID)
	if err != nil || stream.UUID == "" {
		return nil, errors.New("error: invalid stream UUID")
	}
	isRemote := flow.IsRemote

	var producer *model.Producer
	if isRemote {
		cli := client.NewSessionWithToken(flow.FlowToken, flow.FlowIP, flow.FlowPort)
		p, err := cli.GetProducer(producerUUID)
		if err != nil {
			return nil, errors.New("error: issue on get producer over rest client")
		}
		producer = p
	} else {
		p, err := d.GetProducer(producerUUID)
		if err != nil {
			return nil, errors.New("error: issue on get producer")
		}
		producer = p
	}

	if producer.UUID == "" {
		return nil, errors.New("error: no producer producer found with that UUID")
	}

	consumerModel.StreamUUID = stream.UUID
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ProducerThingName = producer.ProducerThingName
	consumerModel.ProducerThingUUID = producer.ProducerThingUUID
	consumerModel.ProducerThingClass = producer.ProducerThingClass
	consumerModel.ProducerThingType = producer.ProducerThingType

	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumer, err := d.CreateConsumer(consumerModel)
	if err != nil {
		return nil, errors.New("error: issue on create consumer")
	}
	// writer
	writerModel.ConsumerUUID = consumer.UUID
	writerModel.ConsumerThingUUID = consumerModel.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return nil, errors.New("error: issue on create writer")
	}
	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.WriterUUID = writer.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API

	if !isRemote {
		writerClone, err := d.CreateWriterClone(&writerCloneModel)
		if err != nil {
			return nil, errors.New("error: issue on create writer clone over rest")
		}
		//update writerCloneUUID to writer
		writerModel.CloneUUID = writerClone.UUID
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			return nil, errors.New("error: issue on update writer over rest")
		}
	} else {
		cli := client.NewSessionWithToken(flow.FlowToken, flow.FlowIP, flow.FlowPort)
		clone, err := cli.CreateWriterClone(writerCloneModel)
		if err != nil {
			return nil, errors.New("error: issue on create writer clone")
		}
		writerModel.CloneUUID = clone.UUID
		_, err = cli.EditWriter(writerModel.UUID, writerModel, false)
		if err != nil {
			return nil, errors.New("error: issue on update writer")
		}
	}
	return consumerModel, nil
}
