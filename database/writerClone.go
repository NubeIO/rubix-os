package database

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

type WriterClone struct {
	*model.WriterClone
}

func (d *GormDatabase) GetWriterClones(args api.Args) ([]*model.WriterClone, error) {
	var writerClones []*model.WriterClone
	query := d.buildWriterCloneQuery(args)
	err := query.Find(&writerClones).Error
	if err != nil {
		return nil, query.Error
	}
	return writerClones, nil
}

func (d *GormDatabase) CreateWriterClone(body *model.WriterClone) (*model.WriterClone, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.WriterClone)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

func (d *GormDatabase) GetWriterClone(uuid string) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

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
		pro.CurrentWriterUUID = uuid
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
			_, err = d.UpdatePointValue(p.ProducerThingUUID, pnt, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return wcm, nil
}

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
