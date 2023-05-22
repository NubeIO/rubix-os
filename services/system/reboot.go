package system

import "os/exec"

func (inst *System) RebootHost() (*Message, error) {
	cmd := exec.Command("shutdown", "-r", "now")
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: "restarted ok",
	}, err
}
