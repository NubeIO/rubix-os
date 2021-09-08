package database

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/patrickmn/go-cache"
)


type Node struct {
	*model.Node
}

// GetNodesList get all of them
func (d *GormDatabase) GetNodesList() ([]*model.Node, error) {
	var producersModel []*model.Node
	query := d.DB.Preload("Out1Connections").Preload("In1Connections").Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateNode make it
func (d *GormDatabase) CreateNode(body *model.Node) (*model.Node, error) {
	body.UUID = utils.MakeTopicUUID("")
	body.Name = nameIsNil(body.Name)
	body.NodeType = typeIsNil(body.NodeType, "add")
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	eventbus.NodeContext.Set(body.UUID, body, cache.NoExpiration)
	return body, nil
}

// GetNode get it
func (d *GormDatabase) GetNode(uuid string) (*model.Node, error) {
	var wcm *model.Node
	query := d.DB.Preload("Out1Connections").Preload("In1Connections").Where("uuid = ? ", uuid).First(&wcm); if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DeleteNode deletes it
func (d *GormDatabase) DeleteNode(uuid string) (bool, error) {
	var wcm *model.Node
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

// UpdateNode  update it
func (d *GormDatabase) UpdateNode(uuid string, body *model.Node) (*model.Node, error) {
	var wcm *model.Node
	query := d.DB.Where("uuid = ?", uuid).Find(&wcm);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	list, err := d.GetNode(uuid)
	if err != nil {
		return nil, err
	}
	eventbus.NodeContext.Set(list.UUID, list, cache.NoExpiration)
	d.Bus.Emit(eventbus.CTX(), eventbus.NodeEventIn,  list)
	return wcm, nil
}

// DropNodesList delete all.
func (d *GormDatabase) DropNodesList() (bool, error) {
	var wcm *model.Node
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

