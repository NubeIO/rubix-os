package edgebioscli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/src/cli/edgebioscli/ebmodel"
)

func (inst *BiosClient) GetRubixEdgeVersion() (*ebmodel.Version, error) {
	installLocation := fmt.Sprintf("/data/rubix-service/apps/install/%s", rubixEdgeName)
	url := fmt.Sprintf("/api/files/list?path=%s", installLocation)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	versions := resp.Result().(*[]fileutils.FileDetails)
	if versions != nil && len(*versions) > 0 {
		return &ebmodel.Version{Version: (*versions)[0].Name}, nil
	}
	return nil, errors.New("doesn't found the installation file")
}
