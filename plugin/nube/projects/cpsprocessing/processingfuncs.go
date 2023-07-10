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
		dateTimestamp := date.Format(time.RFC3339)
		// dateTimestamp := date.UTC().Format(time.RFC3339)
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

// TODO: add door type and use thresholds lookup.  Probably pass in thresholds DF as it is per site, not per sensor

// CalculateDoorUses calculates the totalUses, currentUses, cubicleOccupancy, pendingStatus, toClean, toPending of a door position sensors. doorPosDF must have door position.  lastValuesDF must have the last value for door position, occupancy, totalUses, currentUses, pendingStatus, and applicable use thresholds.
func (inst *Instance) CalculateDoorUses(dfRawDoorSensorHistories, resetsDF, thresholdsDF dataframe.DataFrame, pointLastProcessedData *LastProcessedData, pointDoorInfo *DoorInfo) (*dataframe.DataFrame, error) {
	var err error

	joinedDF := dfRawDoorSensorHistories.OuterJoin(resetsDF, string(timestampColName))
	joinedDF = joinedDF.Arrange(dataframe.Sort(string(timestampColName)))
	joinedDF = joinedDF.Rename(string(doorPositionColName), "value")

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

	// Extract the door position column as a series
	doorPositionSeries := joinedDF.Col(string(doorPositionColName))

	inst.cpsDebugMsg("CalculateDoorUses() doorPositionSeries.Type: ", doorPositionSeries.Type())

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

	lastToPending, _ := time.Parse(time.RFC3339, pointLastProcessedData.LastToPendingTimestamp)
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

	for i, v := range doorPositionSeries.Records() {
		entryTime, _ := time.Parse(time.RFC3339, timestampSeries.Elem(i).String())
		inst.cpsDebugMsg("CalculateDoorUses() entryTime: ", entryTime)

		toPending := 0
		toClean := 0
		// cleaningTime := ""
		cleaningTime := 0

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

			if v == "NaN" { // this pushes series values if there is no data in the door position column
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

		if v == "NaN" { // no door data, could be a reset, or bad data
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

		doorStateFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			inst.cpsDebugMsg("CalculateDoorUses() doorStateFloat, err := strconv.ParseFloat(v, 64) error: ", err)
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
	resultDF = resultDF.Select([]string{string(timestampColName), string(doorPositionColName), string(cubicleOccupancyColName), string(totalUsesColName), string(currentUsesColName), string(pendingStatusColName), string(toCleanColName), string(toPendingColName), string(areaResetColName), string(cleaningTimeColName)})
	resultDF = resultDF.Rename(string(cleaningResetColName), string(areaResetColName))
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
		// dateTimestamp := date.UTC().Format(time.RFC3339)
		dateTimestamp := date.In(timeZoneLocation).Format(time.RFC3339)
		timestampsArray = append(timestampsArray, dateTimestamp)
	}
	inst.cpsDebugMsg("Calculate15MinUsageRollup() timestampsArray: ", timestampsArray)
	rollupTimestampsDF := dataframe.New(
		series.New(timestampsArray, series.String, string(timestampColName)),
	)

	// join the processed data DF with the 15 min rollup timestampsArray.  Now we have all the timestamps that we need
	joinedDF := dfUsageCalculationResults.OuterJoin(rollupTimestampsDF, string(timestampColName))
	joinedDF = joinedDF.Arrange(dataframe.Sort(string(timestampColName)))
	inst.cpsDebugMsg("Calculate15MinUsageRollup() joinedDF:")
	inst.cpsDebugMsg(joinedDF)

	// Extract the timestamp column as a series
	timestampSeries := joinedDF.Col(string(timestampColName))
	inst.cpsDebugMsg("Calculate15MinUsageRollup() timestampSeries:")
	inst.cpsDebugMsg(timestampSeries)

	// Extract the totalUses column as a series
	totalUseCountSeries := joinedDF.Col(string(totalUsesColName))
	inst.cpsDebugMsg("Calculate15MinUsageRollup() totalUseCountSeries:")
	inst.cpsDebugMsg(totalUseCountSeries)

	lastEntryTotalUses := last15MinRollupTotalUseCount
	totalUseNaNArray := totalUseCountSeries.IsNaN() // boolean slice indicating which values are NaN
	resultTimestampsArray := make([]string, 0)
	resultUsesRollupArray := make([]int, 0)
	for i, v := range timestampSeries.Records() {
		entryTime, _ := time.Parse(time.RFC3339, v)
		inst.cpsDebugMsg("Calculate15MinUsageRollup() entryTime: ", entryTime)
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

		inst.cpsDebugMsg("Calculate15MinUsageRollup() last15MinRollupTotalUseCount: ", last15MinRollupTotalUseCount)
		inst.cpsDebugMsg("Calculate15MinUsageRollup() lastEntryTotalUses: ", lastEntryTotalUses)
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
	lastToPending, _ := time.Parse(time.RFC3339, pointLastProcessedData.LastToPendingTimestamp)
	lastToClean, _ := time.Parse(time.RFC3339, pointLastProcessedData.LastToCleanTimestamp)

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

	for i, v := range existingTimestampSeries.Records() {
		entryTimestampSaved := false
		entryTime, _ := time.Parse(time.RFC3339, v)
		inst.cpsDebugMsg("CalculateOverdueCubicles() entryTime: ", entryTime)
		pendingStatus, _ := pendingStatusSeries.Elem(i).Int() // check the pending status of each entry
		if overdueEventPending {                              // looking for an overdue event
			if entryTime.Before(nextOverdueTime) && !pendingStatusNaNArray[i] && pendingStatus != 0 {
				inst.cpsDebugMsg("CalculateOverdueCubicles() 1")
				// not overdue yet, just push 0 values for existing timestamp
				newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339))
				toOverdueArray = append(toOverdueArray, 0)
				overdueStatusArray = append(overdueStatusArray, 0)
				pendingStatusArray = append(pendingStatusArray, 1)
				entryTimestampSaved = true
			} else if (entryTime.After(nextOverdueTime) || entryTime.Equal(nextOverdueTime)) && pendingStatus != 0 {
				inst.cpsDebugMsg("CalculateOverdueCubicles() 2")
				// overdue delay has expired, so create a new timestamp and toOverdue entry
				newTimestampsArray = append(newTimestampsArray, nextOverdueTime.Format(time.RFC3339))
				toOverdueArray = append(toOverdueArray, 1)
				overdueStatusArray = append(overdueStatusArray, 1)
				pendingStatusArray = append(pendingStatusArray, 1)
				overdueEventPending = false
				overdueEventInProgress = true
				entryTimestampSaved = true
				if !entryTime.Equal(nextOverdueTime) {
					inst.cpsDebugMsg("CalculateOverdueCubicles() 3")
					// also add data for the existing timestamp
					newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339))
					toOverdueArray = append(toOverdueArray, 0)
					overdueStatusArray = append(overdueStatusArray, 1)
					pendingStatusArray = append(pendingStatusArray, 1)
				}
			}
		} else if overdueEventInProgress && pendingStatus != 0 {
			inst.cpsDebugMsg("CalculateOverdueCubicles() 4")
			// still overdue
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 1)
			pendingStatusArray = append(pendingStatusArray, 1)
			entryTimestampSaved = true
		}
		if pendingStatus == 0 {
			inst.cpsDebugMsg("CalculateOverdueCubicles() 5")
			// has been reset/cleaned reset overdueStatus
			overdueEventPending = false
			overdueEventInProgress = false
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 0)
			pendingStatusArray = append(pendingStatusArray, 0)
			entryTimestampSaved = true
		}
		if !entryTimestampSaved && pendingStatus != 0 {
			inst.cpsDebugMsg("CalculateOverdueCubicles() 6")
			// has become pending
			overdueEventPending = true
			nextOverdueTime = entryTime.Add(cleaningOverdueAlertDelay)
			overdueEventInProgress = false
			newTimestampsArray = append(newTimestampsArray, entryTime.Format(time.RFC3339))
			toOverdueArray = append(toOverdueArray, 0)
			overdueStatusArray = append(overdueStatusArray, 0)
			pendingStatusArray = append(pendingStatusArray, 1)
			entryTimestampSaved = true
		}
	}

	// check for overdue between the last existing record and the end of the period
	if nextOverdueTime.Before(end) {
		inst.cpsDebugMsg("CalculateOverdueCubicles() 7")
		newTimestampsArray = append(newTimestampsArray, nextOverdueTime.Format(time.RFC3339))
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
func (inst *Instance) PackageProcessedHistories(dfProcessingResults dataframe.DataFrame, thisAssetProcessedDataPoints []DoorProcessingPoint) ([]*pgmodel.History, error) {
	processedHistories := make([]*pgmodel.History, 0)

	// TODO: deal with history IDs for last sync

	// Extract the timestamp column as a series
	timestampSeries := dfProcessingResults.Col(string(timestampColName))
	inst.cpsDebugMsg("PackageProcessedHistories() timestampSeries:")
	inst.cpsDebugMsg(timestampSeries)

	// loop through the processed results DF and make histories
	for _, pdp := range thisAssetProcessedDataPoints {
		switch pdp.Name {
		case string(cubicleOccupancyColName):
			occupancySeries := dfProcessingResults.Col(string(cubicleOccupancyColName))
			for i, ts := range timestampSeries.Records() {
				element := occupancySeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(totalUsesColName):
			totalUsesSeries := dfProcessingResults.Col(string(totalUsesColName))
			for i, ts := range timestampSeries.Records() {
				element := totalUsesSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(currentUsesColName):
			currentUsesSeries := dfProcessingResults.Col(string(currentUsesColName))
			lastVal := 0
			for i, ts := range timestampSeries.Records() {
				element := currentUsesSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == lastVal {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(pendingStatusColName):
			pendingStatusSeries := dfProcessingResults.Col(string(pendingStatusColName))
			lastVal := 0
			for i, ts := range timestampSeries.Records() {
				element := pendingStatusSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == lastVal {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(overdueStatusColName):
			overdueStatusSeries := dfProcessingResults.Col(string(overdueStatusColName))
			lastVal := 0
			for i, ts := range timestampSeries.Records() {
				element := overdueStatusSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				if value == lastVal {
					continue
				}
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(toPendingColName):
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
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(toCleanColName):
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
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}
		case string(toOverdueColName):
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
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(cleaningResetColName):
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
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(cleaningTimeColName):
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
				timestamp, _ := time.Parse(time.RFC3339, ts)
				newHist := pgmodel.History{
					PointUUID: pdp.UUID,
					HostUUID:  pdp.HostUUID,
					Value:     float.New(float64(value)),
					Timestamp: timestamp,
				}
				processedHistories = append(processedHistories, &newHist)
			}

		case string(fifteenMinRollupUsesColName):
			fifteenMinRollupSeries := dfProcessingResults.Col(string(fifteenMinRollupUsesColName))
			for i, ts := range timestampSeries.Records() {
				element := fifteenMinRollupSeries.Elem(i)
				if element.IsNA() {
					continue
				}
				value, _ := element.Int()
				timestamp, _ := time.Parse(time.RFC3339, ts)
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

	return processedHistories, nil
}
