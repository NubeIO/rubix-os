package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) CreateAutoMapping(autoMapping *interfaces.AutoMapping) *interfaces.AutoMappingResponse {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	// level => network (it has all networks)
	//    doesn't exist:
	//       - delete stream_clones (it does cascade delete of consumers and writers)
	//       - delete networks
	//       - delete devices
	//       - delete points
	//    exists:
	//       - stream_clones.enable = <boolean>
	//       - consumers.enable = <boolean>
	//       - network.enable = <boolean>
	//       - device.enable = <boolean>
	//       - point.enable = <boolean>
	//
	// level => device (it has all devices)
	//    doesn't exist:
	//       - delete stream_clones (it does cascade delete of consumers and writers)
	//       - delete streams writers on other side (not needed, coz it already gets deleted)
	//       - delete devices
	//       - delete points
	//    exists:
	//       - stream_clones.enable = <boolean>
	//       - consumers.enable = <boolean>
	//       - device.enable = <boolean>
	//       - point.enable = <boolean>
	//
	// level => point (it has all points)
	//    doesn't exist:
	//       - delete points
	//       - delete consumers (it does cascade delete of writer)
	//       - delete writer_clones on fn side
	//    exists:
	//       - stream_clones.enable = <boolean>
	//       - consumers.enable = <boolean>
	//       - device.enable = <boolean>
	//       - point.enable = <boolean>

	d.cleanAutoMappedModels(tx, autoMapping)
	d.clearStreamClonesAndConsumers(tx)

	var syncWriters []*interfaces.SyncWriter
	for _, amNetwork := range autoMapping.Networks {
		amRes := d.createNetworkAutoMapping(tx, amNetwork, autoMapping.FlowNetworkUUID, autoMapping.GlobalUUID)
		if amRes.HasError {
			tx.Rollback()
			return amRes
		}
		syncWriters = append(syncWriters, amRes.SyncWriters...)
	}

	tx.Commit()

	d.PublishPointsList("")

	return &interfaces.AutoMappingResponse{
		HasError:    false,
		SyncWriters: syncWriters,
	}
}

func (d *GormDatabase) cleanAutoMappedModels(tx *gorm.DB, autoMapping *interfaces.AutoMapping) {
	// delete those which is not deleted when we delete edge
	// level => network (it has all networks)
	//    doesn't exist:
	//       - delete stream_clones (it does cascade delete of consumers and writers)
	//       - delete networks
	//       - delete devices
	//       - delete points
	var edgeNetworks []string
	var edgeDevices []string
	var edgePoints []string
	for _, amNetwork := range autoMapping.Networks {
		edgeNetworks = append(edgeNetworks, amNetwork.UUID)
		for _, amDevice := range amNetwork.Devices {
			edgeDevices = append(edgeDevices, amDevice.UUID)
			for _, amPoints := range amDevice.Points {
				edgePoints = append(edgePoints, amPoints.UUID)
			}
		}
	}

	networks, _ := d.GetNetworks(api.Args{GlobalUUID: &autoMapping.GlobalUUID, WithDevices: true, WithPoints: true})
	for _, network := range networks {
		if autoMapping.Level == interfaces.Network {
			if boolean.IsTrue(network.CreatedFromAutoMapping) &&
				network.AutoMappingUUID != nil && !nstring.ContainsString(edgeNetworks, *network.AutoMappingUUID) {
				tx.Delete(&network)
			}
		}

		for _, device := range network.Devices {
			if autoMapping.Level == interfaces.Network || autoMapping.Level == interfaces.Device {
				if boolean.IsTrue(device.CreatedFromAutoMapping) &&
					device.AutoMappingUUID != nil && !nstring.ContainsString(edgeDevices, *device.AutoMappingUUID) {
					tx.Delete(&device)
				}
			}

			for _, point := range device.Points {
				if boolean.IsTrue(point.CreatedFromAutoMapping) &&
					point.AutoMappingUUID != nil && !nstring.ContainsString(edgePoints, *point.AutoMappingUUID) {
					tx.Delete(&point)
				}
			}
		}
	}
}

func (d *GormDatabase) clearStreamClonesAndConsumers(tx *gorm.DB) {
	// delete those which is not deleted when we delete network, device & points
	tx.Where("created_from_auto_mapping IS TRUE AND auto_mapping_network_uuid NOT IN (?)",
		tx.Where("created_from_auto_mapping IS TRUE").Model(&model.Network{}).Select("uuid")).
		Delete(&model.StreamClone{})
	tx.Where("created_from_auto_mapping IS TRUE AND auto_mapping_device_uuid NOT IN (?)",
		tx.Where("created_from_auto_mapping IS TRUE").Model(&model.Device{}).Select("uuid")).
		Delete(&model.StreamClone{})
	tx.Where("created_from_auto_mapping IS TRUE AND producer_thing_uuid NOT IN (?)",
		tx.Where("created_from_auto_mapping IS TRUE").Model(&model.Point{}).Select("auto_mapping_uuid")).
		Delete(&model.Consumer{})
}

func (d *GormDatabase) createNetworkAutoMapping(tx *gorm.DB, amNetwork *interfaces.AutoMappingNetwork, fnUUID, globalUUID string) *interfaces.AutoMappingResponse {
	amRes := &interfaces.AutoMappingResponse{
		NetworkUUID: amNetwork.UUID,
		HasError:    true,
		Level:       interfaces.Network,
	}

	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nstring.New(fnUUID)})
	if err != nil {
		amRes.Error = err.Error()
		return amRes
	}

	networkName := getAutoMappedNetworkName(fnc.Name, amNetwork.Name)

	network, err := d.GetNetworkByName(networkName, api.Args{})
	if network != nil {
		if network.GlobalUUID != globalUUID {
			amRes.Error = fmt.Sprintf("network.name %s already exists in fnc side with different global_uuid", network.Name)
			return amRes
		} else if boolean.IsFalse(network.CreatedFromAutoMapping) {
			amRes.Error = fmt.Sprintf("manually created network.name %s already exists in fnc side", network.Name)
			return amRes
		}
	}

	network, _ = d.GetOneNetworkByArgs(api.Args{AutoMappingUUID: nstring.New(amNetwork.UUID), GlobalUUID: nstring.New(globalUUID)})
	if network == nil {
		network = &model.Network{}
		network.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, network, globalUUID)
		network, err = d.CreateNetworkTransaction(tx, network) //todo meta-tags
		if err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	} else {
		network.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, network, globalUUID)
		network, err = d.UpdateNetworkTransaction(tx, network.UUID, network, true) //todo meta-tags
		if err != nil {
			amRes.Error = err.Error()
			return amRes
		}
	}

	var syncWriters []*interfaces.SyncWriter
	for _, amDevice := range amNetwork.Devices {
		amRes.DeviceUUID = amDevice.UUID
		amRes.Level = interfaces.Device

		device, _ := d.GetOneDeviceByArgs(api.Args{AutoMappingUUID: nstring.New(amDevice.UUID)})
		if device == nil {
			device = &model.Device{}
			device.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, device) //todo meta-tags
			if device, err = d.CreateDeviceTransaction(tx, device, true); err != nil {
				amRes.Error = err.Error()
				return amRes
			}
		} else {
			device.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, device) //todo meta-tags
			if device, err = d.UpdateDeviceTransaction(tx, device.UUID, device, true); err != nil {
				amRes.Error = err.Error()
				return amRes
			}
			//_, _ = d.CreateDeviceMetaTags(device.UUID, amDevice.MetaTags)//todo meta-tags
		}

		streamClone, _ := GetOneStreamCloneByArgTransaction(tx, api.Args{SourceUUID: nstring.New(amDevice.StreamUUID)})
		streamCloneName := getAutoMappedStreamName(fnc.Name, amNetwork.Name, amDevice.Name)

		if streamClone == nil {
			streamClone = &model.StreamClone{}
			streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, device, streamClone)
			if err = tx.Create(&streamClone).Error; err != nil {
				return amRes
			}
		} else {
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, device, streamClone)
			if err = tx.Model(&streamClone).Where("uuid = ?", streamClone.UUID).Updates(streamClone).Error; err != nil {
				return amRes
			}
		}

		for _, amPoint := range amDevice.Points {
			amRes.PointUUID = amPoint.UUID
			amRes.Level = interfaces.Point
			point, _ := d.GetOnePointByArgsTransaction(tx, api.Args{AutoMappingUUID: nstring.New(amPoint.UUID)})
			if point == nil {
				point = &model.Point{}
				d.setPointModel(device.UUID, amPoint, point) //todo meta-tags
				if _, err = d.CreatePointTransaction(tx, point, true); err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setPointModel(device.UUID, amPoint, point) //todo meta-tags
				if point, err = d.UpdatePointTransactionForAutoMapping(tx, point.UUID, point); err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			}

			consumer, _ := GetOneConsumerByArgsTransaction(tx, api.Args{ProducerThingUUID: nstring.New(amPoint.UUID)})
			if consumer == nil {
				consumer = &model.Consumer{}
				consumer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Consumer)
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Create(&consumer).Error; err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Model(&consumer).Where("uuid = ?", consumer.UUID).Updates(consumer).Error; err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			}

			writer, _ := GetOneWriterByArgsTransaction(tx, api.Args{WriterThingUUID: nstring.New(point.UUID)})
			if writer == nil {
				writer = &model.Writer{}
				writer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Writer)
				d.setWriterModel(amPoint.Name, point.UUID, consumer.UUID, writer)
				if err = tx.Create(&writer).Error; err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setWriterModel(amPoint.Name, point.UUID, consumer.UUID, writer)
				if err = tx.Model(&writer).Where("uuid = ?", consumer.UUID).Updates(writer).Error; err != nil {
					amRes.Error = err.Error()
					return amRes
				}
			}

			syncWriters = append(syncWriters, &interfaces.SyncWriter{
				ProducerUUID:      amPoint.ProducerUUID,
				WriterUUID:        writer.UUID,
				FlowFrameworkUUID: fnc.SourceUUID,
				PointUUID:         amPoint.UUID,
				PointName:         amPoint.Name,
			})
		}
	}

	amRes_ := d.swapMapperNames(tx, amNetwork, fnc.Name, networkName)
	if amRes_ != nil {
		return amRes_
	}

	return &interfaces.AutoMappingResponse{
		HasError:    false,
		SyncWriters: syncWriters,
	}
}

func (d *GormDatabase) swapMapperNames(db *gorm.DB, amNetwork *interfaces.AutoMappingNetwork, fncName, networkName string) *interfaces.AutoMappingResponse {
	for _, amDevice := range amNetwork.Devices {
		if err := db.Model(&model.StreamClone{}).
			Where("source_uuid = ?", amDevice.StreamUUID).
			Update("name", getAutoMappedStreamName(fncName, amNetwork.Name, amDevice.Name)).
			Error; err != nil {
			return &interfaces.AutoMappingResponse{
				NetworkUUID: amNetwork.UUID,
				DeviceUUID:  amDevice.UUID,
				HasError:    true,
				Error:       err.Error(),
				Level:       interfaces.Device,
			}
		}
	}

	if err := db.Model(&model.Network{}).
		Where("auto_mapping_uuid = ?", amNetwork.UUID).
		Update("name", networkName).
		Error; err != nil {
		return &interfaces.AutoMappingResponse{
			NetworkUUID: amNetwork.UUID,
			HasError:    true,
			Error:       err.Error(),
			Level:       interfaces.Network,
		}
	}

	for _, amDevice := range amNetwork.Devices {
		if err := db.Model(&model.Device{}).
			Where("auto_mapping_uuid = ?", amDevice.UUID).
			Update("name", amDevice.Name).
			Error; err != nil {
			return &interfaces.AutoMappingResponse{
				NetworkUUID: amNetwork.UUID,
				DeviceUUID:  amDevice.UUID,
				HasError:    true,
				Error:       err.Error(),
				Level:       interfaces.Device,
			}
		}

		for _, amPoint := range amDevice.Points {
			pointModel := model.Point{}
			if err := db.Model(&pointModel).
				Where("auto_mapping_uuid = ?", amPoint.UUID).
				Update("name", amPoint.Name).
				Error; err != nil {
				return &interfaces.AutoMappingResponse{
					NetworkUUID: amNetwork.UUID,
					DeviceUUID:  amDevice.UUID,
					PointUUID:   amPoint.UUID,
					HasError:    true,
					Error:       err.Error(),
					Level:       interfaces.Point,
				}
			}

			if err := db.Model(&model.Consumer{}).
				Where("producer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", amPoint.UUID).
				Update("name", amPoint.Name).
				Error; err != nil {
				return &interfaces.AutoMappingResponse{
					NetworkUUID: amNetwork.UUID,
					DeviceUUID:  amDevice.UUID,
					PointUUID:   amPoint.UUID,
					HasError:    true,
					Error:       err.Error(),
					Level:       interfaces.Point,
				}
			}

			writer := model.Writer{}
			if err := db.Model(&writer).
				Where("writer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", pointModel.UUID).
				Update("writer_thing_name", amPoint.Name).
				Error; err != nil {
				return &interfaces.AutoMappingResponse{
					NetworkUUID: amNetwork.UUID,
					DeviceUUID:  amDevice.UUID,
					PointUUID:   amPoint.UUID,
					HasError:    true,
					Error:       err.Error(),
					Level:       interfaces.Point,
				}
			}
		}
	}
	return nil
}

func (d *GormDatabase) setStreamCloneModel(streamCloneName, fncUUID, sourceUUID string, device *model.Device, streamClone *model.StreamClone) {
	streamClone.Name = getTempAutoMappedName(streamCloneName)
	streamClone.Enable = boolean.NewTrue()
	streamClone.SourceUUID = sourceUUID
	streamClone.FlowNetworkCloneUUID = fncUUID
	streamClone.CreatedFromAutoMapping = boolean.NewTrue()
	streamClone.AutoMappingNetworkUUID = nstring.New(device.NetworkUUID)
	streamClone.AutoMappingDeviceUUID = nstring.New(device.UUID)
}

func (d *GormDatabase) setNetworkModel(fnc *model.FlowNetworkClone, amNetwork *interfaces.AutoMappingNetwork, networkModel *model.Network, globalUUID string) {
	networkModel.Enable = boolean.NewTrue()
	networkModel.PluginPath = "system"
	networkModel.GlobalUUID = globalUUID
	networkModel.AutoMappingFlowNetworkName = fnc.Name
	networkModel.CreatedFromAutoMapping = boolean.NewTrue()
	networkModel.AutoMappingUUID = &amNetwork.UUID
	networkModel.Tags = amNetwork.Tags
	networkModel.MetaTags = amNetwork.MetaTags
}

func (d *GormDatabase) setDeviceModel(networkUUID string, amDevice *interfaces.AutoMappingDevice, deviceModel *model.Device) {
	deviceModel.Enable = boolean.NewTrue()
	deviceModel.NetworkUUID = networkUUID
	deviceModel.CreatedFromAutoMapping = boolean.NewTrue()
	deviceModel.AutoMappingUUID = &amDevice.UUID
	deviceModel.Tags = amDevice.Tags
	deviceModel.MetaTags = amDevice.MetaTags
}

func (d *GormDatabase) setPointModel(deviceUUID string, amPoint *interfaces.AutoMappingPoint, pointModel *model.Point) {
	pointModel.Name = getTempAutoMappedName(amPoint.Name)
	pointModel.DeviceUUID = deviceUUID
	pointModel.EnableWriteable = boolean.NewTrue()
	pointModel.CreatedFromAutoMapping = boolean.NewTrue()
	pointModel.AutoMappingUUID = &amPoint.UUID
}

func (d *GormDatabase) setConsumerModel(amPoint *interfaces.AutoMappingPoint, stcUUID, pointName string, consumerModel *model.Consumer) {
	consumerModel.Name = getTempAutoMappedName(pointName)
	consumerModel.Enable = boolean.NewTrue()
	consumerModel.StreamCloneUUID = stcUUID
	consumerModel.ProducerUUID = amPoint.ProducerUUID
	consumerModel.ProducerThingName = pointName
	consumerModel.ProducerThingUUID = amPoint.UUID
	consumerModel.ProducerThingClass = "point"
	consumerModel.CreatedFromAutoMapping = boolean.NewTrue()
}

func (d *GormDatabase) setWriterModel(pointName, pointUUID, consumerUUID string, writer *model.Writer) {
	writer.WriterThingName = pointName
	writer.WriterThingClass = "point"
	writer.WriterThingUUID = pointUUID
	writer.ConsumerUUID = consumerUUID
	writer.CreatedFromAutoMapping = boolean.NewTrue()
}
