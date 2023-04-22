package main

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes/priority"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *Instance) doRead(point *model.Point, deviceUUID, networkUUID string) (currentBACServPriority *priority.Float32, highestPriorityValue *float64, readSuccess, writeSuccess bool, err error) {

	fmt.Println(111111, deviceUUID, networkUUID)

	currentBACServPriority = &priority.Float32{
		P1:  nil,
		P2:  nil,
		P3:  nil,
		P4:  nil,
		P5:  nil,
		P6:  nil,
		P7:  nil,
		P8:  nil,
		P9:  nil,
		P10: nil,
		P11: nil,
		P12: nil,
		P13: nil,
		P14: nil,
		P15: nil,
		P16: nil,
	}

	return nil, nil, false, false, err

}
