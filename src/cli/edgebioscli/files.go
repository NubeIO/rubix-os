package edgebioscli

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *BiosClient) ListFilesV2(path string) ([]fileutils.FileDetails, error, error) {
	url := fmt.Sprintf("/api/files/list?path=%s", path)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return *resp.Result().(*[]fileutils.FileDetails), nil, nil
}

func (inst *BiosClient) DeleteFiles(path string) (*interfaces.Message, error, error) {
	url := fmt.Sprintf("/api/files/delete-all?path=%s", path)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Delete(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*interfaces.Message), nil, nil
}
