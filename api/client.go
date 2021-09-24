package api

import (
	"fmt"
	auth2 "github.com/NubeDev/flow-framework/src/auth"

	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ClientDatabase interface for encapsulating database access.
type ClientDatabase interface {
	CreateClient(client *model.Client) error
	GetClientByToken(token string) (*model.Client, error)
	GetClientByID(id uint) (*model.Client, error)
	GetClientsByUser(userID uint) ([]*model.Client, error)
	DeleteClientByID(id uint) error
	UpdateClient(client *model.Client) error
}

// The ClientAPI provides handlers for managing clients and applications.
type ClientAPI struct {
	DB            ClientDatabase
	ImageDir      string
	NotifyDeleted func(uint, string)
}

// UpdateClient updates a client by its id.
// swagger:operation PUT /client/{id} client updateClient
//
// Update a client.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// parameters:
// - name: body
//   in: body
//   description: the client to update
//   required: true
//   schema:
//     $ref: "#/definitions/Client"
// - name: id
//   in: path
//   description: the client id
//   required: true
//   type: integer
//   format: int64
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Client"
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) UpdateClient(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		client, err := a.DB.GetClientByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if client != nil && client.UserID == auth2.GetUserID(ctx) {
			newValues := &model.Client{}
			if err := ctx.Bind(newValues); err == nil {
				client.Name = newValues.Name

				if success := successOrAbort(ctx, 500, a.DB.UpdateClient(client)); !success {
					return
				}
				ctx.JSON(200, client)
			}
		} else {
			ctx.AbortWithError(404, fmt.Errorf("client with id %d doesn't exists", id))
		}
	})
}

// CreateClient creates a client and returns the access token.
// swagger:operation POST /client client createClient
//
// Create a client.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// parameters:
// - name: body
//   in: body
//   description: the client to add
//   required: true
//   schema:
//     $ref: "#/definitions/Client"
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Client"
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) CreateClient(ctx *gin.Context) {
	client := model.Client{}
	if err := ctx.Bind(&client); err == nil {
		client.Token = auth2.GenerateNotExistingToken(generateClientToken, a.clientExists)
		client.UserID = auth2.GetUserID(ctx)
		if success := successOrAbort(ctx, 500, a.DB.CreateClient(&client)); !success {
			return
		}
		ctx.JSON(200, client)
	}
}

// GetClients returns all clients a user has.
// swagger:operation GET /client client getClients
//
// Return all clients.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//     schema:
//       type: array
//       items:
//         $ref: "#/definitions/Client"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) GetClients(ctx *gin.Context) {
	userID := auth2.GetUserID(ctx)
	clients, err := a.DB.GetClientsByUser(userID)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	ctx.JSON(200, clients)
}

// DeleteClient deletes a client by its id.
// swagger:operation DELETE /client/{id} client deleteClient
//
// Delete a client.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: id
//   in: path
//   description: the client id
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) DeleteClient(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		client, err := a.DB.GetClientByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if client != nil && client.UserID == auth2.GetUserID(ctx) {
			a.NotifyDeleted(client.UserID, client.Token)
			successOrAbort(ctx, 500, a.DB.DeleteClientByID(id))
		} else {
			ctx.AbortWithError(404, fmt.Errorf("client with id %d doesn't exists", id))
		}
	})
}

func (a *ClientAPI) clientExists(token string) bool {
	client, _ := a.DB.GetClientByToken(token)
	return client != nil
}
