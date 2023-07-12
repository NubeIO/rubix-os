package shared

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/module/common"
)

type Marshaller interface {
	GetNetwork(uuid string, args argspkg.Args) (*model.Network, error)
	GetDevice(uuid string, args argspkg.Args) (*model.Device, error)
	GetPoint(uuid string, args argspkg.Args) (*model.Point, error)

	GetNetworksByPluginName(pluginName string, args argspkg.Args) ([]*model.Network, error)
	GetNetworkByName(pluginName string, args argspkg.Args) (*model.Network, error)

	GetOneNetworkByArgs(args argspkg.Args) (*model.Network, error)
	GetOneDeviceByArgs(args argspkg.Args) (*model.Device, error)
	GetOnePointByArgs(args argspkg.Args) (*model.Point, error)

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

	GetPlugin(pluginUUID string, args argspkg.Args) (*model.PluginConf, error)
	GetPluginByPath(name string, args argspkg.Args) (*model.PluginConf, error)
	SetErrorsForAllDevicesOnNetwork(networkUUID, message, messageLevel, messageCode string, doPoints bool) error
	ClearErrorsForAllDevicesOnNetwork(networkUUID string, doPoints bool) error
	SetErrorsForAllPointsOnDevice(deviceUUID, message, messageLevel, messageCode string) error
	ClearErrorsForAllPointsOnDevice(deviceUUID string) error
	WizardNewNetworkDevicePoint(plugin string, network *model.Network, device *model.Device, point *model.Point) (bool, error)
	DeviceNameExistsInNetwork(deviceName, networkUUID string) (*model.Device, bool)

	GetHosts(args argspkg.Args) ([]*model.Host, error)
	GetHistoryLogByHostUUID(hostUUID string) (*model.HistoryLog, error)
	CloneEdge(host *model.Host) error
	CreateBulkHistory(histories []*model.History) (bool, error)
	UpdateBulkHistoryLogs(logs []*model.HistoryLog) (bool, error)
}

type GRPCMarshaller struct {
	DbHelper DBHelper
}

func (g *GRPCMarshaller) GetNetwork(uuid string, args argspkg.Args) (*model.Network, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("networks", uuid, serializedArgs)
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

func (g *GRPCMarshaller) GetDevice(uuid string, args argspkg.Args) (*model.Device, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("devices", uuid, serializedArgs)
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

func (g *GRPCMarshaller) GetPoint(uuid string, args argspkg.Args) (*model.Point, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("points", uuid, serializedArgs)
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

func (g *GRPCMarshaller) GetNetworksByPluginName(pluginName string, args argspkg.Args) ([]*model.Network, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("networks_by_plugin_name", pluginName, serializedArgs)
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

func (g *GRPCMarshaller) GetNetworkByName(networkName string, args argspkg.Args) (*model.Network, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("network_by_name", networkName, serializedArgs)
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

func (g *GRPCMarshaller) GetOneNetworkByArgs(args argspkg.Args) (*model.Network, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.GetWithoutParam("one_network_by_args", serializedArgs)
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

func (g *GRPCMarshaller) GetOneDeviceByArgs(args argspkg.Args) (*model.Device, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.GetWithoutParam("one_device_by_args", serializedArgs)
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

func (g *GRPCMarshaller) GetOnePointByArgs(args argspkg.Args) (*model.Point, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.GetWithoutParam("one_point_by_args", serializedArgs)
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

func (g *GRPCMarshaller) GetPlugin(pluginUUID string, args argspkg.Args) (*model.PluginConf, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("plugin_by_id", pluginUUID, serializedArgs)
	if err != nil {
		return nil, err
	}
	var pluginConf *model.PluginConf
	err = json.Unmarshal(res, &pluginConf)
	if err != nil {
		return nil, err
	}
	return pluginConf, nil
}

func (g *GRPCMarshaller) GetPluginByPath(name string, args argspkg.Args) (*model.PluginConf, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.Get("plugin_by_path", name, serializedArgs)
	if err != nil {
		return nil, err
	}
	var pluginConf *model.PluginConf
	err = json.Unmarshal(res, &pluginConf)
	if err != nil {
		return nil, err
	}
	return pluginConf, nil
}

func (g *GRPCMarshaller) SetErrorsForAllDevicesOnNetwork(networkUUID, message, messageLevel, messageCode string, doPoints bool) error {
	err := g.DbHelper.SetErrorsForAll("devices_on_network", networkUUID, message, messageLevel, messageCode, doPoints)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) ClearErrorsForAllDevicesOnNetwork(networkUUID string, doPoints bool) error {
	err := g.DbHelper.ClearErrorsForAll("devices_on_network", networkUUID, doPoints)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) SetErrorsForAllPointsOnDevice(deviceUUID, message, messageLevel, messageCode string) error {
	err := g.DbHelper.SetErrorsForAll("points_on_device", deviceUUID, message, messageLevel, messageCode, false)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) ClearErrorsForAllPointsOnDevice(deviceUUID string) error {
	err := g.DbHelper.ClearErrorsForAll("points_on_device", deviceUUID, false)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) WizardNewNetworkDevicePoint(plugin string, network *model.Network, device *model.Device, point *model.Point) (bool, error) {
	net, err := json.Marshal(network)
	if err != nil {
		return false, err
	}
	dev, err := json.Marshal(device)
	if err != nil {
		return false, err
	}
	pnt, err := json.Marshal(point)
	if err != nil {
		return false, err
	}
	chk, err := g.DbHelper.WizardNewNetworkDevicePoint(plugin, net, dev, pnt)
	return chk, err
}

func (g *GRPCMarshaller) DeviceNameExistsInNetwork(deviceName, networkUUID string) (*model.Device, bool) {
	network, err := g.GetNetwork(networkUUID, argspkg.Args{})
	if err != nil {
		return nil, false
	}
	for _, dev := range network.Devices {
		if dev.Name == deviceName {
			return dev, true
		}
	}
	return nil, false
}

func (g *GRPCMarshaller) GetHosts(args argspkg.Args) ([]*model.Host, error) {
	serializedArgs, err := args.SerializeArgs(args)
	if err != nil {
		return nil, err
	}
	res, err := g.DbHelper.GetWithoutParam("hosts", serializedArgs)
	if err != nil {
		return nil, err
	}
	var hosts []*model.Host
	err = json.Unmarshal(res, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (g *GRPCMarshaller) GetHistoryLogByHostUUID(hostUUID string) (*model.HistoryLog, error) {
	res, err := g.DbHelper.Get("history_log_by_id", hostUUID, "")
	if err != nil {
		return nil, err
	}
	var historyLog *model.HistoryLog
	err = json.Unmarshal(res, &historyLog)
	if err != nil {
		return nil, err
	}
	return historyLog, nil
}

func (g *GRPCMarshaller) CloneEdge(host *model.Host) error {
	hst, err := json.Marshal(host)
	if err != nil {
		return err
	}
	_, err = g.DbHelper.Post("clone_edge", hst)
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCMarshaller) CreateBulkHistory(histories []*model.History) (bool, error) {
	hist, err := json.Marshal(histories)
	if err != nil {
		return false, err
	}
	_, err = g.DbHelper.Post("bulk_history", hist)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *GRPCMarshaller) UpdateBulkHistoryLogs(logs []*model.HistoryLog) (bool, error) {
	histLog, err := json.Marshal(logs)
	if err != nil {
		return false, err
	}
	_, err = g.DbHelper.Patch("bulk_history_logs", "", histLog)
	if err != nil {
		return false, err
	}
	return true, nil
}
