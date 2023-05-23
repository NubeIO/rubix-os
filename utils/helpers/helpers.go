package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/gin-gonic/gin"
)

func truncateString(str string, num int) string {
	ret := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		ret = str[0:num] + ""
	}
	return ret
}

func NameIsNil() string {
	uuid := nuuid.MakeTopicUUID("")
	return fmt.Sprintf("n_%s", truncateString(uuid, 8))
}

func CloneRequest(ctx *gin.Context) *http.Request {
	r := ctx.Request
	r2 := r.Clone(ctx)
	*r2 = *r
	var b bytes.Buffer
	b.ReadFrom(r.Body)
	r.Body = ioutil.NopCloser(&b)
	r2.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	return r2
}

func CheckVersionBool(version string) bool {
	var hasV bool
	var correctLen bool
	if version[0] == 'v' { // make sure have a v at the start v0.1.1
		hasV = true
	}
	p := strings.Split(version, ".")
	if len(p) == 3 {
		correctLen = true
	}
	if hasV && correctLen {
		return true
	}
	return false
}

func CheckVersion(version string) error {
	if version[0:1] != "v" { // make sure have a v at the start v0.1.1
		return errors.New(fmt.Sprintf("incorrect provided: %s version number try: v1.2.3", version))
	}
	p := strings.Split(version, ".")
	if len(p) != 3 {
		return errors.New(fmt.Sprintf("incorrect length provided: %s version number try: v1.2.3", version))
	}
	return nil
}
