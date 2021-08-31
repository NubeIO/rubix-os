package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)




// The WriterDatabase interface for encapsulating database access.
type WriterDatabase interface {
	GetWriter(uuid string) (*model.Writer, error)
	GetWriters() ([]*model.Writer, error)
	CreateWriter(body *model.Writer) (*model.Writer, error)
	UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error)
	DeleteWriter(uuid string) (bool, error)
	RemoteWriterRead(uuid string) (*model.ProducerHistory, error)
	RemoteWriterWrite(uuid string, body *model.Writer, askRefresh bool) (*model.ProducerHistory, error)

}

type WriterAPI struct {
	DB WriterDatabase
}


func (j *WriterAPI) GetWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriter(uuid)
	reposeHandler(q, err, ctx)
}


func (j *WriterAPI) GetWriters(ctx *gin.Context) {
	q, err := j.DB.GetWriters()
	reposeHandler(q, err, ctx)

}

func (j *WriterAPI) CreateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriter(body)
	reposeHandler(q, err, ctx)
}


func (j *WriterAPI) UpdateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateWriter(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *WriterAPI) DeleteWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriter(uuid)
	reposeHandler(q, err, ctx)
}


//withConsumerArgs
func withConsumerArgs(ctx *gin.Context) (askResponse bool, askRefresh bool, write bool, thingType string, flowNetworkUUID string){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault

	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	args.Write = ctx.DefaultQuery(aType.Write, aDefault.Write)
	args.ThingType = ctx.DefaultQuery(aType.ThingType, aDefault.ThingType)
	args.FlowNetworkUUID = ctx.DefaultQuery(aType.FlowNetworkUUID, aDefault.FlowNetworkUUID)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	write, _ = toBool(args.Write)
	return askRefresh, askResponse, write, args.ThingType, args.FlowNetworkUUID
}


//RemoteWriterRead get or update a producer value by using the consumer uuid
//Default will just read the stored value of the consumer (as in don't get the current value from the producer)
//AskRefresh:   "ask_refresh",  // consumer to ask for value from the producer, And producer must resend its value, But don't wait for a response
//AskResponse:  "ask_response", //consumer to ask for value from the producer, And wait for a response
//Write:  "write", //write a new value to the consumer
//thingsType:  "thing_type", //write a new value to the consumer
func (j *WriterAPI) RemoteWriterRead(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteWriterRead(uuid)
	reposeHandler(q, err, ctx)
}


//RemoteWriterWrite get or update a producer value by using the consumer uuid
func (j *WriterAPI) RemoteWriterWrite(ctx *gin.Context) {
	_, _, write, _, _ := withConsumerArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYWriter(ctx)
	q, err := j.DB.RemoteWriterWrite(uuid, body, write)
	reposeHandler(q, err, ctx)
}
