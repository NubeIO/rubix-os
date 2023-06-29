package main

import "time"

const sensorIDMetaTagKey = "sensorRef"
const sensorMakeMetaTagKey = "sensorMake"
const sensorModelMetaTagKey = "sensorModel"

const emptyModuleStorageResyncPeriod = time.Duration(48 * time.Hour)
