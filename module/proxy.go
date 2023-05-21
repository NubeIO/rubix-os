package module

import (
	"github.com/NubeIO/flow-framework/interfaces"
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
	urlPrefix := fullPathParts[3]
	module := modules[urlPrefix]
	var res []byte
	var err error
	var status int
	if method == "GET" {
		status = http.StatusOK
		res, err = (*module).Get(proxyPath)
	} else if method == "POST" {
		status = http.StatusCreated
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Post(proxyPath, requestBody)
	} else if method == "PUT" {
		status = http.StatusOK
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Put(proxyPath, requestBody)
	} else if method == "PATCH" {
		status = http.StatusOK
		requestBody, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Patch(proxyPath, requestBody)
	} else if method == "DELETE" {
		status = http.StatusNoContent
		res, err = (*module).Delete(proxyPath)
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
