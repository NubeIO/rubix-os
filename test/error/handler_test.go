package error

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nerrors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorInternal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(500, errors.New("something went wrong"))

	nerrors.Handler()(ctx)

	assertJSONResponse(t, rec, 500, `{"message":"[500]: something went wrong"}`)
}

func TestBindingErrorDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something")).SetType(gin.ErrorTypeBind)

	nerrors.Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"message":"[400]: you need todo something"}`)
}

func TestDefaultErrorBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something"))

	nerrors.Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"message":"[400]: you need todo something"}`)
}

type testValidate struct {
	Username string `json:"username" binding:"required"`
	Mail     string `json:"mail" binding:"email"`
	Age      int    `json:"age" binding:"max=100"`
	Limit    int    `json:"limit" binding:"min=50"`
}

func TestValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest("GET", "/uri", nil)

	assert.Error(t, ctx.Bind(&testValidate{Age: 150, Limit: 20}))
	nerrors.Handler()(ctx)

	err := new(interfaces.Message)
	_ = json.NewDecoder(rec.Body).Decode(err)
	assert.Contains(t, err.Message, "[400]")
	assert.Contains(t, err.Message, "Field 'username' is required")
	assert.Contains(t, err.Message, "Field 'mail' is not valid")
	assert.Contains(t, err.Message, "Field 'age' must be less or equal to 100")
	assert.Contains(t, err.Message, "Field 'limit' must be more or equal to 50")
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, code int, json string) {
	bytes, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, code, rec.Code)
	assert.JSONEq(t, json, string(bytes))
}
