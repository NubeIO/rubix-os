package main

import (
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
)

//delete point make sure
func (i *Instance) serverDeletePoint(body *pkgmodel.BacnetPoint) (bool, error) {
	return true, nil
}
