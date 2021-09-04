package api

import (

	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (bool, error) //delete all networks, gateways and children
	SyncTopics() //sync all the topics into the event bus
	WizardLocalPointMapping()  (bool, error)
	WizardRemotePointMapping()  (bool, error)
	Wizard2ndFlowNetwork(body *AddNewFlowNetwork)  (bool, error)
	NodeWizard()  (bool, error)

}
type DatabaseAPI struct {
	DB DBDatabase
}

func (a *DatabaseAPI) DropAllFlow(ctx *gin.Context) {
	q, err := a.DB.DropAllFlow()
	reposeHandler(q, err, ctx)
}

func (a *DatabaseAPI) SyncTopics() {
	a.DB.SyncTopics()
}

func (a *DatabaseAPI) WizardLocalPointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardLocalPointMapping()
	reposeHandler(mapping, err, ctx)
}


func (a *DatabaseAPI) WizardRemotePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardRemotePointMapping()
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) NodeWizard(ctx *gin.Context) {
	mapping, err := a.DB.NodeWizard()
	reposeHandler(mapping, err, ctx)
}


type AddNewFlowNetwork struct {
	StreamUUID        	string `json:"stream_uuid"`
	StreamListUUID 		string `json:"stream_list_uuid"`
	ProducerUUID      	string `json:"producer_uuid"`
	ExistingPointUUID 	string `json:"existing_point_uuid"`
	FlowToken 			string 	`json:"flow_token"`

}


func getBODYWizard(ctx *gin.Context) (dto *AddNewFlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}


func (a *DatabaseAPI) Wizard2ndFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYWizard(ctx)
	q, err := a.DB.Wizard2ndFlowNetwork(body)
	reposeHandler(q, err, ctx)
}


