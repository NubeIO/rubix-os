package ufw

import (
	"github.com/NubeDev/flow-framework/src/system/command"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func EnableFirewall() {
	cmd := "sudo ufw enable"
	_, err := command.RunCMD(cmd, false)
	if err != nil {
		log.Error("firewall: Enable Error: ", err)
	}
}

func FirewallStatus() (result string, err error) {
	cmd := "sudo ufw status"
	c, err := command.RunCMD(cmd, false)
	return string(c), err
}

func DisableFirewall() {
	cmd := exec.Command("ufw", "disable")
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Disable Error: ", err)
	}
}

func ProtocolRateLimit(protocol string) {
	cmd := exec.Command("ufw", "limit", protocol)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Protocol Rate Limit Error: ", err)
	}
}

func ResetFirewall() {
	cmd := exec.Command("ufw", "reset")
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Reset Error: ", err)
	}
}

func AllowProtocol(protocol string) {
	cmd := exec.Command("ufw", "allow", protocol)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow Protocol Error: ", err)
	}
}

func AllowSrcIPProtocol(srcip, port string) {
	cmd := exec.Command("ufw", "allow", "from", srcip, "to", "any", "port", port)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow SrcIP Protocol Error: ", err)
	}
}

func DenySrcIP(targetIP string) {
	cmd := exec.Command("ufw", "deny", "from", targetIP)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Configuration Error: ", err)
	}
}

func AllowSrcIPInterface(targetIP, direction, vulInterface string) {
	cmd := exec.Command("ufw", "allow", direction, "on", vulInterface, "from", targetIP)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow SrcIP Interface Error: ", err)
	}
}

func AllowSrcIP(targetIP string) {
	cmd := exec.Command("ufw", "allow", "in", "from", targetIP)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow SrcIP Error: ", err)
	}
}

func AllowInterface(direction, secureInterface string) {
	cmd := exec.Command("ufw", "allow", direction, "on", secureInterface)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow Interface Error: ", err)
	}
}

func AllowIPPortInterface(targetIP, targetPort, vulInterface string) {
	cmd := exec.Command("ufw", "allow", "in", "on", vulInterface, "to", targetIP, "port", targetPort)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Configuration Error: ", err)
	}
}

func AllowForwardInterface(secureInterface, vulInterface string) {
	cmd := exec.Command("ufw", "route", "allow", "in", "on", secureInterface, "out", "on", vulInterface)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow Forward Interface Error: ", err)
	}
}

func AllowForwardInterfacePort(srcInterface, port string) {
	cmd := exec.Command("ufw", "route", "allow", "in", "on", srcInterface, "to", "any", "port", port)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Allow Forward Interface Error: ", err)
	}
}

func AddDefaultGateway(gatewayIP string) {
	cmd := exec.Command("route", "add", "default", "gw", gatewayIP)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Add Default Gateway Error: ", err)
	}
}

func DenySrcIPInterface(targetIP, vulInterface string) {
	cmd := exec.Command("ufw", "deny", "in", "on", vulInterface, "from", targetIP)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Configuration Error: ", err)
	}
}

func AllowSrcIPPortInterface(sourceIP, targetIP, targetPort, vulInterface string) {
	cmd := exec.Command("ufw", "allow", "in", "on", vulInterface, "from", sourceIP, "to", targetIP, "port", targetPort)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Configuration Error: ", err)
	}
}

func AllowSrcIPPort(sourceIP, targetIP, targetPort string) {
	cmd := exec.Command("ufw", "allow", "in", "from", sourceIP, "to", targetIP, "port", targetPort)
	err := cmd.Run()
	if err != nil {
		log.Error("firewall: Configuration Error: ", err)
	}
}
