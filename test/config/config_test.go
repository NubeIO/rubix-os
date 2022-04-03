package config

import (
	"github.com/NubeIO/flow-framework/config"
	"github.com/gin-gonic/gin"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("FLOW_DEFAULTUSER_NAME", "jmattheis")
	os.Setenv("FLOW_SERVER_SSL_LETSENCRYPT_HOSTS", "- push.example.tld\n- push.other.tld")
	os.Setenv("FLOW_SERVER_RESPONSEHEADERS",
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
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Cors.AllowOrigins)
	assert.Equal(t, []string{"GET", "POST"}, conf.Server.Cors.AllowMethods)
	assert.Equal(t, []string{"Authorization", "content-type"}, conf.Server.Cors.AllowHeaders)
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Stream.AllowedOrigins)

	os.Unsetenv("FLOW_DEFAULTUSER_NAME")
	os.Unsetenv("FLOW_SERVER_SSL_LETSENCRYPT_HOSTS")
	os.Unsetenv("FLOW_SERVER_RESPONSEHEADERS")
	os.Unsetenv("FLOW_SERVER_CORS_ALLOWORIGINS")
	os.Unsetenv("FLOW_SERVER_CORS_ALLOWMETHODS")
	os.Unsetenv("FLOW_SERVER_CORS_ALLOWHEADERS")
	os.Unsetenv("FLOW_SERVER_STREAM_ALLOWEDORIGINS")
}
