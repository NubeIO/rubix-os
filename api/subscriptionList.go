package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)




// The SubscriptionListDatabase interface for encapsulating database access.
type SubscriptionListDatabase interface {
	GetSubscriptionList(uuid string) (*model.SubscriptionList, error)
	GetSubscriptionLists() ([]*model.SubscriptionList, error)
	CreateSubscriptionList(body *model.SubscriptionList) (*model.SubscriptionList, error)
	UpdateSubscriptionList(uuid string, body *model.SubscriptionList) (*model.SubscriptionList, error)
	DeleteSubscriptionList(uuid string) (bool, error)
	SubscriptionRead(uuid string, askRefresh bool, askResponse bool) (interface{}, error)
	SubscriptionWrite(uuid string, askRefresh bool, askResponse bool) (*model.SubscriptionList, error)
}

type SubscriptionListAPI struct {
	DB SubscriptionListDatabase
}


func (j *SubscriptionListAPI) GetSubscriptionList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetSubscriptionList(uuid)
	reposeHandler(q, err, ctx)
}


func (j *SubscriptionListAPI) GetSubscriptionLists(ctx *gin.Context) {
	q, err := j.DB.GetSubscriptionLists()
	reposeHandler(q, err, ctx)

}

func (j *SubscriptionListAPI) CreateSubscriptionList(ctx *gin.Context) {
	body, _ := getBODYSubscriptionList(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateSubscriptionList(body)
	reposeHandler(q, err, ctx)
}


func (j *SubscriptionListAPI) UpdateSubscriptionList(ctx *gin.Context) {
	body, _ := getBODYSubscriptionList(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateSubscriptionList(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *SubscriptionListAPI) DeleteSubscriptionList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteSubscriptionList(uuid)
	reposeHandler(q, err, ctx)
}


//withSubscriptionArgs
func withSubscriptionArgs(ctx *gin.Context) (askResponse bool, askRefresh bool){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	return askRefresh, askResponse
}


//SubscriptionRead get the latest value
//Default will just read the stored value of the subscription (as in don't get the current value from the producer)
//AskRefresh:   "ask_refresh",  // subscription to ask for value from the producer, And producer must resend its value, But don't wait for a response
//AskResponse:  "ask_response", //subscription to ask for value from the producer, And wait for a response
func (j *SubscriptionListAPI) SubscriptionRead(ctx *gin.Context) {
	askRefresh, askResponse := withSubscriptionArgs(ctx) //TODO add this in
	uuid := resolveID(ctx)
	q, err := j.DB.SubscriptionRead(uuid, askRefresh, askResponse)
	reposeHandler(q, err, ctx)

}

func (j *SubscriptionListAPI) SubscriptionWrite(ctx *gin.Context) {
	askResponse, askRefresh := withSubscriptionArgs(ctx) //TODO add this in
	uuid := resolveID(ctx)
	q, err := j.DB.SubscriptionWrite(uuid, askRefresh, askResponse)
	reposeHandler(q, err, ctx)

}