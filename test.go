package main

import (
	"fmt"
	"os/exec"
)

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

func main() {
	cmd := "sudo ufw enable"
	o, err := RunCMD(cmd, false)
	fmt.Println(string(o), err)

}
