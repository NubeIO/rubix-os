package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

)




// The SubscriptionsDatabase interface for encapsulating database access.
type SubscriptionsDatabase interface {
	GetSubscription(uuid string) (*model.Subscriptions, error)
	GetSubscriptions() ([]model.Subscriptions, error)
	CreateSubscription(body *model.Subscriptions) error
	UpdateSubscription(uuid string, body *model.Subscriptions) (*model.Subscriptions, error)
	DeleteSubscription(uuid string) (bool, error)


}
type SubscriptionsAPI struct {
	DB SubscriptionsDatabase
}


func (j *SubscriptionsAPI) GetSubscription(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetSubscription(uuid)
	reposeHandler(q, err, ctx)
}


func (j *SubscriptionsAPI) GetSubscriptions(ctx *gin.Context) {
	q, err := j.DB.GetSubscriptions()
	reposeHandler(q, err, ctx)

}

func (j *SubscriptionsAPI) CreateSubscription(ctx *gin.Context) {
	body, _ := getBODYSubscription(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	err = j.DB.CreateSubscription(body)
	reposeHandler(body, err, ctx)
}


func (j *SubscriptionsAPI) UpdateSubscription(ctx *gin.Context) {
	body, _ := getBODYSubscription(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateSubscription(uuid, body)
	reposeHandler(q, err, ctx)
}




func (j *SubscriptionsAPI) DeleteSubscription(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteSubscription(uuid)
	reposeHandler(q, err, ctx)
}

