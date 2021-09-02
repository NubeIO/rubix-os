package utils

import "math"

func cov(new float64, existingData float64, cov float64) (bool, float64) {
	c := new - existingData
	if math.Abs(c) >= cov {
		return false, existingData
	} else {
		return false, new
	}
}
