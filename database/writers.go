package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/rest"
	"github.com/NubeDev/flow-framework/utils"
)

type Writer struct {
	*model.Writer
}

// GetWriters get all of them
func (d *GormDatabase) GetWriters() ([]*model.Writer, error) {
	var consumersModel []*model.Writer
	query := d.DB.Find(&consumersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumersModel, nil
}

// CreateWriter make it
func (d *GormDatabase) CreateWriter(body *model.Writer) (*model.Writer, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Writer)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetWriter get it
func (d *GormDatabase) GetWriter(uuid string) (*model.Writer, error) {
	var consumerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).First(&consumerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil
}

// GetWriterByThing get it by its thing uuid
func (d *GormDatabase) GetWriterByThing(producerThingUUID string) (*model.Writer, error) {
	var consumerModel *model.Writer
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&consumerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil
}

/*

1. update writer
2. write to writerClone
3. producer to decide if it's a valid cov
3. producer to write to history
4. return to consumer and update as required if it's a valid cov (update the writerConeUUID, this could be from another flow-framework instance)

WRITER
get http post:
update the writer and writerHistory

http post to writeClone:
try to write value to writeClone
update the writerClone history
then try and write to the producer, the producer will decide if it will accept the value. example a point write with point cov
if: the producer accepts then.
- update producer with the writerCloneUUID (this is, so we know who wrote the last value to the producer)
- update the consumer
- return message to user

else:
update the writerClone history

*/

func validType(t string, body *model.WriterBody) ([]byte, bool, error) {
	if t == model.CommonNaming.Point {
		var bk model.WriterBody
		if body.Action == model.CommonNaming.Read {
			return nil, true, nil
		} else  if body.Priority == bk.Priority {
			return nil, false, errors.New("error: invalid json on writerBody")
		}
		b, err := json.Marshal(body.Priority);if err != nil {
			return nil, false, errors.New("error: failed to marshal json on writeBody")
		}
		return b, true, err
	} else {
		return nil, false, errors.New("error: invalid data type on writerBody, ie type could be a point")
	}
}


func (d *GormDatabase) RemoteWriterWrite(uuid string, body *model.WriterBody, askRefresh bool) (*model.ProducerHistory, error) {
	writer, err := d.GetWriter(uuid); if err != nil {
		return nil, err
	}
	data, _,  err := validType(writer.WriterType, body);if err != nil {
		return nil, err
	}
	wc := new(model.WriterClone)
	wc.DataStore = data
	writer.DataStore = data
	consumer, err := d.GetConsumer(writer.ConsumerUUID); if err != nil {
		return nil, errors.New("error: on get consumer")
	}
	consumerUUID := consumer.UUID
	producerUUID := consumer.ProducerUUID
	writerCloneUUID := writer.WriteCloneUUID
	flow, err := d.GetFlowNetwork(body.FlowUUID); if err != nil {
		return nil, errors.New("error: invalid flow UUID")
	}

	// update producer
	_, err = rest.WriteClone(writerCloneUUID, flow, wc, true)
	if err != nil {
		return nil, errors.New("error: write new value to writerClone")
	}
	if producerUUID == "" {
		return nil, errors.New("error: producer uuid is none")
	}
	producerFeedback, err := rest.ProducerHistory(flow, producerUUID);if err != nil {
		return nil, errors.New("error: on get feedback from producer history")
	}
	// update the consumer based of the response from the producer
	if askRefresh {
		updateConsumer := new(model.Consumer)
		updateConsumer.DataStore = producerFeedback.DataStore
		updateConsumer.CurrentWriterCloneUUID = producerFeedback.CurrentWriterCloneUUID
		_, _ = d.UpdateConsumer(consumerUUID, updateConsumer)
		if err != nil {
			return nil, errors.New("error: on update consumer feedback")
		}
		return producerFeedback, err
	} else {
		return producerFeedback, err
	}
}



func (d *GormDatabase) RemoteWriterRead(uuid string, body *model.WriterBody) (*model.ProducerHistory, error) {
	writer, err := d.GetWriter(uuid); if err != nil {
		return nil, err
	}
	_, valid, err := validType(writer.WriterType, body);if err != nil || !valid {
		return nil, err
	}

	consumer, err := d.GetConsumer(writer.ConsumerUUID); if err != nil {
		return nil, errors.New("error: on get consumer")
	}
	consumerUUID := consumer.UUID
	producerUUID := consumer.ProducerUUID
	flow, err := d.GetFlowNetwork(body.FlowUUID); if err != nil {
		return nil, errors.New("error: invalid flow UUID")
	}

	if producerUUID == "" {
		return nil, errors.New("error: producer uuid is none")
	}

	producerFeedback, err := rest.ProducerHistory(flow, producerUUID);if err != nil {
		return nil, errors.New("error: on get feedback from producer history")
	}
	// update the consumer based of the response from the producer
	updateConsumer := new(model.Consumer)
	updateConsumer.DataStore = producerFeedback.DataStore
	updateConsumer.CurrentWriterCloneUUID = producerFeedback.CurrentWriterCloneUUID
	_, _ = d.UpdateConsumer(consumerUUID, updateConsumer)
	if err != nil {
		return nil, errors.New("error: on update consumer feedback")
	}
	return producerFeedback, err
}


// DeleteWriter deletes it
func (d *GormDatabase) DeleteWriter(uuid string) (bool, error) {
	var consumerModel *model.Writer
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

// UpdateWriter  update it
func (d *GormDatabase) UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error) {
	var consumerModel *model.Writer
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

// DropConsumersList delete all.
func (d *GormDatabase) DropConsumersList() (bool, error) {
	var consumerModel *model.Writer
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
