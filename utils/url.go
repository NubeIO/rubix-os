package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type URLParts struct {
	Transport string //tcp
	Host      string
	Port      string
}

func SplitURL(url string) URLParts {
	var o URLParts
	u := strings.SplitN(url, "://", 2)
	host := ""
	if len(u) == 2 {
		o.Transport = u[0]
		host = u[1]
	}
	p := strings.Split(host, ":")
	o.Host = p[0]
	o.Port = p[1]
	return o
}

func IsTCP(target string) bool {
	if strings.HasPrefix(target, "tcp://") {
		return true
	} else {
		return false
	}
}

func JoinURL(u URLParts) (url string, err error) {
	t := u.Transport
	h := u.Host
	p := u.Port
	//if !IsTCP(t) {
	//	return "",  errors.New("in valid url prefix try ie: tcp://")
	//}
	if !ValidIP4(h) {
		return "", errors.New("in valid url try ie: 192.168.1.1")
	}
	if !ValidPort(p) {
		return "", errors.New("in valid url try ie: 8080 as a string")
	}
	ip := fmt.Sprintf("%s://%s:%s", t, h, p)
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
