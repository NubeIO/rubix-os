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

func (d *GormDatabase) CreateAutoMapping(amNetwork *interfaces.AutoMappingNetwork) *interfaces.AutoMappingError {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	amError := &interfaces.AutoMappingError{
		NetworkUUID: amNetwork.UUID,
		HasError:    true,
		Level:       interfaces.Network,
	}

	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nstring.New(amNetwork.FlowNetworkUUID)})
	if err != nil {
		amError.Error = err.Error()
		return amError
	}

	networkName := getAutoMappedNetworkName(fnc.Name, amNetwork.Name)

	network, err := d.GetNetworkByName(networkName, api.Args{})
	if network != nil {
		if network.GlobalUUID != amNetwork.GlobalUUID {
			amError.Error = fmt.Sprintf("network.name %s already exists in fnc side with different global_uuid", network.Name)
			return amError
		} else if boolean.IsFalse(network.CreatedFromAutoMapping) {
			amError.Error = fmt.Sprintf("manually created network.name %s already exists in fnc side", network.Name)
			return amError
		}
	}

	network, _ = d.GetOneNetworkByArgs(api.Args{AutoMappingUUID: nstring.New(amNetwork.UUID), GlobalUUID: nstring.New(amNetwork.GlobalUUID)})
	if network == nil {
		networkModel := model.Network{}
		networkModel.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, &networkModel)
		network, err = d.CreateNetworkTransaction(tx, &networkModel)
		if err != nil {
			tx.Rollback()
			amError.Error = err.Error()
			return amError
		}
	} else {
		network.Name = getTempAutoMappedName(networkName)
		d.setNetworkModel(fnc, amNetwork, network)
		network, err = d.UpdateNetworkTransaction(tx, network.UUID, network)
		if err != nil {
			tx.Rollback()
			amError.Error = err.Error()
			return amError
		}
	}

	for _, amDevice := range amNetwork.Devices {
		amError.DeviceUUID = amDevice.UUID
		amError.Level = interfaces.Device

		streamClone, _ := d.GetOneStreamCloneByArg(api.Args{SourceUUID: nstring.New(amDevice.StreamUUID)})
		streamCloneName := getAutoMappedStreamName(fnc.Name, amNetwork.Name, amDevice.Name)

		if streamClone == nil {
			streamClone = &model.StreamClone{}
			streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, streamClone)
			if err = tx.Create(&streamClone).Error; err != nil {
				tx.Rollback()
				return amError
			}
			amDevice.StreamCloneUUID = streamClone.UUID
		} else {
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, streamClone)
			if err = tx.Model(&streamClone).Where("uuid = ?", streamClone.UUID).Updates(streamClone).Error; err != nil {
				tx.Rollback()
				return amError
			}
			amDevice.StreamCloneUUID = streamClone.UUID
		}

		device, _ := d.GetOneDeviceByArgs(api.Args{AutoMappingUUID: nstring.New(amDevice.UUID)})
		if device == nil {
			deviceModel := model.Device{}
			deviceModel.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, &deviceModel)
			if device, err = d.CreateDeviceTransaction(tx, &deviceModel); err != nil {
				tx.Rollback()
				amError.Error = err.Error()
				return amError
			}
		} else {
			device.Name = getTempAutoMappedName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, device)
			if device, err = d.UpdateDeviceTransaction(tx, device.UUID, device); err != nil {
				tx.Rollback()
				amError.Error = err.Error()
				return amError
			}
			//_, _ = d.CreateDeviceMetaTags(device.UUID, amDevice.MetaTags)//todo
		}

		for _, amPoint := range amDevice.Points {
			amError.PointUUID = amPoint.UUID
			amError.Level = interfaces.Point
			point, _ := d.GetOnePointByArgs(api.Args{AutoMappingUUID: nstring.New(amPoint.UUID)})
			if point == nil {
				pointModel := model.Point{}
				d.setPointModel(device.UUID, amPoint, &pointModel) //todo meta-tags
				if _, err = d.CreatePointTransaction(tx, &pointModel); err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			} else {
				d.setPointModel(device.UUID, amPoint, point) //todo meta-tags
				if point, err = d.UpdatePointTransactionForAutoMapping(tx, point.UUID, point); err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			}

			consumer, _ := d.GetOneConsumerByArgs(api.Args{ProducerThingUUID: nstring.New(amPoint.UUID)})
			if consumer == nil {
				consumer = &model.Consumer{}
				consumer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Consumer)
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Create(&consumer).Error; err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			} else {
				d.setConsumerModel(amPoint, streamClone.UUID, amPoint.Name, consumer)
				if err = tx.Model(&consumer).Where("uuid = ?", consumer.UUID).Updates(consumer).Error; err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			}
		}
	}

	mappingError := d.swapMapperNames(tx, amNetwork, fnc.Name, networkName)
	if mappingError != nil {
		tx.Rollback()
		return mappingError
	}
	tx.Commit()
	return &interfaces.AutoMappingError{
		HasError: false,
	}
}

func (d *GormDatabase) swapMapperNames(db *gorm.DB, amNetwork *interfaces.AutoMappingNetwork, fncName, networkName string) *interfaces.AutoMappingError {
	for _, amDevice := range amNetwork.Devices {
		if err := db.Model(&model.StreamClone{}).
			Where("source_uuid = ?", amDevice.StreamUUID).
			Update("name", getAutoMappedStreamName(fncName, amNetwork.Name, amDevice.Name)).
			Error; err != nil {
			return &interfaces.AutoMappingError{
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
		return &interfaces.AutoMappingError{
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
			return &interfaces.AutoMappingError{
				NetworkUUID: amNetwork.UUID,
				DeviceUUID:  amDevice.UUID,
				HasError:    true,
				Error:       err.Error(),
				Level:       interfaces.Device,
			}
		}

		for _, amPoint := range amDevice.Points {
			if err := db.Model(&model.Point{}).
				Where("auto_mapping_uuid = ?", amPoint.UUID).
				Update("name", amPoint.Name).
				Error; err != nil {
				return &interfaces.AutoMappingError{
					NetworkUUID: amNetwork.UUID,
					DeviceUUID:  amDevice.UUID,
					PointUUID:   amPoint.UUID,
					HasError:    true,
					Error:       err.Error(),
					Level:       interfaces.Point,
				}
			}

			if err := db.Model(&model.Consumer{}).
				Where("producer_thing_uuid = ?", amPoint.UUID).
				Update("name", amPoint.Name).
				Error; err != nil {
				return &interfaces.AutoMappingError{
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
}
