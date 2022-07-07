package priorityarray

import (
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

func ParsePriority(pointPriority *model.Priority, priority *map[string]*float64) (*map[string]*float64, *float64, *int, bool) {
	priorityMap := map[string]*float64{}
	priorityValue := reflect.ValueOf(*pointPriority)
	typeOfPriority := priorityValue.Type()
	var highestValue *float64 = nil
	var currentPriority *int = nil
	doesPriorityExist := false
	isFirst := false
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float64)
			var v *float64
			ok := false
			if priority != nil { //TODO: This code looks like it could be optimized
				p := *priority
				v, ok = p[key]
				if ok {
					doesPriorityExist = true
					val = v
				}
			}
			if val == nil {
				priorityMap[key] = nil
			} else {
				if !isFirst {
					currentPriority = integer.New(i) // can't assign address of i, coz it's referencing looping
					highestValue = val
				}
				priorityMap[key] = val
				isFirst = true
			}
		}
	}
	return &priorityMap, highestValue, currentPriority, doesPriorityExist
}
