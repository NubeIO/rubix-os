package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest/v1/rest"
	"github.com/go-resty/resty/v2"
)

// An Error maps to Form3 API error responses
type Error struct {
	Code    int    `json:"error_code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

func failedResponse(err error, resp *resty.Response) error {
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return getAPIError(resp)
	}
	if rest.StatusCodesAllBad(resp.StatusCode()) {
		return getAPIError(resp)
	}
	return nil
}

// Convert error response into error message
func getAPIError(resp *resty.Response) error {
	e := new(Error)
	e.Code = resp.StatusCode()
	e.Message = resp.String()
	return fmt.Errorf("request failed [%d]: %s", e.Code, e.Message)
}
