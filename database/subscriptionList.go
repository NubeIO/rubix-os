package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
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

// pass in writer UUID as what starts everything
// get consumer uuid
// get producer uuid
// get flow network uuid
// get consumer uuid
// update writeValue to writer value
// update writeValue to writerCopy value
// producer to decide to aspect the new write value
// producer to send a COV event to the consumer if the value was updated.
// add histories if enabled



//point
// - presentValue

// consumer needs readonly of producer presentValue
// - presentValue

// writer needs to write to the writerCopy
// - writeValue

//producer
// - presentValue
// - SLWriteUUID //writer UUID

//writerCopy (consumerUUID)
// - writeValue

//producerHist
// - presentValue
// - SLWriteUUID //writer UUID

// ConsumerAction get its value
func (d *GormDatabase) ConsumerAction(uuid string, body *model.Writer, write bool) (*model.Producer, error) {
	var slm *model.Writer
	writer := d.DB.Where("uuid = ? ", uuid).First(&slm); if writer.Error != nil {
		return nil, writer.Error
	}
	if slm == nil {
		return nil, nil
	}
	var sm *model.Consumer
	consumer := d.DB.Where("uuid = ? ", slm.ConsumerUUID).First(&sm); if consumer.Error != nil {
		return nil, consumer.Error
	}
	subType := sm.ConsumerType
	consumerUUID := sm.UUID
	streamUUID := sm.StreamUUID
	producerUUID := sm.ProducerUUID
	writeV := body.WriteValue
 
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
			pm.PresentValue = body.WriteValue
			query = d.DB.Model(&pm).Updates(pm);if query.Error != nil {
				return nil, query.Error
			}
			ph := new(model.ProducerHistory)
			ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			ph.ProducerUUID = producerUUID
			ph.PresentValue = pm.PresentValue
			query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
				return nil, query.Error
			}
			return pm, nil
		} else {
			d.DB.Where("uuid = ? ", uuid).First(&pm); if query.Error != nil {
				return nil, query.Error
			}
			return pm, nil
		}
	} else {
		pm := new(model.Producer)
		pm.UUID = producerUUID
		pm.PresentValue = body.WriteValue
		point, err := eventbus.EventREST(fn, pm, write)
		if err != nil {
			return nil, err
		}
		return point, err
	}
}

// pass in writer UUID as what starts everything
// get consumer uuid
// get producer uuid
// get flow network uuid
// get consumer uuid
// get update writeValue to consumer
// send event to point to update writeValue and update the point
// update the priory array
// update the producer based of what the point returns
// producer to send a COV event to the consumer.
// consumer to update the presentValue and the history if enabled

// ConsumerActionPoint get its value or write
func (d *GormDatabase) ConsumerActionPoint(slUUID string, pointBody *model.Point, write bool) (*model.Producer, error) {
	var slm *model.Writer
	sl := d.DB.Where("uuid = ? ", slUUID).First(&slm); if sl.Error != nil {
		return nil, sl.Error
	}
	if slm == nil {
		return nil, nil
	}

	var sm *model.Consumer
	consumer := d.DB.Where("uuid = ? ", slm.ConsumerUUID).First(&sm); if consumer.Error != nil {
		return nil, consumer.Error
	}
	subType := sm.ConsumerType
	consumerUUID := sm.UUID
	streamUUID := sm.StreamUUID
	producerUUID := sm.ProducerUUID
	writeV := pointBody.WriteValue
	pointUUID := sm.ProducerThingUUID

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
			point, err := eventbus.EventRESTPoint(pointUUID, fn, pointBody, write)
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
			point, err := eventbus.EventRESTPoint(pointUUID, fn, pointBody, write)
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
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
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
