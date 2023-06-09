package module

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/module/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func ProxyModule(c *gin.Context) {
	method := c.Request.Method
	proxyPath := c.Param("proxyPath")
	fullPath := c.FullPath()
	fullPathParts := strings.Split(fullPath, "/")
	if len(fullPathParts) < 4 {
		log.Error("Error on module framework")
		c.JSON(http.StatusInternalServerError, interfaces.Message{Message: "error on module framework!"})
	}
	moduleName := fullPathParts[3]
	module := modules[moduleName]
	if module == nil {
		msg := fmt.Sprintf("we don't have module with module_name %s", moduleName)
		c.JSON(http.StatusBadRequest, interfaces.Message{Message: msg})
		return
	}
	var res []byte
	var err error
	var status int
	fmt.Println("method", method)
	if method == "GET" {
		status = http.StatusOK
		fmt.Println("proxyPath", proxyPath)
		res, err = module.Get(proxyPath)
	} else if method == "POST" {
		status = http.StatusCreated
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = module.Post(proxyPath, requestBody)
	} else if method == "PUT" {
		status = http.StatusOK
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		parentURL, uuid := utils.ParseUUID(proxyPath)
		res, err = module.Put(parentURL, uuid, requestBody)
	} else if method == "PATCH" {
		status = http.StatusOK
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		parentURL, uuid := utils.ParseUUID(proxyPath)
		res, err = module.Patch(parentURL, uuid, requestBody)
	} else if method == "DELETE" {
		status = http.StatusNoContent
		parentURL, uuid := utils.ParseUUID(proxyPath)
		res, err = module.Delete(parentURL, uuid)
	}

	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = not found" {
			c.JSON(http.StatusNotFound, interfaces.Message{Message: err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
	}
	c.Data(status, "application/json", res)
}
