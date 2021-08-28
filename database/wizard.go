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
	var producerModel model.Producer
	var subscriptionModel model.Subscription
	var subscriptionListModel model.SubscriptionList
	var producerListModel model.SubscriberList

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

	flowNetwork.IsRemote = false
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
	pointModel.Name = "is the producer"
	pnt, err := d.CreatePoint(&pointModel)

	// stream
	streamModel.FlowNetworkUUID = f.UUID
	stream, err := d.CreateStreamGateway(&streamModel)

	// producer
	producerModel.StreamUUID = stream.UUID
	fmt.Println(pnt.UUID)
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = "producer stream"
	producer, err := d.CreateProducer(&producerModel)
	fmt.Println(producer.Name)

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
	pointModel2.Name = "is the subscription"
	pnt2, err := d.CreatePoint(&pointModel2)

	// subscription
	subscriptionListModel.SubscriptionUUID = subscriptionModel.UUID
	subscriptionListModel.ProducerThingUUID = pnt.UUID
	subscriptionList, err := d.CreateSubscriptionList(&subscriptionListModel)
	fmt.Println(subscriptionList)

	// add subscription to the producerList
	producerListModel.ProducerUUID = producer.UUID
	producerListModel.SubscriptionUUID = pnt2.UUID
	producerList, err := d.CreateProducerList(&producerListModel)
	fmt.Println(producerList)

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}
