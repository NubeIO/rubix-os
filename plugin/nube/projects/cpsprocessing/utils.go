package main

import (
	"errors"
	"github.com/go-gota/gota/dataframe"
	"time"
)

func GetOverdueDelay(doorType DoorType, thresholdsDF dataframe.DataFrame) (time.Duration, error) {
	switch doorType {
	case facilityEntrance, facilityToilet, facilityDDA:
		cleaningOverdueAlertDelayMins, err := thresholdsDF.Col(string(facilityCleaningOverdueAlertDelayColName)).Elem(0).Int()
		if err != nil {
			return 30 * time.Minute, err
		}
		return time.Duration(cleaningOverdueAlertDelayMins) * time.Minute, nil

	case eotEntrance, eotToilet, eotShower, eotDDA:
		cleaningOverdueAlertDelayMins, err := thresholdsDF.Col(string(eotCleaningOverdueAlertDelayColName)).Elem(0).Int()
		if err != nil {
			return 30 * time.Minute, err
		}
		return time.Duration(cleaningOverdueAlertDelayMins) * time.Minute, nil
	}
	return 30 * time.Minute, errors.New("unknown door type")
}
