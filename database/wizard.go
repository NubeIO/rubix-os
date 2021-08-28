package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

// WizardLocalPointMapping add a local network mapping stream.
func (d *GormDatabase) WizardLocalPointMapping() (bool, error) {
	//delete networks
	var flowNetwork model.FlowNetwork
	//var pluginModel *model.PluginConf
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamModel model.Stream
	var subscriberModel model.Subscriber
	var subscriptionModel model.Subscription
	var subscriptionListModel model.SubscriptionList
	var subscriberListModel model.SubscriberList

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

	flowNetwork.IsRemote = true
	flowNetwork.RemoteUUID =  utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	fmt.Println("CreateFlowNetwork")
	// network
	networkModel.PluginConfId = p.UUID
	n, err := d.CreateNetwork(&networkModel)
	fmt.Println("CreateNetwork")
	// device
	deviceModel.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(&deviceModel)

	// point
	pointModel.DeviceUUID = dev.UUID
	pnt, err := d.CreatePoint(&pointModel)

	// stream
	streamModel.FlowNetworkUUID = f.UUID
	stream, err := d.CreateStreamGateway(&streamModel)

	// subscriber
	subscriberModel.StreamUUID = stream.UUID
	subscriberModel.FromThingUUID = pnt.UUID
	subscriberModel.Name = "subscriber stream"
	subscriber, err := d.CreateSubscriber(&subscriberModel)
	fmt.Println(subscriber.Name)

	// subscription stream
	streamModel.IsSubscription = true
	streamModel.FlowNetworkUUID = f.UUID
	streamSubscription, err := d.CreateStreamGateway(&streamModel)

	// subscription
	subscriptionModel.StreamUUID = streamSubscription.UUID
	subscriptionModel.Name = "subscription stream"
	subscription, err := d.CreateSubscription(&subscriptionModel)
	fmt.Println(subscription.Name)

	// device to be used for subscription list
	deviceModel.NetworkUUID = n.UUID
	dev2, err := d.CreateDevice(&deviceModel)

	// point 2 to add to subscription list
	var pointModel2 model.Point
	pointModel2.DeviceUUID = dev2.UUID
	pnt2, err := d.CreatePoint(&pointModel2)

	// subscription
	subscriptionListModel.SubscriptionUUID = subscriptionModel.UUID
	subscriptionListModel.ToThingUUID = pnt2.UUID
	subscriptionList, err := d.CreateSubscriptionList(&subscriptionListModel)
	fmt.Println(subscriptionList)

	// add subscription to the subscriberList
	subscriberListModel.SubscriberUUID = subscriber.UUID
	subscriberListModel.FromThingUUID = pnt2.UUID
	subscriberList, err := d.CreateSubscriberList(&subscriberListModel)
	fmt.Println(subscriberList)

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}
