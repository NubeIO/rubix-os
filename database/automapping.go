package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

func (d *GormDatabase) CreateNetworkAutoMappings(network *model.Network, level interfaces.Level) error {
	if boolean.IsTrue(network.CreatedFromAutoMapping) {
		err := d.updateNetworkConnectionInCloneSide(network)
		if err != nil {
			return err
		}
	}
	if boolean.IsTrue(network.AutoMappingEnable) {
		amRes := interfaces.AutoMappingResponse{
			NetworkUUID: network.UUID,
			HasError:    true,
			Level:       interfaces.Network,
		}
		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
		if err != nil {
			errMessage := fmt.Sprintf("failed to find flow network with name %s", network.AutoMappingFlowNetworkName)
			network.Connection = connection.Broken.String()
			network.ConnectionMessage = nstring.New(errMessage)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			return errors.New(errMessage)
		} else {
			network.Connection = connection.Connected.String()
			network.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
		}

		var amDevices []*interfaces.AutoMappingDevice
		deviceUUIDToStreamUUIDMap, err := d.createPointAutoMappingStreams(fn, network.Name, network.Devices)
		if err != nil {
			amRes.Error = err.Error()
			d.updateCascadeConnectionError(d.DB, &amRes)
			log.Error(err)
			return err
		}
		for _, device := range network.Devices {
			if boolean.IsTrue(device.AutoMappingEnable) {
				streamUUID, _ := deviceUUIDToStreamUUIDMap[device.UUID]
				pointUUIDToProducerUUIDMap, err := d.createPointsAutoMappingProducers(streamUUID, device.Points)
				if err != nil {
					amRes.DeviceUUID = device.UUID
					amRes.Level = interfaces.Device
					d.updateCascadeConnectionError(d.DB, &amRes)
					log.Error(err)
					return err
				}

				var amPoints []*interfaces.AutoMappingPoint
				for _, point := range device.Points {
					amPoints = append(amPoints, &interfaces.AutoMappingPoint{
						UUID:         point.UUID,
						Name:         point.Name,
						Tags:         point.Tags,
						MetaTags:     point.MetaTags,
						ProducerUUID: pointUUIDToProducerUUIDMap[point.UUID],
					})
				}
				amDevices = append(amDevices, &interfaces.AutoMappingDevice{
					UUID:       device.UUID,
					Name:       device.Name,
					Tags:       device.Tags,
					MetaTags:   device.MetaTags,
					Points:     amPoints,
					StreamUUID: streamUUID,
				})
			}
		}
		cli := client.NewFlowClientCliFromFN(fn)
		deviceInfo, _ := deviceinfo.GetDeviceInfo()
		amNetwork := &interfaces.AutoMappingNetwork{
			GlobalUUID:      deviceInfo.GlobalUUID,
			UUID:            network.UUID,
			Name:            network.Name,
			Tags:            network.Tags,
			MetaTags:        network.MetaTags,
			Devices:         amDevices,
			FlowNetworkUUID: fn.UUID,
		}
		amRes = cli.CreateAutoMapping(amNetwork)
		if amRes.HasError {
			errMsg := fmt.Sprintf("Flow Network Clone side: %s", amRes.Error)
			amRes.Error = errMsg
			d.updateCascadeConnectionError(d.DB, &amRes)
		} else {
			pointUUID, err := d.createWriterClones(amRes.SyncWriters)
			if pointUUID != nil && err != nil {
				amRes.PointUUID = *pointUUID
				amRes.Level = interfaces.Point
				d.updateCascadeConnectionError(d.DB, &amRes)
			}
		}
	}
	return nil
}

func (d *GormDatabase) updateNetworkConnectionInCloneSide(network *model.Network) error {
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: &network.AutoMappingFlowNetworkName})
	if err != nil {
		amRes := interfaces.AutoMappingResponse{
			NetworkUUID: network.UUID,
			HasError:    true,
			Error:       err.Error(),
			Level:       interfaces.Network,
		}
		d.updateCascadeConnectionError(d.DB, &amRes)
		return err
	}

	cli := client.NewFlowClientCliFromFNC(fnc)

	net, connectionErr, _ := cli.GetNetworkV2(*network.AutoMappingUUID)
	if net == nil && connectionErr == nil {
		amRes := interfaces.AutoMappingResponse{
			NetworkUUID: network.UUID,
			HasError:    true,
			Error:       "Its Network creator has been already deleted, manually delete it",
			Level:       interfaces.Network,
		}
		d.updateCascadeConnectionError(d.DB, &amRes)
	} else {
		amRes := interfaces.AutoMappingResponse{
			NetworkUUID: network.UUID,
			HasError:    false,
			Error:       nstring.NotAvailable,
			Level:       interfaces.Network,
		}
		d.updateCascadeConnectionError(d.DB, &amRes)
	}

	d.updateDevicesConnectionsInCloneSide(cli, network)
	return nil
}

func (d *GormDatabase) updateDevicesConnectionsInCloneSide(cli *client.FlowClient, network *model.Network) {
	var wg sync.WaitGroup
	tx := d.DB.Begin()
	for _, device := range network.Devices {
		wg.Add(1)
		go func(device *model.Device, tx *gorm.DB) {
			defer wg.Done()
			dev, connectionErr, _ := cli.GetDeviceV2(*device.AutoMappingUUID)
			if dev == nil && connectionErr == nil {
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: network.UUID,
					DeviceUUID:  device.UUID,
					HasError:    true,
					Error:       "Its Device creator has been already deleted, manually delete it",
					Level:       interfaces.Device,
				}
				d.updateCascadeConnectionError(tx, &amRes)
			} else {
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: network.UUID,
					DeviceUUID:  device.UUID,
					HasError:    false,
					Error:       nstring.NotAvailable,
					Level:       interfaces.Device,
				}
				d.updateCascadeConnectionError(tx, &amRes)

				d.updatePointsConnectionsInCloneSide(tx, cli, network, device)
			}
		}(device, tx)
		wg.Wait()
	}
	tx.Commit()
}

func (d *GormDatabase) updatePointsConnectionsInCloneSide(tx *gorm.DB, cli *client.FlowClient, network *model.Network, device *model.Device) {
	var wg sync.WaitGroup
	for _, point := range device.Points {
		wg.Add(1)
		go func(point *model.Point, tx *gorm.DB) {
			defer wg.Done()
			pnt, connectionErr, _ := cli.GetPointV2(*point.AutoMappingUUID)
			if pnt == nil && connectionErr == nil {
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: network.UUID,
					DeviceUUID:  device.UUID,
					PointUUID:   point.UUID,
					HasError:    true,
					Error:       "Its Point creator has been already deleted, manually delete it",
					Level:       interfaces.Point,
				}
				d.updateCascadeConnectionError(tx, &amRes)
			} else {
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: network.UUID,
					DeviceUUID:  device.UUID,
					PointUUID:   point.UUID,
					HasError:    false,
					Error:       nstring.NotAvailable,
					Level:       interfaces.Point,
				}
				d.updateCascadeConnectionError(tx, &amRes)
			}
		}(point, tx)
	}
	wg.Wait()
}

func (d *GormDatabase) createPointAutoMappingStreams(flowNetwork *model.FlowNetwork, networkName string, devices []*model.Device) (map[string]string, error) {
	deviceUUIDToStreamUUIDMap := map[string]string{}
	tx := d.DB.Begin()
	for _, device := range devices {
		if boolean.IsTrue(device.AutoMappingEnable) {
			streamName := getAutoMappedStreamName(flowNetwork.Name, networkName, device.Name)
			stream, _ := d.GetOneStreamByArgs(api.Args{Name: nils.NewString(streamName)})
			if stream != nil {
				if boolean.IsFalse(stream.CreatedFromAutoMapping) {
					tx.Commit()
					errMsg := fmt.Sprintf("manually created stream_name %s already exists", streamName)
					amRes := interfaces.AutoMappingResponse{
						NetworkUUID: device.NetworkUUID,
						DeviceUUID:  device.UUID,
						HasError:    true,
						Error:       errMsg,
						Level:       interfaces.Device,
					}
					d.updateCascadeConnectionError(d.DB, &amRes)
					return nil, errors.New(errMsg)
				}
			}
			stream, _ = d.GetOneStreamByArgs(api.Args{AutoMappingUUID: nstring.New(device.UUID), WithFlowNetworks: true})
			if stream == nil {
				stream = &model.Stream{}
				stream.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Stream)
				d.setStreamModel(flowNetwork, device, streamName, stream)
				if err := tx.Create(&stream).Error; err != nil {
					tx.Rollback()
					device.Connection = connection.Broken.String()
					errMsg := fmt.Sprintf("create stream: %s", err.Error())
					amRes := interfaces.AutoMappingResponse{
						NetworkUUID: device.NetworkUUID,
						DeviceUUID:  device.UUID,
						HasError:    true,
						Error:       errMsg,
						Level:       interfaces.Device,
					}
					d.updateCascadeConnectionError(d.DB, &amRes)
					return nil, errors.New(errMsg)
				}
				deviceUUIDToStreamUUIDMap[device.UUID] = stream.UUID
			} else {
				device.Connection = connection.Connected.String()
				device.ConnectionMessage = nstring.New(nstring.NotAvailable)
				_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
				if err := tx.Model(&stream).Association("FlowNetworks").Replace([]*model.FlowNetwork{flowNetwork}); err != nil {
					tx.Rollback()
					errMsg := fmt.Sprintf("update flow_networks on stream: %s", err.Error())
					amRes := interfaces.AutoMappingResponse{
						NetworkUUID: device.NetworkUUID,
						DeviceUUID:  device.UUID,
						HasError:    true,
						Error:       errMsg,
						Level:       interfaces.Device,
					}
					d.updateCascadeConnectionError(d.DB, &amRes)
					return nil, err
				}
				d.setStreamModel(flowNetwork, device, streamName, stream)
				if err := tx.Model(&stream).Updates(&stream).Error; err != nil {
					tx.Rollback()
					errMsg := fmt.Sprintf("update stream: %s", err.Error())
					amRes := interfaces.AutoMappingResponse{
						NetworkUUID: device.NetworkUUID,
						DeviceUUID:  device.UUID,
						HasError:    true,
						Error:       errMsg,
						Level:       interfaces.Device,
					}
					d.updateCascadeConnectionError(d.DB, &amRes)
					return nil, err
				}
				deviceUUIDToStreamUUIDMap[device.UUID] = stream.UUID
			}
		} else {
			// todo: disable stream
		}
	}

	// swap back the names
	for _, device := range devices {
		if boolean.IsTrue(device.AutoMappingEnable) {
			streamName := getAutoMappedStreamName(flowNetwork.Name, networkName, device.Name)
			tempStreamName := getTempAutoMappedName(streamName)
			q := tx.Model(&model.Stream{}).Where("name = ?", tempStreamName).Update("name", streamName)
			if q.Error != nil {
				tx.Rollback()
				errMsg := fmt.Sprintf("update stream: %s", q.Error.Error())
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: device.NetworkUUID,
					DeviceUUID:  device.UUID,
					HasError:    true,
					Error:       errMsg,
					Level:       interfaces.Device,
				}
				d.updateCascadeConnectionError(d.DB, &amRes)
				return nil, q.Error
			}
		}
	}
	tx.Commit()
	return deviceUUIDToStreamUUIDMap, nil
}

func (d *GormDatabase) createPointsAutoMappingProducers(streamUUID string, points []*model.Point) (map[string]string, error) {
	pointUUIDToProducerUUIDMap := map[string]string{}
	tx := d.DB.Begin()
	for _, point := range points {
		if boolean.IsTrue(point.AutoMappingEnable) {
			producer, _ := d.GetOneProducerByArgs(api.Args{StreamUUID: nils.NewString(streamUUID), ProducerThingUUID: nils.NewString(point.UUID)})
			if producer == nil {
				producerModel := model.Producer{}
				producerModel.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
				d.setProducerModel(streamUUID, point, &producerModel)
				if err := tx.Create(&producerModel).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				pointUUIDToProducerUUIDMap[point.UUID] = producerModel.UUID
			} else {
				d.setProducerModel(streamUUID, point, producer)
				if err := tx.Save(producer).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				pointUUIDToProducerUUIDMap[point.UUID] = producer.UUID
			}
		}
	}

	// swap back the names
	for _, point := range points {
		if boolean.IsTrue(point.AutoMappingEnable) {
			q := tx.Model(&model.Producer{}).Where("producer_thing_uuid = ?", point.UUID).Update("name", point.Name)
			if q.Error != nil {
				tx.Rollback()
				return nil, q.Error
			}
		}
	}
	tx.Commit()
	return pointUUIDToProducerUUIDMap, nil
}

func (d *GormDatabase) createWriterClones(syncWriters []*interfaces.SyncWriter) (*string, error) {
	tx := d.DB.Begin()
	for _, syncWriter := range syncWriters {
		wc, _ := d.GetOneWriterCloneByArgs(api.Args{SourceUUID: &syncWriter.WriterUUID})
		if wc == nil {
			wc = &model.WriterClone{}
			wc.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setWriterCloneModel(syncWriter, wc)
			if err := tx.Create(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.PointUUID, err
			}
		} else {
			d.setWriterCloneModel(syncWriter, wc)
			if err := tx.Model(&wc).Where("uuid = ?", wc.UUID).Updates(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.PointUUID, err
			}
		}
	}
	tx.Commit()
	return nil, nil
}

func (d *GormDatabase) updateCascadeConnectionError(tx *gorm.DB, amRes *interfaces.AutoMappingResponse) {
	networkModel := model.Network{}
	deviceModel := model.Device{}
	pointModel := model.Point{}

	connection_ := connection.Connected.String()
	if amRes.HasError {
		connection_ = connection.Broken.String()
	}

	switch amRes.Level {
	case interfaces.Point:
		pointModel.Connection = connection_
		pointModel.ConnectionMessage = &amRes.Error
		_ = UpdatePointConnectionErrorsTransaction(tx, amRes.PointUUID, &pointModel)
		fallthrough
	case interfaces.Device:
		deviceModel.Connection = connection_
		deviceModel.ConnectionMessage = &amRes.Error
		_ = UpdateDeviceConnectionErrorsTransaction(tx, amRes.DeviceUUID, &deviceModel)
		fallthrough
	case interfaces.Network:
		networkModel.Connection = connection_
		networkModel.ConnectionMessage = &amRes.Error
		_ = UpdateNetworkConnectionErrorsTransaction(tx, amRes.NetworkUUID, &networkModel)
	}
}

func (d *GormDatabase) setStreamModel(flowNetwork *model.FlowNetwork, device *model.Device, streamName string, streamModel *model.Stream) {
	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	streamModel.Name = getTempAutoMappedName(streamName)
	streamModel.Enable = boolean.NewTrue()
	streamModel.CreatedFromAutoMapping = boolean.NewTrue()
	streamModel.AutoMappingUUID = nstring.New(device.UUID)
}

func (d *GormDatabase) setProducerModel(streamUUID string, point *model.Point, producerModel *model.Producer) {
	producerModel.Name = getTempAutoMappedName(point.Name)
	producerModel.Enable = boolean.NewTrue()
	producerModel.StreamUUID = streamUUID
	producerModel.ProducerThingUUID = point.UUID
	producerModel.ProducerThingName = point.Name
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producerModel.EnableHistory = point.HistoryEnable
	producerModel.HistoryType = point.HistoryType
	producerModel.HistoryInterval = point.HistoryInterval
}

func (d *GormDatabase) setWriterCloneModel(syncWriter *interfaces.SyncWriter, writerClone *model.WriterClone) {
	writerClone.WriterThingName = syncWriter.PointName
	writerClone.WriterThingClass = "point"
	writerClone.WriterThingUUID = syncWriter.PointUUID
	writerClone.ProducerUUID = syncWriter.ProducerUUID
	writerClone.SourceUUID = syncWriter.WriterUUID
}
