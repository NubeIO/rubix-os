package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"strings"
	"unicode"
)

type URLParts struct {
	transport string
	host      string
	port      string
}

func SplitURL(url string) URLParts {
	var o URLParts
	u := strings.SplitN(url, "://", 2)
	host := ""
	if len(u) == 2 {
		o.transport = u[0]
		host = u[1]
	}
	p := strings.Split(host, ":")
	o.host = p[0]
	o.port = p[1]
	return o
}

func JoinURL(u URLParts) (url string) {
	return fmt.Sprintf("%s://%s:%s", u.transport, u.host, u.port)
}
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
func main() {

	a := SplitURL("tcp://192.168.15.202:502")
	fmt.Println(a.port, a.host, a.transport)

	aaa := JoinURL(a)
	fmt.Println(aaa)

	str := utils.NewString("this_is_a test")

	res := str.LcFirstLetter()
	fmt.Println(UcFirst(res))
	fmt.Println(LcFirst("this_is_a test"))

}
