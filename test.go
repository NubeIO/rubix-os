package main

import (
	"fmt"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/NubeDev/flow-framework/utils"
)

func main() {
	_, res, err := unit.Process(1, "c", "c")
	if err != nil {
		return
	}
	fmt.Println(err)
	fmt.Println(res.String())
	fmt.Println(res.AsFloat())

	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandFloat(1, 1011))
	fmt.Println(utils.RandFloat(1, 1011))

}
