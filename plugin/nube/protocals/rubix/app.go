package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/rubix/rubixapi"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/rubix/rubixmodel"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest"
	"net/http"
	"time"
)

func (i *Instance) request(body rubixmodel.TokenBody, integration model.Integration) (res rubixapi.Req, err error) {
	cli := rubixapi.New()
	h := fmt.Sprintf("http://%s:%s", integration.IP, integration.PORT)
	var rb = rest.RequestBuilder{
		Timeout:        5000 * time.Millisecond,
		BaseURL:        h,
		ContentType:    rest.JSON,
		DisableCache:   false,
		DisableTimeout: false,
	}
	var rr rubixapi.Req
	rr.RequestBuilder = &rb
	rr.URL = "/api/users/login"
	rr.Method = rubixapi.POST
	rr.Body = body
	r, _, err := cli.GetToken(rr)
	if err != nil {
		return rr, nil
	}
	rr.Token = r.AccessToken
	rr.RequestBuilder = &rb
	return rr, nil
}

type tokenExp struct {
	Expired bool
}

func (i *Instance) getIntegration(uuid string, name string) (res rubixapi.Req, err error) {
	if name != "" {
		integration, err := i.db.GetIntegrationByName(name)
		if err != nil {
			return rubixapi.Req{}, err
		}
		var t tokenExp
		t.Expired = false
		i.store.Set(integration.UUID, t, -1)
		var tb rubixmodel.TokenBody
		tb.Username = integration.Username
		tb.Password = integration.Password
		req, err := i.request(tb, *integration)
		headers := make(http.Header)
		headers.Add("Authorization", req.Token)
		h := fmt.Sprintf("http://%s:%s", integration.IP, integration.PORT)
		var rb = rest.RequestBuilder{
			Headers:        headers,
			Timeout:        5000 * time.Millisecond,
			BaseURL:        h,
			ContentType:    rest.JSON,
			DisableCache:   false,
			DisableTimeout: false,
		}
		var rr rubixapi.Req
		rr.RequestBuilder = &rb
		return rr, err
	} else {
		_, err := i.db.GetIntegration(uuid)
		if err != nil {
			return rubixapi.Req{}, err
		}
		return rubixapi.Req{}, err
	}

}
