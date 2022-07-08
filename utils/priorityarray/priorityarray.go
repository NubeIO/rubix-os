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

func ParsePriority(pointPriority *model.Priority, priority *map[string]*float64, isTypeBool bool) (*map[string]*float64, *float64, *int, bool, bool) {
	priorityMap := map[string]*float64{}
	priorityValue := reflect.ValueOf(*pointPriority)
	typeOfPriority := priorityValue.Type()
	var highestValue *float64 = nil
	var currentPriority *int = nil
	doesPriorityExist := false
	isPriorityChanged := false
	isFirst := false
	priority_ := map[string]*float64{}
	if priority != nil {
		priority_ = *priority
	}
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float64)

			v, ok := priority_[key]
			// isPriorityChanged calculation:
			// point.priority array and body.priority array values are compared with body nullability check
			// For example:
			//   false => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"_15": 10, "16": 11}`
			//   false => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"16": 11}`
			//   true => point.priority: `{"_15": 10, "16": 11}` and body.priority: `{"_15": null, "16": 11}`
			if !float.ComparePtrValues(v, val) && ok {
				isPriorityChanged = true
			}
			if ok {
				doesPriorityExist = true
				val = v
			}

			if val == nil {
				priorityMap[key] = nil
			} else {
				if isTypeBool {
					val = float.EvalAsBool(val)
				}
				if !isFirst {
					currentPriority = integer.New(i) // can't assign address of i, coz it's referencing looping
					highestValue = val
				}
				priorityMap[key] = val
				isFirst = true
			}
		}
	}
	return &priorityMap, highestValue, currentPriority, doesPriorityExist, isPriorityChanged
}
