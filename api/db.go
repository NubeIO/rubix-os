package api

import (
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (string, error) //delete all networks, gateways and children
	SyncTopics()                  //sync all the topics into the event bus
	WizardLocalPointMapping() (bool, error)
	WizardRemotePointMapping() (bool, error)
	WizardRemoteSchedule() (bool, error)
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

func (a *DatabaseAPI) WizardRemotePointMapping(ctx *gin.Context) {
	mapping, err := a.DB.WizardRemotePointMapping()
	reposeHandler(mapping, err, ctx)
}

func (a *DatabaseAPI) WizardRemoteSchedule(ctx *gin.Context) {
	sch, err := a.DB.WizardRemoteSchedule()
	reposeHandler(sch, err, ctx)
}

func (a *DatabaseAPI) WizardRemotePointMappingOnConsumerSideByProducerSide(ctx *gin.Context) {
	globalUUID := resolveGlobalUUID(ctx)
	sch, err := a.DB.WizardRemotePointMappingOnConsumerSideByProducerSide(globalUUID)
	reposeHandler(sch, err, ctx)
}

type AddNewFlowNetwork struct {
	StreamUUID         string `json:"stream_uuid"`
	ProducerUUID       string `json:"producer_uuid"`
	ProducerThingUUID  string `json:"producer_thing_uuid"` // this is the remote point UUID
	ProducerThingClass string `json:"producer_thing_class"`
	ProducerThingType  string `json:"producer_thing_type"`
	FlowToken          string `json:"flow_token"`
}
