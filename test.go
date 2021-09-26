package main

import (
	"fmt"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/NubeDev/flow-framework/utils"
)

func main() {
	_, res := unit.Process(5000, "meter", "kilometer")
	fmt.Println(res.String())
	fmt.Println(res.AsFloat())

	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandFloat(1, 1011))
	fmt.Println(utils.RandFloat(1, 1011))

}
