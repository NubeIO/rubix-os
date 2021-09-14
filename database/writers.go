package database

import (
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/streams"
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
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

// GetWriterByThing get it by its thing uuid
func (d *GormDatabase) GetWriterByThing(producerThingUUID string) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

// DeleteWriter deletes it
func (d *GormDatabase) DeleteWriter(uuid string) (bool, error) {
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&writerModel)
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
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ?", uuid).Find(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&writerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil

}

// DropWriters delete all.
func (d *GormDatabase) DropWriters() (bool, error) {
	var writerModel *model.Writer
	query := d.DB.Where("1 = 1").Delete(&writerModel)
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
WriterAction read or write a value to the writer and onto the writer clone
1. update writer
2. write to writerClone
3. producer to decide if it's a valid cov
3. producer to write to history
4. return to consumer and update as required if it's a valid cov (update the writerConeUUID, this could be from another flow-framework instance)
*/
func (d *GormDatabase) WriterAction(uuid string, body *model.WriterBody) (*model.ProducerHistory, error) {
	askRefresh := body.AskRefresh
	writer, err := d.GetWriter(uuid)
	if err != nil {
		return nil, err
	}
	data, action, err := streams.ValidateTypes(writer.WriterThingClass, body)
	if err != nil {
		return nil, err
	}
	wc := new(model.WriterClone)
	consumer, err := d.GetConsumer(writer.ConsumerUUID, false)
	if err != nil {
		return nil, errors.New("error: on get consumer")
	}
	consumerUUID := consumer.UUID
	producerUUID := consumer.ProducerUUID
	writerCloneUUID := writer.CloneUUID
	streamUUID := consumer.StreamUUID
	stream, flow, err := d.GetFlowUUID(streamUUID)
	if err != nil || stream.UUID == "nil" {
		return nil, errors.New("error: invalid stream UUID")
	}
	wc.DataStore = data
	writer.DataStore = data
	if flow.IsRemote { //IF IS REMOTE FLOW-NETWORK

		if action == model.CommonNaming.Write {
			_, err = streams.WriteClone(writerCloneUUID, flow, wc, true)
			if err != nil {
				return nil, errors.New("WRITER-REMOTE: on update REMOTE WriteClone feedback")
			}
		}
		producerHistory, err := streams.GetProducerHist(producerUUID, flow)
		if err != nil {
			return nil, errors.New("WRITER-REMOTE: on update GetProducerHist feedback")
		}
		if askRefresh { //THIS WILL MAKE THE CONSUMER REFLECT THE CURRENT STATE OF THE PRODUCER (THIS WOULD BE USED FOR POINT MAPPING ONE TO MANY)
			updateConsumer, err := consumerRefresh(producerHistory)
			if err != nil {
				return nil, err
			}
			_, _ = d.UpdateConsumer(consumerUUID, updateConsumer)
			if err != nil {
				return nil, errors.New("WRITER-REMOTE: on update consumer feedback")
			}
			return producerHistory, err
		} else {
			return producerHistory, err
		}
	} else { //IF IS LOCAL FLOW-NETWORK
		var producerHistory *model.ProducerHistory
		if action == model.CommonNaming.Write {
			producerHistory, err = d.UpdateCloneAndHist(writerCloneUUID, wc, true)
			if err != nil {
				return nil, errors.New("WRITER-LOCAL: error on local WRITE to writer-clone")
			}
		} else {
			producerHistory, err = d.HistoryLatestByProducerUUID(producerUUID)
			if err != nil {
				return nil, errors.New("WRITER-LOCAL: error on local READ to producer history")
			}
		}
		//producer feedback
		if askRefresh {
			updateConsumer, err := consumerRefresh(producerHistory)
			if err != nil {
				return nil, err
			}
			_, _ = d.UpdateConsumer(consumerUUID, updateConsumer)
			if err != nil {
				return nil, errors.New("error: on update consumer feedback")
			}
			return producerHistory, err
		} else {
			return producerHistory, err
		}
	}
}

type hists struct {
	Items []*model.ProducerHistory
}

func buildHists(item *model.ProducerHistory) []*model.ProducerHistory {
	h := new(hists)
	h.Items = append(h.Items, item)
	return h.Items
}

func (d *GormDatabase) WriterBulkAction(body []*model.WriterBulk) ([]*model.ProducerHistory, error) {
	var out []*model.ProducerHistory
	for _, wri := range body {
		b := new(model.WriterBody)
		b.Action = wri.Action
		b.AskRefresh = wri.AskRefresh
		b.Priority = wri.Priority
		action, err := d.WriterAction(wri.WriterUUID, b)
		if err != nil {
			return nil, err
		}
		out = buildHists(action)
	}
	return out, nil

}

func consumerRefresh(producerFeedback *model.ProducerHistory) (*model.Consumer, error) {
	updateConsumer := new(model.Consumer)
	updateConsumer.DataStore = producerFeedback.DataStore
	updateConsumer.CurrentWriterUUID = producerFeedback.ThingWriterUUID
	return updateConsumer, nil
}
