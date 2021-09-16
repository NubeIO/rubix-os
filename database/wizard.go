package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

// WizardLocalPointMapping add a local network mapping stream.
func (d *GormDatabase) WizardLocalPointMapping() (bool, error) {
	var flowNetwork model.FlowNetwork
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamModel model.Stream
	var producerModel model.Producer
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	//get plugin
	p, err := d.GetPluginByPath("system")
	if p.UUID == "" {
		return false, errors.New("no valid plugin")
	}

	flowNetwork.IsRemote = false
	flowNetwork.RemoteFlowUUID = utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
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
	fmt.Println("CreateDevice")
	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "is the producer"
	pointModel.IsProducer = true
	pnt, err := d.CreatePoint(&pointModel, "")
	fmt.Println("CreatePoint")

	// stream
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	fmt.Println(streamModel.FlowNetworks, 9898989)
	stream, err := d.CreateStream(&streamModel)
	log.Debug("Created Streams at Producer side: ", stream.Name)

	// producer
	log.Debug("stream.UUID: ", stream.UUID)
	producerModel.StreamUUID = stream.UUID

	// producer
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = "producer stream"
	producerModel.ProducerThingClass = model.ThingClass.Point
	producerModel.ProducerThingType = model.ThingClass.Point
	producerModel.ProducerApplication = model.CommonNaming.Mapping
	producer, err := d.CreateProducer(&producerModel)
	fmt.Println(producer.Name)
	fmt.Println("CreateProducer")

	// consumer
	consumerModel.StreamUUID = streamModel.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID = producerModel.UUID
	consumerModel.ProducerThingClass = model.ThingClass.Point
	consumerModel.ProducerThingType = model.ThingClass.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	fmt.Println(consumer.Name)
	fmt.Println("CreateConsumer")
	// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.Point
	writerModel.ConsumerThingUUID = consumerModel.UUID //itself
	writer, err := d.CreateWriter(&writerModel)
	fmt.Println(writer)
	fmt.Println("CreateWriter")
	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.ThingClass = model.ThingClass.Point
	writerCloneModel.ThingType = model.ThingClass.Point
	writerCloneModel.WriterUUID = writer.UUID
	fmt.Println(writer.UUID, 1, 1, 1, 1)
	writerClone, err := d.CreateWriterClone(&writerCloneModel)
	fmt.Println(writerClone)
	fmt.Println("CreateWriterClone")
	writerModel.CloneUUID = writerClone.UUID
	_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
	if err != nil {
		return false, err
	}
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
	var flowNetwork model.FlowNetwork
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamModel model.Stream
	var producerModel model.Producer
	var consumerFlowNetwork model.FlowNetwork
	var consumerStreamModel model.Stream
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	//get plugin
	p, err := d.GetPluginByPath("system")
	if err != nil {
		return false, errors.New("not valid plugin found")
	}

	//in writer add writeCloneUUID and same in writerClone
	flowNetwork.IsRemote = true
	flowNetwork.FlowIP = "0.0.0.0"
	flowNetwork.FlowPort = "1660"
	flowNetwork.RemoteFlowUUID = "ID-" + utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)

	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	log.Debug("Created a FlowNetwork: ", f.UUID)

	// network
	networkModel.PluginConfId = p.UUID
	n, err := d.CreateNetwork(&networkModel)
	log.Debug("Created a Network")
	// device
	deviceModel.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(&deviceModel)
	log.Debug("Created a Device")
	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "is the producer"
	pointModel.IsProducer = true
	pnt, err := d.CreatePoint(&pointModel, "")
	log.Debug("Created a Point")

	// stream
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	stream, err := d.CreateStream(&streamModel)
	log.Debug("Created Streams at Producer side: ", stream.Name)

	// producer
	log.Debug("stream.UUID: ", stream.UUID)
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = "producer stream"
	producerModel.ProducerThingClass = model.ThingClass.Point
	producerModel.ProducerThingType = model.ThingClass.Point
	producerModel.ProducerApplication = model.CommonNaming.Mapping
	producer, err := d.CreateProducer(&producerModel)
	log.Debug("Created Producer: ", producer.Name)

	consumerFlowNetwork.Name = "Consumer flow network"
	consumerFlowNetwork.IsRemote = true
	consumerFlowNetwork.FlowIP = "0.0.0.0"
	consumerFlowNetwork.FlowPort = "1660"
	consumerFlowNetwork.FlowToken = "fakeToken123"
	cfn, err := d.CreateFlowNetwork(&consumerFlowNetwork)
	log.Debug("Created Consumer FlowNetwork: ", cfn.UUID)

	// consumer stream (edge-2)
	consumerStreamModel.IsConsumer = true
	consumerStreamModel.FlowNetworks = []*model.FlowNetwork{&consumerFlowNetwork}
	consumerStream, err := d.CreateStream(&consumerStreamModel)

	// consumer (edge-2)
	consumerModel.StreamUUID = consumerStream.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID = producerModel.UUID
	consumerModel.ProducerThingClass = model.ThingClass.Point
	consumerModel.ProducerThingType = model.ThingClass.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	log.Debug("Created Consumer: ", consumer.Name)

	// device to be used for consumer list (edge-2)
	deviceModel.NetworkUUID = n.UUID
	dev2, err := d.CreateDevice(&deviceModel)
	if err != nil {
		return false, err
	}

	// point 2 to add to consumer list (edge-2)
	var pointModel2 model.Point
	pointModel2.DeviceUUID = dev2.UUID
	pointModel2.Name = "is the consumer"
	pointModel2.IsConsumer = true
	pnt2, err := d.CreatePoint(&pointModel2, "")
	if err != nil {
		return false, err
	}

	// writer (edge-2)
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = pnt2.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.Point
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {

		return false, err
	}

	// add consumer to the writerClone (edge-1)
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.WriterUUID = writerModel.UUID
	writerCloneModel.ThingClass = model.ThingClass.Point
	writerCloneModel.ThingType = model.ThingClass.Point
	writerClone, err := d.CreateWriterClone(&writerCloneModel)
	if err != nil {
		return false, err
	}

	// Update write_clone_uuid on consumer side (edge-2)
	writerModel.CloneUUID = writerClone.UUID
	_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
	if err != nil {
		return false, err
	}
	log.Debug("Updated write_clone_uuid on consumer side (edge-2): ", writerCloneModel.UUID)
	writerCloneModel.WriterUUID = writer.UUID
	_, err = d.UpdateWriterClone(writerCloneModel.UUID, &writerCloneModel, false)
	if err != nil {
		return false, err
	}
	return true, nil
}

//add a new flow network
// need existing streamListUUID
// need an existing point and producerUUID
// add a new consumer
// add a new point
// add a new writer

// Wizard2ndFlowNetwork add a local network mapping stream.
func (d *GormDatabase) Wizard2ndFlowNetwork(body *api.AddNewFlowNetwork) (bool, error) {
	//delete networks
	var flowNetwork model.FlowNetwork
	var consumerModel model.Consumer
	var streamModel model.Stream
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	isRemote := true
	url := "165.227.72.56" //165.227.72.56
	token := "fakeToken123"

	//in writer add writeCloneUUID and same in writerClone
	flowNetwork.IsRemote = isRemote
	flowNetwork.FlowIP = url
	flowNetwork.FlowPort = "1660"
	flowNetwork.FlowToken = token

	flowNetwork.RemoteFlowUUID = "ID-" + utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)

	flowNetwork.Name = "NAME 2nd network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	if err != nil {
		fmt.Println("Error on wizard CreateFlowNetwork")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	fmt.Println("CreateFlowNetwork", f.UUID)

	// consumer stream (edge-2)
	streamModel.IsConsumer = true
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	consumerStream, err := d.CreateStream(&streamModel)
	if err != nil {
		fmt.Println("Error on wizard CreateStream")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	fmt.Println(consumerStream.Name, consumerStream.Name)

	// consumer
	consumerModel.StreamUUID = consumerStream.UUID
	consumerModel.Name = "consumer-2"
	consumerModel.ProducerUUID = body.ProducerUUID
	consumerModel.ProducerThingClass = body.ProducerThingClass
	consumerModel.ProducerThingType = body.ProducerThingType
	consumerModel.ProducerThingUUID = body.ProducerThingUUID
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		fmt.Println("Error on wizard CreateConsumer")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	fmt.Println(consumer.Name, consumer.UUID)

	// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = consumerModel.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		fmt.Println("Error on wizard CreateWriter")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	fmt.Println(writer.WriterThingClass, writer.UUID)

	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = body.ProducerUUID
	writerCloneModel.WriterUUID = writer.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API

	if !isRemote {
		writerClone, err := d.CreateWriterClone(&writerCloneModel)
		if err != nil {
			fmt.Println("Error on wizard CreateWriterClone")
			fmt.Println(err)
			fmt.Println("Error on wizard")
			return false, err
		}
		fmt.Println(writerClone)
		//update writerCloneUUID to writer
		writerModel.CloneUUID = writerClone.UUID
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			fmt.Println("Error on wizard UpdateWriter")
			fmt.Println(err)
			fmt.Println("Error on wizard")
			return false, err
		}
		fmt.Println(writer)
	} else {
		fmt.Println(writerCloneModel.ProducerUUID, writerCloneModel.WriterUUID, "rest add writerClone")
		ap := client.NewSessionWithToken(token, url, "1660")
		clone, err := ap.CreateWriterClone(writerCloneModel)
		if err != nil {
			fmt.Println("Error on wizard CreateWriterClone", err)
			return false, err
		}
		fmt.Println(clone.UUID, clone.ProducerUUID)
		writerModel.CloneUUID = clone.UUID
		fmt.Println(writerModel.CloneUUID, writerModel.UUID, "rest add EditWriter")
		_, err = ap.EditWriter(writerModel.UUID, writerModel, false)
		if err != nil {
			fmt.Println("Error on wizard EditWriter", err)
			return false, err
		}

	}

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) NodeWizard() (bool, error) {
	//delete networks
	var nm1 model.Node
	nm1.Name = "NODE-1"
	n1, err := d.CreateNode(&nm1)

	var nm2 model.Node
	nm2.Name = "NODE-2"
	n2, err := d.CreateNode(&nm2)

	var nm3 model.Node
	nm3.Name = "NODE-3"
	n3, err := d.CreateNode(&nm3)

	var nm4 model.Node
	nm4.Name = "NODE-4"
	n4, err := d.CreateNode(&nm4)

	var nm5 model.Node
	nm5.Name = "NODE-5"
	nm5.NodeType = "add"
	n5, err := d.CreateNode(&nm5)

	var nm6 model.Node
	nm6.Name = "NODE-6"
	nm6.NodeType = "add"
	n6, err := d.CreateNode(&nm6)

	var out1m model.Out1Connections
	out1m.UUID = utils.MakeTopicUUID("")
	out1m.NodeUUID = n1.UUID
	out1m.ToUUID = n2.UUID
	out1m.Connection = "in1"

	query := d.DB.Create(out1m)
	if query.Error != nil {
		return false, query.Error
	}

	var out2m model.Out1Connections
	out2m.UUID = utils.MakeTopicUUID("")
	out2m.NodeUUID = n1.UUID
	out2m.ToUUID = n2.UUID
	out2m.Connection = "in2"

	query = d.DB.Create(out2m)
	if query.Error != nil {
		return false, query.Error
	}

	var out3m model.Out1Connections
	out3m.UUID = utils.MakeTopicUUID("")
	out3m.NodeUUID = n2.UUID
	out3m.ToUUID = n3.UUID
	out3m.Connection = "in1"

	query = d.DB.Create(out3m)
	if query.Error != nil {
		return false, query.Error
	}

	//out of 4 goes into node-3 in-1
	var out4m model.Out1Connections
	out4m.UUID = utils.MakeTopicUUID("")
	out4m.NodeUUID = n4.UUID
	out4m.ToUUID = n3.UUID
	out4m.Connection = "in1"

	query = d.DB.Create(out4m)
	if query.Error != nil {
		return false, query.Error
	}

	//out of 5 goes into node-4 in-1
	var out5m model.Out1Connections
	out5m.UUID = utils.MakeTopicUUID("")
	out5m.NodeUUID = n5.UUID
	out5m.ToUUID = n4.UUID
	out5m.Connection = "in1"

	query = d.DB.Create(out5m)
	if query.Error != nil {
		return false, query.Error
	}

	//out of 3 goes into node-6 in-1
	var out6m model.Out1Connections
	out6m.UUID = utils.MakeTopicUUID("")
	out6m.NodeUUID = n3.UUID
	out6m.ToUUID = n6.UUID
	out6m.Connection = "in1"

	query = d.DB.Create(out6m)
	if query.Error != nil {
		return false, query.Error
	}

	var in1m model.In1Connections
	in1m.UUID = utils.MakeTopicUUID("")
	in1m.NodeUUID = n2.UUID
	in1m.FromUUID = n1.UUID
	in1m.Connection = "out1"

	query = d.DB.Create(in1m)
	if query.Error != nil {
		return false, query.Error
	}

	var in2m model.In1Connections
	in2m.UUID = utils.MakeTopicUUID("")
	in2m.NodeUUID = n3.UUID
	in2m.FromUUID = n2.UUID
	in2m.Connection = "out1"

	query = d.DB.Create(in2m)
	if query.Error != nil {
		return false, query.Error
	}

	var in3m model.In1Connections
	in3m.UUID = utils.MakeTopicUUID("")
	in3m.NodeUUID = n3.UUID
	in3m.FromUUID = n4.UUID
	in3m.Connection = "out1"

	query = d.DB.Create(in3m)
	if query.Error != nil {
		return false, query.Error
	}

	var in4m model.In1Connections
	in4m.UUID = utils.MakeTopicUUID("")
	in4m.NodeUUID = n4.UUID
	in4m.FromUUID = n5.UUID
	in4m.Connection = "out1"

	query = d.DB.Create(in4m)
	if query.Error != nil {
		return false, query.Error
	}

	//node-6 in1 from node-3 out1
	var in5m model.In1Connections
	in5m.UUID = utils.MakeTopicUUID("")
	in5m.NodeUUID = n6.UUID
	in5m.FromUUID = n3.UUID
	in5m.Connection = "out1"

	query = d.DB.Create(in5m)
	if query.Error != nil {
		return false, query.Error
	}

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}


// NetworkDevicePoint add a local network mapping stream.
func (d *GormDatabase) NetworkDevicePoint() (bool, error) {

	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point

	//get plugin
	p, err := d.GetPluginByPath("bacnetserver")
	if p.UUID == "" {
		return false, errors.New("no valid plugin")
	}

	// network
	networkModel.PluginConfId = p.UUID
	n, err := d.CreateNetwork(&networkModel)
	fmt.Println("CreateNetwork")
	// device
	deviceModel.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(&deviceModel)
	fmt.Println("CreateDevice")
	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.IsProducer = true
	_, err = d.CreatePoint(&pointModel, "")

	fmt.Println("CreatePoint")
	if err != nil {
		return false, err
	}
	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}
