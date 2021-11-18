package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type TagDatabase interface {
	GetTags(args Args) ([]*model.Tag, error)
	GetTag(tag string, args Args) (*model.Tag, error)
	CreateTag(body *model.Tag) (*model.Tag, error)
	DeleteTag(tag string) (bool, error)
}

type TagAPI struct {
	DB TagDatabase
}

func (j *TagAPI) GetTags(ctx *gin.Context) {
	args := buildTagArgs(ctx)
	q, err := j.DB.GetTags(args)
	responseHandler(q, err, ctx)
}

func (j *TagAPI) GetTag(ctx *gin.Context) {
	args := buildTagArgs(ctx)
	tag := getTagParam(ctx)
	q, err := j.DB.GetTag(tag, args)
	responseHandler(q, err, ctx)
}

func (j *TagAPI) CreateTag(ctx *gin.Context) {
	body, err := getBodyTag(ctx)
	if err != nil {
		responseHandler(nil, err, ctx)
	} else {
		q, e := j.DB.CreateTag(body)
		responseHandler(q, e, ctx)
	}
}

func (j *TagAPI) DeleteTag(ctx *gin.Context) {
	tag := getTagParam(ctx)
	q, err := j.DB.DeleteTag(tag)
	responseHandler(q, err, ctx)
}
