package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/units"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nmath"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/PaesslerAG/gval"
)

func (d *GormDatabase) pointNameExists(body *model.Point) (nameExist, addressIDExist bool) {
	var arg api.Args
	arg.WithPoints = true
	device, err := d.GetDevice(body.DeviceUUID, arg)
	if err != nil {
		return false, false
	}

	nameExist = false
	addressIDExist = false

	for _, pnt := range device.Points {
		if pnt.Name == body.Name {
			if pnt.UUID == body.UUID {
				nameExist = false
			} else {
				nameExist = true
			}
		}
	}

	// check body with existing devices that a point with the same objectType do not have the same addrID
	for _, pnt := range device.Points {
		if pnt.ObjectType == body.ObjectType {
			if integer.NonNil(pnt.AddressID) == integer.NonNil(body.AddressID) {
				if pnt.UUID == body.UUID {
					addressIDExist = false
				} else {
					addressIDExist = true
				}
			}
		}
	}
	return nameExist, addressIDExist
}

func (d *GormDatabase) pointNameExistsInDevice(pointName, deviceUUID string) (existing bool) {
	device, err := d.GetDevice(deviceUUID, api.Args{WithPoints: true})
	if err != nil {
		return false
	}
	for _, pnt := range device.Points {
		if pnt.Name == pointName {
			return true
		}
	}
	return false
}

// PointDeviceByAddressID will query by device_uuid = ? AND object_type = ? AND address_id = ?
func (d *GormDatabase) PointDeviceByAddressID(pointUUID string, body *model.Point) (*model.Point, bool) {
	var pointModel *model.Point
	deviceUUID := body.DeviceUUID
	objType := body.ObjectType
	addressID := body.AddressID
	f := fmt.Sprintf("device_uuid = ? AND object_type = ? AND address_id = ?")
	query := d.DB.Where(f, deviceUUID, objType, addressID, pointUUID).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false
	}
	return pointModel, true
}

func pointUnits(presentValue *float64, unitFrom, unitTo *string) (value *float64, err error) {
	if presentValue == nil || unitFrom == nil || unitTo == nil || *unitFrom == "" || *unitTo == "" {
		return presentValue, nil
	}
	_, res, err := units.Process(*presentValue, *unitFrom, *unitTo)
	if err != nil {
		return float.New(0), err
	}
	return float.New(res.AsFloat()), err
}

func PointRange(presentValue, limitMin, limitMax *float64) (value *float64) {
	if !float.IsNil(presentValue) && !float.IsNil(limitMin) && !float.IsNil(limitMax) {
		if *limitMin == 0 && *limitMax == 0 {
			return presentValue
		}
		out := nmath.LimitToRange(*presentValue, *limitMin, *limitMax)
		return &out
	}
	return presentValue
}

func PointScale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax *float64) (value *float64) {
	if !float.IsNil(presentValue) && !float.IsNil(scaleInMin) && !float.IsNil(scaleInMax) &&
		!float.IsNil(scaleOutMin) && !float.IsNil(scaleOutMax) {
		if *scaleInMin == 0 && *scaleInMax == 0 && *scaleOutMin == 0 && *scaleOutMax == 0 {
			return presentValue
		}
		out := nmath.Scale(*presentValue, *scaleInMin, *scaleInMax, *scaleOutMin, *scaleOutMax)
		return &out
	}
	return presentValue
}

func PointEval(val *float64, evalString string) (value *float64, err error) {
	if evalString != "" {
		eval, err := gval.Full().NewEvaluable(evalString)
		if err != nil || val == nil {
			return nil, err
		}
		v, err := eval.EvalFloat64(context.Background(), map[string]interface{}{"x": *val})
		if err != nil {
			return nil, err
		}
		_v := v
		return &_v, nil
	}
	return val, nil
}

func PointValueTransformOnRead(originalValue *float64, scaleEnable, limitEnable *bool, factor, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax, offset *float64) (transformedValue *float64, err error) {
	if originalValue == nil {
		return nil, errors.New("input value is undefined")
	}
	ov := float.NonNil(originalValue)

	// perform factor operation
	factored := ov
	if float.NonNil(factor) != 0 {
		factored = ov * float.NonNil(factor)
	}

	// perform scaling and limit operations
	scaledAndLimited := factored
	if boolean.IsTrue(scaleEnable) {
		scaledAndLimited = float.NonNil(PointScale(float.New(factored), scaleInMin, scaleInMax, scaleOutMin, scaleOutMax))
	} else if boolean.IsTrue(limitEnable) {
		scaledAndLimited = float.NonNil(PointRange(float.New(factored), scaleOutMin, scaleOutMax))
	}

	// perform offset operation
	offsetted := scaledAndLimited + float.NonNil(offset)

	return float.New(offsetted), nil
}

func PointValueTransformOnWrite(originalValue *float64, scaleEnable, limitEnable *bool, factor, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax, offset *float64) (transformedValue *float64, err error) {
	if originalValue == nil {
		return nil, errors.New("input value is undefined")
	}
	ov := float.NonNil(originalValue)

	// reverse offset operation
	unoffsetted := ov - float.NonNil(offset)

	// reverse scaling and limit operations
	unscaledAndUnlimited := unoffsetted
	if boolean.IsTrue(scaleEnable) {
		unscaledAndUnlimited = float.NonNil(PointScale(float.New(unoffsetted), scaleOutMin, scaleOutMax, scaleInMin, scaleInMax))
	} else if boolean.IsTrue(limitEnable) {
		unscaledAndUnlimited = float.NonNil(PointRange(float.New(unoffsetted), scaleOutMin, scaleOutMax))
	}

	// reverse factoring operation
	unfactored := unscaledAndUnlimited
	if float.NonNil(factor) != 0 {
		unfactored = unscaledAndUnlimited * float.NonNil(factor)
	}

	return float.New(unfactored), nil
}
