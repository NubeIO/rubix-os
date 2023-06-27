package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	log "github.com/sirupsen/logrus"
)

// SyncTopics sync all the topics
func (d *GormDatabase) SyncTopics() {
	s, err := d.GetPlugins()
	for _, obj := range s {
		d.Bus.RegisterTopicParent(model.CommonNaming.Plugin, obj.UUID)
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
