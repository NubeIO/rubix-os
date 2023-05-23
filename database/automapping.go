package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/interfaces/connection"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/deviceinfo"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

func (d *GormDatabase) CreateNetworksAutoMappings(fnName string, networks []*model.Network, level interfaces.Level) error {
	d.clearStreamsAndProducers()

	err := d.createNetworksAutoMappings(fnName, networks, level)
	if err != nil {
		return err
	}

	err = d.updateNetworksConnectionInCloneSide(fnName)
	if err != nil {
		return err
	}

	go d.PublishPointsList("")

	return nil
}

func (d *GormDatabase) createNetworksAutoMappings(fnName string, networks []*model.Network, level interfaces.Level) error {
	doAutoMapping := false
	var amNetworks []*interfaces.AutoMappingNetwork
	fn, fnError := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(fnName)})

	for _, network := range networks {
		if boolean.IsTrue(network.CreatedFromAutoMapping) {
			continue
		}

		// we are sending extra networks, devices, points to make sure whether it's available or not in fn side
		if network.AutoMappingFlowNetworkName != fnName || boolean.IsFalse(network.AutoMappingEnable) {
			var amDevices []*interfaces.AutoMappingDevice
			for _, device := range network.Devices {
				var amPoints []*interfaces.AutoMappingPoint
				for _, point := range device.Points {
					amPoints = append(amPoints, &interfaces.AutoMappingPoint{
						Enable:            boolean.IsTrue(point.Enable),
						AutoMappingEnable: boolean.IsTrue(point.AutoMappingEnable),
						EnableWriteable:   boolean.IsTrue(point.EnableWriteable),
						UUID:              point.UUID,
						Name:              point.Name,
						Tags:              point.Tags,
						MetaTags:          point.MetaTags,
						Priority:          *point.Priority,
					})
				}
				amDevices = append(amDevices, &interfaces.AutoMappingDevice{
					Enable:            boolean.IsTrue(device.Enable),
					AutoMappingEnable: boolean.IsTrue(device.AutoMappingEnable),
					UUID:              device.UUID,
					Name:              device.Name,
					Tags:              device.Tags,
					MetaTags:          device.MetaTags,
					Points:            amPoints,
				})
			}
			amNetwork := &interfaces.AutoMappingNetwork{
				Enable:            boolean.IsTrue(network.Enable),
				AutoMappingEnable: boolean.IsTrue(network.AutoMappingEnable),
				UUID:              network.UUID,
				Name:              network.Name,
				Tags:              network.Tags,
				MetaTags:          network.MetaTags,
				Devices:           amDevices,
				CreateNetwork:     false,
			}
			amNetworks = append(amNetworks, amNetwork)

			if boolean.IsFalse(network.AutoMappingEnable) {
				network.Connection = connection.Connected.String()
				network.ConnectionMessage = nstring.New(nstring.NotAvailable)
				_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			}
			continue
		}

		// if fnError has issue then return that just right away
		if fnError != nil {
			var msg string
			if network.AutoMappingFlowNetworkName == "" {
				msg = fmt.Sprintf("No flow-network has been selected for the enabled auto-mapping network.")
			} else {
				msg = fmt.Sprintf("The flow network with the name '%s' could not be found.", network.AutoMappingFlowNetworkName)
			}
			network.Connection = connection.Broken.String()
			network.ConnectionMessage = nstring.New(msg)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			return errors.New(msg)
		} else {
			network.Connection = connection.Connected.String()
			network.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
		}

		doAutoMapping = true // this is the case where it has auto_mapping creator with valid flow_network

		tx := d.DB.Begin()
		deviceUUIDToStreamUUIDMap, err := createAutoMappingStreamsTransaction(tx, fn, network, network.Devices)
		if err != nil {
			tx.Rollback()
			log.Error(err)
			return err
		}
		tx.Commit()

		var amDevices []*interfaces.AutoMappingDevice
		for _, device := range network.Devices {
			streamUUID, ok := deviceUUIDToStreamUUIDMap[device.UUID]
			if !ok {
				log.Warn("device is auto-mapping already been disabled, we can't go further on depth")
				continue
			}
			pointUUIDToProducerUUIDMap, err := d.createPointsAutoMappingProducers(streamUUID, device.Points)
			if err != nil {
				amRes := interfaces.AutoMappingResponse{
					HasError:    true,
					NetworkUUID: network.UUID,
					DeviceUUID:  device.UUID,
					Level:       interfaces.Device,
				}
				updateCascadeConnectionError(d.DB, &amRes)
				log.Error(err)
				return err
			}

			var amPoints []*interfaces.AutoMappingPoint
			for _, point := range device.Points {
				producerUUID, ok := pointUUIDToProducerUUIDMap[point.UUID]
				if !ok {
					log.Warn("point is auto-mapping already been disabled, we can't go further on depth")
					continue
				}
				amPoints = append(amPoints, &interfaces.AutoMappingPoint{
					Enable:            boolean.IsTrue(point.Enable),
					AutoMappingEnable: boolean.IsTrue(point.AutoMappingEnable),
					EnableWriteable:   boolean.IsTrue(point.EnableWriteable),
					UUID:              point.UUID,
					Name:              point.Name,
					Tags:              point.Tags,
					MetaTags:          point.MetaTags,
					ProducerUUID:      producerUUID,
					Priority:          *point.Priority,
				})
			}
			amDevices = append(amDevices, &interfaces.AutoMappingDevice{
				Enable:            boolean.IsTrue(device.Enable),
				AutoMappingEnable: boolean.IsTrue(device.AutoMappingEnable),
				UUID:              device.UUID,
				Name:              device.Name,
				Tags:              device.Tags,
				MetaTags:          device.MetaTags,
				Points:            amPoints,
				StreamUUID:        streamUUID,
			})
		}
		amNetwork := &interfaces.AutoMappingNetwork{
			Enable:            boolean.IsTrue(network.Enable),
			AutoMappingEnable: boolean.IsTrue(network.AutoMappingEnable),
			UUID:              network.UUID,
			Name:              network.Name,
			Tags:              network.Tags,
			MetaTags:          network.MetaTags,
			Devices:           amDevices,
			CreateNetwork:     true,
		}
		amNetworks = append(amNetworks, amNetwork)
	}

	if !doAutoMapping {
		return nil
	}

	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	autoMapping := &interfaces.AutoMapping{
		GlobalUUID:      deviceInfo.GlobalUUID,
		FlowNetworkUUID: fn.UUID,
		Level:           level,
		Networks:        amNetworks,
	}

	cli := client.NewFlowClientCliFromFN(fn)
	amRes := cli.CreateAutoMapping(autoMapping)
	if amRes.HasError {
		errMsg := fmt.Sprintf("Flow network clone side: %s", amRes.Error)
		log.Error(errMsg)
		amRes.Error = errMsg
		updateCascadeConnectionError(d.DB, &amRes)
	} else {
		for _, amNetwork := range autoMapping.Networks {
			if amNetwork.CreateNetwork { // just update its own network
				err := d.clearConnectionErrorTransaction(d.DB, amNetwork)
				if err != nil {
					return err
				}
			}
		}

		pointUUID, err := d.createWriterClones(amRes.SyncWriters)
		if pointUUID != nil && err != nil {
			amRes.PointUUID = *pointUUID
			amRes.Level = interfaces.Point
			updateCascadeConnectionError(d.DB, &amRes)
		}
	}
	return nil
}

func (d *GormDatabase) clearStreamsAndProducers() {
	// delete those which is not deleted when we delete network, device & points
	d.DB.Where("created_from_auto_mapping IS TRUE AND IFNULL(auto_mapping_schedule_uuid,'') = '' AND "+
		"auto_mapping_network_uuid NOT IN (?)", d.DB.Model(&model.Network{}).Select("uuid")).
		Delete(&model.Stream{})
	d.DB.Where("created_from_auto_mapping IS TRUE AND IFNULL(auto_mapping_schedule_uuid,'') = '' AND "+
		"auto_mapping_device_uuid NOT IN (?)", d.DB.Model(&model.Device{}).Select("uuid")).Delete(&model.Stream{})
	d.DB.Where("created_from_auto_mapping IS TRUE AND producer_thing_class = ? AND producer_thing_uuid NOT IN (?)",
		model.ThingClass.Point, d.DB.Model(&model.Point{}).Select("uuid")).
		Delete(&model.Producer{})
}

func (d *GormDatabase) updateNetworksConnectionInCloneSide(fnName string) error {
	networks, err := d.GetNetworksTransaction(d.DB, api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	for _, network := range networks {
		if boolean.IsTrue(network.CreatedFromAutoMapping) && network.AutoMappingFlowNetworkName == fnName {
			err = d.updateNetworkConnectionInCloneSide(network)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *GormDatabase) updateNetworkConnectionInCloneSide(network *model.Network) error {
	fnc, err := d.GetOneFlowNetworkCloneByArgsTransaction(d.DB, api.Args{Name: &network.AutoMappingFlowNetworkName})
	if err != nil {
		network.Connection = connection.Broken.String()
		network.ConnectionMessage = nstring.New(err.Error())
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
		return err
	}

	cli := client.NewFlowClientCliFromFNC(fnc)

	net, connectionErr, _ := cli.GetNetworkV2(*network.AutoMappingUUID)
	if connectionErr != nil {
		network.Connection = connection.Broken.String()
		network.ConnectionMessage = nstring.New("connection error: " + connectionErr.Error())
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
	} else if net == nil {
		network.Connection = connection.Broken.String()
		network.ConnectionMessage = nstring.New("The network creator has already been deleted. Delete manually if needed.")
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
	} else if boolean.IsFalse(net.AutoMappingEnable) && boolean.IsTrue(network.CreatedFromAutoMapping) {
		// here we use 'net' instead 'network' because network wouldn't get updated in auto-mapped disabled case
		// where on 'devices' and 'points', it is fine to use its own model
		network.Connection = connection.Broken.String()
		msg := fmt.Sprintf("The auto-mapping feature for the network creator '%s' is currently disabled. Delete manually if needed.", net.Name)
		network.ConnectionMessage = nstring.New(msg)
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
	} else if net.AutoMappingFlowNetworkName != network.AutoMappingFlowNetworkName {
		network.Connection = connection.Broken.String()
		var msg string
		if net.AutoMappingFlowNetworkName != "" {
			msg = fmt.Sprintf("The network creator '%s' is attached to a different flow network named '%s'. Delete manually if needed.", net.Name, net.AutoMappingFlowNetworkName)
		} else {
			msg = fmt.Sprintf("The network creator '%s' isn't attached with any flow network. Delete manually if needed.", net.Name)
		}
		network.ConnectionMessage = nstring.New(msg)
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
	} else {
		network.Connection = connection.Connected.String()
		network.ConnectionMessage = nstring.New(nstring.NotAvailable)
		_ = UpdateNetworkConnectionErrorsTransaction(d.DB, network.UUID, network)
	}
	updateDevicesConnectionsInCloneSide(d.DB, cli, network)
	return nil
}

func updateDevicesConnectionsInCloneSide(tx *gorm.DB, cli *client.FlowClient, network *model.Network) {
	var wg sync.WaitGroup
	pointsUUIDs, pointConnectionErr, _ := cli.GetPointsBulkUUIDs()
	for _, device := range network.Devices {
		wg.Add(1)
		go func(device *model.Device, tx *gorm.DB) {
			defer wg.Done()
			dev, connectionErr, _ := cli.GetDeviceV2(*device.AutoMappingUUID)
			if connectionErr != nil {
				device.Connection = connection.Broken.String()
				device.ConnectionMessage = nstring.New("connection error: " + connectionErr.Error())
				_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
			} else if dev == nil {
				device.Connection = connection.Broken.String()
				device.ConnectionMessage = nstring.New("The device creator has already been deleted. Delete manually if needed.")
				_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
			} else if boolean.IsFalse(dev.AutoMappingEnable) && boolean.IsTrue(device.CreatedFromAutoMapping) {
				device.Connection = connection.Broken.String()
				msg := fmt.Sprintf("The auto-mapping feature for the device creator '%s' is currently disabled. Delete manually if needed.", dev.Name)
				device.ConnectionMessage = nstring.New(msg)
				_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
			} else {
				device.Connection = connection.Connected.String()
				device.ConnectionMessage = nstring.New(nstring.NotAvailable)
				_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)

				updatePointsConnectionsInCloneSide(tx, device, pointsUUIDs, pointConnectionErr)
			}
		}(device, tx)
		wg.Wait()
	}

	// update parent with its child error
	var networkModel []*model.Network
	var deviceModel []*model.Device
	tx.Where("uuid = ? AND connection = ?", network.UUID, connection.Connected.String()).Find(&networkModel)
	tx.Where("network_uuid = ? AND connection = ?", network.UUID, connection.Broken.String()).Find(&deviceModel)
	if len(networkModel) > 0 && len(deviceModel) > 0 {
		networkModel[0].Connection = connection.Broken.String()
		networkModel[0].ConnectionMessage = deviceModel[0].ConnectionMessage
		_ = UpdateNetworkConnectionErrorsTransaction(tx, network.UUID, networkModel[0])
	}
}

func updatePointsConnectionsInCloneSide(tx *gorm.DB, device *model.Device, pointsUUIDs *[]string, pointConnectionErr error) {
	for _, point := range device.Points {
		if pointConnectionErr != nil {
			point.Connection = connection.Broken.String()
			point.ConnectionMessage = nstring.New("connection error: " + pointConnectionErr.Error())
		} else if pointsUUIDs == nil || !nstring.ContainsString(*pointsUUIDs, *point.AutoMappingUUID) {
			point.Connection = connection.Broken.String()
			point.ConnectionMessage = nstring.New("The point creator has already been deleted. Delete manually if needed.")
		} else if boolean.IsFalse(point.AutoMappingEnable) && boolean.IsTrue(point.CreatedFromAutoMapping) {
			point.Connection = connection.Broken.String()
			msg := fmt.Sprintf("The auto-mapping feature for the point creator '%s' is currently disabled. Delete manually if needed.", point.Name)
			point.ConnectionMessage = nstring.New(msg)
		} else {
			point.Connection = connection.Connected.String()
			point.ConnectionMessage = nstring.New(nstring.NotAvailable)
		}
		_ = UpdatePointConnectionErrorsTransaction(tx, point.UUID, point)
	}

	// update parent with its child error
	var deviceModel []*model.Device
	var pointModel []*model.Point
	tx.Where("uuid = ? AND connection = ?", device.UUID, connection.Connected.String()).Find(&deviceModel)
	tx.Where("device_uuid = ? AND connection = ?", device.UUID, connection.Broken.String()).Find(&pointModel)
	if len(deviceModel) > 0 && len(pointModel) > 0 {
		deviceModel[0].Connection = connection.Broken.String()
		deviceModel[0].ConnectionMessage = pointModel[0].ConnectionMessage
		_ = UpdateDeviceConnectionErrorsTransaction(tx, deviceModel[0].UUID, deviceModel[0])
	}
}

func createAutoMappingStreamsTransaction(tx *gorm.DB, flowNetwork *model.FlowNetwork, network *model.Network,
	devices []*model.Device) (map[string]string, error) {

	deviceUUIDToStreamUUIDMap := map[string]string{}
	for _, device := range devices {
		streamName := getAutoMappedStreamName(flowNetwork.Name, network.Name, device.Name)
		stream, _ := GetOneStreamByArgsTransaction(tx, api.Args{Name: nils.NewString(streamName)})
		if stream != nil {
			if boolean.IsFalse(stream.CreatedFromAutoMapping) {
				errMsg := fmt.Sprintf("manually created stream_name %s already exists", streamName)
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: device.NetworkUUID,
					DeviceUUID:  device.UUID,
					HasError:    true,
					Error:       errMsg,
					Level:       interfaces.Device,
				}
				updateCascadeConnectionError(tx, &amRes)
				return nil, errors.New(errMsg)
			}
		}

		stream, _ = GetOneStreamByArgsTransaction(tx, api.Args{AutoMappingDeviceUUID: nstring.New(device.UUID), WithFlowNetworks: true})
		if stream == nil {
			if boolean.IsTrue(device.AutoMappingEnable) { // create stream only when auto_mapping is enabled
				stream = &model.Stream{}
				stream.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Stream)
				setStreamModel(flowNetwork, network, device, stream)
				if err := tx.Create(&stream).Error; err != nil {
					device.Connection = connection.Broken.String()
					errMsg := fmt.Sprintf("create stream: %s", err.Error())
					amRes := interfaces.AutoMappingResponse{
						NetworkUUID: device.NetworkUUID,
						DeviceUUID:  device.UUID,
						HasError:    true,
						Error:       errMsg,
						Level:       interfaces.Device,
					}
					updateCascadeConnectionError(tx, &amRes)
					return nil, errors.New(errMsg)
				}
				deviceUUIDToStreamUUIDMap[device.UUID] = stream.UUID
			} else {
				continue
			}
		} else {
			device.Connection = connection.Connected.String()
			device.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
			if err := tx.Model(&stream).Association("FlowNetworks").Replace([]*model.FlowNetwork{flowNetwork}); err != nil {
				errMsg := fmt.Sprintf("update flow_networks on stream: %s", err.Error())
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: device.NetworkUUID,
					DeviceUUID:  device.UUID,
					HasError:    true,
					Error:       errMsg,
					Level:       interfaces.Device,
				}
				updateCascadeConnectionError(tx, &amRes)
				return nil, err
			}
			setStreamModel(flowNetwork, network, device, stream)
			if err := tx.Model(&stream).Updates(&stream).Error; err != nil {
				errMsg := fmt.Sprintf("update stream: %s", err.Error())
				amRes := interfaces.AutoMappingResponse{
					NetworkUUID: device.NetworkUUID,
					DeviceUUID:  device.UUID,
					HasError:    true,
					Error:       errMsg,
					Level:       interfaces.Device,
				}
				updateCascadeConnectionError(tx, &amRes)
				return nil, err
			}
			deviceUUIDToStreamUUIDMap[device.UUID] = stream.UUID
		}
	}

	// swap back the names
	for _, device := range devices {
		streamName := getAutoMappedStreamName(flowNetwork.Name, network.Name, device.Name)
		if err := tx.Model(&model.Stream{}).
			Where("auto_mapping_device_uuid = ? AND created_from_auto_mapping IS TRUE", device.UUID).
			Update("name", streamName).
			Error; err != nil {
			errMsg := fmt.Sprintf("update stream: %s", err.Error())
			amRes := interfaces.AutoMappingResponse{
				NetworkUUID: device.NetworkUUID,
				DeviceUUID:  device.UUID,
				HasError:    true,
				Error:       errMsg,
				Level:       interfaces.Device,
			}
			updateCascadeConnectionError(tx, &amRes)
			return nil, err
		}
	}
	return deviceUUIDToStreamUUIDMap, nil
}

func (d *GormDatabase) createPointsAutoMappingProducers(streamUUID string, points []*model.Point) (map[string]string, error) {
	pointUUIDToProducerUUIDMap := map[string]string{}
	tx := d.DB.Begin()
	for _, point := range points {
		producer, _ := d.GetOneProducerByArgsTransaction(tx, api.Args{StreamUUID: nils.NewString(streamUUID), ProducerThingUUID: nils.NewString(point.UUID)})
		if producer == nil {
			if boolean.IsTrue(point.AutoMappingEnable) { // create stream only when auto_mapping is enabled
				producer = &model.Producer{}
				producer.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
				d.setProducerModel(streamUUID, point, producer)
				if err := tx.Create(&producer).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				pointUUIDToProducerUUIDMap[point.UUID] = producer.UUID
			}
			continue
		} else {
			d.setProducerModel(streamUUID, point, producer)
			if err := tx.Save(producer).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			pointUUIDToProducerUUIDMap[point.UUID] = producer.UUID
		}
	}

	// swap back the names
	for _, point := range points {
		if err := tx.Model(&model.Producer{}).
			Where("producer_thing_uuid = ? AND created_from_auto_mapping IS TRUE", point.UUID).
			Update("name", point.Name).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return pointUUIDToProducerUUIDMap, nil
}

func (d *GormDatabase) createWriterClones(syncWriters []*interfaces.SyncWriter) (*string, error) {
	tx := d.DB.Begin()
	for _, syncWriter := range syncWriters {
		// it will restrict duplicate creation of writer_clone
		wc, _ := d.GetOneWriterCloneByArgsTransaction(tx, api.Args{ProducerUUID: &syncWriter.ProducerUUID, CreatedFromAutoMapping: boolean.NewTrue()})
		if wc == nil {
			wc = &model.WriterClone{}
			wc.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
			d.setWriterCloneModel(syncWriter, wc)
			if err := tx.Create(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.UUID, err
			}
		} else {
			d.setWriterCloneModel(syncWriter, wc)
			if err := tx.Model(&wc).Updates(&wc).Error; err != nil {
				tx.Rollback()
				return &syncWriter.UUID, err
			}
		}
	}
	tx.Commit()
	return nil, nil
}

func (d *GormDatabase) clearConnectionErrorTransaction(db *gorm.DB, amNetwork *interfaces.AutoMappingNetwork) error {
	tx := db.Begin()
	networkModel := model.Network{
		Connection:        connection.Connected.String(),
		ConnectionMessage: nstring.New(nstring.NotAvailable),
	}
	deviceModel := model.Device{
		Connection:        connection.Connected.String(),
		ConnectionMessage: nstring.New(nstring.NotAvailable),
	}
	pointModel := model.Point{
		Connection:        connection.Connected.String(),
		ConnectionMessage: nstring.New(nstring.NotAvailable),
	}
	err := UpdateNetworkConnectionErrorsTransaction(tx, amNetwork.UUID, &networkModel)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, amDevice := range amNetwork.Devices {
		err = UpdateDeviceConnectionErrorsTransaction(tx, amDevice.UUID, &deviceModel)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, amPoint := range amDevice.Points {
			err = UpdatePointConnectionErrorsTransaction(tx, amPoint.UUID, &pointModel)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func updateCascadeConnectionError(tx *gorm.DB, amRes *interfaces.AutoMappingResponse) {
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

func setStreamModel(flowNetwork *model.FlowNetwork, network *model.Network, device *model.Device, streamModel *model.Stream) {
	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	streamModel.Name = getTempAutoMappedName(getAutoMappedStreamName(flowNetwork.Name, network.Name, device.Name))
	streamModel.Enable = boolean.New(boolean.IsTrue(network.Enable) && boolean.IsTrue(network.AutoMappingEnable) &&
		boolean.IsTrue(device.Enable) && boolean.IsTrue(device.AutoMappingEnable))
	streamModel.CreatedFromAutoMapping = boolean.NewTrue()
	streamModel.AutoMappingNetworkUUID = nstring.New(device.NetworkUUID)
	streamModel.AutoMappingDeviceUUID = nstring.New(device.UUID)
}

func (d *GormDatabase) setProducerModel(streamUUID string, point *model.Point, producerModel *model.Producer) {
	producerModel.Name = getTempAutoMappedName(point.Name)
	producerModel.Enable = boolean.New(boolean.IsTrue(point.Enable) && boolean.IsTrue(point.AutoMappingEnable))
	producerModel.StreamUUID = streamUUID
	producerModel.ProducerThingUUID = point.UUID
	producerModel.ProducerThingName = point.Name
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producerModel.EnableHistory = point.HistoryEnable
	producerModel.HistoryType = point.HistoryType
	producerModel.HistoryInterval = point.HistoryInterval
	producerModel.CreatedFromAutoMapping = boolean.NewTrue()
}

func (d *GormDatabase) setWriterCloneModel(syncWriter *interfaces.SyncWriter, writerClone *model.WriterClone) {
	writerClone.WriterThingName = syncWriter.Name
	writerClone.WriterThingClass = "point"
	writerClone.FlowFrameworkUUID = syncWriter.FlowFrameworkUUID
	writerClone.WriterThingUUID = syncWriter.UUID
	writerClone.ProducerUUID = syncWriter.ProducerUUID
	writerClone.SourceUUID = syncWriter.WriterUUID
	writerClone.CreatedFromAutoMapping = boolean.NewTrue()
}
