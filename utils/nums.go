package utils

import "math"

func Round(n, unit float64) float64 {
	if unit <= 0 {
		return n
	}
	return math.Round(n/unit) * unit
}

func RoundTo(n float64, decimals uint32) float64 {
	if decimals <= 0 {
		return n
	}
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}
