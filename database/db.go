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
