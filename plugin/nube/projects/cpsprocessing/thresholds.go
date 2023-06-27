package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"time"
)

type Threshold struct {
	// gorm.Model
	ID                                uint           `json:"id" gorm:"primarykey"`
	CreatedAt                         time.Time      `json:"created_at"`
	UpdatedAt                         time.Time      `json:"updated_at"`
	DeletedAt                         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	SiteRef                           string         `json:"site_ref" gorm:"type:varchar(255)"`
	Timezone                          string         `json:"timezone"`
	OccupancyStartTime                string         `json:"occupancyStartTime"`
	OccupancyStopTime                 string         `json:"occupancyStopTime"`
	AllAreaResetTime                  string         `json:"allAreaResetTime"`
	FacilityCleaningOverdueAlertDelay int            `json:"facilityCleaningOverdueAlertDelay"`
	EOTCleaningOverdueAlertDelay      int            `json:"eotCleaningOverdueAlertDelay"`
	FacilityToiletUseThreshold        int            `json:"facilityToiletUseThreshold"`
	FacilityEntranceUseThreshold      int            `json:"facilityEntranceUseThreshold"`
	FacilityDDAUseThreshold           int            `json:"facilityDDAUseThreshold"`
	EOTToiletUseThreshold             int            `json:"eotToiletUseThreshold"`
	EOTShowerUseThreshold             int            `json:"eotShowerUseThreshold"`
	EOTEntranceUseThreshold           int            `json:"eotEntranceUseThreshold"`
	EOTDDAUseThreshold                int            `json:"eotDDAUseThreshold"`
	LowBatteryAlertDelay              int            `json:"lowBatteryAlertDelay"`
	TemperatureAlertDelay             int            `json:"temperatureAlertDelay"`
	HumidityAlertDelay                int            `json:"humidityAlertDelay"`
	CO2AlertDelay                     int            `json:"co2AlertDelay"`
	VOCAlertDelay                     int            `json:"vocAlertDelay"`
	SensorOfflineAlertDelay           int            `json:"sensorOfflineAlertDelay"`
	GatewayOfflineAlertDelay          int            `json:"gatewayOfflineAlertDelay"`
	LowBatteryVoltageThreshold        float64        `json:"lowBatteryVoltageThreshold"`
	LowBatteryPercentThreshold        int            `json:"lowBatteryPercentThreshold"`
	HighTemperatureAlertThreshold     float64        `json:"highTemperatureAlertThreshold"`
	LowTemperatureAlertThreshold      float64        `json:"lowTemperatureAlertThreshold"`
	HighHumidityAlertThreshold        int            `json:"highHumidityAlertThreshold"`
	HighCo2AlertThreshold             int            `json:"highCo2AlertThreshold"`
	HighVocAlertThreshold             int            `json:"highVocAlertThreshold"`
	LowShowerAvailabilityThreshold    int            `json:"lowShowerAvailabilityThreshold"`
	LowToiletAvailabilityThreshold    int            `json:"lowToiletAvailabilityThreshold"`
}

/*
{
  "site_ref": "cps_b49e0c73919c47ef",
  "timezone": "Australia/Sydney",
  "occupancyStartTime": "7:00",
  "occupancyStopTime": "18:00",
  "allAreaResetTime": "22:00",
  "facilityCleaningOverdueAlertDelay": 30,
  "eotCleaningOverdueAlertDelay": 30,
  "facilityToiletUseThreshold": 3,
  "facilityEntranceUseThreshold": 100,
  "facilityDDAUseThreshold": 50,
  "eotToiletUseThreshold": 10,
  "eotShowerUseThreshold": 1,
  "eotEntranceUseThreshold": 20,
  "eotDDAUseThreshold": 20,
  "lowBatteryAlertDelay": 60,
  "temperatureAlertDelay": 60,
  "humidityAlertDelay": 60,
  "co2AlertDelay": 60,
  "vocAlertDelay": 60,
  "sensorOfflineAlertDelay": 120,
  "gatewayOfflineAlertDelay": 30,
  "lowBatteryVoltageThreshold": 2.5,
  "lowBatteryPercentThreshold": 10,
  "highTemperatureAlertThreshold": 30.0,
  "lowTemperatureAlertThreshold": 10.0,
  "highHumidityAlertThreshold": 85,
  "highCo2AlertThreshold": 1500,
  "highVocAlertThreshold": 2500,
  "lowShowerAvailabilityThreshold": 80,
  "lowToiletAvailabilityThreshold": 80
}
*/

// CreateThreshold creates a new threshold entry
func (inst *Instance) CreateThreshold(c *gin.Context) {
	var threshold Threshold
	var err error
	err = c.ShouldBindJSON(&threshold)
	if err != nil {
		inst.cpsErrorMsg("CreateThreshold() ShouldBindJSON() error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for zero values in the Threshold struct
	if thresholdHasZeroValues(threshold) {
		errMsg := "threshold json contains zero values"
		inst.cpsErrorMsg("CreateThreshold() validation error: ", errors.New(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("CreateThreshold() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	siteRef := threshold.SiteRef
	var site Site

	err = postgresSetting.postgresConnectionInstance.db.Where("site_ref = ?", siteRef).First(&site).Error
	if err != nil {
		inst.cpsErrorMsg("CreateThreshold() db.Where(site_ref = ?, siteRef).First(&site) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	if err := postgresSetting.postgresConnectionInstance.db.Create(&threshold).Error; err != nil {
		inst.cpsErrorMsg("CreateThreshold() db.Create(&threshold) error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, threshold)
}

// GetLastThresholdBySiteRef retrieves a site threshold entry by site_ref
func (inst *Instance) GetLastThresholdBySiteRef(c *gin.Context) {
	siteRef := c.Param("site_ref")

	var threshold Threshold
	var err error

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("GetLastThresholdBySiteRef() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.Last(&threshold, siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("GetLastThresholdBySiteRef() db.First(&site, siteRef) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, threshold)
}

// thresholdHasZeroValues checks if the given thresholds struct contains any zero values
func thresholdHasZeroValues(s interface{}) bool {
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		if fieldName == "Model" || fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt" {
			continue
		}
		field := v.Field(i)
		zeroValue := reflect.Zero(field.Type()).Interface()
		if reflect.DeepEqual(field.Interface(), zeroValue) {
			return true
		}
	}
	return false
}
