package system

import (
	"github.com/NubeIO/lib-date/datectl"
	"github.com/NubeIO/lib-dhcpd/dhcpd"
	"github.com/NubeIO/lib-networking/networking"
	systats "github.com/NubeIO/lib-system"
	"github.com/NubeIO/lib-ufw/ufw"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type System struct {
	ufw     *ufw.System
	datectl *datectl.DateCTL
	dhcp    *dhcpd.DHCP
	syStats systats.SyStats
}

var debug = true
var nets = networking.New()

type Message struct {
	Message string `json:"message"`
}

func New(inst *System) *System {
	if inst == nil {
		inst = &System{}
	}
	inst.ufw = ufw.New(&ufw.System{})
	inst.datectl = datectl.New(&datectl.DateCTL{})
	inst.dhcp = dhcpd.New(&dhcpd.DHCP{})
	inst.syStats = systats.New()
	return inst
}

func cleanCommand(resp string, cmd *exec.Cmd, err error, debug ...bool) string {
	outAsString := strings.TrimRight(resp, "\n")
	if len(debug) > 0 {
		if debug[0] {
			log.Infof("cmd %s", cmd.String())
			log.Infof("path %s", cmd.Path)
			if err != nil {
				log.Errorf("err: %s", err.Error())
			} else {
				log.Infof("resp: %s", outAsString)
			}
		}
	}
	return outAsString
}
