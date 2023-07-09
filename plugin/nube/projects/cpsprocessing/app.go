package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/go-gota/gota/dataframe"
	"os"
	"strings"
	"time"
)

func (inst *Instance) CPSProcessing() {
	inst.cpsDebugMsg("CPSProcessing()")

	periodStart, _ := time.Parse(time.RFC3339, "2023-06-25T07:50:00Z")
	periodEnd, _ := time.Parse(time.RFC3339, "2023-06-26T14:00:00Z")

	// get site thresholds from thresholds table in pg database
	var allSiteThresholds []Threshold
	err := postgresSetting.postgresConnectionInstance.db.Raw(`
			SELECT DISTINCT ON (site_ref) *
			FROM thresholds
			ORDER BY site_ref, updated_at DESC
		`).
		Scan(&allSiteThresholds).Error
	if err != nil {
		inst.cpsErrorMsg("CPSProcessing() allSiteThresholds error: ", err)
	}
	inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() allSiteThresholds: %+v", allSiteThresholds))
	// TODO: just for debug
	dfAllSiteThresholds := dataframe.LoadStructs(allSiteThresholds)
	// dfSiteThresholds := dataframe.ReadCSV(strings.NewReader(csvSiteThresholds))
	fmt.Println("dfAllSiteThresholds")
	fmt.Println(dfAllSiteThresholds)

	// get a list of point UUIDs for door sensors
	var doorPointUUIDs []string
	_ = postgresSetting.postgresConnectionInstance.db.Table("point_meta_tags").
		Select("point_uuid").
		Where("key = ? AND value IN ?", "assetFunc", []string{string(managedCubicleDoorSensorAssetFunctionTagValue), string(managedFacilityEntranceDoorSensorAssetFunctionTagValue), string(usageCountDoorSensorAssetFunctionTagValue)}).
		Find(&doorPointUUIDs)

	// doorPointUUIDs = []string{"pnt_3dc020ef688e4ae1", "pnt_d6a078410aae48c1", "pnt_75610d0f08ab41b2"} // TODO: just for testing
	inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() doorPointUUIDs: %+v", doorPointUUIDs))

	// get the required meta tags grouped by point_uuid and host uuid
	var doorPointsAndTags []DoorProcessingPoint
	err = postgresSetting.postgresConnectionInstance.db.Raw(`
			SELECT
			p.uuid AS point_uuid,
			p.name,
			p.host_uuid,
			MAX(CASE WHEN m.key = 'siteRef' THEN m.value END) AS site_ref,
			MAX(CASE WHEN m.key = 'assetRef' THEN m.value END) AS asset_ref,
			MAX(CASE WHEN m.key = 'assetFunc' THEN m.value END) AS asset_func,
			MAX(CASE WHEN m.key = 'floorRef' THEN m.value END) AS floor_ref,
			MAX(CASE WHEN m.key = 'genderRef' THEN m.value END) AS gender_ref,
			MAX(CASE WHEN m.key = 'locationRef' THEN m.value END) AS location_ref,
			MAX(CASE WHEN m.key = 'pointFunction' THEN m.value END) AS point_function,
			MAX(CASE WHEN m.key = 'measurementRef' THEN m.value END) AS measurement_ref,
			MAX(CASE WHEN m.key = 'doorType' THEN m.value END) AS door_type,
			MAX(CASE WHEN m.key = 'normalPosition' THEN m.value END) AS normal_position,
			MAX(CASE WHEN m.key = 'enableCleaningTracking' THEN m.value END) AS enable_cleaning_cracking,
			MAX(CASE WHEN m.key = 'enableUseCounting' THEN m.value END) AS enable_use_counting,
			MAX(CASE WHEN m.key = 'isEOT' THEN m.value END) AS is_eot,
			MAX(CASE WHEN m.key = 'availabilityID' THEN m.value END) AS availability_id,
			MAX(CASE WHEN m.key = 'resetID' THEN m.value END) AS reset_id
		FROM
			points AS p
		JOIN
			point_meta_tags AS m ON p.uuid = m.point_uuid
		WHERE
			p.uuid IN (?)
		GROUP BY
			p.uuid, p.host_uuid
		`, doorPointUUIDs).
		Find(&doorPointsAndTags).Error
	if err != nil {
		inst.cpsErrorMsg("CPSProcessing() doorPointsAndTags error: ", err)
	}
	// -- dataframe implementation --
	inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() doorPointsAndTags: %+v", doorPointsAndTags))
	// TODO: just for debug
	dfDoorPointsAndTags := dataframe.LoadStructs(doorPointsAndTags)
	fmt.Println("dfDoorPointsAndTags")
	fmt.Println(dfDoorPointsAndTags)

	// get reset points and histories
	var resetPointUUIDs []string
	_ = postgresSetting.postgresConnectionInstance.db.Table("point_meta_tags").
		Select("point_uuid").
		Where("(key = ? AND value = ?)", "pointFunction", string(doorResetPointFunctionTagValue)).
		Find(&resetPointUUIDs)

	// resetPointUUIDs = []string{"", "", ""} // TODO: just for testing
	inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() resetPointUUIDs: %+v", resetPointUUIDs))

	// get the required reset point meta tags grouped by point_uuid and host uuid
	var doorResetPointsAndTags []DoorResetPoint
	err = postgresSetting.postgresConnectionInstance.db.Raw(`
			SELECT
				p.uuid AS point_uuid,
				p.name,
				p.host_uuid,
				MAX(CASE WHEN m.key = 'siteRef' THEN m.value END) AS site_ref,
				MAX(CASE WHEN m.key = 'pointFunction' THEN m.value END) AS point_function,
				MAX(CASE WHEN m.key = 'measurementRef' THEN m.value END) AS measurement_ref,
				MAX(CASE WHEN m.key = 'isEOT' THEN m.value END) AS is_eot,
				MAX(CASE WHEN m.key = 'resetID' THEN m.value END) AS reset_id
			FROM
				points AS p
			JOIN
				point_meta_tags AS m ON p.uuid = m.point_uuid
			WHERE
				p.uuid IN (?)
			GROUP BY
				p.uuid, p.host_uuid
				`, resetPointUUIDs).
		Find(&doorResetPointsAndTags).Error
	if err != nil {
		inst.cpsErrorMsg("CPSProcessing() doorResetPointsAndTags error: ", err)
	}
	// -- dataframe implementation --
	inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() doorResetPointsAndTags: %+v", doorResetPointsAndTags))
	// TODO: just for debug
	dfDoorResetPointsAndTags := dataframe.LoadStructs(doorResetPointsAndTags)
	fmt.Println("dfDoorResetPointsAndTags")
	fmt.Println(dfDoorResetPointsAndTags)

	// do processing steps for each site
	// TODO: allow for only processing select sites based on tags
	for _, siteThresholds := range allSiteThresholds {
		siteRef := siteThresholds.SiteRef
		inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() siteRef: %+v", siteRef))

		// this site thresholds as a dataframe
		siteThresholdsToDF := make([]Threshold, 0)
		siteThresholdsToDF = append(siteThresholdsToDF, siteThresholds)
		dfSiteThresholds := dataframe.LoadStructs(siteThresholdsToDF)
		fmt.Println("dfSiteThresholds")
		fmt.Println(dfSiteThresholds)

		// get the points for this site
		thisSiteDoorPointsAndTags := make([]DoorProcessingPoint, 0)
		for _, point := range doorPointsAndTags {
			if point.SiteRef == siteRef {
				thisSiteDoorPointsAndTags = append(thisSiteDoorPointsAndTags, point)
			}
		}
		// -- dataframe implementation --
		// dfThisSitePointsAndTags := dfDoorPointsAndTags.Filter(dataframe.F{Colname: "site_ref", Comparator: series.Eq, Comparando: siteRef})
		// TODO: just for debug
		dfThisSiteDoorPointsAndTags := dataframe.LoadStructs(thisSiteDoorPointsAndTags)
		fmt.Println("dfThisSiteDoorPointsAndTags")
		fmt.Println(dfThisSiteDoorPointsAndTags)

		// get the raw door sensor points for this site (meta-tag = pointFunction: "sensor")
		// TODO: allow for only processing select assets based on tags
		thisSiteDoorSensorPointsAndTags := make([]DoorProcessingPoint, 0)
		for _, point := range doorPointsAndTags {
			if point.PointFunction == string(rawDoorSensorPointFunctionTagValue) && point.MeasurementRef == string(doorSensorMeasurementRefTagValue) {
				thisSiteDoorSensorPointsAndTags = append(thisSiteDoorSensorPointsAndTags, point)
			}
		}
		// -- dataframe implementation --
		// dfThisSiteSensorPointsAndTags := dfThisSitePointsAndTags.Filter(dataframe.F{Colname: pointFunctionTag, Comparator: series.Eq, Comparando: rawDoorSensorPointFunctionTagValue})
		// TODO: just for debug
		dfThisSiteDoorSensorPointsAndTags := dataframe.LoadStructs(thisSiteDoorSensorPointsAndTags)
		fmt.Println("dfThisSiteDoorSensorPointsAndTags")
		fmt.Println(dfThisSiteDoorSensorPointsAndTags)

		// iterate through the raw points, push histories (to points that exist), and update the relevant point values (for app)
		for i, doorSensorPoint := range thisSiteDoorSensorPointsAndTags {
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() point: %+v, uuid: %v, host: %v", doorSensorPoint.Name, doorSensorPoint.UUID, doorSensorPoint.HostUUID))
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() doorSensorPoint: %+v", doorSensorPoint))

			if doorSensorPoint.AssetRef == "" {
				inst.cpsErrorMsg(fmt.Sprintf("CPSProcessing() no assetRef tag on point: %v - %v", doorSensorPoint.Name, doorSensorPoint.UUID))
				continue
			}
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() assetRef: %+v", doorSensorPoint.AssetRef))

			// collect the processed data points for this asset
			thisAssetProcessedDataPoints := make([]DoorProcessingPoint, 0)
			for _, p := range thisSiteDoorPointsAndTags {
				if p.AssetRef == doorSensorPoint.AssetRef && p.PointFunction == string(processedDataPointFunctionTagValue) {
					thisAssetProcessedDataPoints = append(thisAssetProcessedDataPoints, p)
				}
			}
			if len(thisAssetProcessedDataPoints) <= 0 {
				inst.cpsErrorMsg(fmt.Sprintf("CPSProcessing() no processed data points for asset: (%v) %v - %v - %v", doorSensorPoint.AssetRef, doorSensorPoint.FloorRef, doorSensorPoint.GenderRef, doorSensorPoint.LocationRef))
				continue
			}
			// add in the current sensor point, because we need to get the last door position value along with the last value of the rest of the processed data points
			thisAssetProcessedDataPoints = append(thisAssetProcessedDataPoints, doorSensorPoint)
			// TODO: just for debug
			dfThisAssetProcessedDataPoints := dataframe.LoadStructs(thisAssetProcessedDataPoints)
			fmt.Println("dfThisAssetProcessedDataPoints")
			fmt.Println(dfThisAssetProcessedDataPoints)

			// verify that the point has all the required information and calculations should actually be done
			processedDataPointUUIDs := make([]string, 0)
			var cubicleOccupancyPoint, totalUsesPoint, currentUsesPoint, fifteenMinUsesPoint *DoorProcessingPoint
			// , pendingStatusPoint, overdueStatusPoint, toPendingPoint, toCleanPoint, toOverduePoint, cleaningTimePoint *DoorProcessingPoint
			for _, pdp := range thisAssetProcessedDataPoints {
				switch pdp.Name {
				case string(cubicleOccupancyColName):
					cubicleOccupancyPoint = &pdp
				case string(totalUsesColName):
					totalUsesPoint = &pdp
				case string(currentUsesColName):
					currentUsesPoint = &pdp
				case string(fifteenMinRollupUsesColName):
					fifteenMinUsesPoint = &pdp
					/*
						case string(pendingStatusColName):
							pendingStatusPoint = &pdp
						case string(overdueStatusColName):
							overdueStatusPoint = &pdp
						case string(toPendingColName):
							toPendingPoint = &pdp
						case string(toCleanColName):
							toCleanPoint = &pdp
						case string(toOverdueColName):
							toOverduePoint = &pdp
						case string(cleaningTimeColName):
							cleaningTimePoint = &pdp
					*/
				}
				// add point_uuid to the list of the processed data point uuids
				processedDataPointUUIDs = append(processedDataPointUUIDs, pdp.UUID)
			}
			// first verify the points exist for `usageCount` processing
			if cubicleOccupancyPoint == nil || totalUsesPoint == nil || currentUsesPoint == nil || fifteenMinUsesPoint == nil {
				inst.cpsErrorMsg(fmt.Sprintf("CPSProcessing() missing `usageCount` proccessed data point for asset: (%v) %v - %v - %v", doorSensorPoint.AssetRef, doorSensorPoint.FloorRef, doorSensorPoint.GenderRef, doorSensorPoint.LocationRef))
				continue
			}

			// pull data for sensor for the given time range
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() periodStart: %+v", periodStart))
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() periodEnd: %+v", periodEnd))

			var rawDoorSensorHistories []History
			err = postgresSetting.postgresConnectionInstance.db.Model(&model.History{}).Where("point_uuid = ? AND host_uuid = ? AND (timestamp AT TIME ZONE 'UTC' BETWEEN ? AND ?)", doorSensorPoint.UUID, doorSensorPoint.HostUUID, periodStart, periodEnd).Scan(&rawDoorSensorHistories).Error
			if err != nil {
				inst.cpsErrorMsg("CPSProcessing() rawSensorData error: ", err)
			}
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() rawSensorData: %+v", rawDoorSensorHistories))
			// dfRawDoor := dataframe.ReadCSV(strings.NewReader(csvRawDoor))
			dfRawDoorSensorHistories := dataframe.LoadStructs(rawDoorSensorHistories)
			fmt.Println("dfRawDoorSensorHistories")
			fmt.Println(dfRawDoorSensorHistories)

			// get last stored processed data values
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() hostUUID: %+v", doorSensorPoint.HostUUID))
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() processedDataPointUUIDs: %+v", processedDataPointUUIDs))
			var lastProcessedDataHistories []History
			// TODO: ensure that this query gets the LAST value from each processed history point
			err = postgresSetting.postgresConnectionInstance.db.Raw(`
					SELECT DISTINCT ON (point_uuid, host_uuid) *
					FROM histories
					WHERE host_uuid = ? AND point_uuid IN (?) AND timestamp AT TIME ZONE 'UTC' < ?
					ORDER BY point_uuid, host_uuid, timestamp DESC
				`, doorSensorPoint.HostUUID, processedDataPointUUIDs, periodEnd).
				Scan(&lastProcessedDataHistories).Error
			if err != nil {
				inst.cpsErrorMsg("CPSProcessing() lastProcessedData error: ", err)
			}
			/*  // for viewing the resulting SQL (DEBUG)
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() lastProcessedData SQL: %+v", postgresSetting.postgresConnectionInstance.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx.Raw(`
					SELECT DISTINCT ON (point_uuid, host_uuid) *
					FROM histories
					WHERE host_uuid = ? AND point_uuid IN (?) AND timestamp AT TIME ZONE 'UTC' < ?
					ORDER BY point_uuid, host_uuid, timestamp DESC
				`, doorSensorPoint.HostUUID, processedDataPointUUIDs, periodEnd).
					Scan(&lastProcessedDataHistories)
			})))
			*/
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() lastProcessedDataHistories: %+v", lastProcessedDataHistories))
			// dfLastProcessedDoor := dataframe.ReadCSV(strings.NewReader(csvLastProNODoor))
			dfLastProcessedDataHistories := dataframe.LoadStructs(lastProcessedDataHistories)
			fmt.Println("dfLastProcessedDataHistories")
			fmt.Println(dfLastProcessedDataHistories)

			// join the last values with the processed data points
			// var lastProcessedData LastProcessedData
			var dfJoinedLastProcessedValuesAndPoints dataframe.DataFrame
			if len(lastProcessedDataHistories) > 0 {
				dfJoinedLastProcessedValuesAndPoints = dfThisAssetProcessedDataPoints.OuterJoin(dfLastProcessedDataHistories, "point_uuid", "host_uuid")
				// for _, hist := range lastProcessedDataHistories {
			} else {
				dfJoinedLastProcessedValuesAndPoints = dfThisAssetProcessedDataPoints
				inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() lastProcessedDataHistories error: no last values were found, using zero values for processing"))
			}
			fmt.Println("dfJoinedLastProcessedValuesAndPoints")
			fmt.Println(dfJoinedLastProcessedValuesAndPoints)

			// TODO: DELETE ME (just for debug)
			if i == 0 {
				ResultFile, err := os.Create("/home/marc/Documents/Nube/CPS/Development/Data_Processing/5_dfJoinedLastProcessedValuesAndPoints.csv")
				if err != nil {
					inst.cpsErrorMsg(err)
				}
				defer ResultFile.Close()
				dfJoinedLastProcessedValuesAndPoints.WriteCSV(ResultFile)
			}

			// get reset point and histories applicable to this asset/point
			thisAssetResetDataPoints := make([]DoorResetPoint, 0)
			for _, rp := range doorResetPointsAndTags {
				if rp.ResetID == doorSensorPoint.ResetID {
					thisAssetResetDataPoints = append(thisAssetResetDataPoints, rp)
				}
			}
			var thisAssetResetDataPoint DoorResetPoint
			thisAssetResetDataPoint = thisAssetResetDataPoints[0] // TODO: consider the case where the point has been deleted and re-added (so there would be multiple points that match)
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() thisAssetResetDataPoints: %+v", thisAssetResetDataPoints))
			// -- dataframe implementation --
			// TODO: just for debug
			dfThisAssetResetDataPoints := dataframe.LoadStructs(thisAssetResetDataPoints)
			fmt.Println("dfThisAssetResetDataPoints")
			fmt.Println(dfThisAssetResetDataPoints)

			// get the history logs for the reset point for the calculation period
			var resetHistoryData []History
			err = postgresSetting.postgresConnectionInstance.db.Model(&model.History{}).Where("point_uuid = ? AND host_uuid = ? AND (timestamp AT TIME ZONE 'UTC' BETWEEN ? AND ?)", thisAssetResetDataPoint.UUID, thisAssetResetDataPoint.HostUUID, periodStart, periodEnd).Scan(&resetHistoryData).Error
			if err != nil {
				inst.cpsErrorMsg("CPSProcessing() resetHistoryData error: ", err)
			}
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() resetHistoryData: %+v", resetHistoryData))
			// dfRawDoor := dataframe.ReadCSV(strings.NewReader(csvRawDoor))
			dfResetHistoryData := dataframe.LoadStructs(resetHistoryData)
			fmt.Println("dfResetHistoryData")
			fmt.Println(dfResetHistoryData)

			// create a dataframe of the daily reset times
			dfDailyResets, err := inst.MakeDailyResetsDF(periodStart, periodEnd, dfSiteThresholds)
			if err != nil {
				inst.cpsErrorMsg("MakeDailyResetsDF() error: ", err)
				return
			}
			fmt.Println("dfDailyResets")
			fmt.Println(dfDailyResets)

			// join daily reset timestamps with the manual resets
			var dfAllResets dataframe.DataFrame
			if len(resetHistoryData) > 0 {
				if dfDailyResets.Nrow() > 0 {
					dfAllResets = dfResetHistoryData.Concat(*dfDailyResets)
					dfAllResets = dfAllResets.Arrange(dataframe.Sort(string(timestampColName)))
				} else {
					dfAllResets = dfResetHistoryData
				}
			} else if dfDailyResets.Nrow() > 0 {
				dfAllResets = *dfDailyResets
			}
			fmt.Println("dfAllResets")
			fmt.Println(dfAllResets)

			// extract the last processed data values and the door type info from the point tags and values
			pointLastProcessedData, pointDoorInfo, err := inst.GetLastProcessedDataAndDoorType(&dfJoinedLastProcessedValuesAndPoints, &doorSensorPoint)
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() pointLastProcessedData: %+v", pointLastProcessedData))
			inst.cpsDebugMsg(fmt.Sprintf("CPSProcessing() pointDoorInfo: %+v", pointDoorInfo))

			// now do the door usage calculations
			dfDoorResults, err := inst.CalculateDoorUses(dfRawDoorSensorHistories, dfJoinedLastProcessedValuesAndPoints, dfAllResets, dfSiteThresholds, pointLastProcessedData, pointDoorInfo)
			if err != nil {
				inst.cpsErrorMsg("CalculateDoorUses() error: ", err)
				return
			}
			fmt.Println("dfDoorResults")
			fmt.Println(dfDoorResults)

			// TODO: DELETE ME (just for debug)
			ResultFile, err := os.Create(fmt.Sprintf("/home/marc/Documents/Nube/CPS/Development/Data_Processing/7_DoorResultDF.csv"))
			if err != nil {
				inst.cpsErrorMsg(err)
			}
			defer ResultFile.Close()
			dfDoorResults.WriteCSV(ResultFile)

			// now calculate the 15 minute usage rollup
			dfDoorResults15Min, err := inst.Calculate15MinUsageRollup(periodStart, periodEnd, *dfDoorResults, dfJoinedLastProcessedValuesAndPoints)
			if err != nil {
				inst.cpsErrorMsg("CalculateDoorUses() error: ", err)
				return
			}
			fmt.Println("dfDoorResults15Min")
			fmt.Println(dfDoorResults15Min)

			// TODO: DELETE ME (just for debug)
			ResultFile2, err := os.Create(fmt.Sprintf("/home/marc/Documents/Nube/CPS/Development/Data_Processing/8_dfDoorResults15Min.csv"))
			if err != nil {
				inst.cpsErrorMsg(err)
			}
			defer ResultFile2.Close()
			dfDoorResults15Min.WriteCSV(ResultFile2)

			/*
				// TODO: need to get last totalUses at the time when the last 15MinRollup value was stored (prior to the processing time range)
				// get last totalUses at the time when the last 15MinRollup value was stored (prior to the processing time range)
				dfLastTotalUsesAt15Min := dataframe.ReadCSV(strings.NewReader(csvLastTotalUsesAt15Min))
				fmt.Println("dfLastTotalUsesAt15Min")
				fmt.Println(dfLastTotalUsesAt15Min)

				dfResets := dataframe.ReadCSV(strings.NewReader(csvRawResets))
				// fmt.Println("dfResets")
				// fmt.Println(dfResets)

				dfDailyResets, err := inst.MakeDailyResetsDF(periodStart, periodEnd, dfSiteThresholds)
				if err != nil {
					inst.cpsErrorMsg("MakeDailyResetsDF() error: ", err)
					return
				}
				fmt.Println("dfDailyResets")
				fmt.Println(dfDailyResets)

				// join daily reset timestamps with the manual resets
				dfAllResets := dfResets.Concat(*dfDailyResets)
				dfAllResets = dfAllResets.Arrange(dataframe.Sort(string(timestampColName)))
				fmt.Println("dfAllResets")
				fmt.Println(dfAllResets)

				lastToPending := "2023-06-19T07:42:00Z"
				lastToClean := "2023-06-17T07:15:00Z"

				// Normally OPEN Door Usage Count and Occupancy DF
				DoorResultDF, err := inst.CalculateDoorUses(facilityToilet, normallyOpen, dfRawDoor, dfLastProcessedDoor, dfAllResets, dfSiteThresholds, lastToPending)
				if err != nil {
					inst.cpsErrorMsg("CalculateDoorUses() error: ", err)
					return
				}
				fmt.Println("DoorResultDF")
				fmt.Println(DoorResultDF)

				RollupResultDF, err := inst.Calculate15MinUsageRollup(periodStart, periodEnd, *DoorResultDF, dfLastTotalUsesAt15Min)
				if err != nil {
					inst.cpsErrorMsg("Calculate15MinUsageRollup() error: ", err)
					return
				}
				fmt.Println("RollupResultDF")
				fmt.Println(RollupResultDF)

				OverdueResultDF, err := inst.CalculateOverdueCubicles(facilityToilet, periodStart, periodEnd, *DoorResultDF, dfLastProcessedDoor, dfSiteThresholds, lastToPending, lastToClean)
				if err != nil {
					inst.cpsErrorMsg("CalculateOverdueCubicles() error: ", err)
					return
				}
				fmt.Println("OverdueResultDF")
				fmt.Println(OverdueResultDF)
			*/
		}
	}

	/*
		ResultFile, err := os.Create("/home/marc/Documents/Nube/CPS/Development/Data_Processing/1_Results.csv")
		if err != nil {
			inst.cpsErrorMsg(err)
		}
		defer ResultFile.Close()
		OverdueResultDF.WriteCSV(ResultFile)

	*/
}

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.cpsDebugMsg("addNetwork(): ", body.Name)

	body.HistoryEnable = boolean.NewTrue()

	network, err = inst.db.CreateNetwork(body)
	if network == nil || err != nil {
		inst.cpsErrorMsg("addNetwork(): failed to create cps network: ", body.Name)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("failed to create cps network")
	}
	return network, err
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.cpsDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.cpsDebugMsg("addDevice(): ", body.Name)

	if body.Name == "Cleaning Resets" {
		body.HistoryEnable = boolean.NewTrue()
	} else {
		body.HistoryEnable = boolean.NewFalse()
	}

	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.cpsDebugMsg("addDevice(): failed to create cps device: ", body.Name)
		return nil, errors.New("failed to create cps device")
	}

	if device.Name == "Cleaning Resets" || device.Name == "Availability" {
		return device, nil
	}

	// Automatically create door sensor processed data points
	createThesePoints := make([]model.Point, 0)
	occupancyPoint := model.Point{
		Name: string(cubicleOccupancyColName),
	}
	createThesePoints = append(createThesePoints, occupancyPoint)

	totalUsesPoint := model.Point{
		Name: string(totalUsesColName),
	}
	createThesePoints = append(createThesePoints, totalUsesPoint)

	currentUsesPoint := model.Point{
		Name: string(currentUsesColName),
	}
	createThesePoints = append(createThesePoints, currentUsesPoint)

	fifteenMinUsesPoint := model.Point{
		Name: string(fifteenMinRollupUsesColName),
	}
	createThesePoints = append(createThesePoints, fifteenMinUsesPoint)

	pendingStatusPoint := model.Point{
		Name: string(pendingStatusColName),
	}
	createThesePoints = append(createThesePoints, pendingStatusPoint)

	overdueStatusPoint := model.Point{
		Name: string(overdueStatusColName),
	}
	createThesePoints = append(createThesePoints, overdueStatusPoint)

	toPendingPoint := model.Point{
		Name: string(toPendingColName),
	}
	createThesePoints = append(createThesePoints, toPendingPoint)

	toCleanPoint := model.Point{
		Name: string(toCleanColName),
	}
	createThesePoints = append(createThesePoints, toCleanPoint)

	toOverduePoint := model.Point{
		Name: string(toOverdueColName),
	}
	createThesePoints = append(createThesePoints, toOverduePoint)

	cleaningTimePoint := model.Point{
		Name: string(cleaningTimeColName),
	}
	createThesePoints = append(createThesePoints, cleaningTimePoint)

	devNameSplit := strings.Split(device.Name, "-")
	if len(devNameSplit) < 3 {
		inst.cpsErrorMsg("addDevice(): device name should be of the form Level-Gender-Location")
		return device, nil
	}
	level := devNameSplit[0]
	gender := strings.ToLower(devNameSplit[1])
	location := devNameSplit[2]

	for _, point := range createThesePoints {
		point.DeviceUUID = device.UUID
		point.MetaTags = make([]*model.PointMetaTag, 0)

		metaTag1 := model.PointMetaTag{Key: "assetFunc", Value: "managedCubicle"}
		point.MetaTags = append(point.MetaTags, &metaTag1)
		metaTag2 := model.PointMetaTag{Key: "doorType", Value: "toilet"}
		point.MetaTags = append(point.MetaTags, &metaTag2)
		metaTag3 := model.PointMetaTag{Key: "enableCleaningTracking", Value: "true"}
		point.MetaTags = append(point.MetaTags, &metaTag3)
		metaTag4 := model.PointMetaTag{Key: "enableUseCounting", Value: "true"}
		point.MetaTags = append(point.MetaTags, &metaTag4)
		metaTag5 := model.PointMetaTag{Key: "floorRef", Value: level}
		point.MetaTags = append(point.MetaTags, &metaTag5)
		metaTag6 := model.PointMetaTag{Key: "genderRef", Value: gender}
		point.MetaTags = append(point.MetaTags, &metaTag6)
		metaTag7 := model.PointMetaTag{Key: "isEOT", Value: "false"}
		point.MetaTags = append(point.MetaTags, &metaTag7)
		metaTag8 := model.PointMetaTag{Key: "locationRef", Value: location}
		point.MetaTags = append(point.MetaTags, &metaTag8)
		metaTag9 := model.PointMetaTag{Key: "measurementRef", Value: "door_position"}
		point.MetaTags = append(point.MetaTags, &metaTag9)
		metaTag10 := model.PointMetaTag{Key: "normalPosition", Value: "NO"}
		point.MetaTags = append(point.MetaTags, &metaTag10)
		metaTag11 := model.PointMetaTag{Key: "siteRef", Value: "cps_b49e0c73919c47ef"}
		point.MetaTags = append(point.MetaTags, &metaTag11)
		metaTag12 := model.PointMetaTag{Key: "assetRef", Value: device.Description}
		point.MetaTags = append(point.MetaTags, &metaTag12)
		metaTag13 := model.PointMetaTag{Key: "resetID", Value: "rst_1"}
		point.MetaTags = append(point.MetaTags, &metaTag13)
		metaTag14 := model.PointMetaTag{Key: "availabilityID", Value: "avl_1"}
		point.MetaTags = append(point.MetaTags, &metaTag14)

		inst.addPoint(&point)
	}

	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil || body.DeviceUUID == "" {
		inst.cpsDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.cpsDebugMsg("addPoint(): ", body.Name)

	device, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if device == nil || err != nil {
		inst.cpsDebugMsg("addPoint(): failed to find device", body.DeviceUUID)
		return nil, err
	}

	if device.Name == "Cleaning Resets" {
		body.HistoryEnable = boolean.NewTrue()
		body.HistoryType = model.HistoryTypeCov
		body.HistoryCOVThreshold = float.New(0.1)
	} else {
		body.HistoryEnable = boolean.NewFalse()
	}

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.cpsDebugMsg("addPoint(): failed to create cps point: ", body.Name)
		return nil, err
	}
	// inst.cpsDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))
	return point, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.cpsDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updateNetwork():  nil network object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "network disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}

	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.cpsDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
		body.CommonFault.Message = "device disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(body.UUID, body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.cpsDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("updatePoint(): nil point object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}
	body.CommonFault.InFault = false
	body.CommonFault.MessageLevel = model.MessageLevel.Info
	body.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	body.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	body.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.cpsErrorMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.cpsDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deleteNetwork(): nil network object")
		return
	}

	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.cpsDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deleteDevice(): nil device object")
		return
	}
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.cpsDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.cpsDebugMsg("deletePoint(): nil point object")
		return
	}
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}
