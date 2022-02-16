package utils

import (
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstNotNilTest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var a *float64 = nil
	b := utils.NewFloat64(2.3)
	result := utils.FirstNotNilFloat(a, b)
	assert.Equal(t, 2.3, *result)

	var c *float64 = nil
	result2 := utils.FirstNotNilFloat(a, c)
	assert.Nil(t, result2)
}
