package utils

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	event, _ := COV(2, 0, 1)
	fmt.Println(event)
	event, _ = COV(0, 0, 1)
	fmt.Println(event)
	event, _ = COV(10, 9.9, 1)
	fmt.Println(event)
	event, _ = COV(-2, -1, 1)
	fmt.Println(event)
}
