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

func (d *GormDatabase) CreateAutoMapping(amNetwork *interfaces.AutoMappingNetwork) *interfaces.AutoMappingResponse {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	amRes := &interfaces.AutoMappingResponse{
		NetworkUUID: amNetwork.UUID,
		HasError:    true,
		Level:       interfaces.Network,
	}

	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nstring.New(amNetwork.FlowNetworkUUID)})
	if err != nil {
		amRes.Error = err.Error()
		return amRes
	}

	networkName := getAutoMappedNetworkName(fnc.Name, amNetwork.Name)

	network, err := d.GetNetworkByName(networkName, api.Args{})
	if network != nil {
		if network.GlobalUUID != amNetwork.GlobalUUID {
			amRes.Error = fmt.Sprintf("network.name %s already exists in fnc side with different global_uuid", network.Name)
			return amRes
		} else if boolean.IsFalse(network.CreatedFromAutoMapping) {
			amRes.Error = fmt.Sprintf("manually created network.name %s already exists in fnc side", network.Name)
			return amRes
		}
	}

	network, _ = d.GetOneNetworkByArgs(api.Args{AutoMappingUUID: nstring.New(amNetwork.UUID), GlobalUUID: nstring.New(amNetwork.GlobalUUID)})
	if network == nil {
		network = &model.Network{}
		network.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, network)
		network, err = d.CreateNetworkTransaction(tx, network) //todo meta-tags
		if err != nil {
			tx.Rollback()
			amRes.Error = err.Error()
			return amRes
		}
	} else {
		network.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, network)
		network, err = d.UpdateNetworkTransaction(tx, network.UUID, network) //todo meta-tags
		if err != nil {
			tx.Rollback()
			amRes.Error = err.Error()
			return amRes
		}
	}

	var syncWriters []*interfaces.SyncWriter
	for _, amDevice := range amNetwork.Devices {
		amRes.DeviceUUID = amDevice.UUID
		amRes.Level = interfaces.Device

		streamClone, _ := d.GetOneStreamCloneByArg(api.Args{SourceUUID: nstring.New(amDevice.StreamUUID)})
		streamCloneName := getAutoMappedStreamName(fnc.Name, amNetwork.Name, amDevice.Name)

		if streamClone == nil {
			streamClone = &model.StreamClone{}
			streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, streamClone)
			if err = tx.Create(&streamClone).Error; err != nil {
				tx.Rollback()
				return amRes
			}
			amDevice.StreamCloneUUID = streamClone.UUID
		} else {
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, streamClone)
			if err = tx.Model(&streamClone).Where("uuid = ?", streamClone.UUID).Updates(streamClone).Error; err != nil {
				tx.Rollback()
				return amRes
			}
			amDevice.StreamCloneUUID = streamClone.UUID
		}

		device, _ := d.GetOneDeviceByArgs(api.Args{AutoMappingUUID: nstring.New(amDevice.UUID)})
		if device == nil {
			device = &model.Device{}
			device.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, device) //todo meta-tags
			if device, err = d.CreateDeviceTransaction(tx, device); err != nil {
				tx.Rollback()
				amRes.Error = err.Error()
				return amRes
			}
		} else {
			device.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, device) //todo meta-tags
			if device, err = d.UpdateDeviceTransaction(tx, device.UUID, device); err != nil {
				tx.Rollback()
				amRes.Error = err.Error()
				return amRes
			}
			//_, _ = d.CreateDeviceMetaTags(device.UUID, amDevice.MetaTags)//todo meta-tags
		}

		for _, amPoint := range amDevice.Points {
			amRes.PointUUID = amPoint.UUID
			amRes.Level = interfaces.Point
			point, _ := d.GetOnePointByArgs(api.Args{AutoMappingUUID: nstring.New(amPoint.UUID)})
			if point == nil {
				point = &model.Point{}
				d.setPointModel(device.UUID, amPoint, point) //todo meta-tags
				if _, err = d.CreatePointTransaction(tx, point); err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setPointModel(device.UUID, amPoint, point) //todo meta-tags
				if point, err = d.UpdatePointTransactionForAutoMapping(tx, point.UUID, point); err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			}

			consumer, _ := d.GetOneConsumerByArgs(api.Args{ProducerThingUUID: nstring.New(amPoint.UUID)})
			if consumer == nil {
				consumer = &model.Consumer{}
				consumer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Consumer)
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Create(&consumer).Error; err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Model(&consumer).Where("uuid = ?", consumer.UUID).Updates(consumer).Error; err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			}

			writer, _ := d.GetOneWriterByArgs(api.Args{WriterThingUUID: nstring.New(point.UUID)})
			if writer == nil {
				writer = &model.Writer{}
				writer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Writer)
				d.setWriterModel(amPoint.Name, point.UUID, consumer.UUID, writer)
				if err = tx.Create(&writer).Error; err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			} else {
				d.setWriterModel(amPoint.Name, point.UUID, consumer.UUID, writer)
				if err = tx.Model(&writer).Where("uuid = ?", consumer.UUID).Updates(writer).Error; err != nil {
					tx.Rollback()
					amRes.Error = err.Error()
					return amRes
				}
			}

			syncWriters = append(syncWriters, &interfaces.SyncWriter{
				ProducerUUID: amPoint.ProducerUUID,
				WriterUUID:   writer.UUID,
				PointUUID:    amPoint.UUID,
				PointName:    amPoint.Name,
			})
		}
	}

	amRes_ := d.swapMapperNames(tx, amNetwork, fnc.Name, networkName)
	if amRes_ != nil {
		tx.Rollback()
		return amRes_
	}

	tx.Commit()

	d.PublishPointsList("")

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

func (d *GormDatabase) setStreamCloneModel(streamCloneName, fncUUID, sourceUUID string, streamClone *model.StreamClone) {
	streamClone.Name = getTempAutoMappedName(streamCloneName)
	streamClone.Enable = boolean.NewTrue()
	streamClone.SourceUUID = sourceUUID
	streamClone.FlowNetworkCloneUUID = fncUUID
	streamClone.CreatedFromAutoMapping = boolean.NewTrue()
}

func (d *GormDatabase) setNetworkModel(fnc *model.FlowNetworkClone, amNetwork *interfaces.AutoMappingNetwork, networkModel *model.Network) {
	networkModel.Enable = boolean.NewTrue()
	networkModel.PluginPath = "system"
	networkModel.GlobalUUID = amNetwork.GlobalUUID
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
	writer.WriterThingName = getTempAutoMappedName(pointName)
	writer.WriterThingClass = "point"
	writer.WriterThingUUID = pointUUID
	writer.ConsumerUUID = consumerUUID
	writer.CreatedFromAutoMapping = boolean.NewTrue()
}
