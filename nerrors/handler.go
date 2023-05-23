package nerrors

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Handler creates a gin middleware for handling errors.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypeBind:
					errs, ok := e.Err.(validator.ValidationErrors)
					if !ok {
						writeError(c, e.Error())
						return
					}
					var stringErrors []string
					for _, err := range errs {
						stringErrors = append(stringErrors, validationErrorToText(err))
					}
					writeError(c, strings.Join(stringErrors, "; "))
				default:
					writeError(c, e.Err.Error())
				}
			}
		}
	}
}

func validationErrorToText(e validator.FieldError) string {
	runes := []rune(e.Field())
	runes[0] = unicode.ToLower(runes[0])
	fieldName := string(runes)
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", fieldName)
	case "max":
		return fmt.Sprintf("Field '%s' must be less or equal to %s", fieldName, e.Param())
	case "min":
		return fmt.Sprintf("Field '%s' must be more or equal to %s", fieldName, e.Param())
	}
	return fmt.Sprintf("Field '%s' is not valid", fieldName)
}

func writeError(ctx *gin.Context, err string) {
	statusCode := http.StatusBadRequest
	if ctx.Writer.Status() != http.StatusOK {
		statusCode = ctx.Writer.Status()
	}
	message := fmt.Sprintf("[%d]: %s", statusCode, err)
	ctx.JSON(statusCode, interfaces.Message{Message: message})
}
