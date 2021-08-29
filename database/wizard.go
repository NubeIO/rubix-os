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
	var streamListModel model.StreamList
	var streamModel model.Stream
	var producerModel model.Producer
	var subscriptionModel model.Subscription
	var subscriptionListModel model.SubscriptionList
	var producerListModel model.ProducerSubscriptionList

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

	// make a stream list
	// use the stream_list UUID and add a flow network
	// use the stream_list UUID and add a stream
	// use the plugin name to add a network then add dev/pnt
	// make a producer with pnt uuid
	// make a 2nd point
	// use pnt1 uuid and pnt2 uuid to make a subscription
	// add point2 uuid to the producerList so the producer has a record of who is subscribing to it


	// streamList
	streamList, err := d.CreateStreamList(&streamListModel)

	flowNetwork.IsRemote = false
	flowNetwork.StreamListUUID = streamList.UUID
	flowNetwork.RemoteFlowUUID =  utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.GlobalFlowID =  "ID-" + utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	fmt.Println("CreateFlowNetwork", f.UUID)
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
	//streamModel.StreamListUUID = f.UUID
	streamModel.StreamListUUID = streamList.UUID
	stream, err := d.CreateStreamGateway(&streamModel)

	// producer
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = "producer stream"
	producerModel.ProducerType = model.CommonNaming.Point
	producerModel.ProducerApplication = model.CommonNaming.Mapping
	producer, err := d.CreateProducer(&producerModel)
	fmt.Println(producer.Name)

	// subscription stream
	var streamModel2 model.Stream
	streamModel2.IsSubscription = true
	streamModel2.StreamListUUID = streamList.UUID
	streamSubscription, err := d.CreateStreamGateway(&streamModel2)

	// subscription
	subscriptionModel.StreamUUID = streamSubscription.UUID
	subscriptionModel.Name = "subscription stream"
	subscriptionModel.SubscriptionType = model.CommonNaming.Point
	subscriptionModel.SubscriptionApplication = model.CommonNaming.Mapping
	subscriptionModel.ProducerThingUUID = pnt.UUID
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

	// subscriptionList
	subscriptionListModel.SubscriptionUUID = subscriptionModel.UUID
	subscriptionListModel.SubscriptionThingUUID = pnt2.UUID
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

// WizardRemotePointMapping add a local network mapping stream.
func (d *GormDatabase) WizardRemotePointMapping() (bool, error) {
	//delete networks
	var flowNetwork model.FlowNetwork
	//var pluginModel *model.PluginConf
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamListModel model.StreamList
	var streamModel model.Stream
	var producerModel model.Producer
	var subscriptionModel model.Subscription
	var subscriptionListModel model.SubscriptionList
	var producerListModel model.ProducerSubscriptionList

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

	// make a stream list
	// use the stream_list UUID and add a flow network
	// use the stream_list UUID and add a stream
	// use the plugin name to add a network then add dev/pnt
	// make a producer with pnt uuid
	// make a 2nd point
	// use pnt1 uuid and pnt2 uuid to make a subscription
	// add point2 uuid to the producerList so the producer has a record of who is subscribing to it


	// streamList
	streamList, err := d.CreateStreamList(&streamListModel)

	flowNetwork.IsRemote = true
	flowNetwork.FlowIP = "0.0.0.0"
	flowNetwork.FlowPort = "1660"
	flowNetwork.StreamListUUID = streamList.UUID
	flowNetwork.RemoteFlowUUID =  utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.GlobalFlowID =  "ID-" + utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	fmt.Println("CreateFlowNetwork", f.UUID)
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
	//streamModel.StreamListUUID = f.UUID
	streamModel.StreamListUUID = streamList.UUID
	stream, err := d.CreateStreamGateway(&streamModel)

	// producer
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = "producer stream"
	producerModel.ProducerType = model.CommonNaming.Point
	producerModel.ProducerApplication = model.CommonNaming.Mapping
	producer, err := d.CreateProducer(&producerModel)
	fmt.Println(producer.Name)

	// subscription stream
	var streamModel2 model.Stream
	streamModel2.IsSubscription = true
	streamModel2.StreamListUUID = streamList.UUID
	streamSubscription, err := d.CreateStreamGateway(&streamModel2)

	// subscription
	subscriptionModel.StreamUUID = streamSubscription.UUID
	subscriptionModel.Name = "subscription stream"
	subscriptionModel.SubscriptionType = model.CommonNaming.Point
	subscriptionModel.SubscriptionApplication = model.CommonNaming.Mapping
	subscriptionModel.ProducerThingUUID = pnt.UUID
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

	// subscriptionList
	subscriptionListModel.SubscriptionUUID = subscriptionModel.UUID
	subscriptionListModel.SubscriptionThingUUID = pnt2.UUID
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
