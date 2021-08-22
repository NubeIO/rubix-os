package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// An Error maps to Form3 API error responses
type Error struct {
	Code    string `json:"error_code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

// Convert error response into error message
func getAPIError(resp *resty.Response) error {
	apiError := resp.Error().(*Error)
	return fmt.Errorf("request failed [%s]: %s", apiError.Code, apiError.Message)
}
