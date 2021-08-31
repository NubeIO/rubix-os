package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)




// The WriterCopyDatabase interface for encapsulating database access.
type WriterCopyDatabase interface {
	GetWriterCopy(uuid string) (*model.WriterClone, error)
	GetWriterCopys() ([]*model.WriterClone, error)
	CreateWriterCopy(body *model.WriterClone) (*model.WriterClone, error)
	UpdateWriterCopy(uuid string, body *model.WriterClone) (*model.WriterClone, error)
	DeleteWriterCopy(uuid string) (bool, error)


}
type WriterCopyAPI struct {
	DB WriterCopyDatabase
}


func (j *WriterCopyAPI) GetWriterCopy(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriterCopy(uuid)
	reposeHandler(q, err, ctx)
}


func (j *WriterCopyAPI) GetWriterCopys(ctx *gin.Context) {
	q, err := j.DB.GetWriterCopys()
	reposeHandler(q, err, ctx)

}

func (j *WriterCopyAPI) CreateWriterCopy(ctx *gin.Context) {
	body, _ := getBODYWriterCopy(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriterCopy(body)
	reposeHandler(q, err, ctx)
}


func (j *WriterCopyAPI) UpdateWriterCopy(ctx *gin.Context) {
	body, _ := getBODYWriterCopy(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateWriterCopy(uuid, body)
	reposeHandler(q, err, ctx)
}




func (j *WriterCopyAPI) DeleteWriterCopy(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriterCopy(uuid)
	reposeHandler(q, err, ctx)
}

