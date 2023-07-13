package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-os/plugin/nube/database/postgres/pgmodel"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"os"
	"strconv"
	"strings"
	"time"
)

// MakeDailyResetsDF creates a DF of daily reset triggers for the given time period.
func (inst *Instance) MakeDailyResetsDF(start, end time.Time, thresholdsDF dataframe.DataFrame) (*dataframe.DataFrame, error) {
	// Extract the allAreaResetTime from thresholdsDF
	allAreaResetTime := thresholdsDF.Col(string(allAreaResetTimeColName)).Elem(0).String()
	allAreaResetTimeSplit := strings.Split(allAreaResetTime, ":")
	allAreaResetTimeHour, err := strconv.Atoi(allAreaResetTimeSplit[0])
	if err != nil {
		return nil, err
	}
	allAreaResetTimeMins, err := strconv.Atoi(allAreaResetTimeSplit[1])
	if err != nil {
		return nil, err
	}

	timeZone := thresholdsDF.Col(string(timeZoneColName)).Elem(0).String()
	timeZoneLocation, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, err
	}

	// Create an empty slices to store the timestamps and areaReset
	timestampsArray := make([]string, 0)
	areaResetArray := make([]int, 0)

	// Set the time to 22:00 for the start date
	start = time.Date(start.Year(), start.Month(), start.Day(), allAreaResetTimeHour, allAreaResetTimeMins, 0, 0, timeZoneLocation)

	// Iterate from the start date until the end date, adding 24 hours each iteration
	for date := start; date.Before(end); date = date.Add(24 * time.Hour) {
		dateTimestamp := date.Format(time.RFC3339Nano)
		// dateTimestamp := date.UTC().Format(time.RFC3339Nano)
		timestampsArray = append(timestampsArray, dateTimestamp)
		areaResetArray = append(areaResetArray, 1)
	}

	// Convert the slice of timestamps to a series
	areaResetDF := dataframe.New(
		series.New(timestampsArray, series.String, string(timestampColName)),
		series.New(areaResetArray, series.Int, string(areaResetColName)),
	)

	return &areaResetDF, nil
}

// CalculateDoorUses calculates the totalUses, currentUses, cubicleOccupancy, pendingStatus, toClean, toPending of a door position sensors. doorPosDF must have door position.  lastValuesDF must have the last value for door position, occupancy, totalUses, currentUses, pendingStatus, and applicable use thresholds.
func (inst *Instance) CalculateDoorUses(dfRawDoorSensorHistories, resetsDF, thresholdsDF dataframe.DataFrame, pointLastProcessedData *LastProcessedData, pointDoorInfo *DoorInfo) (*dataframe.DataFrame, error) {
	var err error

	var joinedDF dataframe.DataFrame
	if dfRawDoorSensorHistories.Nrow() > 0 && resetsDF.Nrow() > 0 { // if both input dataframes aren't empty combine them
		dfRawDoorSensorHistories = dfRawDoorSensorHistories.Rename(string(doorPositionColName), "value")
		dfRawDoorSensorHistories = dfRawDoorSensorHistories.Drop([]string{string(pointUUIDColName), string(hostUUIDColName), "id"})
		joinedDF = dfRawDoorSensorHistories.OuterJoin(resetsDF, string(timestampColName))
		joinedDF = joinedDF.Arrange(dataframe.Sort(string(timestampColName)))
	} else if resetsDF.Nrow() > 0 {
		joinedDF = resetsDF.Arrange(dataframe.Sort(string(timestampColName)))
	} else if dfRawDoorSensorHistories.Nrow() > 0 {
		joinedDF = dfRawDoorSensorHistories.Arrange(dataframe.Sort(string(timestampColName)))
		joinedDF = joinedDF.Rename(string(doorPositionColName), "value")
		joinedDF = joinedDF.Drop([]string{string(pointUUIDColName), string(hostUUIDColName)})
	}
	fmt.Println("CalculateDoorUses() joinedDF")
	fmt.Println(joinedDF)

	// TODO: DELETE ME (just for debug)
	ResultFile, err := os.Create(fmt.Sprintf("/home/marc/Documents/Nube/CPS/Development/Data_Processing/6_joinedDF.csv"))
	if err != nil {
		inst.cpsErrorMsg(err)
	}
	defer ResultFile.Close()
	joinedDF.WriteCSV(ResultFile)

	// Extract the timestamp column as a series
	timestampSeries := joinedDF.Col(string(timestampColName))

	// check which columns exist
	columnNames := joinedDF.Names()
	doorPositionColumnExists := false
	resetColumnExists := false
	for _, cName := range columnNames {
		if cName == string(doorPositionColName) {
			doorPositionColumnExists = true
			continue
		}
		if cName == string(areaResetColName) {
			resetColumnExists = true
			continue
		}
	}

	// Extract the door position column as a series (if it exists)
	doorPositionSeries := joinedDF.Col(string(doorPositionColName))

	// Extract the reset column as a series
	resetSeries := joinedDF.Col(string(areaResetColName))

	// get the use threshold from the site thresholds df
	useThreshold := 10
	switch pointDoorInfo.DoorTypeID {
	case facilityEntrance:
		useThreshold, err = thresholdsDF.Col(string(facilityEntranceUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case facilityToilet:
		useThreshold, err = thresholdsDF.Col(string(facilityToiletUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case facilityDDA:
		useThreshold, err = thresholdsDF.Col(string(facilityDDAUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case eotEntrance:
		useThreshold, err = thresholdsDF.Col(string(eotEntranceUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case eotToilet:
		useThreshold, err = thresholdsDF.Col(string(eotToiletUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case eotShower:
		useThreshold, err = thresholdsDF.Col(string(eotShowerUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	case eotDDA:
		useThreshold, err = thresholdsDF.Col(string(eotDDAUseThresholdColName)).Elem(0).Int()
		if err != nil {
			return nil, err
		}
	}
	inst.cpsDebugMsg("CalculateDoorUses() useThreshold: ", useThreshold)

	totalUsesArray := make([]int, 0)
	currentUsesArray := make([]int, 0)
	occupancyArray := make([]int, 0)
	pendingStatusArray := make([]int, 0)
	toCleanArray := make([]int, 0)
	toPendingArray := make([]int, 0)
	// cleaningTimeArray := make([]string, 0)
	cleaningTimeArray := make([]int, 0)

	lastToPending, _ := time.Parse(time.RFC3339Nano, pointLastProcessedData.LastToPendingTimestamp)
	inst.cpsDebugMsg("CalculateDoorUses() lastToPending: ", lastToPending)
	invalidLastToPending := false
	if pointLastProcessedData.LastToPendingTimestamp == "" || lastToPending.After(time.Now()) {
		invalidLastToPending = true
	}

	lastPendingStatus := pointLastProcessedData.PendingStatus
	lastCurrentUseCount := pointLastProcessedData.CurrentUses
	lastTotalUseCount := pointLastProcessedData.TotalUses
	occupancy := pointLastProcessedData.CubicleOccupancy
	lastPosition := pointLastProcessedData.DoorPosition

	inst.cpsDebugMsg("CalculateDoorUses() pointDoorInfo.NormalPosition == normallyOpen: ", pointDoorInfo.NormalPosition == normallyOpen)

	for i, v := range timestampSeries.Records() {
		entryTime, _ := time.Parse(time.RFC3339Nano, v)
		// inst.cpsDebugMsg("CalculateDoorUses() entryTime: ", entryTime)

		doorPosition := "NaN"
		if doorPositionColumnExists {
			doorPosition = doorPositionSeries.Elem(i).String()
		}

		toPending := 0
		toClean := 0
		// cleaningTime := ""
		cleaningTime := 0

		if resetColumnExists {
			if resetVal, _ := resetSeries.Elem(i).Int(); resetVal == 1 { // This is a reset row.
				inst.cpsDebugMsg("CalculateDoorUses() RESET ROW")
				if lastPendingStatus == 1 && !invalidLastToPending {
					toClean = 1
					// cleaningTime = entryTime.Sub(lastToPending).String()
					cleaningTime = int(entryTime.Sub(lastToPending).Seconds())
					inst.cpsDebugMsg("CalculateDoorUses() cleaningTime: ", cleaningTime)
					// lastToClean = entryTime
				}
				lastCurrentUseCount = 0
				lastPendingStatus = 0

				if doorPosition == "NaN" { // this pushes series values if there is no data in the door position column
					totalUsesArray = append(totalUsesArray, lastTotalUseCount)
					occupancyArray = append(occupancyArray, occupancy)
					currentUsesArray = append(currentUsesArray, lastCurrentUseCount)
					pendingStatusArray = append(pendingStatusArray, lastPendingStatus)
					toCleanArray = append(toCleanArray, toClean)
					toPendingArray = append(toPendingArray, toPending)
					cleaningTimeArray = append(cleaningTimeArray, cleaningTime)
					continue
				}
			}
		}

		if doorPosition == "NaN" { // no door data, could be a reset, or bad data
			// still need to push values to the series arrays
			totalUsesArray = append(totalUsesArray, lastTotalUseCount)
			occupancyArray = append(occupancyArray, occupancy)
			currentUsesArray = append(currentUsesArray, lastCurrentUseCount)
			pendingStatusArray = append(pendingStatusArray, lastPendingStatus)
			toCleanArray = append(toCleanArray, toClean)
			toPendingArray = append(toPendingArray, toPending)
			cleaningTimeArray = append(cleaningTimeArray, cleaningTime)
			continue
		}

		doorStateFloat, err := strconv.ParseFloat(doorPosition, 64)
		if err != nil {
			inst.cpsDebugMsg("CalculateDoorUses() doorStateFloat, err := strconv.ParseFloat(doorPosition, 64) error: ", err)
		}
		doorState := int(doorStateFloat)
		if pointDoorInfo.NormalPosition == normallyOpen {
			if doorState == open && lastPosition == closed {
				lastTotalUseCount++
				lastCurrentUseCount++
			}
			if doorState == open {
				occupancy = vacant
			} else {
				occupancy = occupied
			}
		} else if pointDoorInfo.NormalPosition == normallyClosed {
			if doorState == closed && lastPosition == open && occupancy == vacant {
				occupancy = occupied
			} else if doorState == closed && lastPosition == open && occupancy == occupied {
				lastTotalUseCount++
				lastCurrentUseCount++
				occupancy = vacant
			}
		}

		lastPosition = doorState
		if lastCurrentUseCount >= useThreshold {
			if lastPendingStatus == 0 {
				toPending = 1
				lastToPending = entryTime
				invalidLastToPending = false
			}
			lastPendingStatus = 1
		}

		// append the new values to their respective series (to be joined later)
		totalUsesArray = append(totalUsesArray, lastTotalUseCount)
		occupancyArray = append(occupancyArray, occupancy)
		currentUsesArray = append(currentUsesArray, lastCurrentUseCount)
		pendingStatusArray = append(pendingStatusArray, lastPendingStatus)
		toCleanArray = append(toCleanArray, toClean)
		toPendingArray = append(toPendingArray, toPending)
		cleaningTimeArray = append(cleaningTimeArray, cleaningTime)
	}

	// Add count column to the dataframe
	resultDF := joinedDF.Mutate(series.New(totalUsesArray, series.Int, string(totalUsesColName)))
	resultDF = resultDF.Mutate(series.New(currentUsesArray, series.Int, string(currentUsesColName)))
	resultDF = resultDF.Mutate(series.New(occupancyArray, series.Int, string(cubicleOccupancyColName)))
	resultDF = resultDF.Mutate(series.New(pendingStatusArray, series.Int, string(pendingStatusColName)))
	resultDF = resultDF.Mutate(series.New(toCleanArray, series.Int, string(toCleanColName)))
	resultDF = resultDF.Mutate(series.New(toPendingArray, series.Int, string(toPendingColName)))
	// resultDF = resultDF.Mutate(series.New(cleaningTimeArray, series.String, string(cleaningTimeColName)))
	resultDF = resultDF.Mutate(series.New(cleaningTimeArray, series.Int, string(cleaningTimeColName)))
	// resultDF = resultDF.Select([]string{string(timestampColName), string(doorPositionColName), string(cubicleOccupancyColName), string(totalUsesColName), string(currentUsesColName), string(pendingStatusColName), string(toCleanColName), string(toPendingColName), string(areaResetColName), string(cleaningTimeColName)})
	if resetsDF.Nrow() > 0 {
		resultDF = resultDF.Rename(string(cleaningResetColName), string(areaResetColName))
	}
	return &resultDF, nil
}

// Calculate15MinUsageRollup creates a DF that adds timestamps at 0, 15, 30, 45 which rollup the sensor usage counts by 15 min periods.
func (inst *Instance) Calculate15MinUsageRollup(start, end time.Time, dfUsageCalculationResults dataframe.DataFrame, last15MinIntervalTotalUsesHistory *History, lastTotalUsesHistoryFound bool, timeZone string) (*dataframe.DataFrame, error) {
	last15MinRollupTotalUseCount := 0
	if lastTotalUsesHistoryFound {
		last15MinRollupTotalUseCount = int(last15MinIntervalTotalUsesHistory.Value)
	}

	// Create an empty slices to store the new timestamps
	timestampsArray := make([]string, 0)
	timeZoneLocation, err := time.LoadLocation(timeZone)
	if err != nil {
		inst.cpsErrorMsg("Calculate15MinUsageRollup() timezone error: ", err)
	}

	// Set the time for the first entry
	startRounded := start.Round(time.Minute * 15)
	if startRounded.Before(start) {
		startRounded = startRounded.Add(time.Minute * 15)
	}

	// Iterate from the start date until the end date, adding 15 mins each iteration
	for date := startRounded; date.Before(end) || date.Equal(end); date = date.Add(time.Minute * 15) {
		// dateTimestamp := date.UTC().Format(time.RFC3339Nano)
		dateTimestamp := date.In(timeZoneLocation).Format(time.RFC3339Nano)
		timestampsArray = append(timestampsArray, dateTimestamp)
	}
	rollupTimestampsDF := dataframe.New(
		series.New(timestampsArray, series.String, string(timestampColName)),
	)

	// join the processed data DF with the 15 min rollup timestampsArray.  Now we have all the timestamps that we need
	joinedDF := dfUsageCalculationResults.OuterJoin(rollupTimestampsDF, string(timestampColName))
	joinedDF = joinedDF.Arrange(dataframe.Sort(string(timestampColName)))

	// Extract the timestamp column as a series
	timestampSeries := joinedDF.Col(string(timestampColName))
	// inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() Calculate15MinUsageRollup() timestampSeries: %+v", timestampSeries))

	// Extract the totalUses column as a series
	totalUseCountSeries := joinedDF.Col(string(totalUsesColName))
	// inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() Calculate15MinUsageRollup() totalUseCountSeries: %+v", totalUseCountSeries))

	lastEntryTotalUses := last15MinRollupTotalUseCount
	totalUseNaNArray := totalUseCountSeries.IsNaN() // boolean slice indicating which values are NaN
	resultTimestampsArray := make([]string, 0)
	resultUsesRollupArray := make([]int, 0)
	for i, v := range timestampSeries.Records() {
		entryTime, _ := time.Parse(time.RFC3339Nano, v)
		// inst.cpsDebugMsg("Calculate15MinUsageRollup() entryTime: ", entryTime)
		if totalUseNaNArray[i] != true { // there is a totalUses count stored on this timestamp
			lastEntryTotalUses, _ = totalUseCountSeries.Elem(i).Int()
			if !lastTotalUsesHistoryFound {
				last15MinRollupTotalUseCount = lastEntryTotalUses
				lastTotalUsesHistoryFound = true
			}
		}
		if entryTime.Minute()%15 != 0 || entryTime.Second() != 0 {
			continue
		}

		// inst.cpsDebugMsg("Calculate15MinUsageRollup() last15MinRollupTotalUseCount: ", last15MinRollupTotalUseCount)
		// inst.cpsDebugMsg("Calculate15MinUsageRollup() lastEntryTotalUses: ", lastEntryTotalUses)
		usageRollup := lastEntryTotalUses - last15MinRollupTotalUseCount
		last15MinRollupTotalUseCount = lastEntryTotalUses
		if usageRollup < 0 {
			inst.cpsErrorMsg("Calculate15MinUsageRollup(): totalUses has decreased")
			return nil, errors.New("totalUses has decreased")
		}
		resultTimestampsArray = append(resultTimestampsArray, entryTime.Format(time.RFC3339))
		resultUsesRollupArray = append(resultUsesRollupArray, usageRollup)
	}

	// Convert the slice of timestamps to a series
	usesRollupDF := dataframe.New(
		series.New(resultTimestampsArray, series.String, string(timestampColName)),
		series.New(resultUsesRollupArray, series.Int, string(fifteenMinRollupUsesColName)),
	)

	// join the 15min usage rollup DF with existing DF.
	resultDF := joinedDF.OuterJoin(usesRollupDF, string(timestampColName))
	resultDF = resultDF.Arrange(dataframe.Sort(string(timestampColName)))

	return &resultDF, nil
}

// CalculateOverdueCubicles creates a DF that adds timestamps for overdueStatus and toOverdue
func (inst *Instance) CalculateOverdueCubicles(start, end time.Time, dfUsageCalculationResults, thresholdsDF dataframe.DataFrame, pointLastProcessedData *LastProcessedData, pointDoorInfo *DoorInfo) (*dataframe.DataFrame, error) {
	lastToPending, _ := time.Parse(time.RFC3339Nano, pointLastProcessedData.LastToPendingTimestamp)
	lastToClean, _ := time.Parse(time.RFC3339Nano, pointLastProcessedData.LastToCleanTimestamp)

	cleaningOverdueAlertDelay, err := GetOverdueDelay(pointDoorInfo.DoorTypeID, thresholdsDF)
	if err != nil {
		return nil, err
	}

	lastOverdueStatus := pointLastProcessedData.OverdueStatus

	// Check if there is an overdue event from the last pending time
	overdueEventPending := false
	var nextOverdueTime time.Time
	if pointLastProcessedData.LastToPendingTimestamp != "" {
		nextOverdueTime = lastToPending.Add(cleaningOverdueAlertDelay)
		if pointLastProcessedData.LastToCleanTimestamp != "" {
			if lastToPending.After(lastToClean) && nextOverdueTime.After(start) && nextOverdueTime.Before(end) { // checks that the cubicle is still pending, and that the overdue time would fall within the calculation time range
				overdueEventPending = true
			}
		} else {
			overdueEventPending = true
		}
	}
	// inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() CalculateOverdueCubicles() nextOverdueTime: %+v", nextOverdueTime))

	// Extract the timestamp column as a series
	existingTimestampSeries := dfUsageCalculationResults.Col(string(timestampColName))

	// Extract the pending status column as a series
	pendingStatusSeries := dfUsageCalculationResults.Col(string(pendingStatusColName))

	// Create an empty slices to store the new timestamps
	newTimestampsArray := make([]string, 0)
	toOverdueArray := make([]int, 0)
	overdueStatusArray := make([]int, 0)
	pendingStatusArray := make([]int, 0)

	pendingStatusNaNArray := pendingStatusSeries.IsNaN() // boolean slice indicating which values are NaN

	overdueEventInProgress := lastOverdueStatus != 0

	lastPendingStatus := pointLastProcessedData.PendingStatus
	for i, v := range existingTimestampSeries.Records() {
		entryTimestampSaved := false
		entryTime, _ := time.Parse(time.RFC3339Nano, v)
		// inst.cpsDebugMsg("CalculateOverdueCubicles() entryTime: ", entryTime)
		pendingStatus, _ := pendingStatusSeries.Elem(i).Int() // check the pending status of each entry
		// account for the first row, and other NaN rows
		if pendingStatusNaNArray[i] == true {
			pendingStatus = lastPendingStatus
		}
		if overdueEventPending { // looking for an overdue event
			if entryTime.Before(nextOverdueTime) && pendingStatus != 0 {
				// not overdue yet, just push 0 values for existing timestamp
				newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339Nano))
				toOverdueArray = append(toOverdueArray, 0)
				overdueStatusArray = append(overdueStatusArray, 0)
				pendingStatusArray = append(pendingStatusArray, 1)
				entryTimestampSaved = true
			} else if (entryTime.After(nextOverdueTime) || entryTime.Equal(nextOverdueTime)) && pendingStatus != 0 {
				// overdue delay has expired, so create a new timestamp and toOverdue entry
				newTimestampsArray = append(newTimestampsArray, nextOverdueTime.Format(time.RFC3339Nano))
				toOverdueArray = append(toOverdueArray, 1)
				overdueStatusArray = append(overdueStatusArray, 1)
				pendingStatusArray = append(pendingStatusArray, 1)
				overdueEventPending = false
				overdueEventInProgress = true
				entryTimestampSaved = true
				if !entryTime.Equal(nextOverdueTime) {
					// also add data for the existing timestamp
					newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339Nano))
					toOverdueArray = append(toOverdueArray, 0)
					overdueStatusArray = append(overdueStatusArray, 1)
					pendingStatusArray = append(pendingStatusArray, 1)
				}
			}
		} else if overdueEventInProgress && pendingStatus != 0 {
			// still overdue
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339Nano))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 1)
			pendingStatusArray = append(pendingStatusArray, 1)
			entryTimestampSaved = true
		}
		if pendingStatus == 0 {
			// not pending, or has been reset/cleaned. reset overdueStatus.
			overdueEventPending = false
			overdueEventInProgress = false
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339Nano))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 0)
			pendingStatusArray = append(pendingStatusArray, 0)
			entryTimestampSaved = true
		}
		if !entryTimestampSaved && pendingStatus != 0 {
			// has become pending
			overdueEventPending = true
			nextOverdueTime = entryTime.Add(cleaningOverdueAlertDelay)
			overdueEventInProgress = false
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339Nano))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 0)
			pendingStatusArray = append(pendingStatusArray, 1)
			entryTimestampSaved = true
		}
		lastPendingStatus = pendingStatus
	}

	// check for overdue between the last existing record and the end of the period
	if overdueEventPending && nextOverdueTime.After(start) && nextOverdueTime.Before(end) {
		newTimestampsArray = append(newTimestampsArray, nextOverdueTime.Format(time.RFC3339Nano))
		toOverdueArray = append(toOverdueArray, 1)
		overdueStatusArray = append(overdueStatusArray, 1)
		pendingStatusArray = append(pendingStatusArray, 1)
	}

	// Convert the slice of timestamps to a series
	overdueDF := dataframe.New(
		series.New(newTimestampsArray, series.String, string(timestampColName)),
		series.New(toOverdueArray, series.Int, string(toOverdueColName)),
		series.New(overdueStatusArray, series.Int, string(overdueStatusColName)),
		series.New(pendingStatusArray, series.Int, string(pendingStatusColName)),
	)

	// join overdue DF with existing DF.
	dfUsageCalculationResults = dfUsageCalculationResults.Drop(string(pendingStatusColName))
	resultDF := dfUsageCalculationResults.OuterJoin(overdueDF, string(timestampColName))
	resultDF = resultDF.Arrange(dataframe.Sort(string(timestampColName)))
	return &resultDF, nil
}

// PackageProcessedHistories ingests processed data DF, and outputs histores to be sent to the CPS postgres database
func (inst *Instance) PackageProcessedHistories(dfProcessingResults dataframe.DataFrame, thisAssetProcessedDataPoints []DoorProcessingPoint) (processedHistories []*pgmodel.History, latestPendingStatus, latestOverdueStatus *int, err error) {
	processedHistories = make([]*pgmodel.History, 0)

	// TODO: deal with history IDs for last sync

	// Extract the timestamp column as a series
	timestampSeries := dfProcessingResults.Col(string(timestampColName))
	// inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() PackageProcessedHistories() timestampSeries: %+v", timestampSeries))

	// check which series are available in the processing results dataframe
	columnNames := dfProcessingResults.Names()
	cubicleOccupancyColumnExists := false
	totalUsesColumnExists := false
	currentUsesColumnExists := false
	pendingStatusColumnExists := false
	overdueStatusColumnExists := false
	toPendingColumnExists := false
	toCleanColumnExists := false
	toOverdueColumnExists := false
	cleaningResetColumnExists := false
	cleaningTimeColumnExists := false
	fifteenMinRollupUsesColumnExists := false

	for _, cName := range columnNames {
		switch cName {
		case string(cubicleOccupancyColName):
			cubicleOccupancyColumnExists = true
		case string(totalUsesColName):
			totalUsesColumnExists = true
		case string(currentUsesColName):
			currentUsesColumnExists = true
		case string(pendingStatusColName):
			pendingStatusColumnExists = true
		case string(overdueStatusColName):
			overdueStatusColumnExists = true
		case string(toPendingColName):
			toPendingColumnExists = true
		case string(toCleanColName):
			toCleanColumnExists = true
		case string(toOverdueColName):
			toOverdueColumnExists = true
		case string(cleaningResetColName):
			cleaningResetColumnExists = true
		case string(cleaningTimeColName):
			cleaningTimeColumnExists = true
		case string(fifteenMinRollupUsesColName):
			fifteenMinRollupUsesColumnExists = true
		}
	}

	// loop through the processed results DF and make histories
	for _, pdp := range thisAssetProcessedDataPoints {
		switch pdp.Name {
		case string(cubicleOccupancyColName):
			if !cubicleOccupancyColumnExists {
				continue
			}
			occupancySeries := dfProcessingResults.Col(string(cubicleOccupancyColName))
			for i, ts := range timestampSeries.Records() {
				element := occupancySeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(totalUsesColName):
			if !totalUsesColumnExists {
				continue
			}
			totalUsesSeries := dfProcessingResults.Col(string(totalUsesColName))
			// TODO: could implement a last value check
			for i, ts := range timestampSeries.Records() {
				element := totalUsesSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(currentUsesColName):
			if !currentUsesColumnExists {
				continue
			}
			currentUsesSeries := dfProcessingResults.Col(string(currentUsesColName))
			lastVal := 0
			lastValSet := false
			for i, ts := range timestampSeries.Records() {
				element := currentUsesSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if lastValSet && value == lastVal {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
				lastVal = value
			}
		case string(pendingStatusColName):
			if !pendingStatusColumnExists {
				continue
			}
			pendingStatusSeries := dfProcessingResults.Col(string(pendingStatusColName))
			lastVal := 0
			lastValSet := false
			for i, ts := range timestampSeries.Records() {
				element := pendingStatusSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if lastValSet && value == lastVal {
					continue
				}
				latestPendingStatus = &value
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
				lastVal = value
			}
		case string(overdueStatusColName):
			if !overdueStatusColumnExists {
				continue
			}
			overdueStatusSeries := dfProcessingResults.Col(string(overdueStatusColName))
			lastVal := 0
			lastValSet := false
			for i, ts := range timestampSeries.Records() {
				element := overdueStatusSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if lastValSet && value == lastVal {
					continue
				}
				latestOverdueStatus = &value
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
				lastVal = value
			}
		case string(toPendingColName):
			if !toPendingColumnExists {
				continue
			}
			toPendingSeries := dfProcessingResults.Col(string(toPendingColName))
			for i, ts := range timestampSeries.Records() {
				element := toPendingSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == 0 {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(toCleanColName):
			if !toCleanColumnExists {
				continue
			}
			toCleanSeries := dfProcessingResults.Col(string(toCleanColName))
			for i, ts := range timestampSeries.Records() {
				element := toCleanSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == 0 {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(toOverdueColName):
			if !toOverdueColumnExists {
				continue
			}
			toOverdueSeries := dfProcessingResults.Col(string(toOverdueColName))
			for i, ts := range timestampSeries.Records() {
				element := toOverdueSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == 0 {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(cleaningResetColName):
			if !cleaningResetColumnExists {
				continue
			}
			cleaningResetSeries := dfProcessingResults.Col(string(cleaningResetColName))
			for i, ts := range timestampSeries.Records() {
				element := cleaningResetSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == 0 {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(cleaningTimeColName):
			if !cleaningTimeColumnExists {
				continue
			}
			cleaningTimeSeries := dfProcessingResults.Col(string(cleaningTimeColName))
			for i, ts := range timestampSeries.Records() {
				element := cleaningTimeSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == 0 {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(fifteenMinRollupUsesColName):
			if !fifteenMinRollupUsesColumnExists {
				continue
			}
			fifteenMinRollupSeries := dfProcessingResults.Col(string(fifteenMinRollupUsesColName))
			for i, ts := range timestampSeries.Records() {
				element := fifteenMinRollupSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339Nano, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		}
	}

	return
}
