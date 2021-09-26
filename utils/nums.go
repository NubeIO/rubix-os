package utils

import (
	"math"
	"math/rand"
	"time"
)

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

//RandInt with min and max range
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

//RandFloat with min and max range
func RandFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}
