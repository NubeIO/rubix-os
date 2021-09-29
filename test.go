package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
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
	fmt.Println(unit.Exists("length1"))

	out := utils.NewArray()
	objs := utils.ArrayValues(model.ObjectTypes)
	for _, obj := range objs {
		switch obj {
		case model.ObjectTypes.AnalogInput:
			out.Add(obj)
		case model.ObjectTypes.AnalogOutput:
		case model.ObjectTypes.AnalogValue:
		case model.ObjectTypes.BinaryInput:
		case model.ObjectTypes.BinaryOutput:
		case model.ObjectTypes.BinaryValue:
		default:
			//out.Add(obj)
		}

	}
	fmt.Println(out)

}
