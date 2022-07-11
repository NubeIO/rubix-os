package priorityarray

import (
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
