package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type WriterCopy struct {
	*model.WriterClone
}

// GetWriterCopys get all of them
func (d *GormDatabase) GetWriterCopys() ([]*model.WriterClone, error) {
	var producersModel []*model.WriterClone

	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateWriterCopy make it
func (d *GormDatabase) CreateWriterCopy(body *model.WriterClone) (*model.WriterClone, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Writer)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetWriterCopy get it
func (d *GormDatabase) GetWriterCopy(uuid string) (*model.WriterClone, error) {
	var producerModel *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// GetWriterCopyBySubUUID get it by its
func (d *GormDatabase) GetWriterCopyBySubUUID(consumerUUID string) (*model.WriterClone, error) {
	var producerModel *model.WriterClone
	query := d.DB.Where("consumer_uuid = ? ", consumerUUID).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// DeleteWriterCopy deletes it
func (d *GormDatabase) DeleteWriterCopy(uuid string) (bool, error) {
	var producerModel *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&producerModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateWriterCopy  update it
func (d *GormDatabase) UpdateWriterCopy(uuid string, body *model.WriterClone) (*model.WriterClone, error) {
	var producerModel *model.WriterClone
	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
	query = d.DB.Model(&producerModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil

}

// DropWriterCopy delete all.
func (d *GormDatabase) DropWriterCopy() (bool, error) {
	var producerModel *model.WriterClone
	query := d.DB.Where("1 = 1").Delete(&producerModel)
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
// get
 */

//// AddHistory  add a history record
//func (d *GormDatabase) AddHistory(uuid string, body *model.WriterClone) (*model.ProducerHistory, error) {
//	var producerModel *model.WriterClone
//	var producerHist *model.ProducerHistory
//	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel);if query.Error != nil {
//		return nil, query.Error
//	}
//	body.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
//	query = d.DB.Model(&producerModel).Updates(body);if query.Error != nil {
//		return nil, query.Error
//	}
//	return producerModel, nil
//
//}