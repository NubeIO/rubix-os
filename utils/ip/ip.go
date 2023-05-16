package ip

import (
	"fmt"
	"net/url"
)

func Builder(https *bool, ip string, port int) (*url.URL, error) {
	if https != nil && *https == true {
		return url.ParseRequestURI(fmt.Sprintf("https://%s:%d", ip, port))
	} else {
		return url.ParseRequestURI(fmt.Sprintf("http://%s:%d", ip, port))
	}
}
