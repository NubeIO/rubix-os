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
	SubscriptionAction(uuid string, body interface{}, askRefresh bool, askResponse bool, write bool, thingType string) (interface{}, error)
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
func withSubscriptionArgs(ctx *gin.Context) (askResponse bool, askRefresh bool, write bool, thingType string){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	args.Write = ctx.DefaultQuery(aType.Write, aDefault.Write)
	args.ThingType = ctx.DefaultQuery(aType.ThingType, aDefault.ThingType)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	write, _ = toBool(args.Write)
	return askRefresh, askResponse, write, args.ThingType
}


//SubscriptionAction get or update a producer value by using the subscription uuid
//Default will just read the stored value of the subscription (as in don't get the current value from the producer)
//AskRefresh:   "ask_refresh",  // subscription to ask for value from the producer, And producer must resend its value, But don't wait for a response
//AskResponse:  "ask_response", //subscription to ask for value from the producer, And wait for a response
//Write:  "write", //write a new value to the subscription
//thingsType:  "thing_type", //write a new value to the subscription
func (j *SubscriptionListAPI) SubscriptionAction(ctx *gin.Context) {
	askRefresh, askResponse, write, thingType := withSubscriptionArgs(ctx)
	uuid := resolveID(ctx)
	//TODO is a remote subscriber then logic needs to be added
	if thingType == model.CommonNaming.Point{
		body, _ := getBODYPoint(ctx) //TODO add in support for other types
		q, err := j.DB.SubscriptionAction(uuid, body ,askRefresh, askResponse, write, thingType)
		reposeHandler(q, err, ctx)
	} else {
		reposeHandler(nil, nil, ctx)
	}
}