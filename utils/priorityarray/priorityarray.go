package priorityarray

import (
	"errors"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func ConvertToMap(priority model.Priority) map[string]*float64 {
	priorityMap := map[string]*float64{}
	priorityValue := reflect.ValueOf(priority)
	typeOfPriority := priorityValue.Type()
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float64)
			priorityMap[key] = val
		}
	}
	return priorityMap
}

func ApplyMapToPriorityArray(pnt *model.Point, updatedPriorityMap *map[string]*float64) (*model.Point, error) {
	if updatedPriorityMap != nil && pnt.Priority != nil {
		for key, val := range *updatedPriorityMap {
			newVal := float.Copy(val)

			switch key {
			case "_1":
				pnt.Priority.P1 = newVal
			case "_2":
				pnt.Priority.P2 = newVal
			case "_3":
				pnt.Priority.P3 = newVal
			case "_4":
				pnt.Priority.P4 = newVal
			case "_5":
				pnt.Priority.P5 = newVal
			case "_6":
				pnt.Priority.P6 = newVal
			case "_7":
				pnt.Priority.P7 = newVal
			case "_8":
				pnt.Priority.P8 = newVal
			case "_9":
				pnt.Priority.P9 = newVal
			case "_10":
				pnt.Priority.P10 = newVal
			case "_11":
				pnt.Priority.P11 = newVal
			case "_12":
				pnt.Priority.P12 = newVal
			case "_13":
				pnt.Priority.P13 = newVal
			case "_14":
				pnt.Priority.P14 = newVal
			case "_15":
				pnt.Priority.P15 = newVal
			case "_16":
				pnt.Priority.P16 = newVal
			}
		}
		return pnt, nil
	}
	return pnt, errors.New("Invalid priority map. not applied to point.")
}

func ParsePriority(originalPointPriority *model.Priority, newPriorityMapPtr *map[string]*float64, isTypeBool bool) (
	*map[string]*float64, *float64, *int, bool, bool) {
	resultPriorityMap := map[string]*float64{}
	priorityValue := reflect.ValueOf(*originalPointPriority)
	typeOfPriority := priorityValue.Type()
	var highestValue *float64 = nil
	var currentPriority *int = nil
	doesPriorityExist := false
	isPriorityChanged := false
	isFirst := false
	newPriorityMap := map[string]*float64{}
	if newPriorityMapPtr != nil {
		newPriorityMap = *newPriorityMapPtr
	}
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float64)

			v, ok := newPriorityMap[key]
			// isPriorityChanged calculation:
			// point.priority array and body.priority array values are compared with body nullability check
			// For example:
			//   false => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"_15": 10, "16": 11}`
			//   false => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"16": 11}`
			//   true => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"_15": null, "16": 11}`
			if ok {
				doesPriorityExist = true
				val = v
				if !float.ComparePtrValues(v, val) {
					isPriorityChanged = true
				}
			}
			if val == nil {
				resultPriorityMap[key] = nil
			} else {
				if isTypeBool {
					val = float.EvalAsBoolOnlyOneIsTrue(val)
				}
				if !isFirst {
					currentPriority = integer.New(i) // can't assign address of i, coz it's referencing looping
					highestValue = val
				}
				resultPriorityMap[key] = val
				isFirst = true
			}
		}
	}
	return &resultPriorityMap, highestValue, currentPriority, doesPriorityExist, isPriorityChanged
}
