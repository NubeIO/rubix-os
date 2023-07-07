package main

import (
	"errors"
	"github.com/go-gota/gota/dataframe"
	"strconv"
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
func GetLastProcessedDataAndDoorType(dfJoinedLastProcessedValuesAndPoints *dataframe.DataFrame, doorSensorPoint *DoorProcessingPoint) (pointLastProcessedData *LastProcessedData, pointDoorInfo *DoorInfo, err error) {
	pointDoorInfo = &DoorInfo{}
	pointLastProcessedData = &LastProcessedData{}
	var lastOccupancy, lastPendingStatus, lastOverdueStatus, lastCurrentUseCount, lastTotalUseCount int

	// convert the dataframe to a map and iterate through the rows to pick out the applicable values
	lastValuesAndPointsMap := dfJoinedLastProcessedValuesAndPoints.Maps()
	ok := false
	for _, row := range lastValuesAndPointsMap {

		switch row["name"].(string) {
		/*  // TODO: get the door info from the doorSensorPoint instead.
		case string(doorTypeTag):
			_, ok = row["value"]
			if ok {
				pointDoorInfo.DoorTypeTag = row["value"].(string)
			}
		case string(normalPositionTag):
			_, ok = row["value"]
			if ok {
				doorNormalPositionString = row["value"].(string)
				if doorNormalPositionString == string(normallyOpenNormalPositionTagValue) {
					pointDoorInfo.NormalPosition = normallyOpen
				} else if doorNormalPositionString == string(normallyClosedNormalPositionTagValue) {
					pointDoorInfo.NormalPosition = normallyClosed
				}
			}
		case string(assetFuncTag):
			_, ok = row["value"]
			if ok && row["value"].(string) == string(usageCountDoorSensorAssetFunctionTagValue) {
				pointDoorInfo.AssetFunction = row["value"].(string)
			}
		case string(enableCleaningTrackingTag):
			_, ok = row["value"]
			if ok && row["value"].(string) == string(enabledEnableCleaningTrackingTagValue) {
				pointDoorInfo.EnableCleaningTracking = true
			}
		case string(enableUseCountingTag):
			_, ok = row["value"]
			if ok && row["value"].(string) == string(enabledEnableUseCountingTagValue) {
				pointDoorInfo.EnableUseCounting = true
			}
		case string(isEOTTag):
			_, ok = row["value"]
			if ok && row["value"].(string) == string(EOTisEOTTagValue) {
				pointDoorInfo.IsEOT = true
			}
		case string(availabilityIDTag):
			_, ok = row["value"]
			if ok {
				pointDoorInfo.AvailabilityID = row["value"].(string)
			}
		case string(resetIDTag):
			_, ok = row["value"]
			if ok {
				pointDoorInfo.ResetID = row["value"].(string)
			}
		*/
		case string(currentUsesColName):
			_, ok = row["value"]
			if ok {
				lastCurrentUseCount, _ = strconv.Atoi(row["value"].(string))
				pointLastProcessedData.CurrentUses = lastCurrentUseCount
			}
		case string(totalUsesColName):
			_, ok = row["value"]
			if ok {
				lastTotalUseCount, _ = strconv.Atoi(row["value"].(string))
				pointLastProcessedData.TotalUses = lastTotalUseCount
			}
		case string(cubicleOccupancyColName):
			_, ok = row["value"]
			if ok {
				lastOccupancy, _ = strconv.Atoi(row["value"].(string))
				pointLastProcessedData.CubicleOccupancy = lastOccupancy
			}
		case string(pendingStatusColName):
			_, ok = row["value"]
			if ok {
				lastPendingStatus, _ = strconv.Atoi(row["value"].(string))
				pointLastProcessedData.PendingStatus = lastPendingStatus
			}
		case string(overdueStatusColName):
			_, ok = row["value"]
			if ok {
				lastOverdueStatus, _ = strconv.Atoi(row["value"].(string))
				pointLastProcessedData.OverdueStatus = lastOverdueStatus
			}
		case string(toPendingColName):
			_, ok = row["timestamp"]
			if ok {
				pointLastProcessedData.LastToPendingTimestamp = row["timestamp"].(string)
			}
		}
	}

	// get the door type
	if pointDoorInfo.IsEOT {
		switch pointDoorInfo.DoorTypeTag {
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
		}
	} else {
		switch pointDoorInfo.DoorTypeTag {
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
		}
	}
	return
}
