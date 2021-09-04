package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Node struct {
	*model.NodeList
}

// GetNodesList get all of them
func (d *GormDatabase) GetNodesList() ([]*model.NodeList, error) {
	var producersModel []*model.NodeList
	query := d.DB.Preload("NodeOut1").Preload("NodeIn1").Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateNodeList make it
func (d *GormDatabase) CreateNodeList(body *model.NodeList) (*model.NodeList, error) {
	body.UUID = utils.MakeTopicUUID("")
	body.Name = nameIsNil(body.Name)
	body.NodeType = typeIsNil(body.NodeType, "add")
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetNodeList get it
func (d *GormDatabase) GetNodeList(uuid string) (*model.NodeList, error) {
	var wcm *model.NodeList
	query := d.DB.Preload("NodeOut1").Preload("NodeIn1").Where("uuid = ? ", uuid).First(&wcm); if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DeleteNodeList deletes it
func (d *GormDatabase) DeleteNodeList(uuid string) (bool, error) {
	var wcm *model.NodeList
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wcm);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// UpdateNodeList  update it
func (d *GormDatabase) UpdateNodeList(uuid string, body *model.NodeList) (*model.NodeList, error) {
	var wcm *model.NodeList
	query := d.DB.Where("uuid = ?", uuid).Find(&wcm);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	list, err := d.GetNodeList(uuid)
	if err != nil {
		return nil, err
	}
	busNodes(list.UUID,  list)
	return wcm, nil
}

// DropNodesList delete all.
func (d *GormDatabase) DropNodesList() (bool, error) {
	var wcm *model.NodeList
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

func busNodes(UUID string, body *model.NodeList){
	notificationService := eventbus.NewNotificationService(eventbus.BUS)
	body.UUID = UUID
	notificationService.Emit(eventbus.BusContext,eventbus.NodeEventIn, body)
	fmt.Println("topics", eventbus.BUS.Topics())
}

