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

