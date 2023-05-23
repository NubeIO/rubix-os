package main

import (
	"fmt"
	"github.com/NubeIO/rubix-os/plugin/nube/projects/thresholdalerts/ffhistoryrest"
)

func (inst *Instance) ThresholdAnalysis(ffHistoryArray []ffhistoryrest.FFHistory, jobConfig Job) (highThresholdAlerts, lowThresholdAlerts []string, err error) {
	historyDf, err := inst.ConvertFFHistoriesToDataframe(ffHistoryArray)
	if err != nil {
		inst.thresholdalertsErrorMsg("ThresholdAnalysis() ConvertFFHistoriesToDataframe(): ", err)
		return nil, nil, err
	}
	inst.thresholdalertsDebugMsg("ThresholdAnalysis() historyDf: ", historyDf)

	// deviceNameFilter := dataframe.F{Colname: "rubix_device_name", Comparator: series.Eq, Comparando: "Greenhouse"}
	// inst.thresholdalertsDebugMsg("ThresholdAnalysis() historyDf.Filter(deviceNameFilter): ", historyDf.Filter(deviceNameFilter))
	// inst.thresholdalertsDebugMsg("ThresholdAnalysis() historyDf.GroupBy(rubix_device_name).GetGroups(): ", historyDf.GroupBy("rubix_device_name").GetGroups())

	// Take the whole History DataFrame and split it up into groupings to be analyzed

	// groupings := []string{"client_id", "client_name", "site_id", "site_name", "device_id", "device_name", "rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}
	groupings := []string{"rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}

	inst.thresholdalertsDebugMsg("ThresholdAnalysis() groupings: ", groupings)

	mapOfHistoryGroups := historyDf.GroupBy(groupings...).GetGroups()
	inst.thresholdalertsDebugMsg("ThresholdAnalysis() mapOfHistoryGroups: ", mapOfHistoryGroups)

	for key, groupedDf := range mapOfHistoryGroups {
		maxValue := groupedDf.Col("value").Max()
		minValue := groupedDf.Col("value").Min()
		inst.thresholdalertsDebugMsg(fmt.Sprintf("ThresholdAnalysis() df: %s, max: %f, min: %f", key, maxValue, minValue))
		if jobConfig.HighLimitThreshold != nil && minValue >= *jobConfig.HighLimitThreshold { // All values are above the high limit threshold
			highThresholdAlerts = append(highThresholdAlerts, key)
		}
		if jobConfig.LowLimitThreshold != nil && maxValue <= *jobConfig.LowLimitThreshold { // All values are below the low limit threshold
			lowThresholdAlerts = append(lowThresholdAlerts, key)
		}
	}

	return highThresholdAlerts, lowThresholdAlerts, nil
}
