package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

func (d *GormDatabase) GetMemberDevicesByMemberUUID(memberUUID string) ([]*model.MemberDevice, error) {
	var memberDevicesModel []*model.MemberDevice
	if err := d.DB.Where("member_uuid = ?", memberUUID).Find(&memberDevicesModel).Error; err != nil {
		return nil, err
	}
	return memberDevicesModel, nil
}

func (d *GormDatabase) GetMemberDevicesByArgs(args argspkg.Args) ([]*model.MemberDevice, error) {
	var memberDevicesModel []*model.MemberDevice
	query := d.buildMemberDeviceQuery(args)
	if err := query.Find(&memberDevicesModel).Error; err != nil {
		return nil, err
	}
	return memberDevicesModel, nil
}

func (d *GormDatabase) GetOneMemberDeviceByArgs(args argspkg.Args) (*model.MemberDevice, error) {
	var memberDeviceModel *model.MemberDevice
	query := d.buildMemberDeviceQuery(args)
	if err := query.First(&memberDeviceModel).Error; err != nil {
		return nil, err
	}
	return memberDeviceModel, nil
}

func (d *GormDatabase) CreateMemberDevice(body *model.MemberDevice) (*model.MemberDevice, error) {
	obj, err := checkMemberDevicePlatform(*body.Platform)
	if err != nil {
		return nil, err
	}
	body.Platform = nstring.New(string(obj))
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.MemberDevice)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateMemberDevice(uuid string, body *model.MemberDevice) (*model.MemberDevice, error) {
	if body.Platform != nil {
		obj, err := checkMemberDevicePlatform(*body.Platform)
		if err != nil {
			return nil, err
		}
		body.Platform = nstring.New(string(obj))
	}
	var memberDeviceModel *model.MemberDevice
	if err := d.DB.Where("uuid = ?", uuid).Find(&memberDeviceModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return memberDeviceModel, nil
}

func (d *GormDatabase) DeleteMemberDevicesByArgs(args argspkg.Args) (bool, error) {
	query := d.buildMemberDeviceQuery(args)
	query = query.Delete(&model.MemberDevice{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SendNotificationByMemberUUID(memberUUID string, data map[string]interface{}) {
	key := d.GetFcmServerKey()
	cli := cligetter.GetFcmServerClient(key)
	memberDevices, _ := d.GetMemberDevicesByMemberUUID(memberUUID)
	wg := &sync.WaitGroup{}
	for _, memberDevice := range memberDevices {
		wg.Add(1)
		go func(memberDevice *model.MemberDevice, data map[string]interface{}) {
			defer wg.Done()
			log.Infof(">>>>>>>>>>>> Sending data to device: %s", *memberDevice.DeviceName)
			if data["to"] != nil {
				data["to"] = nstring.New(memberDevice.DeviceID)
			}
			content := cli.SendNotification(data)
			if len(content) > 0 {
				failure := content["failure"].(bool)
				results := content["results"].([]interface{})
				if failure && len(results) > 0 {
					errorMsg := results[0].(map[string]interface{})["error"].(string)
					if errorMsg == "InvalidRegistration" || errorMsg == "NotRegistered" {
						log.Warnf(">>>>>>>>>>>>>>> Removing device: %s from list!", *memberDevice.DeviceName)
						_, _ = d.DeleteMemberDevicesByArgs(argspkg.Args{DeviceId: nstring.New(memberDevice.DeviceID)})
					}
				}
			}
		}(memberDevice, data)
	}
	wg.Wait()
}
