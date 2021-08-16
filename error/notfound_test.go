package error

import (
	"net/http/httptest"
	"testing"

	"github.com/NubeDev/flow-framework/mode"
	"github.com/gin-gonic/gin"
)

func TestNotFound(t *testing.T) {
	mode.Set(mode.TestDev)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	NotFound()(ctx)

	assertJSONResponse(t, rec, 404, `{"errorCode":404, "errorDescription":"page not found", "error":"Not Found"}`)
}
