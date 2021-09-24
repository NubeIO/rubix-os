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
//make a flow network, stream network/device/point
//use the point uuid to make a producer
//use the writerWizard to make consumer, writer and writer clone api.WriterWizard
func (d *GormDatabase) WizardLocalPointMapping(body *api.WizardLocalMapping) (bool, error) {
	var flowNetwork model.FlowNetwork
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var streamModel model.Stream
	var producerModel model.Producer

	var writerWizard api.WriterWizard

	if body.PluginName == "" {
		body.PluginName = "system"
	}

	//get plugin
	p, err := d.GetPluginByPath(body.PluginName)
	if p.UUID == "" {
		return false, errors.New("no valid plugin")
	}

	flowNetwork.IsRemote = utils.NewFalse()
	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	fmt.Println("CreateFlowNetwork", f.UUID)
	if err != nil {
		log.Errorf("wizard:  CreateFlowNetwork: %v\n", err)
		return false, err
	}
	// network
	networkModel.PluginConfId = p.UUID
	networkModel.TransportType = model.TransType.IP
	n, err := d.CreateNetwork(&networkModel)
	fmt.Println("CreateNetwork")
	if err != nil {
		log.Errorf("wizard:  CreateNetwork: %v\n", err)
		return false, err
	}
	// device
	deviceModel.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(&deviceModel)
	fmt.Println("CreateDevice")
	if err != nil {
		log.Errorf("wizard:  CreateDevice: %v\n", err)
		return false, err
	}
	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "is the producer"
	pointModel.IsProducer = utils.NewTrue()
	pnt, err := d.CreatePoint(&pointModel, "")
	fmt.Println("CreatePoint")
	if err != nil {
		log.Errorf("wizard:  CreatePoint: %v\n", err)
		return false, err
	}
	// stream
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		log.Errorf("wizard:  CreateStream: %v\n", err)
		return false, err
	}
	log.Debug("Created Streams at Producer side: ", stream.Name)
	log.Debug("stream.UUID: ", stream.UUID)
	producerModel.StreamUUID = stream.UUID

	// producer
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pnt.UUID
	producerModel.Name = pnt.Name
	producerModel.ProducerThingClass = pnt.ThingClass
	producerModel.ProducerThingType = pnt.ThingType
	producerModel.ProducerApplication = model.CommonNaming.Mapping
	producer, err := d.CreateProducer(&producerModel)
	fmt.Println("CreateProducer")
	if err != nil {
		log.Errorf("wizard:  CreateProducer: %v\n", err)
		return false, err
	}
	writerWizard.ConsumerFlowUUID = flowNetwork.UUID
	writerWizard.ConsumerStreamUUID = stream.UUID
	writerWizard.ProducerUUID = producer.UUID
	_, err = d.CreateWriterWizard(&writerWizard)
	fmt.Println("CreateConsumer")
	if err != nil {
		log.Errorf("wizard:  CreateWriterWizard: %v\n", err)
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
	flowNetwork.IsRemote = utils.NewFalse()
	//flowNetwork.RemoteFlowUUID = "ID-" + utils.MakeTopicUUID(model.CommonNaming.RemoteFlowNetwork)

	flowNetwork.Name = "flow network"
	f, err := d.CreateFlowNetwork(&flowNetwork)
	if err != nil {
		return false, err
	}
	log.Debug("Created a FlowNetwork: ", f.UUID)

	// network
	networkModel.PluginConfId = p.UUID
	networkModel.TransportType = "ip"
	n, err := d.CreateNetwork(&networkModel)
	log.Debug("Created a Network")
	// device
	deviceModel.NetworkUUID = n.UUID
	deviceModel.TransportType = "serial"
	dev, err := d.CreateDevice(&deviceModel)
	if err != nil {
		return false, err
	}
	log.Debug("Created a Device")

	// point
	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "is the producer"
	pointModel.IsProducer = utils.NewTrue()
	pointModel.ObjectType = "analogInput" //TODO: check
	pnt, err := d.CreatePoint(&pointModel, "")
	if err != nil {
		return false, err
	}
	log.Debug("Created a Point", err, pnt)

	// stream
	streamModel.FlowNetworks = []*model.FlowNetwork{&flowNetwork}
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		return false, err
	}
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
	if err != nil {
		return false, err
	}
	log.Debug("Created Producer: ", producer.Name)

	consumerFlowNetwork.Name = "Consumer flow network"
	consumerFlowNetwork.IsRemote = utils.NewTrue()
	consumerFlowNetwork.FlowIP = "0.0.0.0"
	consumerFlowNetwork.FlowPort = "1660"
	consumerFlowNetwork.FlowToken = "fakeToken123"
	cfn, err := d.CreateFlowNetwork(&consumerFlowNetwork)
	if err != nil {
		return false, err
	}
	log.Debug("Created Consumer FlowNetwork: ", cfn.UUID)

	// consumer stream (edge-2)
	consumerStreamModel.IsConsumer = true
	consumerStreamModel.FlowNetworks = []*model.FlowNetwork{&consumerFlowNetwork}
	consumerStream, err := d.CreateStream(&consumerStreamModel)
	if err != nil {
		return false, err
	}

	// consumer (edge-2)
	consumerModel.StreamUUID = consumerStream.UUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID = producerModel.UUID
	consumerModel.ProducerThingClass = model.ThingClass.Point
	consumerModel.ProducerThingType = model.ThingClass.Point
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = pnt.UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		return false, err
	}
	log.Debug("Created Consumer: ", consumer.Name)
	// writer (edge-2)
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = consumerModel.UUID
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
	url := "0.0.0.0" //165.227.72.56
	//in writer add writeCloneUUID and same in writerClone
	flowNetwork.IsRemote = utils.NewTrue()
	flowNetwork.FlowIP = url
	flowNetwork.FlowPort = "1660"

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
		ap := client.NewSessionWithToken("token", url, "1660")
		clone, err := ap.CreateWriterClone(writerCloneModel)
		if err != nil {
			fmt.Println("Error on wizard CreateWriterClone", err)
			return false, err
		}
		fmt.Println(clone.UUID, clone.ProducerUUID)
		writerModel.CloneUUID = clone.UUID
		fmt.Println(writerModel.CloneUUID, writerModel.UUID, "rest add EditWriter")
		writerModel.CloneUUID = clone.UUID
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			fmt.Println("Error on wizard UpdateWriter")
			fmt.Println(err)
			fmt.Println("Error on wizard")
			return false, err
		}
		fmt.Println(writer)
	}

	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}

// WizardNewNetDevPnt add a local network mapping stream.
func (d *GormDatabase) WizardNewNetDevPnt(plugin string, net *model.Network, dev *model.Device, pnt *model.Point) (bool, error) {

	//get plugin
	p, err := d.GetPluginByPath(plugin)
	if p.UUID == "" {
		return false, errors.New("no valid plugin")
	}
	// network
	net.PluginConfId = p.UUID
	n, err := d.CreateNetwork(net)
	fmt.Println("CreateNetwork")
	// device
	dev.NetworkUUID = n.UUID
	de, err := d.CreateDevice(dev)
	fmt.Println("CreateDevice")
	// point
	pnt.DeviceUUID = de.UUID
	pnt.IsProducer = utils.NewTrue()
	_, err = d.CreatePoint(pnt, "")
	fmt.Println("CreatePoint")
	if err != nil {
		fmt.Println("Error on wizard")
		fmt.Println(err)
		fmt.Println("Error on wizard")
		return false, err
	}
	return true, nil
}
