package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/go-resty/resty/v2"
	"strings"
)

// RestClient is used to invoke Form3 Accounts API.
type RestClient struct {
	client      *resty.Client
	ClientToken string
}

type AzureDeviceConnectionDetails struct {
	HostName              string `json:"azure_host_name"`
	DeviceId              string `json:"azure_device_id"`
	SharedAccessSignature string `json:"azure_shared_access_signature"`
}

func (inst *Instance) checkDeviceConnectionDetails() (*AzureDeviceConnectionDetails, error) {
	inst.edgeazureDebugMsg("makeDeviceConnectionDetails()")
	acdc := AzureDeviceConnectionDetails{}
	acdc.HostName = inst.config.Azure.HostName
	acdc.DeviceId = inst.config.Azure.DeviceId
	acdc.SharedAccessSignature = inst.config.Azure.SharedAccessSignature

	if acdc.HostName == "" || acdc.DeviceId == "" || !strings.HasPrefix(acdc.SharedAccessSignature, "SharedAccessSignature") {
		return nil, errors.New("invalid azure connection details")
	}
	return &acdc, nil
}

// NewNoAuth returns a new instance
func (inst *Instance) NewAzureRestClient(azureDetails *AzureDeviceConnectionDetails) *RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("https://%s/devices/%s/", azureDetails.HostName, azureDetails.DeviceId)
	apiURL := url
	client.SetBaseURL(apiURL)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", azureDetails.SharedAccessSignature)
	return &RestClient{client: client}
}

func (a *RestClient) sendAzureDeviceEventHttp(history *History) (*resty.Response, error) {
	// sas := "SharedAccessSignature sr=Nube-Test-Hub-1.azure-devices.net%2Fdevices%2FTest_Device_1&sig=GsEyASD9jRFqZnEPrYPSFawxn2v%2Fgxx%2BgGqoevv6nlg%3D&se=1668141526"

	/*
		// FOR SENDING PAYLOAD AS BUFFER
		buf := bytes.Buffer{}
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(history)
		historyAsByteArray := buf.Bytes()
	*/

	type AzureHistory struct {
		Device string  `json:"device"`
		Point  string  `json:"point"`
		Value  float64 `json:"value"`
	}

	deviceName, ok := history.Tags["rubix_device_name"]
	if !ok {
		return nil, errors.New("no rubix_device_name tag in history")
	}

	pointName, ok := history.Tags["rubix_point_name"]
	if !ok {
		return nil, errors.New("no rubix_point_name tag in history")
	}

	data := AzureHistory{
		Device: deviceName,
		Point:  pointName,
		Value:  history.Value,
	}

	fmt.Printf("data: %+v", data)

	resp, err := nresty.FormatRestyResponse(a.client.R().
		// SetBody(history).
		SetBody(data).
		// SetHeader("Authorization", sas).
		Post("messages/events?api-version=2020-03-13"))
	if err != nil {
		return resp, err
	}
	return resp, nil
}

/*
func (inst *Instance) getAzureClient(azureConnectionDetails *AzureDeviceConnectionDetails) (*iotdevice.Client, error) {
	inst.edgeazureDebugMsg("getAzureClient()")
	c, err := iotdevice.NewFromConnectionString(iotmqtt.New(), azureConnectionDetails.DeviceConnectionString)
	if err != nil {
		inst.edgeazureErrorMsg("getAzureClient() err:", err)
		return nil, err
	}
	return c, nil
}
*/
