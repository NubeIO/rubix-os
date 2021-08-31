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
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

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

	// consumer stream
	var streamModel2 model.Stream
	streamModel2.IsConsumer = true
	streamModel2.StreamListUUID = streamList.UUID
	streamConsumer, err := d.CreateStreamGateway(&streamModel2)

	// consumer
	consumerModel.StreamUUID = streamConsumer.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID =producerModel.UUID
	consumerModel.ConsumerType = model.CommonNaming.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	fmt.Println(consumer.Name)

	// device to be used for consumer list
	deviceModel.NetworkUUID = n.UUID
	dev2, err := d.CreateDevice(&deviceModel)

	// point 2 to add to consumer list
	var pointModel2 model.Point
	pointModel2.DeviceUUID = dev2.UUID
	pointModel2.Name = "is the consumer"
	pnt2, err := d.CreatePoint(&pointModel2)

	// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = pnt2.UUID
	writer, err := d.CreateWriter(&writerModel)
	fmt.Println(writer)

	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.ConsumerUUID = pnt2.UUID
	writerClone, err := d.CreateWriterClone(&writerCloneModel)
	fmt.Println(writerClone)

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
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	//get plugin
	p, err := d.GetPluginByPath("system")
	fmt.Println("GetPluginByPath", p.UUID)

	//in writer add writeCloneUUID and same in writerClone

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

	// consumer stream
	var streamModel2 model.Stream
	streamModel2.IsConsumer = true
	streamModel2.StreamListUUID = streamList.UUID
	streamConsumer, err := d.CreateStreamGateway(&streamModel2)

	// consumer
	consumerModel.StreamUUID = streamConsumer.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID =producerModel.UUID
	consumerModel.ConsumerType = model.CommonNaming.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	fmt.Println(consumer.Name)

	// device to be used for consumer list
	deviceModel.NetworkUUID = n.UUID
	dev2, err := d.CreateDevice(&deviceModel)

	// point 2 to add to consumer list
	var pointModel2 model.Point
	pointModel2.DeviceUUID = dev2.UUID
	pointModel2.Name = "is the consumer"
	pnt2, err := d.CreatePoint(&pointModel2)

	// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = pnt2.UUID
	writer, err := d.CreateWriter(&writerModel)
	fmt.Println(writer)

	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.ConsumerUUID = pnt2.UUID
	writerCloneModel.WriterUUID = writerModel.UUID
	writerClone, err := d.CreateWriterClone(&writerCloneModel)
	fmt.Println(writerClone)

	writerModel.WriteCloneUUID = writerClone.UUID
	_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
	fmt.Println(writer)

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}
