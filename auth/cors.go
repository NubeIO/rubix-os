package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
)

// CorsConfig generates a config to use in gin cors middleware based on server configuration.
func CorsConfig() cors.Config {
	corsConf := cors.Config{
		MaxAge:                 12 * time.Hour,
		AllowBrowserExtensions: true,
	}
	corsConf.AllowAllOrigins = true
	corsConf.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	corsConf.AllowHeaders = []string{
		"X-ROS-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Referer", "User-Agent", "Accept",
		"Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Host-Uuid",
	}
	return corsConf
}

func HostProxyOptions() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			// Set the necessary headers to allow all requests
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers",
				"X-ROS-Key, Authorization, Content-Type, Upgrade, Origin, "+
					"Connection, Accept-Encoding, Accept-Language, Host, Referer, User-Agent, Accept, "+
					"Access-Control-Allow-Origin, Access-Control-Allow-Headers, host-uuid, Host-Uuid")
			c.Header("Access-Control-Max-Age", "43200") // 12 hours
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
