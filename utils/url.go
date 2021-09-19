package utils

import (
	"fmt"
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

func JoinURL(u URLParts) (url string) {
	return fmt.Sprintf("%s://%s:%s", u.Transport, u.Host, u.Port)
}
