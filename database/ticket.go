package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func checkTicketStatus(status model.TicketStatus) error {
	switch status {
	case model.TicketStatusNew:
		return nil
	case model.TicketStatusReplied:
		return nil
	case model.TicketStatusResolved:
		return nil
	case model.TicketStatusClosed:
		return nil
	case model.TicketStatusBlocked:
		return nil
	}
	return errors.New("invalid ticket status, try NEW, REPLIED, RESOLVED, CLOSED, BLOCKED")
}

func checkTicketPriority(priority model.TicketPriority) error {
	switch priority {
	case model.TicketPriorityLow:
		return nil
	case model.TicketPriorityMedium:
		return nil
	case model.TicketPriorityHigh:
		return nil
	case model.TicketPriorityCritical:
		return nil
	}
	return errors.New("invalid ticket priority, try LOW, MEDIUM, HIGH, CRITICAL")
}

func (d *GormDatabase) GetTickets(args parentArgs.Args) ([]*model.Ticket, error) {
	var ticketsModel []*model.Ticket
	query := d.buildTicketQuery(args)
	if err := query.Find(&ticketsModel).Error; err != nil {
		return nil, err
	}
	return ticketsModel, nil
}

func (d *GormDatabase) GetTicket(uuid string, args parentArgs.Args) (*model.Ticket, error) {
	var ticketModel *model.Ticket
	query := d.buildTicketQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&ticketModel).Error; err != nil {
		return nil, err
	}
	return ticketModel, nil
}

func (d *GormDatabase) CreateTicket(body *model.Ticket) (*model.Ticket, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Ticket)
	if body.Status == "" {
		body.Status = model.TicketStatusNew
	} else {
		if err := checkTicketStatus(body.Status); err != nil {
			return nil, err
		}
	}
	if body.Priority == "" {
		body.Priority = model.TicketPriorityLow
	} else {
		if err := checkTicketPriority(body.Priority); err != nil {
			return nil, err
		}
	}
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
	if err := checkTicketPriority(model.TicketPriority(priority)); err != nil {
		return false, err
	}
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
	if err := checkTicketStatus(model.TicketStatus(status)); err != nil {
		return false, err
	}
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
