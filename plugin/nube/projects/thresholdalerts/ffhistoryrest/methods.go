package ffhistoryrest

import (
	"github.com/NubeIO/rubix-os/nresty"
)

type FFCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FFToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type FFHistory struct {
	Value            float64 `json:"value" dataframe:"value"`
	Timestamp        string  `json:"timestamp" dataframe:"timestamp"`
	RubixNetworkUUID string  `json:"rubix_network_uuid" dataframe:"rubix_network_uuid"`
	RubixNetworkName string  `json:"rubix_network_name" dataframe:"rubix_network_name"`
	RubixDeviceUUID  string  `json:"rubix_device_uuid" dataframe:"rubix_device_uuid"`
	RubixDeviceName  string  `json:"rubix_device_name" dataframe:"rubix_device_name"`
	RubixPointUUID   string  `json:"rubix_point_uuid" dataframe:"rubix_point_uuid"`
	RubixPointName   string  `json:"rubix_point_name" dataframe:"rubix_point_name"`
}

func (a *RestClient) GetFFToken(user, pass string) (*FFToken, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(FFCredentials{
			Username: user,
			Password: pass,
		}).
		SetResult(FFToken{}).
		Post("/api/users/login"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*FFToken), nil
}

func (a *RestClient) GetFFHistories(token FFToken, queryParams string) (*[]FFHistory, error) {
	url := "/ff/api/plugins/api/postgres/histories" + queryParams
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Authorization", token.AccessToken).
		SetHeader("Accept", "application/json").
		SetResult([]FFHistory{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]FFHistory), nil
}
