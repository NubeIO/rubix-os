package client

import (
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) CreateAutoMappingSchedule(body *interfaces.AutoMapping) interfaces.AutoMappingScheduleResponse {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingScheduleResponse{}).
		SetBody(body).
		Post("/api/auto_mapping_schedules"))
	if err != nil {
		scheduleUUID := "" // pick first valid schedule
		for _, schedule := range body.Schedules {
			if schedule.CreateSchedule {
				scheduleUUID = schedule.UUID
				break
			}
		}
		return interfaces.AutoMappingScheduleResponse{
			HasError:     true,
			ScheduleUUID: scheduleUUID,
			Error:        err.Error(),
		}
	}
	return *resp.Result().(*interfaces.AutoMappingScheduleResponse)
}
