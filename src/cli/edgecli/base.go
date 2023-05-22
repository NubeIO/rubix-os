package edgecli

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex              = &sync.RWMutex{}
	clients            = map[string]*Client{}
	clientsFastTimeout = map[string]*Client{}
)

type Client struct {
	Rest          *resty.Client
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	HTTPS         *bool  `json:"https"`
	ExternalToken string `json:"external_token"`
}

// The dialTimeout normally catches: when the server is unreachable and returns i/o timeout within 2 seconds.
// Otherwise, the i/o timeout takes 1.3 minutes on default; which is a very long time for waiting.
// It uses the DialTimeout function of the net package which connects to a server address on a named network before
// a specified timeout.
func dialTimeout(_ context.Context, network, addr string) (net.Conn, error) {
	timeout := 2 * time.Second
	return net.DialTimeout(network, addr, timeout)
}

var transport = http.Transport{
	DialContext: dialTimeout,
}

func New(cli *Client) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if cli == nil {
		log.Fatal("edge client cli can not be empty")
		return nil
	}
	baseURL := getBaseUrl(cli)
	if client, found := clients[baseURL]; found {
		client.Rest.SetHeader("Authorization", composeToken(cli.ExternalToken))
		return client
	}
	rest := resty.New()
	rest.SetBaseURL(baseURL)
	rest.SetHeader("Authorization", composeToken(cli.ExternalToken))
	cli.Rest = rest
	clients[baseURL] = cli
	return cli
}

func NewFastTimeout(cli *Client) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if cli == nil {
		log.Fatal("edge client cli can not be empty")
		return nil
	}
	baseURL := getBaseUrl(cli)
	if client, found := clientsFastTimeout[baseURL]; found {
		client.Rest.SetHeader("Authorization", composeToken(cli.ExternalToken))
		return client
	}
	rest := resty.New()
	rest.SetBaseURL(baseURL)
	rest.Header.Set("Authorization", composeToken(cli.ExternalToken))
	rest.SetTransport(&transport)
	cli.Rest = rest
	clientsFastTimeout[baseURL] = cli
	return cli
}

func getBaseUrl(cli *Client) string {
	cli.Rest = resty.New()
	if cli.Ip == "" {
		cli.Ip = "0.0.0.0"
	}
	if cli.Port == 0 {
		cli.Port = 1661
	}
	var baseURL string
	if cli.HTTPS != nil && *cli.HTTPS {
		baseURL = fmt.Sprintf("https://%s:%d", cli.Ip, cli.Port)
	} else {
		baseURL = fmt.Sprintf("http://%s:%d", cli.Ip, cli.Port)
	}
	return baseURL
}

func composeToken(token string) string {
	return fmt.Sprintf("External %s", token)
}
