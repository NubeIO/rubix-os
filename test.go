package main

import (
	"fmt"
	unit "github.com/NubeDev/flow-framework/src/units"
)

func main() {

	aa, bb := unit.Process(5000, "s", "min")
	fmt.Println(aa)
	fmt.Println(unit.SupportedUnits2)
	fmt.Println(bb.String())
	fmt.Println(bb.AsFloat())

}
