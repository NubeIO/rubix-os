package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetMemberDevicesByMemberUUID(memberUUID string) ([]*model.MemberDevice, error) {
	var memberDevicesModel []*model.MemberDevice
	if err := d.DB.Where("member_uuid = ?", memberUUID).Find(&memberDevicesModel).Error; err != nil {
		return nil, err
	}
	return memberDevicesModel, nil
}

func (d *GormDatabase) GetMemberDevicesByArgs(args api.Args) ([]*model.MemberDevice, error) {
	var memberDevicesModel []*model.MemberDevice
	query := d.buildMemberDeviceQuery(args)
	if err := query.Find(&memberDevicesModel).Error; err != nil {
		return nil, err
	}
	return memberDevicesModel, nil
}

func (d *GormDatabase) GetOneMemberDeviceByArgs(args api.Args) (*model.MemberDevice, error) {
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

func (d *GormDatabase) DeleteMemberDevicesByArgs(args api.Args) (bool, error) {
	query := d.buildMemberDeviceQuery(args)
	query = query.Delete(&model.MemberDevice{})
	return d.deleteResponseBuilder(query)
}
