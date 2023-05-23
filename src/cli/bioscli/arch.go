package bioscli

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *BiosClient) GetArch() (*interfaces.Arch, error) {
	url := fmt.Sprintf("/api/system/arch")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.Arch{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*interfaces.Arch), nil
}
