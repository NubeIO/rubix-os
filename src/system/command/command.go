package command

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run runs given command with parameters and return combined output
func Run(cmdAndParams ...string) (string, error) {
	if len(cmdAndParams) <= 0 {
		return "", fmt.Errorf("no command provided")
	}

	output, err := exec.Command(cmdAndParams[0], cmdAndParams[1:]...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

// SudoRun runs given command with parameters and return combined output (with sudo)
func SudoRun(cmdAndParams ...string) (string, error) {
	if len(cmdAndParams) <= 0 {
		return "", fmt.Errorf("no command provided")
	}
	output, err := exec.Command("sudo", cmdAndParams...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

func RunCMD(sh string, debug bool) ([]byte, error) {
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
