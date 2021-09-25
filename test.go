package main

import (
	"fmt"
	unit "github.com/NubeDev/flow-framework/src/units"
)

func main() {
	_, res := unit.Process(5000, "meter", "kilometer")
	fmt.Println(res.String())
	fmt.Println(res.AsFloat())

}
