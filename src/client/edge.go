package client

import "github.com/go-resty/resty/v2"

type Client struct {
	Rest          *resty.Client
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	HTTPS         *bool  `json:"https"`
	ExternalToken string `json:"external_token"`
}
