package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/projects/flatlinealerts/ffhistoryrest"
)

func (inst *Instance) FlatlineAnalysis(ffHistoryArray []ffhistoryrest.FFHistory, jobConfig Job) (flatlineAlerts []string, err error) {
	historyDf, err := inst.ConvertFFHistoriesToDataframe(ffHistoryArray)
	if err != nil {
		inst.flatlinealertsErrorMsg("FlatlineAnalysis() ConvertFFHistoriesToDataframe(): ", err)
		return nil, err
	}
	inst.flatlinealertsDebugMsg("FlatlineAnalysis() historyDf: ", historyDf)

	// deviceNameFilter := dataframe.F{Colname: "rubix_device_name", Comparator: series.Eq, Comparando: "Greenhouse"}
	// inst.flatlinealertsDebugMsg("FlatlineAnalysis() historyDf.Filter(deviceNameFilter): ", historyDf.Filter(deviceNameFilter))
	// inst.flatlinealertsDebugMsg("FlatlineAnalysis() historyDf.GroupBy(rubix_device_name).GetGroups(): ", historyDf.GroupBy("rubix_device_name").GetGroups())

	// Take the whole History DataFrame and split it up into groupings to be analyzed

	// groupings := []string{"client_id", "client_name", "site_id", "site_name", "device_id", "device_name", "rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}
	groupings := []string{"rubix_network_uuid", "rubix_network_name", "rubix_device_uuid", "rubix_device_name", "rubix_point_uuid", "rubix_point_name"}

	inst.flatlinealertsDebugMsg("FlatlineAnalysis() groupings: ", groupings)

	mapOfHistoryGroups := historyDf.GroupBy(groupings...).GetGroups()
	inst.flatlinealertsDebugMsg("FlatlineAnalysis() mapOfHistoryGroups: ", mapOfHistoryGroups)

	for key, groupedDf := range mapOfHistoryGroups {
		maxValue := groupedDf.Col("value").Max()
		minValue := groupedDf.Col("value").Min()
		inst.flatlinealertsDebugMsg(fmt.Sprintf("FlatlineAnalysis() df: %s, max: %f, min: %f", key, maxValue, minValue))
		if minValue == maxValue { // All values are the same
			flatlineAlerts = append(flatlineAlerts, key)
		}
	}

	return flatlineAlerts, nil
}
