package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Site struct {
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	SiteRef               string         `json:"site_ref" gorm:"type:varchar(255);unique;primaryKey"`
	Name                  string         `json:"name" gorm:"type:varchar(255);unique"`
	Address               string         `json:"address"`
	State                 string         `json:"state"`
	Region                uint8          `json:"region"`
	Operation             string         `json:"operation"`
	AssetOwner            string         `json:"asset_owner"`
	ManagingAgent         string         `json:"managing_agent"`
	BuildingType          string         `json:"building_type"`
	BuildingGrade         uint8          `json:"building_grade"`
	ContractStartDate     string         `json:"contract_start_date"`
	ContractEndDate       string         `json:"contract_end_date"`
	ContractTimeRemaining *float64       `json:"contract_time_remaining,omitempty"`
	ContractFlag          *uint8         `json:"contract_flag,omitempty"`
	IsCurrentFlag         *uint8         `json:"contract_flag,omitempty"`
}

// CreateSite creates a new site entry
func (inst *Instance) CreateSite(c *gin.Context) {
	var site Site
	var err error
	err = c.ShouldBindJSON(&site)
	if err != nil {
		inst.cpsErrorMsg("CreateSite() ShouldBindJSON() error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("CreateSite() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// create a new cps uuid for the new site
	uuid, _ := uuid.MakeUUID()
	site.SiteRef = fmt.Sprintf("cps_%s", uuid)

	if err := postgresSetting.postgresConnectionInstance.db.Create(&site).Error; err != nil {
		inst.cpsErrorMsg("CreateSite() db.Create(&site) error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, site)
}

// GetAllSites retrieves all sites
func (inst *Instance) GetAllSites(c *gin.Context) {

	var sites []Site
	var err error

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("GetAllSites() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.Find(&sites).Error
	if err != nil {
		inst.cpsErrorMsg("GetAllSites() db.Find(&sites) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, sites)
}

// GetSite retrieves a site entry by site_ref
func (inst *Instance) GetSite(c *gin.Context) {
	siteRef := c.Param("site_ref")

	var site Site
	var err error

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("GetSite() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.First(&site, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("GetSite() db.First(&site, siteRef) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, site)
}

// GetSiteByName retrieves a site entry by its name
func (inst *Instance) GetSiteByName(c *gin.Context) {

	var site Site
	var err error

	err = c.ShouldBindJSON(&site)
	if err != nil {
		inst.cpsErrorMsg("GetSiteByName() ShouldBindJSON() error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	siteName := site.Name
	inst.cpsDebugMsg("GetSiteByName() siteName: ", siteName)
	if siteName == "" {
		inst.cpsErrorMsg("GetSiteByName() error: site 'name' is required in body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "site 'name' is required in body"})
		return
	}

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("GetSiteByName() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.Where("name = ?", siteName).First(&site).Error
	if err != nil {
		inst.cpsErrorMsg("GetSiteByName() db.Where(name = ?, siteName).First(&site) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, site)
}

// GetSiteByAddress retrieves a site entry by its address
func (inst *Instance) GetSiteByAddress(c *gin.Context) {

	var site Site
	var sitesResult []Site
	var err error

	err = c.ShouldBindJSON(&site)
	if err != nil {
		inst.cpsErrorMsg("GetSiteByAddress() ShouldBindJSON() error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	siteAddress := site.Address
	inst.cpsDebugMsg("GetSiteByAddress() siteAddress: ", siteAddress)
	if siteAddress == "" {
		inst.cpsErrorMsg("GetSiteByAddress() error: site 'address' is required in body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "site 'address' is required in body"})
		return
	}

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("GetSiteByAddress() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.Where("address = ?", siteAddress).Find(&sitesResult).Error
	if err != nil {
		inst.cpsErrorMsg("GetSiteByAddress() db.Where(address = ?, siteAddress).First(&site) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, sitesResult)
}

// UpdateSite updates a site entry
func (inst *Instance) UpdateSite(c *gin.Context) {
	siteRef := c.Param("site_ref")

	var site Site
	var existingSite Site
	var err error
	err = c.ShouldBindJSON(&site)
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() ShouldBindJSON() error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if site.SiteRef != "site_ref" {
		inst.cpsErrorMsg("UpdateSite() attempt to modify site_ref is not permitted")
		c.JSON(http.StatusNotFound, gin.H{"error": "modifying site_ref is not permitted"})
		return
	}

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.First(&existingSite, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() db.First(&existingSite, \"site_ref = ?\", siteRef) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	err = postgresSetting.postgresConnectionInstance.db.Model(&existingSite).Updates(site).Error
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() db.Model(&existingSite).Updates(site) error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the updated site from the database
	err = postgresSetting.postgresConnectionInstance.db.First(&site, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() db.First(&site, \"site_ref = ?\", siteRef) error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, site)
}

// DeleteSite deletes a site entry by site_ref
func (inst *Instance) DeleteSite(c *gin.Context) {
	siteRef := c.Param("site_ref")

	var site Site
	var err error

	_, err = inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("UpdateSite() initializePostgresDBConnection() error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = postgresSetting.postgresConnectionInstance.db.First(&site, "site_ref = ?", siteRef).Error
	if err != nil {
		inst.cpsErrorMsg("DeleteSite() db.First(&site, \"site_ref = ?\", siteRef) error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	err = postgresSetting.postgresConnectionInstance.db.Delete(&site).Error
	if err != nil {
		inst.cpsErrorMsg("DeleteSite() db.Delete(&site) error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
