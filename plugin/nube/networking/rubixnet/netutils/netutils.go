package netutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	fmt.Println(getGatewayIP("enp9s0"))
	fmt.Println(HasStaticIP("enp9s0", true))
	fmt.Println(SetStaticIP("enp9s0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))
	//fmt.Println(HasStaticIP("enp9s0", true))
	//fmt.Println(GetCurrentHardwarePortInfo("enp9s0"))
	//removeLine("/etc/dhcpcd.conf", 62)
	//updateStaticIPDhcpcdConf("eth0", "192.168.15.109", "192.168.15.1", "8.8.8.8")

	//func updateStaticIPDhcpcdConf(iFaceName, ip, gatewayIP, dnsIP string) string {
}

func RunCommand(command string, arguments ...string) (int, string, error) {
	cmd := exec.Command(command, arguments...)
	out, err := cmd.Output()
	if err != nil {
		return 1, "", fmt.Errorf("exec.Command(%s) failed: %v: %s", command, err, string(out))
	}
	return cmd.ProcessState.ExitCode(), string(out), nil
}

type NetInterface struct {
	Name         string   // Network interface name
	MTU          int      // MTU
	HardwareAddr string   // Hardware address
	Addresses    []string // Array with the network interface addresses
	Subnets      []string // Array with CIDR addresses of this network interface
	Flags        string   // Network interface flags (up, broadcast, etc)
}

func GetValidNetInterfaces() ([]net.Interface, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("couldn't get list of interfaces: %s", err)
	}
	netIfaces := []net.Interface{}
	for i := range iFaces {
		iface := iFaces[i]
		netIfaces = append(netIfaces, iface)
	}
	return netIfaces, nil
}

func removeLine(path string, lineNumber int) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	info, _ := os.Stat(path)
	mode := info.Mode()
	array := strings.Split(string(file), "\n")
	array = append(array[:lineNumber], array[lineNumber+1:]...)
	err = ioutil.WriteFile(path, []byte(strings.Join(array, "\n")), mode)
	return err
}

func GetValidNetInterfacesForWeb() ([]NetInterface, error) {
	ifaces, err := GetValidNetInterfaces()
	if err != nil {
		return nil, errors.New("couldn't get interfaces")
	}
	if len(ifaces) == 0 {
		return nil, errors.New("couldn't find any legible interface")
	}
	var netInterfaces []NetInterface
	for _, iface := range ifaces {
		addrs, e := iface.Addrs()
		if e != nil {
			return nil, errors.New("failed to get addresses for interface")
		}
		netIface := NetInterface{
			Name:         iface.Name,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr.String(),
		}
		if iface.Flags != 0 {
			netIface.Flags = iface.Flags.String()
		}
		// Collect network interface addresses
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				// not an IPNet, should not happen
				return nil, fmt.Errorf("got iface.Addrs() element %s that is not net.IPNet, it is %T", addr, addr)
			}
			// ignore link-local
			if ipNet.IP.IsLinkLocalUnicast() {
				continue
			}
			netIface.Addresses = append(netIface.Addresses, ipNet.IP.String())
			netIface.Subnets = append(netIface.Subnets, ipNet.String())
		}
		// Discard interfaces with no addresses
		if len(netIface.Addresses) != 0 {
			netInterfaces = append(netInterfaces, netIface)
		}
	}
	return netInterfaces, nil
}

func GetSubnet(iFaceName string) string {
	netIfaces, err := GetValidNetInterfacesForWeb()
	if err != nil {
		fmt.Printf("Could not get network interfaces info: %v", err)
		return ""
	}
	for _, netIface := range netIfaces {
		if netIface.Name == iFaceName && len(netIface.Subnets) > 0 {
			return netIface.Subnets[0]
		}
	}
	return ""
}

// HasStaticIP Check if network interface has a static IP configured
// Supports: Raspbian.
func HasStaticIP(iFaceName string, delete bool) (bool, error) {
	if runtime.GOOS == "linux" {
		body, err := ioutil.ReadFile("/etc/dhcpcd.conf")
		if err != nil {
			return false, err
		}
		return hasStaticIPDhcpcdConf(string(body), iFaceName, delete), nil
	}
	return false, fmt.Errorf("cannot check if IP is static: not supported on %s", runtime.GOOS)
}

//SetStaticIP Set a static IP for the specified network interface
func SetStaticIP(iFaceName, ip, gatewayIP, dnsIP string) error {
	if runtime.GOOS == "linux" {
		return setStaticIPDhcpdConf(iFaceName, ip, gatewayIP, dnsIP)
	}
	return fmt.Errorf("cannot set static IP on %s", runtime.GOOS)
}

// GetCurrentHardwarePortInfo gets information the specified network interface
func GetCurrentHardwarePortInfo(iFaceName string) (HardwarePortInfo, error) {

	m := getNetworkSetupHardwareReports()
	hardwarePort, ok := m[iFaceName]
	if !ok {
		return HardwarePortInfo{}, fmt.Errorf("could not find hardware port for %s", iFaceName)
	}

	return getHardwarePortInfo(hardwarePort)
}

// for dhcpcd.conf
func hasStaticIPDhcpcdConf(dhcpConf, iFaceName string, delete bool) bool {
	lines := strings.Split(dhcpConf, "\n")
	nameLine := fmt.Sprintf("interface %s", iFaceName)
	withinInterfaceCtx := false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if withinInterfaceCtx && len(line) == 0 {
			// an empty line resets our state
			withinInterfaceCtx = false
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		if !withinInterfaceCtx {
			if line == nameLine {
				// we found our interface
				withinInterfaceCtx = true
				if delete {
					for ii := 0; ii < 4; ii++ {
						fmt.Println("line number ", i, line, ii)
						err := removeLine("/etc/dhcpcd.conf", i)
						if err != nil {
							return false
						}
					}
				}
			}
		} else {
			if strings.HasPrefix(line, "interface ") {
				// we found another interface - reset our state
				withinInterfaceCtx = false
				continue
			}
			if strings.HasPrefix(line, "static ip_address=") {
				return true
			}
		}
	}
	return false
}

// Get gateway IP address
func getGatewayIP(iFaceName string) string {
	cmd := exec.Command("ip", "route", "show", "dev", iFaceName)
	d, err := cmd.Output()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		return ""
	}
	fields := strings.Fields(string(d))
	if len(fields) < 3 || fields[0] != "default" {
		return ""
	}
	ip := net.ParseIP(fields[2])
	if ip == nil {
		return ""
	}
	return fields[2]
}

// setStaticIPDhcpdConf - updates /etc/dhcpd.conf and sets the current IP address to be static
func setStaticIPDhcpdConf(iFaceName, ip, gatewayIP, dnsIP string) error {
	fmt.Println(iFaceName, ip, gatewayIP, dnsIP)
	_ip := ip
	if ip == "" {
		_ip = GetSubnet(iFaceName)
	}
	if len(_ip) == 0 {
		return errors.New("can't get IP address")
	}
	ip4, _, err := net.ParseCIDR(_ip)
	if err != nil {
		return err
	}
	gateIP := gatewayIP
	if gatewayIP == "" {
		gateIP = getGatewayIP(iFaceName)
	}
	_dnsIP := dnsIP
	if dnsIP == "" {
		_dnsIP = ip4.String()
	}
	fmt.Println(iFaceName, ip, _ip, gateIP, _dnsIP)
	add := updateStaticIPDhcpcdConf(iFaceName, _ip, gateIP, _dnsIP)
	body, err := ioutil.ReadFile("/etc/dhcpcd.conf")
	if err != nil {
		return err
	}
	body = append(body, []byte(add)...)
	err = os.WriteFile("/etc/dhcpcd.conf", body, 0755)
	if err != nil {
		return err
	}
	return nil
}

// updates dhcpd.conf content -- sets static IP address there
// for dhcpcd.conf
func updateStaticIPDhcpcdConf(iFaceName, ip, gatewayIP, dnsIP string) string {
	var body []byte
	add := fmt.Sprintf("\ninterface %s\nstatic ip_address=%s\n",
		iFaceName, ip)
	body = append(body, []byte(add)...)

	if len(gatewayIP) != 0 {
		add = fmt.Sprintf("static routers=%s\n",
			gatewayIP)
		body = append(body, []byte(add)...)
	}
	add = fmt.Sprintf("static domain_name_servers=%s\n\n",
		dnsIP)
	body = append(body, []byte(add)...)
	return string(body)
}

// getNetworkSetupHardwareReports parses the output of the `networksetup -listallhardwareports` command
// it returns a map where the key is the interface name, and the value is the "hardware port"
// returns nil if it fails to parse the output
func getNetworkSetupHardwareReports() map[string]string {
	_, out, err := RunCommand("networksetup", "-listallhardwareports")
	if err != nil {
		return nil
	}

	re, err := regexp.Compile("Hardware Port: (.*?)\nDevice: (.*?)\n")
	if err != nil {
		return nil
	}

	m := make(map[string]string, 0)

	matches := re.FindAllStringSubmatch(out, -1)
	for i := range matches {
		port := matches[i][1]
		device := matches[i][2]
		m[device] = port
	}

	return m
}

// HardwarePortInfo - information obtained using MacOS networksetup
// about the current state of the internet connection
type HardwarePortInfo struct {
	name      string
	ip        string
	subnet    string
	gatewayIP string
	static    bool
}

func getHardwarePortInfo(hardwarePort string) (HardwarePortInfo, error) {
	h := HardwarePortInfo{}

	_, out, err := RunCommand("networksetup", "-getinfo", hardwarePort)
	if err != nil {
		return h, err
	}

	re := regexp.MustCompile("IP address: (.*?)\nSubnet mask: (.*?)\nRouter: (.*?)\n")

	match := re.FindStringSubmatch(out)
	if len(match) == 0 {
		return h, errors.New("could not find hardware port info")
	}

	h.name = hardwarePort
	h.ip = match[1]
	h.subnet = match[2]
	h.gatewayIP = match[3]

	if strings.Index(out, "Manual Configuration") == 0 {
		h.static = true
	}

	return h, nil
}

// Gets a list of nameservers currently configured in the /etc/resolv.conf
func getEtcResolvConfServers() ([]string, error) {
	body, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("nameserver ([a-zA-Z0-9.:]+)")

	matches := re.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return nil, errors.New("found no DNS servers in /etc/resolv.conf")
	}

	addrs := make([]string, 0)
	for i := range matches {
		addrs = append(addrs, matches[i][1])
	}

	return addrs, nil
}
