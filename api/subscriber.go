package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

)

/*

add subscriber
- user needs to pass a valid producer uuid (for example a point uuid) and type model.SubscriberType (Point, Job, Alarm) and  model.SubscriberApplication (Plugin, Remote, Local)

example for workflow for a point (Point 1 Has a Subscription to Point 2):
Point 1
-- subscription table -> point 2 uuid
-- subscriber table -> nil

Point 2
-- subscription table -> nil
-- subscriber table -> point 1 uuid

remote point will subscribe to cov events
- the local db will store a copy of the subscriber to know where to publish the data to
- the remote device will store a copy of its subscription in the subscriptions table, these will be the details of the remote producer

remote subscriber
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)

*/


// The SubscriberDatabase interface for encapsulating database access.
type SubscriberDatabase interface {
	GetSubscriber(uuid string) (*model.Subscriber, error)
	GetSubscribers() ([]model.Subscriber, error)
	CreateSubscriber(body *model.Subscriber) error
	UpdateSubscriber(uuid string, body *model.Subscriber) (*model.Subscriber, error)
	DeleteSubscriber(uuid string) (bool, error)
	//CreateJobSubscriber(body *model.JobSubscriber, jobUUID string) error
	//UpdateJobSubscriber(uuid string, body *model.JobSubscriber) (*model.JobSubscriber, error)
	//GetJobSubscribers() ([]model.JobSubscriber, error)
	//DeleteJobSubscriber(uuid string) (bool, error)

}
type SubscriberAPI struct {
	DB SubscriberDatabase
}


func (j *SubscriberAPI) GetSubscriber(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetSubscriber(uuid)
	reposeHandler(q, err, ctx)
}


func (j *SubscriberAPI) GetSubscribers(ctx *gin.Context) {
	q, err := j.DB.GetSubscribers()
	reposeHandler(q, err, ctx)

}

func (j *SubscriberAPI) CreateSubscriber(ctx *gin.Context) {
	body, _ := getBODYSubscriber(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	err = j.DB.CreateSubscriber(body)
	reposeHandler(body, err, ctx)
}


func (j *SubscriberAPI) UpdateSubscriber(ctx *gin.Context) {
	body, _ := getBODYSubscriber(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateSubscriber(uuid, body)
	reposeHandler(q, err, ctx)
}




func (j *SubscriberAPI) DeleteSubscriber(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteSubscriber(uuid)
	reposeHandler(q, err, ctx)
}


//func (j *SubscriberAPI) CreateJobSubscriber(ctx *gin.Context) {
//	body, _ := getBODYJobSubscriber(ctx)
//	err := j.DB.CreateJobSubscriber(body, body.JobUUID)
//	reposeHandler(body, err, ctx)
//}

//
//func (j *JobAPI) UpdateJobSubscriber(ctx *gin.Context) {
//	body, _ := getBODYJobSubscriber(ctx)
//	uuid := resolveID(ctx)
//	q, err := j.DB.UpdateJobSubscriber(uuid, body)
//	reposeHandler(q, err, ctx)
//}
//
//func (j *JobAPI) GetJobSubscriber(ctx *gin.Context) {
//	q, err := j.DB.GetJobSubscribers()
//	reposeHandler(q, err, ctx)
//}
//
//
//func (j *JobAPI) DeleteJobSubscriber(ctx *gin.Context) {
//	uuid := resolveID(ctx)
//	q, err := j.DB.DeleteJobSubscriber(uuid)
//	reposeHandler(q, err, ctx)
//}
//

/*
add job
- don't start job until it has one or more subscribers

edit job
- if is set disable then notify all subscriber's and for enable do the same

delete job
- notify all subscriber's, and they will unsubscribe

job subscriber
- On add: make sure the job uuid is valid
- On delete: update the subscriber's list and unsubscribe

remote subscriber
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)

*/

