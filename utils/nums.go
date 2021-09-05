package utils

import "math"

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func RoundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}