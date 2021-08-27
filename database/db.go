package database

import (
	"github.com/NubeDev/flow-framework/model"
)



// DropAllFlow networks, gateways, commandGroup, subscriptions, jobs and children.
func (d *GormDatabase) DropAllFlow() (bool, error) {
	//delete networks
	var networkModel *model.Network
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	//delete networks
	var gatewaysModel *model.Stream
	query = d.DB.Where("1 = 1").Delete(&gatewaysModel)
	if query.Error != nil {
		return false, query.Error
	}
	//delete jobs
	var jobModel *model.Job
	query = d.DB.Where("1 = 1").Delete(&jobModel)
	if query.Error != nil {
		return false, query.Error
	}
	//delete subscriptions
	var subscriptionModel *model.Subscription
	query = d.DB.Where("1 = 1").Delete(&subscriptionModel)
	if query.Error != nil {
		return false, query.Error
	}
	//delete commands
	var commandsModel *model.CommandGroup
	query = d.DB.Where("1 = 1").Delete(&commandsModel)
	if query.Error != nil {
		return false, query.Error
	}
	return true, nil
}

//SyncTopics sync all the topics TODO add more
func (d *GormDatabase) SyncTopics()  {

	g, err := d.GetStreamGateways(false)
	for _, obj := range g {
		GetDatabaseBus.RegisterTopicParent(model.CommonNaming.Network, obj.UUID)
	}

	//s, err := d.GetPlugins()
	//for _, obj := range s {
	//	//GetDatabaseBus.RegisterTopicParent(model.CommonNaming.Network, obj.ID)
	//}

	n, err := d.GetNetworks(false, false)
	for _, obj := range n {
		GetDatabaseBus.RegisterTopicParent(model.CommonNaming.Network, obj.UUID)
	}
	de, err := d.GetDevices(false)
	for _, obj := range de {
		GetDatabaseBus.RegisterTopicParent(model.CommonNaming.Network, obj.UUID)
	}
	p, err := d.GetPoints(false)
	for _, obj := range p {
		GetDatabaseBus.RegisterTopicParent(model.CommonNaming.Point, obj.UUID)
	}

	if err != nil {

	}
}
