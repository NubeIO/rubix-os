package host

import (
	"fmt"
	"os/exec"
)

func CMDRun(sh string, debug bool) ([]byte, error) {
	cmd := exec.Command("bash", "-c", sh)
	res, e := cmd.Output()

	if debug {
		fmt.Printf("[cmd debug] %s\n", cmd.String())
	}
	if e != nil {
		defer cmd.Process.Kill()
		return nil, e
	}

	defer cmd.Process.Kill()
	return res, e
}
