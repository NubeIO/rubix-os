package main

const open = 0
const closed = 1

const vacant = 0
const occupied = 1

type PointTags string

const (
	pointFunctionTag                                       PointTags = "pointFunction"
	rawDoorSensorPointFunctionTagValue                     PointTags = "sensor"
	processedDataPointFunctionTagValue                     PointTags = "processedData"
	doorResetPointFunctionTagValue                         PointTags = "doorReset"
	measurementRefTag                                      PointTags = "measurementRef"
	availabilityMeasurementRefTag                          PointTags = "availability"
	availabilityAlertMeasurementRefTag                     PointTags = "availabilityAlert"
	doorSensorMeasurementRefTagValue                       PointTags = "door_position"
	assetRefTag                                            PointTags = "assetRef"
	assetFuncTag                                           PointTags = "assetFunc"
	managedCubicleDoorSensorAssetFunctionTagValue          PointTags = "managedCubicle"
	managedFacilityEntranceDoorSensorAssetFunctionTagValue PointTags = "managedFacilityEntrance"
	usageCountDoorSensorAssetFunctionTagValue              PointTags = "usageCount"
	doorTypeTag                                            PointTags = "doorType"
	toiletDoorTypeTagValue                                 PointTags = "toilet"
	showerDoorTypeTagValue                                 PointTags = "shower"
	ddaDoorTypeTagValue                                    PointTags = "dda"
	entranceDoorTypeTagValue                               PointTags = "entrance"
	doorDoorTypeTagValue                                   PointTags = "door"
	normalPositionTag                                      PointTags = "normalPosition"
	normallyOpenNormalPositionTagValue                     PointTags = "NO"
	normallyClosedNormalPositionTagValue                   PointTags = "NC"
	enableCleaningTrackingTag                              PointTags = "enableCleaningTracking"
	enabledEnableCleaningTrackingTagValue                  PointTags = "true"
	disabledEnableCleaningTrackingTagValue                 PointTags = "false"
	enableUseCountingTag                                   PointTags = "enableUseCounting"
	enabledEnableUseCountingTagValue                       PointTags = "true"
	disabledEnableUseCountingTagValue                      PointTags = "false"
	isEOTTag                                               PointTags = "isEOT"
	EOTisEOTTagValue                                       PointTags = "true"
	notEOTisEOTTagValue                                    PointTags = "false"
	availabilityIDTag                                      PointTags = "availabilityID"
	resetIDTag                                             PointTags = "resetID"
)

type DoorNormalPosition int

const (
	normallyOpen   DoorNormalPosition = 0
	normallyClosed DoorNormalPosition = 1
)

type DoorType int

const (
	facilityEntrance DoorType = iota
	facilityToilet
	facilityDDA
	facilityDoor
	eotEntrance
	eotToilet
	eotShower
	eotDDA
	eotDoor
)

type RawDataColumnName string

const (
	pointUUIDColName    RawDataColumnName = "point_uuid"
	hostUUIDColName     RawDataColumnName = "host_uuid"
	timestampColName    RawDataColumnName = "timestamp"
	doorPositionColName RawDataColumnName = "door"
	areaResetColName    RawDataColumnName = "areaReset"
	temperatureColName  RawDataColumnName = "temp"
	humidityColName     RawDataColumnName = "humidity"
	lightColName        RawDataColumnName = "light"
	co2ColName          RawDataColumnName = "co2"
	vocColName          RawDataColumnName = "voc"
	motionColName       RawDataColumnName = "motion"
	pccColName          RawDataColumnName = "pcc"
	deskColName         RawDataColumnName = "desk"
	voltageColName      RawDataColumnName = "voltage"
	batteryColName      RawDataColumnName = "battery"
	rssiColName         RawDataColumnName = "rssi"
)

type ProcessedDataColumnName string

const (
	cubicleOccupancyColName             ProcessedDataColumnName = "cubicleOccupancy"
	totalUsesColName                    ProcessedDataColumnName = "totalUses"
	currentUsesColName                  ProcessedDataColumnName = "currentUses"
	fifteenMinRollupUsesColName         ProcessedDataColumnName = "15minUses"
	pendingStatusColName                ProcessedDataColumnName = "pendingStatus"
	overdueStatusColName                ProcessedDataColumnName = "overdueStatus"
	toPendingColName                    ProcessedDataColumnName = "toPending"
	toCleanColName                      ProcessedDataColumnName = "toClean"
	toOverdueColName                    ProcessedDataColumnName = "toOverdue"
	cleaningTimeColName                 ProcessedDataColumnName = "cleaningTime"
	cleaningResetColName                ProcessedDataColumnName = "cleaningReset"
	lowBatteryColName                   ProcessedDataColumnName = "lowBattery"
	highTempColName                     ProcessedDataColumnName = "highTemp"
	lowTempColName                      ProcessedDataColumnName = "lowTemp"
	highCO2ColName                      ProcessedDataColumnName = "highCO2"
	highVOCColName                      ProcessedDataColumnName = "highVOC"
	sensorFlatlineColName               ProcessedDataColumnName = "sensorFlatline"
	gatewayFlatlineColName              ProcessedDataColumnName = "gatewayFlatline"
	lowToiletAvailabilityColName        ProcessedDataColumnName = "lowToiletAvailability"
	lowShowerAvailabilityColName        ProcessedDataColumnName = "lowShowerAvailability"
	lowToiletAvailabilityOverdueColName ProcessedDataColumnName = "lowToiletAvailabilityOverdue"
	lowShowerAvailabilityOverdueColName ProcessedDataColumnName = "lowShowerAvailabilityOverdue"
)

// NOTE: the thresholds struct would also need to be updated if these column names are changed.

type ThresholdColumnName string

const (
	timeZoneColName                               ThresholdColumnName = "timezone"
	occupancyStartTimeColName                     ThresholdColumnName = "occupancyStartTime"
	occupancyStopTimeColName                      ThresholdColumnName = "occupancyStopTime"
	allAreaResetTimeColName                       ThresholdColumnName = "allAreaResetTime"
	facilityCleaningOverdueAlertDelayColName      ThresholdColumnName = "facilityCleaningOverdueAlertDelay"
	eotCleaningOverdueAlertDelayColName           ThresholdColumnName = "eotCleaningOverdueAlertDelay"
	lowBatteryAlertDelayColName                   ThresholdColumnName = "lowBatteryAlertDelay"
	temperatureAlertDelayColName                  ThresholdColumnName = "temperatureAlertDelay"
	humidityAlertDelayColName                     ThresholdColumnName = "humidityAlertDelay"
	co2AlertDelayColName                          ThresholdColumnName = "co2AlertDelay"
	vocAlertDelayColName                          ThresholdColumnName = "vocAlertDelay"
	sensorOfflineAlertDelayColName                ThresholdColumnName = "sensorOfflineAlertDelay"
	gatewayOfflineAlertDelayColName               ThresholdColumnName = "gatewayOfflineAlertDelay"
	facilityToiletUseThresholdColName             ThresholdColumnName = "facilityToiletUseThreshold"
	facilityEntranceUseThresholdColName           ThresholdColumnName = "facilityEntranceUseThreshold"
	facilityDDAUseThresholdColName                ThresholdColumnName = "facilityDDAUseThreshold"
	eotToiletUseThresholdColName                  ThresholdColumnName = "eotToiletUseThreshold"
	eotShowerUseThresholdColName                  ThresholdColumnName = "eotShowerUseThreshold"
	eotEntranceUseThresholdColName                ThresholdColumnName = "eotEntranceUseThreshold"
	eotDDAUseThresholdColName                     ThresholdColumnName = "eotDDAUseThreshold"
	lowBatteryVoltageThresholdColName             ThresholdColumnName = "lowBatteryVoltageThreshold"
	lowBatteryPercentThresholdColName             ThresholdColumnName = "lowBatteryPercentThreshold"
	highTemperatureAlertThresholdColName          ThresholdColumnName = "highTemperatureAlertThreshold"
	lowTemperatureAlertThresholdColName           ThresholdColumnName = "lowTemperatureAlertThreshold"
	highHumidityAlertThresholdColName             ThresholdColumnName = "highHumidityAlertThreshold"
	highCo2AlertThresholdColName                  ThresholdColumnName = "highCo2AlertThreshold"
	highVocAlertThresholdColName                  ThresholdColumnName = "highVocAlertThreshold"
	eotLowShowerAvailabilityThresholdColName      ThresholdColumnName = "eotLowShowerAvailabilityThreshold"
	eotLowToiletAvailabilityThresholdColName      ThresholdColumnName = "eotLowToiletAvailabilityThreshold"
	facilityLowToiletAvailabilityThresholdColName ThresholdColumnName = "facilityLowToiletAvailabilityThreshold"
)
