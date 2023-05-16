package system

import (
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/str"
	"os/exec"
	"strings"
)

type Node struct {
	IsInstalled      bool   `json:"is_installed"`
	InstalledVersion string `json:"installed_version"`
}

func (inst *System) NodeGetVersion() (*Node, error) {
	cmd := exec.Command("/usr/bin/node", "-v")
	output, err := cmd.Output()
	res := cleanCommand(string(output), cmd, err, debug)
	node := &Node{}
	if strings.Contains(res, "v") {
		node.InstalledVersion = str.RemoveNewLine(res)
		node.IsInstalled = true
		return node, err
	} else {
		node.InstalledVersion = res
		node.IsInstalled = false
		return node, err
	}
}
