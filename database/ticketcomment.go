package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetTicketComment(uuid string) (*model.TicketComment, error) {
	var ticketCommentModel *model.TicketComment
	if err := d.DB.Where("uuid = ?", uuid).First(&ticketCommentModel).Error; err != nil {
		return nil, err
	}
	return ticketCommentModel, nil
}

func (d *GormDatabase) CreateTicketComment(body *model.TicketComment) (*model.TicketComment, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.TicketComment)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateTicketComment(uuid string, body *model.TicketComment) (*model.TicketComment, error) {
	var ticketCommentModel *model.TicketComment
	query := d.DB.Where("uuid = ?", uuid).First(&ticketCommentModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Owner != ticketCommentModel.Owner {
		return nil, errors.New("you cannot update this comment")
	}
	if err := d.DB.Model(&ticketCommentModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return ticketCommentModel, nil
}

func (d *GormDatabase) DeleteTicketComment(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.TicketComment{})
	return d.deleteResponseBuilder(query)
}
