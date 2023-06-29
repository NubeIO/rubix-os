package main

import "time"

const sensorIDMetaTagKey = "sensorRef"
const sensorMakeMetaTagKey = "sensorMake"
const sensorModelMetaTagKey = "sensorModel"

const emptyModuleStorageResyncPeriod = time.Duration(72 * time.Hour)
const defaultSensorHistorySyncFrequency = "5m"
const defaultGatewayPayloadSyncFrequency = "15m"
