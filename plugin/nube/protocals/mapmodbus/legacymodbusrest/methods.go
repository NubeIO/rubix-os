package legacymodbusrest

import (
	"errors"
	"github.com/NubeIO/flow-framework/nresty"
)

type ModbusNet struct {
	UUID                         string      `json:"uuid"`
	Name                         string      `json:"name"`
	Enable                       bool        `json:"enable"`
	RTUPort                      string      `json:"rtu_port"`
	RTUSpeed                     int         `json:"rtu_speed"`
	RTUStopBits                  int         `json:"rtu_stop_bits"`
	RTUParity                    string      `json:"rtu_parity"`
	RTUByteSize                  int         `json:"rtu_byte_size"`
	TCPIP                        string      `json:"tcp_ip"`
	TCPPort                      int         `json:"tcp_port"`
	Type                         string      `json:"type"`
	Timeout                      float64     `json:"timeout"`
	PollingIntervalRuntime       float64     `json:"polling_interval_runtime"`
	PointIntervalMSBetweenPoints float64     `json:"point_interval_ms_between_points"`
	Devices                      []ModbusDev `json:"devices"`
}

type ModbusDev struct {
	UUID        string      `json:"uuid"`
	Type        string      `json:"type"`
	NetworkUUID string      `json:"network_uuid"`
	Name        string      `json:"name"`
	Enable      bool        `json:"enable"`
	Address     int         `json:"address"`
	ZeroBased   bool        `json:"zero_based"`
	Points      []ModbusPnt `json:"points"`
}

type ModbusPnt struct {
	UUID           string  `json:"uuid"`
	DeviceUUID     string  `json:"device_uuid"`
	Name           string  `json:"name"`
	Enable         bool    `json:"enable"`
	Writeable      bool    `json:"writeable"`
	Register       int     `json:"register"`
	FunctionCode   string  `json:"function_code"`
	DataType       string  `json:"data_type"`
	DataEndian     string  `json:"data_endian"`
	WriteValueOnce bool    `json:"write_value_once"`
	FallbackValue  float64 `json:"fallback_value"`
}

func (a *RestClient) GetLegacyModbusNetworksAndDevices() (*[]ModbusNet, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		// SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2Njk5MjI4NDYsImlhdCI6MTY2NzMzMDg0Niwic3ViIjoiYWRtaW4ifQ.epFqcUgTj03c7tIU26icpQyOGUkOW4ki5BINbq5rYVE").
		SetHeader("Accept", "*/*").
		SetResult([]ModbusNet{}).
		// Get("/modbus/api/modbus/networks?with_children=true"))
		Get("/api/modbus/networks?with_children=true"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]ModbusNet), nil
}

func (a *RestClient) GetLegacyModbusDeviceAndPoints(modbusDevUUID string) (*ModbusDev, error) {
	if modbusDevUUID == "" {
		return nil, errors.New("GetLegacyModbusDevicesAndPoints(): modbus device is nil")
	}
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult(ModbusDev{}).
		// SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2Njk5MjI4NDYsImlhdCI6MTY2NzMzMDg0Niwic3ViIjoiYWRtaW4ifQ.epFqcUgTj03c7tIU26icpQyOGUkOW4ki5BINbq5rYVE").
		// Get("/modbus/api/modbus/devices/uuid/" + modbusDevUUID + "?with_children=true"))
		Get("/api/modbus/devices/uuid/" + modbusDevUUID + "?with_children=true"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ModbusDev), nil
}
