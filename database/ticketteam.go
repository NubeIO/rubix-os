package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
)

func (d *GormDatabase) UpdateTicketTeams(ticketUUID string, teamUUIDs []*string) ([]*model.TicketTeam, error) {
	teams, _ := d.GetTeamsByUUIDs(teamUUIDs)
	var notInUUIDs []string
	var body []*model.TicketTeam
	for _, team := range teams {
		var ticketTeamModel *model.TicketTeam
		d.DB.Model(&model.TicketTeam{}).Where("ticket_uuid = ?", ticketUUID).Where("team_uuid = ?", team.UUID).
			Find(&ticketTeamModel)
		if ticketTeamModel != nil {
			notInUUIDs = append(notInUUIDs, team.UUID)
		}
		body = append(body, &model.TicketTeam{
			TicketUUID: ticketUUID,
			TeamUUID:   team.UUID,
		})
	}
	condition := d.DB.Where("ticket_uuid = ?", ticketUUID)
	if notInUUIDs != nil {
		condition = condition.Where("team_uuid not in ?", notInUUIDs)
	}
	if err := condition.Delete(&model.TicketTeam{}).Error; err != nil {
		return nil, err
	}
	if len(body) > 0 {
		if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			return nil, err
		}
	}
	if err := d.DB.Where("ticket_uuid = ?", ticketUUID).Find(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}
