package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	log "github.com/sirupsen/logrus"
)

// SyncTopics sync all the topics
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
		log.Error("ERROR sync node topic's at DB")
	}
}
