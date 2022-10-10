package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NubeIO/flow-framework/utils/nuuid"
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
