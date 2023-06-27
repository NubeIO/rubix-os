package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetTickets() ([]*model.Ticket, error) {
	var ticketsModel []*model.Ticket
	query := d.buildTicketQuery()
	if err := query.Find(&ticketsModel).Error; err != nil {
		return nil, err
	}
	return ticketsModel, nil
}

func (d *GormDatabase) GetTicket(uuid string) (*model.Ticket, error) {
	var ticketModel *model.Ticket
	query := d.buildTicketQuery()
	if err := query.Where("uuid = ?", uuid).First(&ticketModel).Error; err != nil {
		return nil, err
	}
	return ticketModel, nil
}

func (d *GormDatabase) CreateTicket(body *model.Ticket) (*model.Ticket, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Ticket)
	body.Status = model.TicketStatusNew
	body.Teams = nil
	body.Comments = nil
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateTicket(uuid string, body *model.Ticket) (*model.Ticket, error) {
	var ticketModel *model.Ticket
	query := d.DB.Where("uuid = ?", uuid).First(&ticketModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Issuer != ticketModel.Issuer {
		return nil, errors.New("you cannot update this ticket")
	}
	body.Issuer = ticketModel.Issuer
	body.Priority = ticketModel.Priority
	body.Status = ticketModel.Status
	if err := d.DB.Model(&ticketModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return ticketModel, nil
}

func (d *GormDatabase) UpdateTicketPriority(uuid string, priority string) (bool, error) {
	var ticketModel *model.Ticket
	query := d.DB.Where("uuid = ?", uuid).First(&ticketModel)
	if query.Error != nil {
		return false, query.Error
	}
	if err := d.DB.Model(&ticketModel).Update("priority", priority).
		Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) UpdateTicketStatus(uuid string, status string) (bool, error) {
	var ticketModel *model.Ticket
	query := d.DB.Where("uuid = ?", uuid).First(&ticketModel)
	if query.Error != nil {
		return false, query.Error
	}
	if err := d.DB.Model(&ticketModel).Update("status", status).
		Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) DeleteTicket(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Ticket{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DropTickets() (bool, error) {
	query := d.DB.Where("1 = 1").Delete(&model.Ticket{})
	return d.deleteResponseBuilder(query)
}
