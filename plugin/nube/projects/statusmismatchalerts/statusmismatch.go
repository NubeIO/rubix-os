package main

import (
	"fmt"
	"github.com/NubeIO/rubix-os/plugin/nube/projects/statusmismatchalerts/ffhistoryrest"
)

func (inst *Instance) StatusMismatchAnalysis(ffHistoryArray []ffhistoryrest.FFHistory, jobConfig Job) (highStatusMismatchAlerts, lowStatusMismatchAlerts []string, err error) {
	historyDf, err := inst.ConvertFFHistoriesToDataframe(ffHistoryArray)
	if err != nil {
		inst.statusmismatchalertsErrorMsg("StatusMismatchAnalysis() ConvertFFHistoriesToDataframe(): ", err)
		return nil, nil, err
	}
	inst.statusmismatchalertsDebugMsg("StatusMismatchAnalysis() historyDf: ", historyDf)

	// deviceNameFilter := dataframe.F{Colname: "rubix_device_name", Comparator: series.Eq, Comparando: "Greenhouse"}
	// inst.statusmismatchalertsDebugMsg("StatusMismatchAnalysis() historyDf.Filter(deviceNameFilter): ", historyDf.Filter(deviceNameFilter))
	// inst.statusmismatchalertsDebugMsg("StatusMismatchAnalysis() historyDf.GroupBy(rubix_device_name).GetGroups(): ", historyDf.GroupBy("rubix_device_name").GetGroups())

	// Take the whole History DataFrame and split it up into groupings to be analyzed

	// groupings := []string{"client_id", "client_name", "site_id", "site_name", "device_id", "device_name", "rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}
	groupings := []string{"rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}

	inst.statusmismatchalertsDebugMsg("StatusMismatchAnalysis() groupings: ", groupings)

	mapOfHistoryGroups := historyDf.GroupBy(groupings...).GetGroups()
	inst.statusmismatchalertsDebugMsg("StatusMismatchAnalysis() mapOfHistoryGroups: ", mapOfHistoryGroups)

	for key, groupedDf := range mapOfHistoryGroups {
		maxValue := groupedDf.Col("value").Max()
		minValue := groupedDf.Col("value").Min()

		inst.statusmismatchalertsDebugMsg(fmt.Sprintf("StatusMismatchAnalysis() df: %s, max: %f, min: %f", key, maxValue, minValue))
		/*
			if jobConfig.HighLimitStatusMismatch != nil && minValue >= *jobConfig.HighLimitStatusMismatch { // All values are above the high limit statusmismatch
				highStatusMismatchAlerts = append(highStatusMismatchAlerts, key)
			}
		*/
	}

	return highStatusMismatchAlerts, lowStatusMismatchAlerts, nil
}
