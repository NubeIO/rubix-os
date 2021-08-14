package test_test

import (
	"testing"

	"github.com/NubeDev/plug-framework/auth"
	"github.com/NubeDev/plug-framework/mode"
	"github.com/NubeDev/plug-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFakeAuth(t *testing.T) {
	mode.Set(mode.TestDev)

	ctx, _ := gin.CreateTestContext(nil)
	test.WithUser(ctx, 5)
	assert.Equal(t, uint(5), auth.GetUserID(ctx))
}
