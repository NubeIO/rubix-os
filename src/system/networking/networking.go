package networking

import (
	"bytes"
	"errors"
	"fmt"
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

// IpAddresses fetches IP addresses
func IpAddresses() ([]string, error) {
	var ips []string
	if ifaces, err := net.Interfaces(); err == nil {
		if err != nil {
			return nil, err
		}
		for _, iface := range ifaces {
			// skip
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

					ips = append(ips, ip.String())
				}
			}
		}
	}
	return ips, nil
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
	return string(content), nil
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
