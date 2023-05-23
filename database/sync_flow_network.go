package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/deviceinfo"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	cli := client.NewFlowClientCliFromFN(body)
	// In this case it's even though FN, this is clone side & need to communicate with slave device from master
	if boolean.IsTrue(body.IsMasterSlave) {
		cli = client.NewMasterToSlaveSession(body.GlobalUUID)
	}
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
	fnc.SyncUUID, _ = nuuid.MakeUUID()
	deviceInfo, err := deviceinfo.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	var flowNetworkClonesModel []*model.FlowNetworkClone
	d.DB.Where("global_uuid = ? ", body.GlobalUUID).Find(&flowNetworkClonesModel)
	if len(flowNetworkClonesModel) == 0 {
		fnc.UUID = nuuid.MakeTopicUUID(model.CommonNaming.FlowNetworkClone)
		if err = d.DB.Create(&fnc).Error; err != nil {
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
