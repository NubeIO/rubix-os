package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetTeams() ([]*model.Team, error) {
	var teamsModel []*model.Team
	query := d.buildTeamQuery()
	if err := query.Find(&teamsModel).Error; err != nil {
		return nil, err
	}
	return teamsModel, nil
}

func (d *GormDatabase) GetTeamsByUUIDs(uuids []*string) ([]*model.Team, error) {
	var teamsModel []*model.Team
	query := d.buildTeamQuery()
	if err := query.Where("uuid IN ?", uuids).Find(&teamsModel).Error; err != nil {
		return nil, err
	}
	return teamsModel, nil
}

func (d *GormDatabase) GetTeam(uuid string) (*model.Team, error) {
	var teamModel *model.Team
	query := d.buildGroupQuery()
	if err := query.Where("uuid = ?", uuid).First(&teamModel).Error; err != nil {
		return nil, err
	}
	return teamModel, nil
}

func (d *GormDatabase) CreateTeam(body *model.Team) (*model.Team, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Team)
	body.Members = nil
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateTeam(uuid string, body *model.Team) (*model.Team, error) {
	var teamModel *model.Team
	if err := d.DB.Where("uuid = ?", uuid).Find(&teamModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return teamModel, nil
}

func (d *GormDatabase) DeleteTeam(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Team{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DropTeams() (bool, error) {
	query := d.DB.Where("1 = 1").Delete(&model.Team{})
	return d.deleteResponseBuilder(query)
}
