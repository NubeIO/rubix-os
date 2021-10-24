package main

import (
	"fmt"
	"net/url"
	"strings"
)

func urlPath(u string) (clean string) {
	_url := fmt.Sprintf("http://%s", u)
	p, _ := url.Parse(_url)
	parts := strings.SplitAfter(p.String(), "any")
	if len(parts) >= 1 {
		return parts[1]
	} else {
		return ""
	}
}

func main() {
	input_url := "e/api/plugins/api/rubix/rubix/system-1/any/ff/api/points?with_priority=true"
	//u, err := url.Parse(input_url)
	//if err != nil {
	//	//log.Fatal(err)
	//}
	//
	//fmt.Println(u.Scheme)
	//fmt.Println(u.User)
	//fmt.Println(u.Hostname())
	//fmt.Println(u.Port())
	//fmt.Println(u.Path)
	//fmt.Println(u.RawQuery)
	//fmt.Println(u.Fragment)
	//fmt.Println(u.String())

	fmt.Println(urlPath(input_url))
}
