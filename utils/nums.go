package utils

import (
	"math"
	"math/rand"
	"time"
)

func Float64IsNil(b *float64) float64 {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func NewInt(value int) *int {
	return &value
}

func NewFloat64(value float64) *float64 {
	return &value
}

func IntIsNil(b *int) int {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func IntNilCheck(b *int) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func Float32IsNil(b *float32) float32 {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func Unit16IsNil(b *uint16) uint16 {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func Unit32IsNil(b *uint32) uint32 {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func Unit32NilCheck(b *uint32) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func FloatIsNilCheck(b *float64) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

//LimitToRange returns the input value clamped within the specified range
func LimitToRange(value float64, range1 float64, range2 float64) float64 {
	if range1 == range2 {
		return range1
	}
	var min, max float64
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	return math.Min(math.Max(value, min), max)
}

//Scale returns the (float64) input value (between inputMin and inputMax) scaled to a value between outputMin and outputMax
func Scale(value float64, inMin float64, inMax float64, outMin float64, outMax float64) float64 {
	scaled := ((value-inMin)/(inMax-inMin))*(outMax-outMin) + outMin
	if scaled > math.Max(outMin, outMax) {
		return math.Max(outMin, outMax)
	} else if scaled < math.Min(outMin, outMax) {
		return math.Min(outMin, outMax)
	} else {
		return scaled
	}
}

//RoundTo returns the input value rounded to the specified number of decimal places.
func RoundTo(value float64, decimals uint32) float64 {
	if decimals < 0 {
		return value
	}
	return math.Round(value*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

//RandInt returns a random int within the specified range.
func RandInt(range1, range2 int) int {
	if range1 == range2 {
		return range1
	}
	var min, max int
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

//RandFloat returns a random float64 within the specified range.
func RandFloat(range1, range2 float64) float64 {
	if range1 == range2 {
		return range1
	}
	var min, max float64
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}
