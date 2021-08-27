package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)




// The SubscriberListDatabase interface for encapsulating database access.
type SubscriberListDatabase interface {
	GetSubscriberList(uuid string) (*model.SubscriberList, error)
	GetSubscriberLists() ([]*model.SubscriberList, error)
	CreateSubscriberList(body *model.SubscriberList) (*model.SubscriberList, error)
	UpdateSubscriberList(uuid string, body *model.SubscriberList) (*model.SubscriberList, error)
	DeleteSubscriberList(uuid string) (bool, error)


}
type SubscriberListAPI struct {
	DB SubscriberListDatabase
}


func (j *SubscriberListAPI) GetSubscriberList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetSubscriberList(uuid)
	reposeHandler(q, err, ctx)
}


func (j *SubscriberListAPI) GetSubscriberLists(ctx *gin.Context) {
	q, err := j.DB.GetSubscriberLists()
	reposeHandler(q, err, ctx)

}

func (j *SubscriberListAPI) CreateSubscriberList(ctx *gin.Context) {
	body, _ := getBODYSubscriberList(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateSubscriberList(body)
	reposeHandler(q, err, ctx)
}


func (j *SubscriberListAPI) UpdateSubscriberList(ctx *gin.Context) {
	body, _ := getBODYSubscriberList(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateSubscriberList(uuid, body)
	reposeHandler(q, err, ctx)
}




func (j *SubscriberListAPI) DeleteSubscriberList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteSubscriberList(uuid)
	reposeHandler(q, err, ctx)
}

