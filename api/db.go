package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (string, error) //delete all networks, gateways and children
	SyncTopics()                  //sync all the topics into the event bus
	WizardP2PMapping(body *model.P2PBody) (bool, error)
	WizardMasterSlavePointMapping() (bool, error)
	WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
	WizardP2PMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
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

func (a *DatabaseAPI) WizardP2PMapping(ctx *gin.Context) {
	body, _ := getP2PBody(ctx)
	mapping, err := a.DB.WizardP2PMapping(body)
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardMasterSlavePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardMasterSlavePointMapping()
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID)
	reposeHandler(sch, err, ctx)
}

func (a *DatabaseAPI) WizardP2PMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardP2PMappingOnConsumerSideByProducerSide(globalUUID)
	reposeHandler(sch, err, ctx)
}
