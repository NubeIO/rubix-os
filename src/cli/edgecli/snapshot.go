package edgecli

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"mime"
	"os"
)

func (inst *Client) CreateSnapshot() ([]byte, string, error) {
	url := fmt.Sprintf("/api/snapshots/create")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		Post(url))
	if err != nil {
		return nil, "", err
	}
	_, param, err := mime.ParseMediaType(resp.RawResponse.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, "", err
	}
	return resp.Body(), param["filename"], nil
}

func (inst *Client) RestoreSnapshot(filename string, reader *os.File) error {
	url := fmt.Sprintf("/api/snapshots/restore")
	_, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetFileReader("file", filename, reader).
		Post(url))
	return err
}
