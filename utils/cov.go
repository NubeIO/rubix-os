package utils

import "math"


func COV(new float64, existingData float64, threshold float64) (bool, float64) {
	if math.Abs(new - existingData) >= threshold {
		return true, existingData
	} else {
		return false, new
	}
}
