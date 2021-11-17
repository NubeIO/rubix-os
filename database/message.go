package database

import (
	"github.com/NubeIO/flow-framework/model"
	"gorm.io/gorm"
)

func (d *GormDatabase) CreateMessage(message *model.Message) error {
	return d.DB.Create(message).Error
}

func (d *GormDatabase) GetMessages() ([]*model.Message, error) {
	var messages []*model.Message
	err := d.DB.Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
}

func (d *GormDatabase) GetMessageByID(id uint) (*model.Message, error) {
	msg := new(model.Message)
	err := d.DB.Find(msg, id).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if msg.ID == id {
		return msg, err
	}
	return nil, err
}

func (d *GormDatabase) GetMessagesSince(limit int, since uint) ([]*model.Message, error) {
	var messages []*model.Message
	db := d.DB.Order("id desc").Limit(limit)
	if since != 0 {
		db = db.Where("messages.id < ?", since)
	}
	err := db.Find(&messages).Error

	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
}

func (d *GormDatabase) DeleteMessages() error {
	return d.DB.Where("1 = 1").Delete(&model.Message{}).Error
}

func (d *GormDatabase) DeleteMessageByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.Message{}).Error
}
