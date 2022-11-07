package legacylorarest

import "github.com/NubeIO/flow-framework/nresty"

type LoRaNet struct {
	Port     string  `json:"port"`
	BaudRate int     `json:"baud_rate"`
	StopBits int     `json:"stop_bits"`
	Parity   string  `json:"parity"`
	ByteSize int     `json:"byte_size"`
	Timeout  float64 `json:"timeout"`
}

type LoRaDev struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	Description string `json:"description"`
	AI1Config   string `json:"ai_1_config"`
	AI2Config   string `json:"ai_2_config"`
	AI3Config   string `json:"ai_3_config"`
}

func (a *RestClient) GetLegacyLoRaNetwork() (*LoRaNet, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		// SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2Njk5MjI4NDYsImlhdCI6MTY2NzMzMDg0Niwic3ViIjoiYWRtaW4ifQ.epFqcUgTj03c7tIU26icpQyOGUkOW4ki5BINbq5rYVE").
		SetHeader("Accept", "*/*").
		SetResult(LoRaNet{}).
		Get("/api/lora/networks"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*LoRaNet), nil
}

func (a *RestClient) GetLegacyLoRaDevices() (*[]LoRaDev, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult([]LoRaDev{}).
		// SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2Njk5MjI4NDYsImlhdCI6MTY2NzMzMDg0Niwic3ViIjoiYWRtaW4ifQ.epFqcUgTj03c7tIU26icpQyOGUkOW4ki5BINbq5rYVE").
		Get("/api/lora/devices"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]LoRaDev), nil
}
