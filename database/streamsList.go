package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type StreamList struct {
	*model.StreamList
}

// GetStreamLists get all of them
func (d *GormDatabase) GetStreamLists() ([]*model.StreamList, error) {
	var streamListsModel []*model.StreamList

	query := d.DB.Find(&streamListsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamListsModel, nil
}

// CreateStreamList make it
func (d *GormDatabase) CreateStreamList(body *model.StreamList) (*model.StreamList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.StreamList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetStreamList get it
func (d *GormDatabase) GetStreamList(uuid string) (*model.StreamList, error) {
	var streamListModel *model.StreamList
	query := d.DB.Where("uuid = ? ", uuid).First(&streamListModel); if query.Error != nil {
		return nil, query.Error
	}
	return streamListModel, nil
}


// GetStreamListBySubUUID get it by its
func (d *GormDatabase) GetStreamListBySubUUID(consumerUUID string) (*model.StreamList, error) {
	var streamListModel *model.StreamList
	query := d.DB.Where("consumer_uuid = ? ", consumerUUID).First(&streamListModel); if query.Error != nil {
		return nil, query.Error
	}
	return streamListModel, nil
}


// DeleteStreamList deletes it
func (d *GormDatabase) DeleteStreamList(uuid string) (bool, error) {
	var streamListModel *model.StreamList
	query := d.DB.Where("uuid = ? ", uuid).Delete(&streamListModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateStreamList  update it
func (d *GormDatabase) UpdateStreamList(uuid string, body *model.StreamList) (*model.StreamList, error) {
	var streamListModel *model.StreamList
	query := d.DB.Where("uuid = ?", uuid).Find(&streamListModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
	query = d.DB.Model(&streamListModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return streamListModel, nil

}

// DropStreamList delete all.
func (d *GormDatabase) DropStreamList() (bool, error) {
	var streamListModel *model.StreamList
	query := d.DB.Where("1 = 1").Delete(&streamListModel)
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
