package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	cli := client.NewSessionWithToken(body.FlowToken, body.FlowIP, body.FlowPort)
	remoteDeviceInfo, err := cli.DeviceInfo()
	if err != nil {
		return nil, errors.New("please change your localstorage device info, we are unable to connect")
	}
	if remoteDeviceInfo.GlobalUUID != body.GlobalUUID {
		return nil, errors.New("please check your flow_ip, flow_port, it's pointing different device")
	}
	mfn, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	fnc := model.FlowNetworkClone{}
	if err = json.Unmarshal(mfn, &fnc); err != nil {
		return nil, err
	}
	fnc.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetworkClone)
	fnc.SourceUUID = body.UUID
	fnc.SyncUUID, _ = utils.MakeUUID()
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	var flowNetworkClonesModel []*model.FlowNetworkClone
	if err = d.DB.Where("global_uuid = ? ", body.GlobalUUID).Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, err
	}
	if len(flowNetworkClonesModel) == 0 {
		if err = d.DB.Create(fnc).Error; err != nil {
			return nil, err
		}
	} else {
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
