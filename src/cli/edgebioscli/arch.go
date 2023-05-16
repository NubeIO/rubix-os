package edgebioscli

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/src/cli/edgebioscli/ebmodel"
)

func (inst *BiosClient) GetArch() (*ebmodel.Arch, error) {
	url := fmt.Sprintf("/api/system/arch")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&ebmodel.Arch{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ebmodel.Arch), nil
}
