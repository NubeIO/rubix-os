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

// GetLastProcessedDataAndDoorType gets the last processed data values and door info from the tags and history values
func (inst *Instance) GetLastProcessedDataAndDoorType(dfJoinedLastProcessedValuesAndPoints *dataframe.DataFrame, doorSensorPoint *DoorProcessingPoint) (pointLastProcessedData *LastProcessedData, pointDoorInfo *DoorInfo, err error) {
	pointDoorInfo = &DoorInfo{}
	pointLastProcessedData = &LastProcessedData{}

	// convert the dataframe to a map and iterate through the rows to pick out the applicable values
	lastValuesAndPointsMap := dfJoinedLastProcessedValuesAndPoints.Maps()
	ok := false
	for _, row := range lastValuesAndPointsMap {
		switch row["name"].(string) {
		case string(doorPositionColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.DoorPosition = int(row["value"].(float64))
			}
		case string(currentUsesColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.CurrentUses = int(row["value"].(float64))
			}
		case string(totalUsesColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.TotalUses = int(row["value"].(float64))
			}
		case string(cubicleOccupancyColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.CubicleOccupancy = int(row["value"].(float64))
			}
		case string(pendingStatusColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.PendingStatus = int(row["value"].(float64))
			}
		case string(overdueStatusColName):
			_, ok = row["value"].(float64)
			if ok {
				pointLastProcessedData.OverdueStatus = int(row["value"].(float64))
			}
		case string(toPendingColName):
			_, ok = row["timestamp"].(string)
			if ok {
				pointLastProcessedData.LastToPendingTimestamp = row["timestamp"].(string)
			}
		}
	}

	// get the door type
	if doorSensorPoint.IsEOT {
		pointDoorInfo.IsEOT = true
		switch doorSensorPoint.DoorType {
		case string(entranceDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = eotEntrance
		case string(toiletDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = eotToilet
		case string(showerDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = eotShower
		case string(ddaDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = eotDDA
		case string(doorDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = eotDoor
		default:
			err = errors.New("doorType tag not recognized")
			inst.cpsErrorMsg("GetLastProcessedDataAndDoorType() error: ", err)
		}
	} else {
		pointDoorInfo.IsEOT = false
		switch doorSensorPoint.DoorType {
		case string(entranceDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = facilityEntrance
		case string(toiletDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = facilityToilet
		case string(ddaDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = facilityDDA
		case string(doorDoorTypeTagValue):
			pointDoorInfo.DoorTypeID = facilityDoor
		default:
			err = errors.New("doorType tag not recognized")
			inst.cpsErrorMsg("GetLastProcessedDataAndDoorType() error: ", err)
		}
	}

	if doorSensorPoint.NormalPosition == string(normallyOpenNormalPositionTagValue) {
		pointDoorInfo.NormalPosition = normallyOpen
	} else {
		pointDoorInfo.NormalPosition = normallyClosed
	}

	if doorSensorPoint.EnableUseCounting == string(enabledEnableUseCountingTagValue) {
		pointDoorInfo.EnableUseCounting = true
	} else {
		pointDoorInfo.EnableUseCounting = false
	}

	if doorSensorPoint.EnableCleaningTracking == string(enabledEnableCleaningTrackingTagValue) {
		pointDoorInfo.EnableCleaningTracking = true
	} else {
		pointDoorInfo.EnableCleaningTracking = false
	}

	if doorSensorPoint.EnableCleaningTracking == string(enabledEnableCleaningTrackingTagValue) {
		pointDoorInfo.EnableCleaningTracking = true
	} else {
		pointDoorInfo.EnableCleaningTracking = false
	}

	if doorSensorPoint.AssetFunc == string(usageCountDoorSensorAssetFunctionTagValue) {
		pointDoorInfo.AssetFunction = string(usageCountDoorSensorAssetFunctionTagValue)
	} else if doorSensorPoint.AssetFunc == string(managedCubicleDoorSensorAssetFunctionTagValue) {
		pointDoorInfo.AssetFunction = string(managedCubicleDoorSensorAssetFunctionTagValue)
	} else if doorSensorPoint.AssetFunc == string(managedFacilityEntranceDoorSensorAssetFunctionTagValue) {
		pointDoorInfo.AssetFunction = string(managedFacilityEntranceDoorSensorAssetFunctionTagValue)
	} else {
		err = errors.New("assetFunc tag not recognized")
		inst.cpsErrorMsg("GetLastProcessedDataAndDoorType() error: ", err)
	}

	pointDoorInfo.AvailabilityID = doorSensorPoint.AvailabilityID
	pointDoorInfo.ResetID = doorSensorPoint.ResetID

	return
}
