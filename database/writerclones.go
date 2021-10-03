package database

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

type WriterClone struct {
	*model.WriterClone
}

// GetWriterClones get all of them
func (d *GormDatabase) GetWriterClones() ([]*model.WriterClone, error) {
	var producersModel []*model.WriterClone
	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateWriterClone make it
func (d *GormDatabase) CreateWriterClone(body *model.WriterClone) (*model.WriterClone, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.WriterClone)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetWriterClone get it
func (d *GormDatabase) GetWriterClone(uuid string) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// GetWriterCloneBySubUUID get it by its
func (d *GormDatabase) GetWriterCloneBySubUUID(consumerUUID string) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("consumer_uuid = ? ", consumerUUID).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DeleteWriterClone deletes it
func (d *GormDatabase) DeleteWriterClone(uuid string) (bool, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wcm)
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

//TODO both functions are doing the same thing UpdateWriterClone UpdateCloneAndHist so merge into one

// UpdateWriterClone  update it
func (d *GormDatabase) UpdateWriterClone(uuid string, body *model.WriterClone, updateProducer bool) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ?", uuid).Find(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	if updateProducer {
		pro := new(model.Producer)
		proUUID := wcm.ProducerUUID
		pro.ThingWriterUUID = uuid
		p, err := d.UpdateProducer(proUUID, pro)
		if err != nil {
			return nil, err
		}
		if body.DataStore != nil {
			_, err := d.ProducerWriteHist(proUUID, wcm.DataStore)
			if err != nil {
				return nil, err
			}
		}
		if p.ProducerThingClass == model.ThingClass.Point {
			pnt := new(model.Point)
			pri := new(model.Priority)
			err := json.Unmarshal(body.DataStore, &pri)
			if err != nil {
				return nil, err
			}
			pnt.Priority = pri
			_, err = d.UpdatePoint(p.ProducerThingUUID, pnt, true, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return wcm, nil
}

// UpdateCloneAndHist  update it
func (d *GormDatabase) UpdateCloneAndHist(uuid string, body *model.WriterClone, updateProducer bool) (*model.ProducerHistory, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ?", uuid).Find(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	var hist *model.ProducerHistory
	if updateProducer {
		pro := new(model.Producer)
		proUUID := wcm.ProducerUUID
		pro.ThingWriterUUID = uuid
		p, err := d.UpdateProducer(proUUID, pro)
		if err != nil {
			return nil, err
		}
		if p.ProducerThingClass == model.ThingClass.Point {
			pnt := new(model.Point)
			pri := new(model.Priority)
			err := json.Unmarshal(body.DataStore, &pri)
			if err != nil {
				return nil, err
			}
			pnt.Priority = pri
			pntReturn, err := d.UpdatePoint(p.ProducerThingUUID, pnt, true, false)
			if err != nil {
				return nil, err
			}
			//get the point presentValue and add it into the history
			marshal, err := json.Marshal(pntReturn)
			if err != nil {
				return nil, err
			}
			if body.DataStore != nil {
				hist, err = d.ProducerWriteHist(proUUID, marshal)
				if err != nil {
					return nil, err
				}
			}
		}
		return hist, err
	}
	return nil, nil
}

// DropWriterClone delete all.
func (d *GormDatabase) DropWriterClone() (bool, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("1 = 1").Delete(&wcm)
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
