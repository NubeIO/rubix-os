package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)

func main() {

	fmt.Println(utils.Round(1.22, -1))
	fmt.Println(utils.RoundTo(1.221111, 1))

}
