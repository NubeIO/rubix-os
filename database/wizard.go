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
func (d *GormDatabase) WizardLocalPointMapping(body *api.WizardLocalMapping) (bool, error) {
	var flowNetwork model.FlowNetwork
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamModel model.Stream
	var producerModel model.Producer
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	if body.PluginName == "" {
		body.PluginName = "system"
	}

	//get plugin
	p, err := d.GetPluginByPath(body.PluginName)
	if p.UUID == "" {
		return false, errors.New("no valid plugin")
	}

	flowNetwork.IsRemote = false
	flowNetwork.RemoteFlowUUID = utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)
	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	fmt.Println("CreateFlowNetwork", f.UUID)
	if err != nil {
		log.Errorf("wizzrad:  CreateFlowNetwork: %v\n", err)
		return false, err
	}
	// network
	networkModel.PluginConfId = p.UUID
	networkModel.TransportType = model.TransType.IP
	n, err := d.CreateNetwork(&networkModel)
	fmt.Println("CreateNetwork")
	if err != nil {
		log.Errorf("wizzrad:  CreateNetwork: %v\n", err)
		return false, err
	}
	// device
	deviceModel.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(&deviceModel)
	fmt.Println("CreateDevice")
	if err != nil {
		log.Errorf("wizzrad:  CreateDevice: %v\n", err)
		return false, err
	}
	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "is the producer"
	*pointModel.IsProducer = true
	pnt, err := d.CreatePoint(&pointModel, "")
	fmt.Println("CreatePoint")
	if err != nil {
		log.Errorf("wizzrad:  CreatePoint: %v\n", err)
		return false, err
	}
	// stream
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	fmt.Println(streamModel.FlowNetworks, 9898989)
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		log.Errorf("wizzrad:  CreateStream: %v\n", err)
		return false, err
	}
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
	fmt.Println("CreateProducer")
	if err != nil {
		log.Errorf("wizzrad:  CreateProducer: %v\n", err)
		return false, err
	}
	// consumer
	consumerModel.StreamUUID = streamModel.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID = producerModel.UUID
	consumerModel.ProducerThingClass = model.ThingClass.Point
	consumerModel.ProducerThingType = model.ThingClass.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	_, err = d.CreateConsumer(&consumerModel)
	fmt.Println("CreateConsumer")
	if err != nil {
		log.Errorf("wizzrad:  CreateConsumer: %v\n", err)
		return false, err
	}
	// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.Point
	writerModel.ConsumerThingUUID = consumerModel.UUID //itself
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		log.Errorf("wizzrad:  CreateWriter: %v\n", err)
		return false, err
	}
	fmt.Println("CreateWriter")
	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.ThingClass = model.ThingClass.Point
	writerCloneModel.ThingType = model.ThingClass.Point
	writerCloneModel.WriterUUID = writer.UUID
	fmt.Println(writer.UUID, 1, 1, 1, 1)
	writerClone, err := d.CreateWriterClone(&writerCloneModel)
	if err != nil {
		log.Errorf("wizzrad:  CreateWriterClone: %v\n", err)
		return false, err
	}
	fmt.Println("CreateWriterClone")
	writerModel.CloneUUID = writerClone.UUID
	_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
	if err != nil {
		log.Errorf("wizzrad:  UpdateWriter: %v\n", err)
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
	*pointModel.IsProducer = true
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
	*pointModel2.IsConsumer = true
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
	*pointModel.IsProducer = true
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
