package database

import (
	"context"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/src/units"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/integer"
	"github.com/NubeIO/rubix-os/utils/nmath"
	"github.com/PaesslerAG/gval"
)

func (d *GormDatabase) pointNameExists(body *model.Point) (nameExist, addressIDExist bool) {
	var arg argspkg.Args
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
	device, err := d.GetDevice(deviceUUID, argspkg.Args{WithPoints: true})
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
		return nil, err
	}
	return float.New(res.AsFloat()), nil
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
		if *scaleInMin == 0 && *scaleInMax == 0 {
			out := nmath.LimitToRange(*presentValue, *scaleOutMin, *scaleOutMax)
			return &out
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

func PointValueTransformOnRead(originalValue *float64, scaleEnable *bool, factor, scaleInMin, scaleInMax, scaleOutMin,
	scaleOutMax, offset *float64) (transformedValue *float64) {
	if originalValue == nil {
		return nil
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
		if float.NonNil(scaleOutMin) != 0 || float.NonNil(scaleOutMax) != 0 {
			if float.NonNil(scaleInMin) != 0 || float.NonNil(scaleInMax) != 0 { // scale with all 4 configs
				scaledAndLimited = float.NonNil(PointScale(float.New(factored), scaleInMin, scaleInMax, scaleOutMin,
					scaleOutMax))
			} else { // do limit with only scaleOutMin and scaleOutMin
				scaledAndLimited = float.NonNil(PointRange(float.New(factored), scaleOutMin, scaleOutMax))
			}
		}
	}
	// perform offset operation
	offsetted := scaledAndLimited + float.NonNil(offset)
	return float.New(offsetted)
}

func PointValueTransformOnWrite(originalValue *float64, scaleEnable *bool, factor, scaleInMin, scaleInMax, scaleOutMin,
	scaleOutMax, offset *float64) (transformedValue *float64) {
	if originalValue == nil {
		return nil
	}
	ov := float.NonNil(originalValue)

	// reverse offset operation
	unoffsetted := ov - float.NonNil(offset)

	// reverse scaling and limit operations
	unscaledAndUnlimited := unoffsetted
	if boolean.IsTrue(scaleEnable) {
		if float.NonNil(scaleOutMin) != 0 || float.NonNil(scaleOutMax) != 0 {
			if float.NonNil(scaleInMin) != 0 || float.NonNil(scaleInMax) != 0 { // scale with all 4 configs
				unscaledAndUnlimited = float.NonNil(PointScale(float.New(unoffsetted), scaleOutMin, scaleOutMax,
					scaleInMin, scaleInMax))
			} else { // do limit with only scaleOutMin and scaleOutMin
				unscaledAndUnlimited = float.NonNil(PointRange(float.New(unoffsetted), scaleOutMin, scaleOutMax))
			}
		}
	}

	// reverse factoring operation
	unfactored := unscaledAndUnlimited
	if float.NonNil(factor) != 0 {
		unfactored = unscaledAndUnlimited / float.NonNil(factor)
	}

	return float.New(unfactored)
}
