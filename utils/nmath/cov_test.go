package nmath

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	event, _ := Cov(2, 0, 1)
	fmt.Println(event)
	event, _ = Cov(0, 0, 1)
	fmt.Println(event)
	event, _ = Cov(10, 9.9, 1)
	fmt.Println(event)
	event, _ = Cov(-2, -1, 1)
	fmt.Println(event)
}
