package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"os"
	"strings"
	"time"
)

func (inst *Instance) CPSProcessing() {
	inst.cpsDebugMsg("CPSProcessing()")

	periodStart, _ := time.Parse(time.RFC3339, "2023-06-19T07:55:00Z")
	periodEnd, _ := time.Parse(time.RFC3339, "2023-06-19T10:00:00Z")

	dfSiteThresholds := dataframe.ReadCSV(strings.NewReader(csvSiteThresholds))
	// fmt.Println("dfSiteThresholds")
	// fmt.Println(dfSiteThresholds)

	// TODO: pull data for each sensor for the given time range.  Do the below functions for each sensor

	// get raw sensor data for time range
	dfRawDoor := dataframe.ReadCSV(strings.NewReader(csvRawDoor))
	// fmt.Println("dfRawDoor")
	// fmt.Println(dfRawDoor)

	// get last stored processed data values
	dfLastProcessedDoor := dataframe.ReadCSV(strings.NewReader(csvLastProNODoor))
	fmt.Println("dfLastProcessedDoor")
	fmt.Println(dfLastProcessedDoor)

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

	ResultFile, err := os.Create("/home/marc/Documents/Nube/CPS/Development/Data_Processing/1_Results.csv")
	if err != nil {
		inst.cpsErrorMsg(err)
	}
	defer ResultFile.Close()
	OverdueResultDF.WriteCSV(ResultFile)
}
