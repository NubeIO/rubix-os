package main

//func (inst *Instance) addNetwork(bacNet *Network) (*model.Network, interface{}, int, error) {
//
//	var ffNetwork model.Network
//	ffNetwork.Name = bacNet.NetworkName
//	ffNetwork.TransportType = model.TransType.IP
//	ffNetwork.PluginPath = "bacnetmaster"
//	if bacNet.InterfaceName != "" {
//		_net, _ := networking.GetInterfaceByName(bacNet.InterfaceName)
//		if _net == nil {
//			return nil, nil, http.StatusBadRequest, errors.New("failed to find a valid network interface")
//		}
//		bacNet.NetworkIp = _net.IP
//		bacNet.NetworkMask = _net.NetMaskLength
//	} else {
//		ffNetwork.IP = bacNet.NetworkIp
//		ffNetwork.Port = nums.NewInt(bacNet.NetworkPort)
//		ffNetwork.NetworkMask = nums.NewInt(bacNet.NetworkMask)
//	}
//
//	rt.Method = nrest.PUT
//	rt.Path = networkBacnet
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacNet})
//	if err != nil {
//		errMsg := fmt.Sprintf("bacnet-master-plugin: ERROR added network over rest response-code:%d", code)
//		log.Error(errMsg)
//		return nil, res.AsJsonNoErr(), http.StatusBadRequest, errors.New(errMsg)
//	}
//
//	res.ToInterfaceNoErr(bacNet)
//	ffNetwork.AddressUUID = bacNet.NetworkUUID
//	if bacNet.NetworkUUID == "" {
//		errMsg := fmt.Sprintf("bacnet-master-plugin: ERROR no bacnet-server network uuid provided")
//		log.Error(errMsg)
//		return nil, nil, http.StatusBadRequest, errors.New(errMsg)
//	}
//	_network, err := inst.db.CreateNetwork(&ffNetwork, true)
//	if err != nil || _network.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR added network: err", err)
//		rt.Method = nrest.DELETE
//		rt.Path = fmt.Sprintf("%s/%s", networkBacnet, bacNet.NetworkUUID)
//		res, code, err = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacNet})
//		if err != nil {
//			errMsg := fmt.Sprintf("bacnet-master-plugin: ERROR delete network over rest response-code:%d", code)
//			log.Error(errMsg)
//			return nil, res.AsJsonNoErr(), code, nil
//		}
//		return nil, nil, http.StatusBadRequest, err
//	}
//	return _network, nil, http.StatusOK, nil
//}
//
//func (inst *Instance) addDevice(bacDevice *Device) (*model.Device, error) {
//	var ffDevice model.Device
//	ffDevice.Name = bacDevice.DeviceName
//	ffDevice.CommonIP.Host = bacDevice.DeviceIp
//	ffDevice.CommonIP.Port = bacDevice.DevicePort
//	ffDevice.DeviceMask = nums.NewInt(bacDevice.DeviceMask)
//	ffDevice.DeviceObjectId = nums.NewInt(bacDevice.DeviceObjectId)
//
//	rt.Method = nrest.PUT
//	rt.Path = deviceBacnet
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacDevice})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR added device over rest response-code:", code)
//		return nil, err
//	}
//	res.ToInterfaceNoErr(bacDevice)
//	ffDevice.AddressUUID = &bacDevice.DeviceUUID
//
//	if bacDevice.DeviceUUID == "" {
//		log.Error("bacnet-master-plugin: ERROR no bacnet-server device uuid provided")
//		return nil, err
//	}
//
//	getNet, err := inst.db.GetNetworkByField("address_uuid", bacDevice.NetworkUUID, false)
//	if err != nil || getNet.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetNetworkByField() failed to find network", err)
//		return nil, err
//	}
//
//	ffDevice.NetworkUUID = getNet.UUID
//	dev, err := inst.db.CreateDevice(&ffDevice)
//	if err != nil || dev.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR added device: err", err)
//		rt.Method = nrest.DELETE
//		rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, bacDevice.DeviceUUID)
//		res, code, err = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacDevice})
//		if err != nil {
//			log.Error("bacnet-master-plugin: ERROR delete device over rest response-code:", code)
//			return nil, err
//		}
//		return nil, err
//	}
//	return dev, nil
//}
//
//func (inst *Instance) addPoint(bacPoint *Point) (*model.Point, error) {
//	var ffPoint model.Point
//	ffPoint.Name = bacPoint.PointName
//	ffPoint.ObjectType = bacPoint.PointObjectType
//	ffPoint.AddressID = nums.NewInt(bacPoint.PointObjectId)
//
//	rt.Method = nrest.PUT
//	rt.Path = pointBacnet
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacPoint})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR added point over rest response-code:", res.AsString())
//		return nil, err
//	}
//	res.ToInterfaceNoErr(bacPoint)
//	ffPoint.AddressUUID = &bacPoint.DeviceUUID
//
//	if bacPoint.DeviceUUID == "" {
//		log.Error("bacnet-master-plugin: ERROR no bacnet-server point uuid provided")
//		return nil, err
//	}
//
//	getDev, err := inst.db.GetOneDeviceByArgs(api.Args{AddressUUID: &bacPoint.DeviceUUID})
//	if err != nil || getDev.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetDeviceByField() failed to find device", err)
//		return nil, err
//	}
//	ffPoint.DeviceUUID = getDev.UUID
//	pnt, err := inst.db.CreatePoint(&ffPoint, false, false)
//	if err != nil || pnt.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR added device: err", err)
//		rt.Method = nrest.DELETE
//		rt.Path = fmt.Sprintf("%s/%s", pointBacnet, bacPoint.PointUUID)
//		res, code, _ = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacPoint})
//		if err != nil {
//			log.Error("bacnet-master-plugin: ERROR delete device over rest response-code:", code)
//		}
//		return nil, err
//	}
//	pnt.AddressUUID = &bacPoint.PointUUID
//	pnt, err = inst.db.UpdatePoint(pnt.UUID, &ffPoint, false)
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR on update point with bacnet point_uuid to address_uuid", code)
//		return nil, err
//	}
//	return pnt, nil
//}
//
//func (inst *Instance) patchNetwork(bacNet *Network, uuid string) (*model.Network, error) {
//	var ffNetwork model.Network
//	ffNetwork.Name = bacNet.NetworkName
//	ffNetwork.TransportType = model.TransType.IP
//	ffNetwork.PluginPath = "bacnetmaster"
//	ffNetwork.IP = bacNet.NetworkIp
//	ffNetwork.Port = nums.NewInt(bacNet.NetworkPort)
//	ffNetwork.NetworkMask = nums.NewInt(bacNet.NetworkMask)
//
//	rt.Method = nrest.GET
//	rt.Path = fmt.Sprintf("%s/%s", networkBacnet, uuid)
//	bacModel := new(Network)
//
//	getBacnetNetwork, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR get network over rest response-code:", code, "response:", getBacnetNetwork.AsString())
//		return nil, err
//	}
//	getBacnetNetwork.ToInterfaceNoErr(bacModel)
//
//	rt.Method = nrest.PATCH
//	rt.Path = fmt.Sprintf("%s/%s", networkBacnet, uuid)
//	_, code, err = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacNet})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR patch network over rest response-code:", code, "response:", getBacnetNetwork.AsString())
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", networkBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch network over rest response-code:", code)
//		return nil, err
//	}
//
//	getNetwork, err := inst.db.GetNetworkByField("address_uuid", bacNet.NetworkUUID, false)
//	if err != nil || getNetwork.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetNetworkByField() failed to find network", err)
//		return nil, err
//	}
//
//	updateNetwork, err := inst.db.UpdateNetwork(getNetwork.UUID, &ffNetwork, true)
//	if err != nil {
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch network over rest response-code:", code)
//		return nil, err
//	}
//	return updateNetwork, nil
//}
//
//func (inst *Instance) patchDevice(bacDevice *Device, uuid string) (*model.Device, error) {
//	var ffDevice model.Device
//	ffDevice.Name = bacDevice.DeviceName
//	ffDevice.CommonIP.Host = bacDevice.DeviceIp
//	ffDevice.CommonIP.Port = bacDevice.DevicePort
//	ffDevice.DeviceMask = nums.NewInt(bacDevice.DeviceMask)
//	ffDevice.DeviceObjectId = nums.NewInt(bacDevice.DeviceObjectId)
//
//	rt.Method = nrest.GET
//	rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//	bacModel := new(Device)
//
//	getBacnetDevice, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR get device over rest response-code:", code, "response:", getBacnetDevice.AsString())
//		return nil, err
//	}
//
//	getBacnetDevice.ToInterfaceNoErr(bacModel)
//
//	rt.Method = nrest.PATCH
//	rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//	_, code, err = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacDevice})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR patch device over rest response-code:", code, "response:", getBacnetDevice.AsString())
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch device over rest response-code:", code)
//		return nil, err
//	}
//
//	getDev, err := inst.db.GetOneDeviceByArgs(api.Args{AddressUUID: &bacDevice.DeviceUUID})
//	if err != nil || getDev.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetDeviceByField() failed to find device", err)
//		return nil, err
//	}
//
//	updateDevice, err := inst.db.UpdateDevice(getDev.UUID, &ffDevice, true)
//	if err != nil {
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch device over rest response-code:", code)
//		return nil, err
//	}
//
//	return updateDevice, nil
//}
//
//func (inst *Instance) patchPoint(bacPoint *Point, uuid string) (*model.Point, error) {
//	var ffPoint model.Point
//	ffPoint.Name = bacPoint.PointName
//	ffPoint.ObjectType = bacPoint.PointObjectType
//	ffPoint.AddressID = nums.NewInt(bacPoint.PointObjectId)
//
//	rt.Method = nrest.GET
//	rt.Path = fmt.Sprintf("%s/%s", pointBacnet, uuid)
//	bacModel := new(Point)
//
//	getBacnetPoint, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR get point over rest response-code:", code, "response:", getBacnetPoint.AsString())
//		return nil, err
//	}
//	getBacnetPoint.ToInterfaceNoErr(bacModel)
//
//	rt.Method = nrest.PATCH
//	rt.Path = fmt.Sprintf("%s/%s", pointBacnet, uuid)
//	_, code, err = nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacPoint})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR patch point over rest response-code:", code, "response:", getBacnetPoint.AsString())
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", pointBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch point over rest response-code:", code)
//		return nil, err
//	}
//
//	getPnt, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: &uuid})
//	if err != nil || getPnt.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetPointByField() failed to find point", err, uuid)
//		return nil, err
//	}
//
//	updatePoint, err := inst.db.UpdatePoint(getPnt.UUID, &ffPoint, true)
//	if err != nil {
//		rt.Method = nrest.PATCH
//		rt.Path = fmt.Sprintf("%s/%s", pointBacnet, uuid)
//		_, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: bacModel})
//		log.Error("bacnet-master-plugin: ERROR re-patch network over rest response-code:", code)
//		return nil, err
//	}
//
//	return updatePoint, nil
//}
//
//func (inst *Instance) deleteNetwork(uuid string) (bool, error) {
//
//	rt.Method = nrest.DELETE
//	rt.Path = fmt.Sprintf("%s/%s", networkBacnet, uuid)
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR delete network over rest response-code:", "response-string:", code, res.AsString())
//		return false, nil
//	}
//
//	getNetwork, err := inst.db.GetNetworkByField("address_uuid", uuid, false)
//	if err != nil || getNetwork.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetNetworkByField() failed to find network", err)
//		return false, err
//	}
//
//	_, err = inst.db.DeleteNetwork(getNetwork.UUID)
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//
//}
//
//func (inst *Instance) deleteDevice(uuid string) (bool, error) {
//	rt.Method = nrest.DELETE
//	rt.Path = fmt.Sprintf("%s/%s", deviceBacnet, uuid)
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR delete device over rest response-code:", "response-string:", code, res.AsString())
//		return false, nil
//	}
//
//	getDev, err := inst.db.GetOneDeviceByArgs(api.Args{AddressUUID: &uuid})
//	if err != nil || getDev.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetDeviceByField() failed to find device", err)
//		return false, err
//	}
//
//	_, err = inst.db.DeleteDevice(getDev.UUID)
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (inst *Instance) deletePoint(uuid string) (bool, error) {
//	rt.Method = nrest.DELETE
//	rt.Path = fmt.Sprintf("%s/%s", pointBacnet, uuid)
//	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
//	if err != nil {
//		log.Error("bacnet-master-plugin: ERROR delete point over rest response-code:", "response-string:", code, res.AsString())
//		return false, nil
//	}
//
//	getPnt, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: &uuid})
//	if err != nil || getPnt.UUID == "" {
//		log.Error("bacnet-master-plugin: ERROR on get GetPointByField() failed to find point", err)
//		return false, err
//	}
//
//	_, err = inst.db.DeletePoint(getPnt.UUID)
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
