package lwrest

import (
	"fmt"
	lwmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/model"
)

const limit = "50"
const orgID = "1"

// GetOrganizations get all
func (a *RestClient) GetOrganizations() (*lwmodel.Organizations, error) {
	q := fmt.Sprintf("/api/organizations?limit=%s", limit)
	resp, err := a.client.R().
		SetResult(lwmodel.Organizations{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetOrganizations %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Organizations), nil
}
