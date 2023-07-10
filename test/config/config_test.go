package config

import (
	"github.com/NubeIO/rubix-os/config"
	"github.com/gin-gonic/gin"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ROS_SERVER_RESPONSEHEADERS",
		"Access-Control-Allow-Origin: \"*\"\nAccess-Control-Allow-Methods: \"GET,POST\"",
	)
	os.Setenv("FLOW_SERVER_CORS_ALLOWORIGINS", "- \".+.example.com\"\n- \"otherdomain.com\"")
	os.Setenv("FLOW_SERVER_CORS_ALLOWMETHODS", "- \"GET\"\n- \"POST\"")
	os.Setenv("FLOW_SERVER_CORS_ALLOWHEADERS", "- \"Authorization\"\n- \"content-type\"")
	os.Setenv("FLOW_SERVER_STREAM_ALLOWEDORIGINS", "- \".+.example.com\"\n- \"otherdomain.com\"")
	conf := config.CreateApp()
	assert.Equal(t, 1660, conf.Server.Port, "should use defaults")
	assert.Equal(t, "*", conf.Server.ResponseHeaders["Access-Control-Allow-Origin"])
	assert.Equal(t, "GET,POST", conf.Server.ResponseHeaders["Access-Control-Allow-Methods"])

	os.Unsetenv("ROS_DEFAULTUSER_NAME")
}
