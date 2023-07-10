package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/gin-gonic/gin"
)

type TicketDatabase interface {
	GetTickets(args argspkg.Args) ([]*model.Ticket, error)
	GetTicket(uuid string, args argspkg.Args) (*model.Ticket, error)
	CreateTicket(body *model.Ticket) (*model.Ticket, error)
	UpdateTicket(uuid string, body *model.Ticket) (*model.Ticket, error)
	UpdateTicketPriority(uuid string, priority string) (bool, error)
	UpdateTicketStatus(uuid string, status string) (bool, error)
	UpdateTicketTeams(ticketUUID string, teamUUIDs []*string) ([]*model.TicketTeam, error)
	DeleteTicket(uuid string) (bool, error)
}

type TicketAPI struct {
	DB TicketDatabase
}

func (a *TicketAPI) GetTickets(ctx *gin.Context) {
	args := buildTicketArgs(ctx)
	q, err := a.DB.GetTickets(args)
	ResponseHandler(q, err, ctx)
}

func (a *TicketAPI) GetTicket(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildTicketArgs(ctx)
	q, err := a.DB.GetTicket(uuid, args)
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
	_, err := a.DB.DeleteTicket(uuid)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "ticket deleted successfully"}, err, ctx)
}
