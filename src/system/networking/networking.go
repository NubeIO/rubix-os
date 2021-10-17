package networking

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/src/system/command"
	"github.com/NubeDev/flow-framework/utils"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type NetInterface struct {
	Name         string   // Network interface name
	MTU          int      // MTU
	HardwareAddr string   // Hardware address
	Addresses    []string // Array with the network interface addresses
	Subnets      []string // Array with CIDR addresses of this network interface
	Flags        string   // Network interface flags (up, broadcast, etc)
}

type NetworkInterfaces struct {
	Interface     string `json:"interface"`
	IP            string `json:"ip"`
	IPMask        string `json:"ip_and_mask"`
	NetMask       string `json:"netmask"`
	NetMaskLength string `json:"net_mask_length"`
	Gateway       string `json:"gateway"`
}

type CheckInternet struct {
	Interface     string `json:"interface"`
	Message       string `json:"message"`
	FoundInternet bool   `json:"found_internet"`
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

func GetInterfacesNames() (*utils.Array, error) {
	i, err := GetValidNetInterfaces()
	if err != nil {
		return nil, errors.New("couldn't get interfaces")
	}
	out := utils.NewArray()
	for _, n := range i {
		out.Add(n.Name)
	}
	return out, nil
}

func GetValidNetInterfaces() ([]net.Interface, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("couldn't get list of interfaces: %s", err)
	}
	var netIfaces []net.Interface
	for i := range iFaces {
		iface := iFaces[i]
		netIfaces = append(netIfaces, iface)
	}
	return netIfaces, nil
}

//GetGatewayIP Get gateway IP address
func GetGatewayIP(iFaceName string) string {
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

//CheckInternetConnection check internet connection for a port (will ping google.com 2 times)
func CheckInternetConnection(iface string) (msg string, connected bool, err error) {
	cmd := fmt.Sprintf("if ping -I %s -c 2 google.com; then echo OK; else echo DEAD ;fi", iface)
	ping, err := command.Run("bash", "-c", cmd)
	if err != nil {
		return "", false, err
	}
	if strings.Contains(ping, "OK") {
		return "PASS", true, nil
	} else if strings.Contains(ping, "unknown iface") {
		return "FAIL: failed to find network interface", false, err
	} else if strings.Contains(ping, "Name or service not known") {
		return "FAIL: Name or service not known", false, err
	}
	return "Unknown Fail", false, err
}

//CheckInternetStatus check internet connection for all ports (will ping google.com 2 times)
func CheckInternetStatus() (msg *utils.Array, err error) {
	_, ifaceNames, _, err := IpAddresses()
	if err != nil {
		return nil, err
	}
	var ci CheckInternet
	out := utils.NewArray()
	for _, iface := range ifaceNames {
		message, ok, err := CheckInternetConnection(iface)
		if err != nil {
			return nil, err
		}
		ci.Interface = iface
		ci.FoundInternet = ok
		ci.Message = message
		out.Add(ci)
	}
	return out, err
}

func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

// IpAddresses fetches IP addresses
func IpAddresses() (ip, ifaceNames []string, ipAndNames *utils.Array, err error) {
	var names []string
	var ips []string
	out := utils.NewArray()
	var interfaces NetworkInterfaces
	if ifaces, err := net.Interfaces(); err == nil {
		if err != nil {
			return nil, nil, nil, err
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}

					ip = ip.To4()
					if ip == nil {
						continue
					}
					mask := strings.Split(addr.String(), "/")
					interfaces.Interface = iface.Name
					interfaces.IP = ip.String()
					interfaces.IPMask = addr.String()
					if len(mask) >= 1 {
						interfaces.NetMaskLength = mask[1]
					}
					interfaces.NetMask = ipv4MaskString(ip.DefaultMask())
					interfaces.Gateway = GetGatewayIP(iface.Name)
					out.Add(interfaces)
					names = append(names, iface.Name)
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return ips, names, out, nil
}

// ExternalIPV4 fetches external IP address in ipv4 format
func ExternalIPV4() (ip string, err error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	out := strings.TrimRight(string(content), "\r\n")
	return out, nil
}

// ExternalIpAddress fetches external IP address
func ExternalIpAddress() (ip string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	// http get request
	var req *http.Request
	if req, err = http.NewRequest("GET", "https://domains.google.com/checkip", nil); err == nil {
		// user-agent
		req.Header.Set("User-Agent", fmt.Sprintf("rpi-tools (golang; %s; %s)", runtime.GOOS, runtime.GOARCH))
		// http get
		var resp *http.Response
		resp, err = client.Do(req)

		if resp != nil {
			defer resp.Body.Close() // in case of http redirects
		}
		if err == nil && resp.StatusCode == 200 {
			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				ip := strings.TrimSpace(string(body))
				return ip, nil
			}
			err = fmt.Errorf("failed to read external ip: %s", err)
		} else {
			err = fmt.Errorf("failed to fetch external ip: %s (http %d)", err, resp.StatusCode)
		}
	}
	return "0.0.0.0", err
}
