package api

import (
	"github.com/gin-gonic/gin"
)

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	DropAllFlow() (bool, error) //delete all networks, gateways and children
	SyncTopics() //sync all the topics into the event bus
	WizardLocalPointMapping()  (bool, error)


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




