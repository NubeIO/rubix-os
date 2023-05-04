package main

import (
	"fmt"
	"time"
)

func newUUID(len int) string {
	t := fmt.Sprint(time.Now().Nanosecond())
	if len < 3 {
		len = 3
	}
	if len >= 6 {
		len = 6
	}
	return t[:len]
}
