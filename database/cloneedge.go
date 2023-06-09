package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

var pluginName = "system"

func (d *GormDatabase) CloneEdge(host *model.Host) error {
	cli := client.NewClient(host.IP, host.Port, host.ExternalToken)
	networks, err := cli.GetNetworksForCloneEdge()
	if err != nil {
		return err
	}
	tx := d.DB.Begin()
	_, _ = d.DeleteNetworkClonesByHostUUIDTransaction(tx, host.UUID)
	plugin, err := d.GetPluginByPathTransaction(tx, pluginName)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, network := range networks {
		d.setNetworkModelClone(host.UUID, host.GlobalUUID, network.UUID, plugin.UUID, network)
		for _, device := range network.Devices {
			d.setDeviceModelClone(network.UUID, device.UUID, device)
			for _, point := range device.Points {
				d.setPointModelClone(device.UUID, point.UUID, point)
			}
		}
	}
	_, err = d.CreateBulkNetworksTransaction(tx, networks)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (d *GormDatabase) setNetworkModelClone(hostUUID, globalUUID, networkUUID, pluginConfId string, network *model.Network) {
	network.HostUUID = nstring.New(hostUUID)
	network.GlobalUUID = nstring.New(globalUUID)
	network.SourceUUID = nstring.New(networkUUID)
	network.UUID = nuuid.MakeTopicUUID(model.ThingClass.Network)
	network.SourcePluginName = nstring.New(network.PluginPath)
	network.IsClone = boolean.NewTrue()
	network.PluginPath = pluginName
	network.PluginConfId = pluginConfId
	network.HistoryEnable = boolean.NewFalse()
	for _, metaTag := range network.MetaTags {
		metaTag.NetworkUUID = network.UUID
	}
}

func (d *GormDatabase) setDeviceModelClone(networkUUID string, deviceUUID string, device *model.Device) {
	device.NetworkUUID = networkUUID
	device.SourceUUID = nstring.New(deviceUUID)
	device.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	device.HistoryEnable = boolean.NewFalse()
	for _, metaTag := range device.MetaTags {
		metaTag.DeviceUUID = device.UUID
	}
}

func (d *GormDatabase) setPointModelClone(deviceUUID string, pointUUID string, point *model.Point) {
	point.DeviceUUID = deviceUUID
	point.SourceUUID = nstring.New(pointUUID)
	point.UUID = nuuid.MakeTopicUUID(model.ThingClass.Point)
	point.HistoryEnable = boolean.NewFalse()
	for _, metaTag := range point.MetaTags {
		metaTag.PointUUID = point.UUID
	}
}
