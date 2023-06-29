package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/gin-gonic/gin"
)

type TagDatabase interface {
	GetTags(args args.Args) ([]*model.Tag, error)
	GetTag(tag string, args args.Args) (*model.Tag, error)
	CreateTag(body *model.Tag) (*model.Tag, error)
	DeleteTag(tag string) (bool, error)
}

type TagAPI struct {
	DB TagDatabase
}

func (j *TagAPI) GetTags(ctx *gin.Context) {
	args := buildTagArgs(ctx)
	q, err := j.DB.GetTags(args)
	ResponseHandler(q, err, ctx)
}

func (j *TagAPI) GetTag(ctx *gin.Context) {
	args := buildTagArgs(ctx)
	tag := getTagParam(ctx)
	q, err := j.DB.GetTag(tag, args)
	ResponseHandler(q, err, ctx)
}

func (j *TagAPI) CreateTag(ctx *gin.Context) {
	body, err := getBodyTag(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	} else {
		q, e := j.DB.CreateTag(body)
		ResponseHandler(q, e, ctx)
	}
}

func (j *TagAPI) DeleteTag(ctx *gin.Context) {
	tag := getTagParam(ctx)
	q, err := j.DB.DeleteTag(tag)
	ResponseHandler(q, err, ctx)
}
