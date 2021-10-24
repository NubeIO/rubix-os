package rubixapi

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixmodel"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest"
	"strings"
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

func splitURL(topic string) (clean, raw *utils.Array) {
	s := strings.SplitAfter(topic, "/")
	clean = utils.NewArray()
	raw = utils.NewArray()
	for _, t := range s {
		if t == "/" {
			clean.Add("EMPTY-TOPIC-SPACE")
		}
		res := strings.ReplaceAll(t, "/", "")
		if res != "" {
			clean.Add(res)
		}
		raw.Add(t)
	}
	return clean, raw
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
	r.URL = fmt.Sprintf("/api/users")
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
	r.URL = fmt.Sprintf("/api/app/control")
	r.Method = POST
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.AppControl)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsInstalled(r Req) (*rubixmodel.AppsInstall, error) {
	r.URL = fmt.Sprintf("/api/app?browser_download_url=true&latest_version=true")
	r.Method = GET
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.AppsInstall)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsLatestVersions(r Req) (*rubixmodel.AppsLatestVersions, error) {
	r.URL = fmt.Sprintf("/api/app/latest_versions")
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

func (a *RestClient) AppsInstall(r Req) (*rubixmodel.AppsInstall, error) {
	r.URL = fmt.Sprintf("/api/app/install")
	r.Method = POST
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.AppsInstall)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsDownloadState(r Req) (*rubixmodel.DownloadState, error) {
	r.URL = fmt.Sprintf("/api/app/install")
	r.Method = GET
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.DownloadState)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AppsDeleteDownloadState(r Req) (*rubixmodel.GeneralResponse, error) {
	r.URL = fmt.Sprintf("/api/app/download_state")
	r.Method = DELETE
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.GeneralResponse)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//SlaveDevices api to discover slaves or slaves added to rubix-service
func (a *RestClient) SlaveDevices(r Req, remoteDevices bool) (*rubixmodel.DiscoveredSlaves, error) {
	if remoteDevices {
		r.URL = fmt.Sprintf("/api/discover/remote_devices")
	} else {
		r.URL = fmt.Sprintf("/api/slaves")
	}
	r.Method = GET
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	body := map[string]json.RawMessage{}
	if err := json.Unmarshal(request.Bytes(), &body); err != nil {
		return nil, err
	}
	s := new(rubixmodel.DiscoveredSlaves)
	for _, e := range body {
		res := new(rubixmodel.Slaves)
		if err := json.Unmarshal(e, &res); err != nil {
			return nil, err
		}
		s.Slaves = append(s.Slaves, *res)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (a *RestClient) SlaveDevicesAddDelete(r Req, delete bool, globalUUID string) (*rubixmodel.GeneralResponse, error) {
	if delete {
		r.Method = DELETE
		r.URL = fmt.Sprintf("/api/slaves/%s", globalUUID)
	} else {
		r.URL = fmt.Sprintf("/api/slaves")
		r.Method = POST
	}
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.GeneralResponse)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) WiresPlat(r Req, update bool) (*rubixmodel.WiresPlat, error) {
	r.URL = fmt.Sprintf("/api/wires/plat")
	if update {
		r.Method = PUT
	} else {
		r.Method = GET
	}
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	res := new(rubixmodel.WiresPlat)
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *RestClient) AnyRequest(r Req) (interface{}, error) {
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	var res interface{}
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//AnyRequestWithBody is for post, put or patch
func (a *RestClient) AnyRequestWithBody(r Req) (interface{}, error) {
	request, err := Request(r)
	if err != nil {
		return nil, err
	}
	var res interface{}
	err = json.Unmarshal(request.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
