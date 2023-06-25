package main

const csvSiteThresholds = `timezone,occupancyStartTime,occupancyStopTime,allAreaResetTime,facilityCleaningOverdueAlertDelay,eotCleaningOverdueAlertDelay,facilityToiletUseThreshold,facilityEntranceUseThreshold,facilityDDAUseThreshold,eotToiletUseThreshold,eotShowerUseThreshold,eotEntranceUseThreshold,eotDDAUseThreshold
Australia/Sydney,7:00,18:00,22:00,30,30,3,100,50,10,1,20,20`

const csvLastProNODoor = `door,cubicleOccupancy,totalUses,currentUses,pendingStatus,overdueStatus
1,0,20,7,1,0`

const csvRawDoor = `timestamp,door
2023-06-19T08:00:00Z,1
2023-06-19T08:05:00Z,0
2023-06-19T08:10:00Z,1
2023-06-19T08:11:00Z,0
2023-06-19T08:13:00Z,1
2023-06-19T08:20:00Z,0
2023-06-19T08:25:00Z,0
2023-06-19T08:31:00Z,1
2023-06-19T08:35:00Z,1
2023-06-19T08:40:00Z,1
2023-06-19T08:41:00Z,0
2023-06-19T08:42:00Z,1
2023-06-19T08:43:00Z,0
2023-06-19T08:45:00Z,1
2023-06-19T08:47:00Z,1
2023-06-19T08:49:00Z,0
2023-06-19T08:51:00Z,1
2023-06-19T08:53:00Z,0
2023-06-19T08:55:00Z,1`

const csvRawResets = `timestamp,areaReset
2023-06-19T08:29:00Z,1`

const csvLastTotalUsesAt15Min = `timestamp,totalUses,15MinUses
2023-06-19T07:45:00Z,18,3`
