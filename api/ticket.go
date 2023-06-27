package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/gin-gonic/gin"
)

type TicketDatabase interface {
	GetTickets() ([]*model.Ticket, error)
	GetTicket(uuid string) (*model.Ticket, error)
	CreateTicket(body *model.Ticket) (*model.Ticket, error)
	UpdateTicket(uuid string, body *model.Ticket) (*model.Ticket, error)
	UpdateTicketPriority(uuid string, priority string) (bool, error)
	UpdateTicketStatus(uuid string, status string) (bool, error)
	UpdateTicketTeams(ticketUUID string, teamUUIDs []*string) ([]*model.TicketTeam, error)
	DeleteTicket(uuid string) (bool, error)
	DropTickets() (bool, error)
}

type TicketAPI struct {
	DB TicketDatabase
}

func (a *TicketAPI) GetTickets(ctx *gin.Context) {
	q, err := a.DB.GetTickets()
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) GetTicket(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetTicket(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) CreateTicket(ctx *gin.Context) {
	username, err := getAuthorizedOrDefaultUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	body, _ := getBodyTicket(ctx)
	body.Issuer = username
	q, err := a.DB.CreateTicket(body)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) UpdateTicket(ctx *gin.Context) {
	username, err := getAuthorizedOrDefaultUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	body, _ := getBodyTicket(ctx)
	uuid := resolveID(ctx)
	body.Issuer = username
	q, err := a.DB.UpdateTicket(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) UpdateTicketPriority(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyTicketPriority(ctx)
	q, err := a.DB.UpdateTicketPriority(uuid, body.Priority)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) UpdateTicketStatus(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyTicketStatus(ctx)
	q, err := a.DB.UpdateTicketStatus(uuid, body.Status)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) UpdateTicketTeams(ctx *gin.Context) {
	body, _ := getBodyTicketTeams(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateTicketTeams(uuid, body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) DeleteTicket(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteTicket(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) DropTickets(ctx *gin.Context) {
	q, err := a.DB.DropTickets()
	ResponseHandler(q, err, ctx)
}
