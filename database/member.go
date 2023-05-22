package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/utils/security"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) GetMembers() ([]*model.Member, error) {
	var membersModel []*model.Member
	query := d.buildTeamQuery(api.Args{})
	query.Find(&membersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return membersModel, nil
}

func (d *GormDatabase) GetMember(uuid string) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildTeamQuery(api.Args{})
	query = query.Where("uuid = ? ", uuid).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) GetMemberByUsername(username string) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildTeamQuery(api.Args{})
	query = query.Where("username = ? ", username).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) GetMemberByEmail(email string) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildTeamQuery(api.Args{})
	query = query.Where("email = ? ", email).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) CreateMember(body *model.Member) (*model.Member, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Member)
	hashedPassword, err := security.GeneratePasswordHash(body.Password)
	if err != nil {
		return nil, err
	}
	body.Password = hashedPassword
	body.State = nstring.New(string(model.UnVerified))
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateMember(uuid string, body *model.Member) (*model.Member, error) {
	if body.State != nil {
		obj, err := checkMemberState(*body.State)
		if err != nil {
			return nil, err
		}
		body.State = nstring.New(string(obj))
	}
	var memberModel *model.Member
	query := d.DB.Where("uuid = ?", uuid).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&memberModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) UpdateMemberGroups(uuid string, body []*string) error {
	member, err := d.GetMember(uuid)
	if err != nil {
		return err
	}
	groups, _ := d.GetGroupsByUUIDs(body)
	if err := d.updateGroups(&member, groups); err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) DeleteMember(uuid string) (bool, error) {
	var memberModel *model.Member
	query := d.DB.Where("uuid = ? ", uuid).Delete(&memberModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteMemberByUsername(username string) (bool, error) {
	var memberModel *model.Member
	query := d.DB.Where("username = ? ", username).Delete(&memberModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) ChangeMemberPassword(uuid string, password string) (bool, error) {
	var memberModel *model.Member
	query := d.DB.Where("uuid = ?", uuid).First(&memberModel)
	if query.Error != nil {
		return false, query.Error
	}
	hashedPassword, err := security.GeneratePasswordHash(password)
	if err != nil {
		return false, err
	}
	query = d.DB.Model(&memberModel).Update("password", hashedPassword)
	if query.Error != nil {
		return false, query.Error
	}
	return true, nil
}
