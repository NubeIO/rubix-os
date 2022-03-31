package model

import (
	"testing"
)

func TestPrintPointValues(*testing.T) {

	var pnt Point
	value := 16.0
	var pri Priority
	pri.P16 = &value
	pnt.Priority = &pri

	pnt.PrintPointValues()

}
