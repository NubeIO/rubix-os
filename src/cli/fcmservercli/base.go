package fcmservercli

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	mutex   = &sync.RWMutex{}
	clients = map[string]*FcmServerClient{}
)

type FcmServerClient struct {
	client  *resty.Client
	BaseUrl string `json:"base_url"`
	Key     string `json:"key"`
}

func New(cli *FcmServerClient) *FcmServerClient {
	mutex.Lock()
	defer mutex.Unlock()

	if cli == nil {
		log.Fatal("fcm cli can not be empty")
		return nil
	}
	cli.BaseUrl = "https://fcm.googleapis.com/fcm"
	if fcmServerClient, found := clients[cli.BaseUrl]; found {
		fcmServerClient.client.SetHeader("Authorization", composeKey(cli.Key))
		return fcmServerClient
	}
	client := resty.New()
	client.SetBaseURL(cli.BaseUrl)
	client.SetHeader("Authorization", composeKey(cli.Key))
	cli.client = client
	clients[cli.BaseUrl] = cli
	return cli
}

func composeKey(key string) string {
	return fmt.Sprintf("key=%s", key)
}
