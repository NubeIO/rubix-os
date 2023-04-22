package main

import (
	"fmt"
	"testing"
)

func Test_decodeMac(t *testing.T) {
	ip, port := decodeMac("192.168.15.10:47808")
	fmt.Println(ip, port)
}
