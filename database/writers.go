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
- update producer with the writerCloneUUID (this is so we know who wrote the last value to the producer)
- update the consumer
- return message to user

else:
update the writerClone history

*/

func (d *GormDatabase) RemoteWriterWrite(uuid string, body *model.Writer, askRefresh bool) (*model.Consumer, error) {

	fmt.Println(body, 222222)
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
	wc.DataStore = body.DataStore
	// update producer
	_, err = rest.WriteClone(writerCloneUUID, fn, wc, true)
	if err != nil {
		return nil, errors.New("error: write new value to writerClone")
	}
	if askRefresh { //get feedback from producer
		producerFeedback, err := rest.ProducerHistory(fn, producerUUID);if err != nil {
			return nil, errors.New("error: on get feedback from producer history")
		}
		// update the consumer based of the response from the producer
		updateConsumer:= new(model.Consumer)
		updateConsumer.DataStore = producerFeedback.DataStore
		updateConsumer.CurrentWriterCloneUUID = producerFeedback.CurrentWriterCloneUUID
		pro, _ := d.UpdateConsumer(consumerUUID, updateConsumer);if err != nil {
			return nil, errors.New("error: on update consumer feedback")
		}

		return pro, err
	}
	return nil, errors.New("error: fail to update the remote producer")
}

func (d *GormDatabase) RemoteWriterRead(uuid string) (*model.Writer, error) {
	//var wm *model.Writer
		fmt.Println("fix me")
		return nil, nil
	//writer := d.DB.Where("uuid = ? ", uuid).First(&wm); if writer.Error != nil {
	//	return nil, writer.Error
	//}
	//if wm == nil {
	//	return nil, errors.New("error: failed to find writer")
	//}
	//var cm *model.Consumer
	//consumer := d.DB.Where("uuid = ? ", wm.ConsumerUUID).First(&cm); if consumer.Error != nil {
	//	return nil, consumer.Error
	//}
	//consumerUUID := cm.UUID
	//streamUUID := cm.StreamUUID
	//producerUUID := cm.ProducerUUID
	//var s *model.Stream
	//stream := d.DB.Where("uuid = ? ", streamUUID).First(&s); if consumer.Error != nil {
	//	return nil, stream.Error
	//}
	//streamListUUID := s.StreamListUUID
	//var fn *model.FlowNetwork
	//flow := d.DB.Where("stream_list_uuid = ? ", streamListUUID).First(&fn); if consumer.Error != nil {
	//	return nil, flow.Error
	//}
	//fmt.Println(111111)
	//producerFeedback, err := rest.ProducerRead(fn, producerUUID);if err != nil {
	//	return nil, errors.New("error: on get feedback from producer")
	//}
	//fmt.Println(111111)
	//// update the consumer based of the response from the producer
	//updateConsumer:= new(model.Consumer)
	//updateConsumer.PresentValue =producerFeedback.PresentValue
	//pro, _ := d.UpdateConsumer(consumerUUID, updateConsumer);if err != nil {
	//	return nil, errors.New("error: on update consumer feedback")
	//}
	//fmt.Println(111111)
	//newWriter:= new(model.Writer)
	//newWriter.PresentValue = pro.PresentValue
	//// now update the writer
	//fmt.Println(111111)
	//updateWriter, _ := d.UpdateWriter(uuid, newWriter);if err != nil {
	//	return nil, errors.New("error: on update consumer feedback")
	//}
	//return updateWriter, err


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
