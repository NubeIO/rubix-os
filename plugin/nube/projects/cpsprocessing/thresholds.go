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
	ID                                     uint           `json:"id" gorm:"primarykey" dataframe:"-"`
	CreatedAt                              time.Time      `json:"created_at" dataframe:"-"`
	UpdatedAt                              time.Time      `json:"updated_at" dataframe:"-"`
	DeletedAt                              gorm.DeletedAt `json:"deleted_at" gorm:"index" dataframe:"-"`
	SiteRef                                string         `json:"site_ref" gorm:"type:varchar(255)" dataframe:"site_ref"`
	Timezone                               string         `json:"timezone" dataframe:"timezone"`
	OccupancyStartTime                     string         `json:"occupancyStartTime" dataframe:"occupancyStartTime"`
	OccupancyStopTime                      string         `json:"occupancyStopTime" dataframe:"occupancyStopTime"`
	AllAreaResetTime                       string         `json:"allAreaResetTime" dataframe:"allAreaResetTime"`
	FacilityCleaningOverdueAlertDelay      int            `json:"facilityCleaningOverdueAlertDelay" dataframe:"facilityCleaningOverdueAlertDelay"`
	EOTCleaningOverdueAlertDelay           int            `json:"eotCleaningOverdueAlertDelay" dataframe:"eotCleaningOverdueAlertDelay"`
	FacilityToiletUseThreshold             int            `json:"facilityToiletUseThreshold" dataframe:"facilityToiletUseThreshold"`
	FacilityEntranceUseThreshold           int            `json:"facilityEntranceUseThreshold" dataframe:"facilityEntranceUseThreshold"`
	FacilityDDAUseThreshold                int            `json:"facilityDDAUseThreshold" dataframe:"facilityDDAUseThreshold"`
	EOTToiletUseThreshold                  int            `json:"eotToiletUseThreshold" dataframe:"eotToiletUseThreshold"`
	EOTShowerUseThreshold                  int            `json:"eotShowerUseThreshold" dataframe:"eotShowerUseThreshold"`
	EOTEntranceUseThreshold                int            `json:"eotEntranceUseThreshold" dataframe:"eotEntranceUseThreshold"`
	EOTDDAUseThreshold                     int            `json:"eotDDAUseThreshold" dataframe:"eotDDAUseThreshold"`
	LowBatteryAlertDelay                   int            `json:"lowBatteryAlertDelay" dataframe:"lowBatteryAlertDelay"`
	TemperatureAlertDelay                  int            `json:"temperatureAlertDelay" dataframe:"temperatureAlertDelay"`
	HumidityAlertDelay                     int            `json:"humidityAlertDelay" dataframe:"humidityAlertDelay"`
	CO2AlertDelay                          int            `json:"co2AlertDelay" dataframe:"co2AlertDelay"`
	VOCAlertDelay                          int            `json:"vocAlertDelay" dataframe:"vocAlertDelay"`
	SensorOfflineAlertDelay                int            `json:"sensorOfflineAlertDelay" dataframe:"sensorOfflineAlertDelay"`
	GatewayOfflineAlertDelay               int            `json:"gatewayOfflineAlertDelay" dataframe:"gatewayOfflineAlertDelay"`
	LowBatteryVoltageThreshold             float64        `json:"lowBatteryVoltageThreshold" dataframe:"lowBatteryVoltageThreshold"`
	LowBatteryPercentThreshold             int            `json:"lowBatteryPercentThreshold" dataframe:"lowBatteryPercentThreshold"`
	HighTemperatureAlertThreshold          float64        `json:"highTemperatureAlertThreshold" dataframe:"highTemperatureAlertThreshold"`
	LowTemperatureAlertThreshold           float64        `json:"lowTemperatureAlertThreshold" dataframe:"lowTemperatureAlertThreshold"`
	HighHumidityAlertThreshold             int            `json:"highHumidityAlertThreshold" dataframe:"highHumidityAlertThreshold"`
	HighCo2AlertThreshold                  int            `json:"highCo2AlertThreshold" dataframe:"highCo2AlertThreshold"`
	HighVocAlertThreshold                  int            `json:"highVocAlertThreshold" dataframe:"highVocAlertThreshold"`
	EOTLowShowerAvailabilityThreshold      int            `json:"eotLowShowerAvailabilityThreshold" dataframe:"eotLowShowerAvailabilityThreshold"`
	EOTLowToiletAvailabilityThreshold      int            `json:"eotLowToiletAvailabilityThreshold" dataframe:"eotLowToiletAvailabilityThreshold"`
	FacilityLowToiletAvailabilityThreshold int            `json:"facilityLowToiletAvailabilityThreshold" dataframe:"facilityLowToiletAvailabilityThreshold"`
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
  "eotLowShowerAvailabilityThreshold": 80,
  "eotlowToiletAvailabilityThreshold": 80,
  "facilityLowToiletAvailabilityThreshold": 80,
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

	err = postgresSetting.postgresConnectionInstance.db.First(&site, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("CreateThreshold() db.First(&site, \"site_ref = ?\", siteRef) error: ", err)
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

	err = postgresSetting.postgresConnectionInstance.db.Last(&threshold, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("GetLastThresholdBySiteRef() db.Last(&thresholds, siteRef) error: ", err)
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
