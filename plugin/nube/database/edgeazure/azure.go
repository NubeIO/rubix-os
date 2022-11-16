package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/go-resty/resty/v2"
)

// RestClient is used to invoke Form3 Accounts API.
type RestClient struct {
	client      *resty.Client
	ClientToken string
}

// NewNoAuth returns a new instance
func (inst *Instance) NewAzureRestClient() *RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("https://%s/devices/%s/", inst.config.Azure.HostName, inst.config.Azure.DeviceId)
	apiURL := url
	client.SetBaseURL(apiURL)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	return &RestClient{client: client}
}

type AzureDeviceConnectionDetails struct {
	HostName               string `json:"azure_host_name"`
	DeviceId               string `json:"azure_device_id"`
	SharedAccessKey        string `json:"azure_shared_access_key"`
	DeviceConnectionString string `json:"azure_device_connection_string"`
}

func (inst *Instance) makeDeviceConnectionDetails() (*AzureDeviceConnectionDetails, error) {
	inst.edgeazureDebugMsg("makeDeviceConnectionDetails()")
	acdc := AzureDeviceConnectionDetails{}
	acdc.HostName = inst.config.Azure.HostName
	acdc.DeviceId = inst.config.Azure.DeviceId
	acdc.SharedAccessKey = inst.config.Azure.SharedAccessKey

	if acdc.HostName == "" || acdc.DeviceId == "" || acdc.SharedAccessKey == "" {
		return nil, errors.New("invalid azure config")
	}
	acdc.DeviceConnectionString = "HostName=" + acdc.HostName + ";DeviceId=" + acdc.DeviceId + ";SharedAccessKey=" + acdc.SharedAccessKey
	return &acdc, nil
}

func (inst *Instance) getAzureClient(azureConnectionDetails *AzureDeviceConnectionDetails) (*iotdevice.Client, error) {
	inst.edgeazureDebugMsg("getAzureClient()")
	c, err := iotdevice.NewFromConnectionString(iotmqtt.New(), azureConnectionDetails.DeviceConnectionString)
	if err != nil {
		inst.edgeazureErrorMsg("getAzureClient() err:", err)
		return nil, err
	}
	return c, nil
}

func (a *RestClient) sendAzureDeviceEventHttp(history *History) (*resty.Response, error) {
	sas := "SharedAccessSignature sr=Nube-Test-Hub-1.azure-devices.net%2Fdevices%2FTest_Device_1&sig=GsEyASD9jRFqZnEPrYPSFawxn2v%2Fgxx%2BgGqoevv6nlg%3D&se=1668141526"

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(history)
	historyAsByteArray := buf.Bytes()

	resp, err := nresty.FormatRestyResponse(a.client.R().
		// SetBody(history).
		SetBody(historyAsByteArray).
		SetHeader("Authorization", sas).
		Post("messages/events?api-version=2020-03-13"))
	if err != nil {
		return resp, err
	}
	return resp, nil
}
