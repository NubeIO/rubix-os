package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"gorm.io/gorm/clause"
)

func (d *GormDatabase) UpdateTeamViews(teamUUID string, viewUUIDs []*string) ([]*model.TeamView, error) {
	views, _ := d.GetViewsByUUIDs(viewUUIDs, argspkg.Args{})
	var notInUUIDs []string
	var body []*model.TeamView
	for _, view := range views {
		var teamViewModel *model.TeamView
		d.DB.Model(&model.TeamView{}).Where("team_uuid = ?", teamUUID).Where("view_uuid = ?", view.UUID).
			Find(&teamViewModel)
		if teamViewModel != nil {
			notInUUIDs = append(notInUUIDs, view.UUID)
		}
		body = append(body, &model.TeamView{
			TeamUUID: teamUUID,
			ViewUUID: view.UUID,
		})
	}
	condition := d.DB.Where("team_uuid = ?", teamUUID)
	if notInUUIDs != nil {
		condition = condition.Where("view_uuid not in ?", notInUUIDs)
	}
	if err := condition.Delete(&model.TeamView{}).Error; err != nil {
		return nil, err
	}
	if len(body) > 0 {
		if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			return nil, err
		}
	}
	if err := d.DB.Where("team_uuid = ?", teamUUID).Find(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}
