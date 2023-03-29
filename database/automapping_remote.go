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

	for _, amDevice := range amNetwork.Devices {
		stc, _ := d.GetOneStreamCloneByArg(api.Args{SourceUUID: nstring.New(amDevice.StreamUUID)})
		streamCloneName := getAutoMappedStreamName(fnc.Name, amNetwork.Name, amDevice.Name)

		if stc == nil {
			streamClone := model.StreamClone{}
			streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, &streamClone)
			if err = tx.Create(&streamClone).Error; err != nil {
				tx.Rollback()
				amError.DeviceUUID = amDevice.UUID
				amError.Level = interfaces.Device
				return amError
			}
		} else {
			d.setStreamCloneModel(streamCloneName, fnc.UUID, amDevice.StreamUUID, stc)
			if err = tx.Model(&stc).Where("uuid = ?", stc.UUID).Updates(stc).Error; err != nil {
				tx.Rollback()
				amError.DeviceUUID = amDevice.UUID
				amError.Level = interfaces.Device
				return amError
			}
		}
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
		networkModel.Name = getAutoMapperName(networkName)
		d.setNetworkModel(fnc, amNetwork, &networkModel)
		network, err = d.CreateNetworkTransaction(tx, &networkModel)
		if err != nil {
			tx.Rollback()
			amError.Error = err.Error()
			return amError
		}
	} else {
		network.Name = getAutoMapperName(networkName)
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
		device, _ := d.GetOneDeviceByArgs(api.Args{AutoMappingUUID: nstring.New(amDevice.UUID)})
		if device == nil {
			deviceModel := model.Device{}
			deviceModel.Name = getAutoMapperName(amDevice.Name)
			d.setDeviceModel(network.UUID, amDevice, &deviceModel)
			if device, err = d.CreateDeviceTransaction(tx, &deviceModel); err != nil {
				tx.Rollback()
				amError.Error = err.Error()
				return amError
			}
		} else {
			device.Name = getAutoMapperName(amDevice.Name)
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
				d.setPointModel(device.UUID, amPoint, &pointModel)
				if _, err = d.CreatePointTransaction(tx, &pointModel); err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			} else {
				d.setPointModel(device.UUID, amPoint, point)
				if _, err = d.UpdatePointTransactionForAutoMapping(tx, point.UUID, point); err != nil {
					tx.Rollback()
					amError.Error = err.Error()
					return amError
				}
			}
		}
	}

	//for _, amDevice := range amNetwork.Devices {
	//	amDeviceError := &interfaces.AutoMappingDeviceError{Name: amDevice.Name}
	//	network, err := d.createPointAutoMappingNetwork(amNetwork)
	//	if err != nil {
	//		amDeviceError.Error = nstring.New(err.Error())
	//		amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
	//		continue
	//	}
	//	device, err := d.createPointAutoMappingDevice(network, amDevice)
	//	if err != nil {
	//		amDeviceError.Error = nstring.New(err.Error())
	//		amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
	//		continue
	//	}
	//	amDeviceError.Consumers = d.createPointAutoMappingConsumers(amDevice)
	//	amDeviceError.Points = d.createPointAutoMappingPoints(network.Name, device.UUID, device.Name, amDevice)
	//	amDeviceError.Writers = d.createPointAutoMappingWriters(network.Name, device.Name, amDevice)
	//	amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
	//}

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
		}
	}
	return nil
}

func (d *GormDatabase) setStreamCloneModel(streamCloneName, fncUUID, sourceUUID string, streamClone *model.StreamClone) {
	streamClone.Name = getAutoMapperName(streamCloneName)
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
	pointModel.Name = getAutoMapperName(amPoint.Name)
	pointModel.DeviceUUID = deviceUUID
	pointModel.CreatedFromAutoMapping = boolean.NewTrue()
	pointModel.AutoMappingUUID = &amPoint.UUID
}
