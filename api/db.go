package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (string, error) //delete all networks, gateways and children
	SyncTopics()                  //sync all the topics into the event bus
	WizardLocalPointMapping() (bool, error)
	WizardMasterSlavePointMapping() (bool, error)
	WizardRemotePointMapping(body *model.FlowNetworkCredential) (bool, error)
	WizardRemoteSchedule() (bool, error)
	WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
	WizardRemotePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error)
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

func (a *DatabaseAPI) WizardMasterSlavePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardMasterSlavePointMapping()
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardRemotePointMapping(ctx *gin.Context) {
	body, _ := getBodyFlowNetworkCredential(ctx)
	mapping, err := a.DB.WizardRemotePointMapping(body)
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardRemoteSchedule(ctx *gin.Context) {
	sch, err := a.DB.WizardRemoteSchedule()
	reposeHandler(sch, err, ctx)
}

func (a *DatabaseAPI) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID)
	reposeHandler(sch, err, ctx)
}

func (a *DatabaseAPI) WizardRemotePointMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardRemotePointMappingOnConsumerSideByProducerSide(globalUUID)
	reposeHandler(sch, err, ctx)
}
