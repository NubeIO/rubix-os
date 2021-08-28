package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)




// The ProducerListDatabase interface for encapsulating database access.
type ProducerListDatabase interface {
	GetProducerList(uuid string) (*model.ProducerList, error)
	GetProducerLists() ([]*model.ProducerList, error)
	CreateProducerList(body *model.ProducerList) (*model.ProducerList, error)
	UpdateProducerList(uuid string, body *model.ProducerList) (*model.ProducerList, error)
	DeleteProducerList(uuid string) (bool, error)


}
type ProducerListAPI struct {
	DB ProducerListDatabase
}


func (j *ProducerListAPI) GetProducerList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetProducerList(uuid)
	reposeHandler(q, err, ctx)
}


func (j *ProducerListAPI) GetProducerLists(ctx *gin.Context) {
	q, err := j.DB.GetProducerLists()
	reposeHandler(q, err, ctx)

}

func (j *ProducerListAPI) CreateProducerList(ctx *gin.Context) {
	body, _ := getBODYProducerList(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateProducerList(body)
	reposeHandler(q, err, ctx)
}


func (j *ProducerListAPI) UpdateProducerList(ctx *gin.Context) {
	body, _ := getBODYProducerList(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateProducerList(uuid, body)
	reposeHandler(q, err, ctx)
}




func (j *ProducerListAPI) DeleteProducerList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteProducerList(uuid)
	reposeHandler(q, err, ctx)
}

