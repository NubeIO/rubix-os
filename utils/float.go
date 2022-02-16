package utils

import "math"

func FirstNotNilFloat(values ...*float64) *float64 {
	for _, n := range values {
		if n != nil {
			return n
		}
	}
	return nil
}

func CompareFloatPtr(value1 *float64, value2 *float64) bool {
	return FloatPtrToFloat(value1) == FloatPtrToFloat(value2)
}

func FloatPtrToFloat(value *float64) float64 {
	if value == nil {
		return math.MaxFloat64
	}
	return *value
}
