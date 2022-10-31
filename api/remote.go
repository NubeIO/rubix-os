package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type RemoteDatabase interface {

	// flow network clones
	RemoteGetFlowNetworkClones(args Args) ([]model.FlowNetworkClone, error)
	RemoteGetFlowNetworkClone(uuid string, args Args) (*model.FlowNetworkClone, error)
	RemoteDeleteFlowNetworkClone(uuid string, args Args) (bool, error)

	// networks
	RemoteGetNetworks(args Args) ([]model.Network, error)
	RemoteGetNetwork(uuid string, args Args) (*model.Network, error)
	RemoteCreateNetwork(body *model.Network, args Args) (*model.Network, error)
	RemoteDeleteNetwork(uuid string, args Args) (bool, error)
	RemoteEditNetwork(uuid string, body *model.Network, args Args) (*model.Network, error)

	// devices
	RemoteGetDevices(args Args) ([]model.Device, error)
	RemoteGetDevice(uuid string, args Args) (*model.Device, error)
	RemoteCreateDevice(body *model.Device, args Args) (*model.Device, error)
	RemoteDeleteDevice(uuid string, args Args) (bool, error)
	RemoteEditDevice(uuid string, body *model.Device, args Args) (*model.Device, error)

	// points
	RemoteGetPoints(args Args) ([]model.Point, error)
	RemoteGetPoint(uuid string, args Args) (*model.Point, error)
	RemoteCreatePoint(body *model.Point, args Args) (*model.Point, error)
	RemoteDeletePoint(uuid string, args Args) (bool, error)
	RemoteEditPoint(uuid string, body *model.Point, args Args) (*model.Point, error)

	// producers
	RemoteGetStreams(args Args) ([]model.Stream, error)
	RemoteGetStream(uuid string, args Args) (*model.Stream, error)
	RemoteCreateStream(body *model.Stream, args Args) (*model.Stream, error)
	RemoteDeleteStream(uuid string, args Args) (bool, error)
	RemoteEditStream(uuid string, body *model.Stream, args Args) (*model.Stream, error)

	// producers
	RemoteGetProducers(args Args) ([]model.Producer, error)
	RemoteGetProducer(uuid string, args Args) (*model.Producer, error)
	RemoteCreateProducer(body *model.Producer, args Args) (*model.Producer, error)
	RemoteDeleteProducer(uuid string, args Args) (bool, error)
	RemoteEditProducer(uuid string, body *model.Producer, args Args) (*model.Producer, error)

	// consumers
	RemoteGetConsumers(args Args) ([]model.Consumer, error)
	RemoteGetConsumer(uuid string, args Args) (*model.Consumer, error)
	RemoteCreateConsumer(body *model.Consumer, args Args) (*model.Consumer, error)
	RemoteDeleteConsumer(uuid string, args Args) (bool, error)
	RemoteEditConsumer(uuid string, body *model.Consumer, args Args) (*model.Consumer, error)

	// writers
	RemoteGetWriters(args Args) ([]model.Writer, error)
	RemoteGetWriter(uuid string, args Args) (*model.Writer, error)
	RemoteCreateWriter(body *model.Writer, args Args) (*model.Writer, error)
	RemoteDeleteWriter(uuid string, args Args) (bool, error)
	RemoteEditWriter(uuid string, body *model.Writer, updateProducer bool, args Args) (*model.Writer, error)
}

type RemoteAPI struct {
	DB RemoteDatabase
}

// FLOW NETWORK CLONES

func (j *RemoteAPI) RemoteGetFlowNetworkClones(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetFlowNetworkClones(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetFlowNetworkClone(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetFlowNetworkClone(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteFlowNetworkClone(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteFlowNetworkClone(uuid, args)
	ResponseHandler(q, err, ctx)
}

// NETWORKS

func (j *RemoteAPI) RemoteGetNetworks(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetNetworks(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetNetwork(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetNetwork(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateNetwork(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYNetwork(ctx)
	q, err := j.DB.RemoteCreateNetwork(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteNetwork(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteNetwork(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditNetwork(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYNetwork(ctx)
	q, err := j.DB.RemoteEditNetwork(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// DEVICES

func (j *RemoteAPI) RemoteGetDevices(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetDevices(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetDevice(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetDevice(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateDevice(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYDevice(ctx)
	q, err := j.DB.RemoteCreateDevice(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteDevice(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteDevice(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditDevice(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYDevice(ctx)
	q, err := j.DB.RemoteEditDevice(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// POINTS

func (j *RemoteAPI) RemoteGetPoints(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetPoints(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetPoint(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetPoint(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreatePoint(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYPoint(ctx)
	q, err := j.DB.RemoteCreatePoint(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeletePoint(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeletePoint(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditPoint(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYPoint(ctx)
	q, err := j.DB.RemoteEditPoint(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// STREAMS

func (j *RemoteAPI) RemoteGetStreams(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetStreams(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetStream(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetStream(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateStream(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYStream(ctx)
	q, err := j.DB.RemoteCreateStream(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteStream(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteStream(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditStream(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYStream(ctx)
	q, err := j.DB.RemoteEditStream(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// PRODUCERS

func (j *RemoteAPI) RemoteGetProducers(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetProducers(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetProducer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetProducer(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateProducer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYProducer(ctx)
	q, err := j.DB.RemoteCreateProducer(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteProducer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteProducer(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditProducer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYProducer(ctx)
	q, err := j.DB.RemoteEditProducer(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// CONSUMERS

func (j *RemoteAPI) RemoteGetConsumers(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetProducers(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetConsumer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetConsumer(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateConsumer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYConsumer(ctx)
	q, err := j.DB.RemoteCreateConsumer(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteConsumer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteConsumer(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditConsumer(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYConsumer(ctx)
	q, err := j.DB.RemoteEditConsumer(uuid, body, args)
	ResponseHandler(q, err, ctx)
}

// WRITERS

func (j *RemoteAPI) RemoteGetWriters(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := j.DB.RemoteGetWriters(args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteGetWriter(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteGetWriter(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteCreateWriter(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	body, _ := getBODYWriter(ctx)
	q, err := j.DB.RemoteCreateWriter(body, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteDeleteWriter(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.RemoteDeleteWriter(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *RemoteAPI) RemoteEditWriter(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	uuid := resolveID(ctx)
	body, _ := getBODYWriter(ctx)
	q, err := j.DB.RemoteEditWriter(uuid, body, false, args)
	ResponseHandler(q, err, ctx)
}
