package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/str"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	if !boolean.IsTrue(body.IsMasterSlave) {
		accessToken, err := client.GetFlowToken(*body.FlowIP, *body.FlowPort, *body.FlowUsername, *body.FlowPassword)
		if err != nil {
			return nil, err
		}
		body.FlowToken = accessToken
	}
	cli := client.NewFlowClientCli(body.FlowIP, body.FlowPort, body.FlowToken, body.IsMasterSlave, body.GlobalUUID, false)
	remoteDeviceInfo, err := cli.DeviceInfo()
	if err != nil {
		return nil, err
	}
	if remoteDeviceInfo.GlobalUUID != body.GlobalUUID {
		return nil, errors.New("please check your flow_ip, flow_port, it's pointing different device")
	}
	fn, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	fnc := model.FlowNetworkClone{}
	if err = json.Unmarshal(fn, &fnc); err != nil {
		return nil, err
	}
	fnc.SourceUUID = body.UUID
	fnc.SyncUUID, _ = utils.MakeUUID()
	if !boolean.IsTrue(fnc.IsRemote) {
		fnc.FlowHTTPS = boolean.NewFalse()
		fnc.FlowIP = str.NewStringAddress("0.0.0.0")
	}
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	var flowNetworkClonesModel []*model.FlowNetworkClone
	d.DB.Where("global_uuid = ? ", body.GlobalUUID).Find(&flowNetworkClonesModel)
	if len(flowNetworkClonesModel) == 0 {
		fnc.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetworkClone)
		if err = d.DB.Create(fnc).Error; err != nil {
			return nil, err
		}
	} else {
		fnc.UUID = flowNetworkClonesModel[0].UUID
		if err = d.DB.Model(&flowNetworkClonesModel[0]).Updates(fnc).Error; err != nil {
			return nil, err
		}
	}
	fnc.GlobalUUID = deviceInfo.GlobalUUID
	fnc.ClientId = deviceInfo.ClientId
	fnc.ClientName = deviceInfo.ClientName
	fnc.SiteId = deviceInfo.SiteId
	fnc.SiteName = deviceInfo.SiteName
	fnc.DeviceId = deviceInfo.DeviceId
	fnc.DeviceName = deviceInfo.DeviceName
	return &fnc, nil
}
