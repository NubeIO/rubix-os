package main

import (
	"github.com/NubeIO/rubix-os/plugin/nube/projects/flatlinealerts/ffhistoryrest"
	"github.com/go-gota/gota/dataframe"
)

func (inst *Instance) ConvertFFHistoriesToDataframe(ffHistoryArray []ffhistoryrest.FFHistory) (*dataframe.DataFrame, error) {

	df := dataframe.LoadStructs(ffHistoryArray)
	if df.Error() != nil && df.Error().Error() != "" {
		return nil, df.Error()
	}
	return &df, nil
}
