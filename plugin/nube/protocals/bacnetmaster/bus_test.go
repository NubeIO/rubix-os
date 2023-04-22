package main

import (
	"fmt"
	"testing"
)

func Test_getReadName(t *testing.T) {

	topic := "bacnet/cmd_result/read_value/ao/1/name"

	readType, ioType, ioNumber, _ := getReadType(topic)
	fmt.Println(readType, ioType, ioNumber)
}
