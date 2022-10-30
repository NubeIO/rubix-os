package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

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
