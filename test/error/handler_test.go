package error

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-framework/error"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorInternal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(500, errors.New("something went wrong"))

	error.Handler()(ctx)

	assertJSONResponse(t, rec, 500, `{"errorCode":500, "errorDescription":"something went wrong", "error":"Internal Server Error"}`)
}

func TestBindingErrorDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something")).SetType(gin.ErrorTypeBind)

	error.Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"errorCode":400, "errorDescription":"you need todo something", "error":"Bad Request"}`)
}

func TestDefaultErrorBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.AbortWithError(400, errors.New("you need todo something"))

	error.Handler()(ctx)

	assertJSONResponse(t, rec, 400, `{"errorCode":400, "errorDescription":"you need todo something", "error":"Bad Request"}`)
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
	error.Handler()(ctx)

	err := new(model.Error)
	json.NewDecoder(rec.Body).Decode(err)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "Bad Request", err.Error)
	assert.Equal(t, 400, err.ErrorCode)
	assert.Contains(t, err.ErrorDescription, "Field 'username' is required")
	assert.Contains(t, err.ErrorDescription, "Field 'mail' is not valid")
	assert.Contains(t, err.ErrorDescription, "Field 'age' must be less or equal to 100")
	assert.Contains(t, err.ErrorDescription, "Field 'limit' must be more or equal to 50")
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, code int, json string) {
	bytes, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, code, rec.Code)
	assert.JSONEq(t, json, string(bytes))
}