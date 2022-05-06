package main

import (
	"github.com/NubeIO/flow-framework/utils/array"
	"strings"
)

func getPointAddr(s string) (objType, addr string) {
	arr := array.NewArray()
	ss := strings.Split(s, "-")
	for _, e := range ss {
		if e != "" {
			arr.Add(e)
		}
	}
	o := arr.Get(0)
	a := arr.Get(1)
	return o.(string), a.(string)
}
