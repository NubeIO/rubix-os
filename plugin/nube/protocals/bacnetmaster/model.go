package main

type readBody struct {
	ObjectType     string `json:"objectType"`
	ObjectInstance string `json:"objectInstance"`
	Property       string `json:"property"`
	DeviceInstance string `json:"deviceInstance"`
	Mac            string `json:"mac"`
	TxnSource      string `json:"txn_source"`
	TxnNumber      string `json:"txn_number"`
}

type writeBody struct {
	ObjectType     string  `json:"object_type"`     // 1 analogue output
	ObjectInstance int     `json:"object_instance"` // 1
	Property       int     `json:"property"`        // 85 presetValue
	DeviceInstance int     `json:"device_instance"` // 2508
	Mac            string  `json:"mac"`             // 192.168.15.10:47808
	Value          float64 `json:"value"`
	TxnSource      string  `json:"txn_source"`
	TxnNumber      int     `json:"txn_number"`
}

type whoIsBody struct {
	NetworkNumber     int    `json:"network_number"`
	DeviceInstanceMin int    `json:"device_instance_min"`
	DeviceInstanceMax int    `json:"device_instance_max"`
	Timeout           int    `json:"timeout"`
	TxnSource         string `json:"txn_source"`
	TxnNumber         string `json:"txn_number"`
}
