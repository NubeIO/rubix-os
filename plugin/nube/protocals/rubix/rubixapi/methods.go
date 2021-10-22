package rubixapi

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixmodel"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest"
)

const (
	GET = iota + 1
	POST
	PATCH
	PUT
	DELETE
)

type Req struct {
	Token          string
	Auth           bool
	RequestBuilder *rest.RequestBuilder
	Method         int
	Body           interface{}
	URL            string
	Repose         interface{}
}

// Request builder. default query is GET, default Req.URL is /
func Request(req Req) (*rest.Response, error) {
	if req.Method == 0 {
		req.Method = GET
	}
	if req.URL == "" {
		req.URL = "/"
	}
	if req.Auth {
		req.RequestBuilder.Headers.Add("Authorization", req.Token)
	}
	var resp *rest.Response
	switch req.Method {
	case GET:
		resp = req.RequestBuilder.Get(req.URL)
	case POST:
		resp = req.RequestBuilder.Post(req.URL, req.Body)
	case PATCH:
		resp = req.RequestBuilder.Patch(req.URL, req.Body)
	case PUT:
		resp = req.RequestBuilder.Put(req.URL, req.Body)
	case DELETE:
		resp = req.RequestBuilder.Delete(req.URL)
	}
	return resp, resp.Err
}

func (a *RestClient) GetToken(r Req) (rubixmodel.TokenResponse, rest.Response, error) {
	request, err := Request(r)
	if err != nil {
		return rubixmodel.TokenResponse{}, rest.Response{}, err
	}
	var res rubixmodel.TokenResponse
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return rubixmodel.TokenResponse{}, rest.Response{}, err
	}
	return res, *request, nil
}

func (a *RestClient) GetUsers(r Req) (*rubixmodel.UserResponse, error) {
	q := fmt.Sprintf("/api/users")
	r.URL = q
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.UserResponse)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppControl(r Req) (*rubixmodel.AppControl, error) {
	q := fmt.Sprintf("/api/app/control")
	r.URL = q
	r.Method = POST
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	fmt.Printf(request.String())
	res := new(rubixmodel.AppControl)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsInstalled(r Req) (*rubixmodel.AppsInstall, error) {
	q := fmt.Sprintf("/api/app?browser_download_url=true&latest_version=true")
	r.URL = q
	r.Method = GET
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	fmt.Printf(request.String())
	res := new(rubixmodel.AppsInstall)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsLatestVersions(r Req) (*rubixmodel.AppsLatestVersions, error) {
	q := fmt.Sprintf("/api/app/latest_versions")
	r.URL = q
	r.Method = GET
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.AppsLatestVersions)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
