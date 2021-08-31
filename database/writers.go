package database

import (
	"errors"
	"fmt"
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
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetWriter get it
func (d *GormDatabase) GetWriter(uuid string) (*model.Writer, error) {
	var consumerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).First(&consumerModel); if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil
}

// GetWriterByThing get it by its thing uuid
func (d *GormDatabase) GetWriterByThing(producerThingUUID string) (*model.Writer, error) {
	var consumerModel *model.Writer
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&consumerModel); if query.Error != nil {
		return nil, query.Error
	}
	return consumerModel, nil
}


//point
// - presentValue

// consumer needs readonly of producer presentValue
// - presentValue

// writer needs to write to the writerClone
// - writeValue

//producer
// - presentValue
// - SLWriteUUID //writer UUID

//writerClone (consumerUUID)
// - writeValue

//producerHist
// - presentValue
// - SLWriteUUID //writer UUID

/*

WRITER
pass in writer UUID as what starts everything
get consumer uuid
get producer uuid
get writerCloneUUID uuid
get flow network uuid
get consumer uuid
update writeValue to writer

HTTP POST to update writerClone
IF 200 then HTTP PRODUCER to update Consumer
BROADCAST to writerClone
try and update writeValue to point
if successful then update the WriterClone writeValue

PRODUCER-BROADCAST
if successful then update the try and update consumer presentValue
if successful then the producer is to broadcast up update event

CONSUMER
then last step is the consumer would update its presentValue

CONSUMER-BROADCAST
to Writers to handle as needed
*/

func (d *GormDatabase) RemoteWriterAction(uuid string, body *model.Writer, write bool) (*model.WriterClone, error) {
	var wm *model.Writer
	writer := d.DB.Where("uuid = ? ", uuid).First(&wm); if writer.Error != nil {
		return nil, writer.Error
	}
	if wm == nil {
		return nil, nil
	}
	_, err := d.UpdateWriter(uuid, body);if err != nil {
		return nil, errors.New("error: on update consumer feedback")
	}
	var cm *model.Consumer
	consumer := d.DB.Where("uuid = ? ", wm.ConsumerUUID).First(&cm); if consumer.Error != nil {
		return nil, consumer.Error
	}
	consumerUUID := cm.UUID
	streamUUID := cm.StreamUUID
	producerUUID := cm.ProducerUUID
	writerCloneUUID := wm.WriteCloneUUID
	var s *model.Stream
	stream := d.DB.Where("uuid = ? ", streamUUID).First(&s); if consumer.Error != nil {
		return nil, stream.Error
	}
	streamListUUID := s.StreamListUUID
	var fn *model.FlowNetwork
	flow := d.DB.Where("stream_list_uuid = ? ", streamListUUID).First(&fn); if consumer.Error != nil {
		return nil, flow.Error
	}
	wc := new(model.WriterClone)
	wc.WriteValue = body.WriteValue
	// update writer clone
	update, err := rest.WriteClone(writerCloneUUID, fn, wc, write)
	if err != nil {
		return nil, errors.New("error: write new value to writerClone")
	}
	if write { //get feedback from producer
		producerFeedback, err := rest.ProducerRead(fn, producerUUID);if err != nil {
			return nil, errors.New("error: on get feedback from producer")
		}
		// update the consumer based of the response from the producer
		//var updateConsumer *model.Consumer
		updateConsumer:= new(model.Consumer)
		updateConsumer.PresentValue =producerFeedback.PresentValue
		pro, _ := d.UpdateConsumer(consumerUUID, updateConsumer);if err != nil {
			return nil, errors.New("error: on update consumer feedback")
		}
		newWriter:= new(model.Writer)
		newWriter.PresentValue = pro.PresentValue
		// now update the writer
		_, _ = d.UpdateWriter(uuid, newWriter);if err != nil {
			return nil, errors.New("error: on update consumer feedback")
		}
	}
	return update, err
}


// WriterActionPoint get its value or write
func (d *GormDatabase) WriterActionPoint(slUUID string, pointBody *model.Point, write bool) (*model.Producer, error) {
	var slm *model.Writer
	sl := d.DB.Where("uuid = ? ", slUUID).First(&slm); if sl.Error != nil {
		return nil, sl.Error
	}
	if slm == nil {
		return nil, nil
	}

	var cm *model.Consumer
	consumer := d.DB.Where("uuid = ? ", slm.ConsumerUUID).First(&cm); if consumer.Error != nil {
		return nil, consumer.Error
	}
	subType := cm.ConsumerType
	consumerUUID := cm.UUID
	streamUUID := cm.StreamUUID
	producerUUID := cm.ProducerUUID
	writeV := pointBody.WriteValue
	pointUUID := cm.ProducerThingUUID

	var s *model.Stream
	stream := d.DB.Where("uuid = ? ", streamUUID).First(&s); if consumer.Error != nil {
		return nil, stream.Error
	}
	streamListUUID := s.StreamListUUID
	var fn *model.FlowNetwork
	flow := d.DB.Where("stream_list_uuid = ? ", streamListUUID).First(&fn); if consumer.Error != nil {
		return nil, flow.Error
	}
	flowUUID := fn.UUID
	isRemote := fn.IsRemote
	fmt.Println("subType", subType, "consumerUUID", consumerUUID, "streamUUID", streamUUID, "producerUUID", producerUUID,"flowUUID", flowUUID, "isRemote", isRemote, writeV, write)
	if !isRemote { // local
		pm := new(model.Producer)
		query := d.DB.Where("uuid = ?", producerUUID).Find(&pm);if query.Error != nil {
			return nil, query.Error
		}
		if query == nil {
			return nil, nil
		}
		if write { //write new value to producer
			//pm.WriteValue = body.WriteValue
			//pm.PresentValue = body.WriteValue
			//query = d.DB.Model(&pm).Updates(pm);if query.Error != nil {
			//	return nil, query.Error
			//}
			//ph := new(model.ProducerHistory)
			//ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			//ph.ProducerUUID = producerUUID
			//ph.PresentValue = pm.PresentValue
			//query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
			//	return nil, query.Error
			//}
			return pm, nil
		} else {
			//d.DB.Where("uuid = ? ", uuid).First(&pm); if query.Error != nil {
			//	return nil, query.Error
			//}
			return pm, nil
		}
	} else {
		if write {
			point, err := rest.EventRESTPoint(pointUUID, fn, pointBody, write)
			if err != nil {
				return nil, err
			}
			pm := new(model.Producer)
			pm.UUID = producerUUID
			pm.PresentValue = point.WriteValue
			query := d.DB.Model(&pm).Updates(pm);if query.Error != nil {
				return nil, query.Error
			}
			ph := new(model.ProducerHistory)
			ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			ph.ProducerUUID = producerUUID
			ph.PresentValue = pm.PresentValue
			query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
				return nil, query.Error
			}
			return pm, err
		} else {
			point, err := rest.EventRESTPoint(pointUUID, fn, pointBody, write)
			if err != nil {
				return nil, err
			}
			pm := new(model.Producer)
			pm.UUID = producerUUID
			pm.PresentValue = point.WriteValue
			return pm, err
		}
	}
}


// DeleteWriter deletes it
func (d *GormDatabase) DeleteWriter(uuid string) (bool, error) {
	var consumerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&consumerModel);if query.Error != nil {
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
	query := d.DB.Where("uuid = ?", uuid).Find(&consumerModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&consumerModel).Updates(body);if query.Error != nil {
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
