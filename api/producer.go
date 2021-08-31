package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

)

/*

add producer
- user needs to pass a valid producer uuid (for example a point uuid) and type model.ProducerType (Point, Job, Alarm) and  model.ProducerApplication (Plugin, Remote, Local)

example for workflow for a point (Point 1 Has a Consumer to Point 2):
Point 1
-- consumer table -> point 2 uuid
-- producer table -> nil

Point 2
-- consumer table -> nil
-- producer table -> point 1 uuid

remote point will subscribe to cov events
- the local db will store a copy of the producer to know where to publish the data to
- the remote device will store a copy of its consumer in the consumers table, these will be the details of the remote producer

remote producer
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)

*/


// The ProducerDatabase interface for encapsulating database access.
type ProducerDatabase interface {
	GetProducer(uuid string) (*model.Producer, error)
	GetProducers() ([]*model.Producer, error)
	CreateProducer(body *model.Producer) (*model.Producer, error)
	UpdateProducer(uuid string, body *model.Producer) (*model.Producer, error)
	DeleteProducer(uuid string) (bool, error)

}
type ProducerAPI struct {
	DB ProducerDatabase
}


func (j *ProducerAPI) GetProducer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetProducer(uuid)
	reposeHandler(q, err, ctx)
}


func (j *ProducerAPI) GetProducers(ctx *gin.Context) {
	q, err := j.DB.GetProducers()
	reposeHandler(q, err, ctx)

}

func (j *ProducerAPI) CreateProducer(ctx *gin.Context) {
	body, _ := getBODYProducer(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	body, err = j.DB.CreateProducer(body)
	reposeHandler(body, err, ctx)
}


func (j *ProducerAPI) UpdateProducer(ctx *gin.Context) {
	body, _ := getBODYProducer(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateProducer(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *ProducerAPI) DeleteProducer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteProducer(uuid)
	reposeHandler(q, err, ctx)
}


//func (j *ProducerAPI) CreateJobProducer(ctx *gin.Context) {
//	body, _ := getBODYJobProducer(ctx)
//	err := j.DB.CreateJobProducer(body, body.JobUUID)
//	reposeHandler(body, err, ctx)
//}

//
//func (j *JobAPI) UpdateJobProducer(ctx *gin.Context) {
//	body, _ := getBODYJobProducer(ctx)
//	uuid := resolveID(ctx)
//	q, err := j.DB.UpdateJobProducer(uuid, body)
//	reposeHandler(q, err, ctx)
//}
//
//func (j *JobAPI) GetJobProducer(ctx *gin.Context) {
//	q, err := j.DB.GetJobProducers()
//	reposeHandler(q, err, ctx)
//}
//
//
//func (j *JobAPI) DeleteJobProducer(ctx *gin.Context) {
//	uuid := resolveID(ctx)
//	q, err := j.DB.DeleteJobProducer(uuid)
//	reposeHandler(q, err, ctx)
//}
//

/*
add job
- don't start job until it has one or more producers

edit job
- if is set disable then notify all producer's and for enable do the same

delete job
- notify all producer's, and they will unsubscribe

job producer
- On add: make sure the job uuid is valid
- On delete: update the producer's list and unsubscribe

remote producer
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)

*/

