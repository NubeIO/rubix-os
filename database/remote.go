package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// FLOW NETWORKS

func (d *GormDatabase) RemoteGetFlowNetworkClones(args api.Args) ([]model.FlowNetworkClone, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetFlowNetworkClones()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetFlowNetworkClone(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteFlowNetworkClone(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteFlowNetworkClone(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

// NETWORKS

func (d *GormDatabase) RemoteGetNetworks(args api.Args) ([]model.Network, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetNetworks()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetNetwork(uuid string, args api.Args) (*model.Network, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	if args.WithPoints {
		resp, err := cli.GetNetworkWithPoints(uuid)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp, err := cli.GetNetwork(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateNetwork(body *model.Network, args api.Args) (*model.Network, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddNetwork(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteNetwork(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteNetwork(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditNetwork(uuid string, body *model.Network, args api.Args) (*model.Network, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditNetwork(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DEVICES

func (d *GormDatabase) RemoteGetDevices(args api.Args) ([]model.Device, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetDevices()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetDevice(uuid string, args api.Args) (*model.Device, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetDevice(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateDevice(body *model.Device, args api.Args) (*model.Device, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddDevice(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteDevice(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteDevice(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditDevice(uuid string, body *model.Device, args api.Args) (*model.Device, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditDevice(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// POINTS

func (d *GormDatabase) RemoteGetPoints(args api.Args) ([]model.Point, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetPoints()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetPoint(uuid string, args api.Args) (*model.Point, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetPoint(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreatePoint(body *model.Point, args api.Args) (*model.Point, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddPoint(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeletePoint(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeletePoint(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditPoint(uuid string, body *model.Point, args api.Args) (*model.Point, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditPoint(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// STREAMS

func (d *GormDatabase) RemoteGetStreams(args api.Args) ([]model.Stream, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetStreams()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetStream(uuid string, args api.Args) (*model.Stream, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetStream(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateStream(body *model.Stream, args api.Args) (*model.Stream, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddStream(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteStream(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteStream(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditStream(uuid string, body *model.Stream, args api.Args) (*model.Stream, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditStream(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// STREAMS CLONES

func (d *GormDatabase) RemoteGetStreamClones(args api.Args) ([]model.StreamClone, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetStreamClones()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteStreamClone(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteStreamClone(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

// PRODUCERS

func (d *GormDatabase) RemoteGetProducers(args api.Args) ([]model.Producer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetProducers()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetProducer(uuid string, args api.Args) (*model.Producer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetProducer(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateProducer(body *model.Producer, args api.Args) (*model.Producer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddProducer(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteProducer(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteProducer(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditProducer(uuid string, body *model.Producer, args api.Args) (*model.Producer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditProducer(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CONSUMERS

func (d *GormDatabase) RemoteGetConsumers(args api.Args) ([]model.Consumer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetConsumers()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetConsumer(uuid string, args api.Args) (*model.Consumer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetConsumer(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateConsumer(body *model.Consumer, args api.Args) (*model.Consumer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.AddConsumer(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteConsumer(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteConsumer(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditConsumer(uuid string, body *model.Consumer, args api.Args) (*model.Consumer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditConsumer(uuid, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WRITERS

func (d *GormDatabase) RemoteGetWriters(args api.Args) ([]model.Writer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetWriters()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteGetWriter(uuid string, args api.Args) (*model.Writer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.GetWriter(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteCreateWriter(body *model.Writer, args api.Args) (*model.Writer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.CreateWriter(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteDeleteWriter(uuid string, args api.Args) (bool, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return false, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.DeleteWriter(uuid)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (d *GormDatabase) RemoteEditWriter(uuid string, body *model.Writer, updateProducer bool, args api.Args) (*model.Writer, error) {
	fn, err := d.GetFlowNetwork(args.FlowNetworkUUID, args)
	if err != nil {
		return nil, err
	}
	cli := client.NewFlowClientCliFromFN(fn)
	resp, err := cli.EditWriter(uuid, body, updateProducer)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
