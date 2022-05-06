package nurl

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parts struct {
	Host string
	Port string
}

func JoinIPPort(u Parts) (url string, err error) {
	h := u.Host
	p := u.Port
	if !ValidIP4(h) {
		return "", errors.New("in valid url try ie: 192.168.1.1")
	}
	if !ValidPort(p) {
		return "", errors.New("in valid url try ie: 8080 as a string")
	}
	ip := fmt.Sprintf("%s:%s", h, p)
	return ip, nil
}

func ValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func ValidPort(port string) bool {
	t := strings.Trim(port, " ")
	re, _ := regexp.Compile(`^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$`)
	if re.MatchString(t) {
		return true
	}
	return false
}

func PortAsString(port int) string {
	t := strconv.Itoa(port)
	return t
}
