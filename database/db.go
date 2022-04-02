package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

// DropAllFlow networks, gateways, commandGroup, consumers, jobs and children.
func (d *GormDatabase) DropAllFlow() (string, error) {

	//delete networks
	var networkModel *model.Network
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete networkModel")
		return "fail networkModel", query.Error
	}
	//delete jobs
	var jobModel *model.Job
	query = d.DB.Where("1 = 1").Delete(&jobModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete jobModel")
		return "fail jobModel", query.Error
	}
	//delete producer
	var producerModel *model.Producer
	query = d.DB.Where("1 = 1").Delete(&producerModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Producer")
		return "fail Producer", query.Error
	}
	////delete consumers
	var consumerModel *model.Consumer
	query = d.DB.Where("1 = 1").Delete(&consumerModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Consumer")
		return "fail Consumer", query.Error
	}
	////delete consumersList
	var consumersList *model.Writer
	query = d.DB.Where("1 = 1").Delete(&consumersList)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Writer")
		return "fail Writer", query.Error
	}
	//delete commands
	var commandsModel *model.CommandGroup
	query = d.DB.Where("1 = 1").Delete(&commandsModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete CommandGroup")
		return "fail CommandGroup", query.Error
	}

	//delete streams
	var streamsModel *model.Stream
	query = d.DB.Where("1 = 1").Delete(&streamsModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Stream")
		return "fail Stream", query.Error
	}

	var StreamClone *model.StreamClone
	query = d.DB.Where("1 = 1").Delete(&StreamClone)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Stream")
		return "fail StreamClone", query.Error
	}

	//delete networks
	var flowNetworkModel *model.FlowNetwork
	query = d.DB.Where("1 = 1").Delete(&flowNetworkModel)
	if query.Error != nil {
		log.Error("DB: fail bulk delete FlowNetwork")
		return "fail FlowNetwork", query.Error
	}

	var FlowNetworkClone *model.FlowNetworkClone
	query = d.DB.Where("1 = 1").Delete(&FlowNetworkClone)
	if query.Error != nil {
		log.Error("DB: fail bulk delete FlowNetworkClone")
		return "fail FlowNetworkClone", query.Error
	}

	var Schedule *model.Schedule
	query = d.DB.Where("1 = 1").Delete(&Schedule)
	if query.Error != nil {
		log.Error("DB: fail bulk delete Schedule")
		return "fail Schedule", query.Error
	}

	var integration *model.Integration
	query = d.DB.Where("1 = 1").Delete(&integration)
	if query.Error != nil {
		log.Error("DB: fail bulk delete integration")
		return "fail integration", query.Error
	}

	return "ok", nil
}

//SyncTopics sync all the topics
func (d *GormDatabase) SyncTopics() {

	g, err := d.GetStreams(api.Args{})
	for _, obj := range g {
		d.Bus.RegisterTopicParent(model.CommonNaming.Stream, obj.UUID)
	}
	s, err := d.GetPlugins()
	for _, obj := range s {
		d.Bus.RegisterTopicParent(model.CommonNaming.Plugin, obj.UUID)
	}
	sub, err := d.GetProducers(api.Args{})
	for _, obj := range sub {
		d.Bus.RegisterTopicParent(model.CommonNaming.Producer, obj.UUID)
	}
	rip, err := d.GetConsumers(api.Args{})
	for _, obj := range rip {
		d.Bus.RegisterTopicParent(model.CommonNaming.Consumer, obj.UUID)
	}
	j, err := d.GetJobs()
	for _, obj := range j {
		d.Bus.RegisterTopicParent(model.CommonNaming.Job, obj.UUID)
	}
	n, err := d.GetNetworks(api.Args{})
	for _, obj := range n {
		d.Bus.RegisterTopicParent(model.ThingClass.Network, obj.UUID)
	}
	de, err := d.GetDevices(api.Args{})
	for _, obj := range de {
		d.Bus.RegisterTopicParent(model.ThingClass.Network, obj.UUID)
	}
	p, err := d.GetPoints(api.Args{})
	for _, obj := range p {
		d.Bus.RegisterTopicParent(model.ThingClass.Point, obj.UUID)
	}

	if err != nil {
		log.Error("ERROR sync node topic's at db")
	}
}
