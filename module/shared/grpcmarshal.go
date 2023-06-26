package shared

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/module/common"
)

type Marshaller interface {
	GetNetwork(uuid, args string) (*model.Network, error)
	GetDevice(uuid, args string) (*model.Device, error)
	GetPoint(uuid, args string) (*model.Point, error)

	GetNetworksByPluginName(pluginName, args string) ([]*model.Network, error)

	GetOneNetworkByArgs(args string) (*model.Network, error)
	GetOneDeviceByArgs(args string) (*model.Device, error)
	GetOnePointByArgs(args string) (*model.Point, error)

	CreateNetwork(body *model.Network) (*model.Network, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	CreatePoint(body *model.Point) (*model.Point, error)

	UpdateNetwork(uuid string, body *model.Network) (*model.Network, error)
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	UpdatePoint(uuid string, body *model.Point) (*model.Point, error)

	UpdateNetworkErrors(uuid string, body *model.Network) error
	UpdateDeviceErrors(uuid string, body *model.Device) error
	UpdatePointErrors(uuid string, body *model.Point) error
	UpdatePointSuccess(uuid string, body *model.Point) error

	DeleteNetwork(uuid string) error
	DeleteDevice(uuid string) error
	DeletePoint(uuid string) error

	PointWrite(uuid string, pointWriter *model.PointWriter) (*common.PointWriteResponse, error)

	GetSchedules() ([]*model.Schedule, error)
	UpdateScheduleAllProps(uuid string, body *model.Schedule) (*model.Schedule, error)
}

type GRPCMarshaller struct {
	DbHelper DBHelper
}

func (g *GRPCMarshaller) GetNetwork(uuid, args string) (*model.Network, error) {
	res, err := g.DbHelper.Get("networks", uuid, args)
	if err != nil {
		return nil, err
	}
	var network *model.Network
	err = json.Unmarshal(res, &network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (g *GRPCMarshaller) GetDevice(uuid, args string) (*model.Device, error) {
	res, err := g.DbHelper.Get("devices", uuid, args)
	if err != nil {
		return nil, err
	}
	var device *model.Device
	err = json.Unmarshal(res, &device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (g *GRPCMarshaller) GetPoint(uuid, args string) (*model.Point, error) {
	res, err := g.DbHelper.Get("points", uuid, args)
	if err != nil {
		return nil, err
	}
	var point *model.Point
	err = json.Unmarshal(res, &point)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (g *GRPCMarshaller) GetNetworksByPluginName(pluginName, args string) ([]*model.Network, error) {
	res, err := g.DbHelper.Get("networks_by_plugin_name", pluginName, args)
	if err != nil {
		return nil, err
	}
	var networks []*model.Network
	err = json.Unmarshal(res, &networks)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

func (g *GRPCMarshaller) GetOneNetworkByArgs(args string) (*model.Network, error) {
	res, err := g.DbHelper.GetWithoutParam("one_network_by_args", args)
	if err != nil {
		return nil, err
	}
	var network *model.Network
	err = json.Unmarshal(res, &network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (g *GRPCMarshaller) GetOneDeviceByArgs(args string) (*model.Device, error) {
	res, err := g.DbHelper.GetWithoutParam("one_device_by_args", args)
	if err != nil {
		return nil, err
	}
	var device *model.Device
	err = json.Unmarshal(res, &device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (g *GRPCMarshaller) GetOnePointByArgs(args string) (*model.Point, error) {
	res, err := g.DbHelper.GetWithoutParam("one_point_by_args", args)
	if err != nil {
		return nil, err
	}
	var point *model.Point
	err = json.Unmarshal(res, &point)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (g *GRPCMarshaller) CreateNetwork(body *model.Network) (*model.Network, error) {
	net, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Post("networks", net)
	if err != nil {
		return nil, err
	}
	var network *model.Network
	err = json.Unmarshal(res, &network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (g *GRPCMarshaller) CreateDevice(body *model.Device) (*model.Device, error) {
	dev, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Post("devices", dev)
	if err != nil {
		return nil, err
	}
	var device *model.Device
	err = json.Unmarshal(res, &device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (g *GRPCMarshaller) CreatePoint(body *model.Point) (*model.Point, error) {
	pnt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Post("points", pnt)
	if err != nil {
		return nil, err
	}
	var point *model.Point
	err = json.Unmarshal(res, &point)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (g *GRPCMarshaller) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	net, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Patch("networks", uuid, net)
	if err != nil {
		return nil, err
	}
	var network *model.Network
	err = json.Unmarshal(res, &network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (g *GRPCMarshaller) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	dev, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Patch("devices", uuid, dev)
	if err != nil {
		return nil, err
	}
	var device *model.Device
	err = json.Unmarshal(res, &device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (g *GRPCMarshaller) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	pnt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Patch("points", uuid, pnt)
	if err != nil {
		return nil, err
	}
	var point *model.Point
	err = json.Unmarshal(res, &point)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (g *GRPCMarshaller) UpdateNetworkErrors(uuid string, body *model.Network) error {
	dev, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = g.DbHelper.Patch("network_errors", uuid, dev)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) UpdateDeviceErrors(uuid string, body *model.Device) error {
	dev, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = g.DbHelper.Patch("device_errors", uuid, dev)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) UpdatePointErrors(uuid string, body *model.Point) error {
	point, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = g.DbHelper.Patch("point_errors", uuid, point)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) UpdatePointSuccess(uuid string, body *model.Point) error {
	point, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = g.DbHelper.Patch("point_success", uuid, point)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) DeleteNetwork(uuid string) error {
	_, err := g.DbHelper.Delete("networks", uuid)
	return err
}

func (g *GRPCMarshaller) DeleteDevice(uuid string) error {
	_, err := g.DbHelper.Delete("devices", uuid)
	return err
}

func (g *GRPCMarshaller) DeletePoint(uuid string) error {
	_, err := g.DbHelper.Delete("points", uuid)
	return err
}

func (g *GRPCMarshaller) PointWrite(uuid string, body *model.PointWriter) (*common.PointWriteResponse, error) {
	pw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Patch("point_write", uuid, pw)
	if err != nil {
		return nil, err
	}
	var pwr *common.PointWriteResponse
	err = json.Unmarshal(res, &pwr)
	if err != nil {
		return nil, err
	}
	return pwr, nil
}

func (g *GRPCMarshaller) GetSchedules() ([]*model.Schedule, error) {
	res, err := g.DbHelper.GetWithoutParam("schedules", "")
	if err != nil {
		return nil, err
	}

	var schedules []*model.Schedule
	if err = json.Unmarshal(res, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (g *GRPCMarshaller) UpdateScheduleAllProps(uuid string, body *model.Schedule) (*model.Schedule, error) {
	sch, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Patch("schedules", uuid, sch)
	if err != nil {
		return nil, err
	}
	var schedule *model.Schedule
	err = json.Unmarshal(res, &schedule)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}
