package float

import (
	"math"
	"math/rand"
	"time"
)

func FirstNotNil(values ...*float64) *float64 {
	for _, n := range values {
		if n != nil {
			return n
		}
	}
	return nil
}

func ComparePtrValues(value1 *float64, value2 *float64) bool {
	return NonNilMax(value1) == NonNilMax(value2)
}

func NonNilMax(value *float64) float64 {
	if value == nil {
		return math.MaxFloat64
	}
	return *value
}

func New(value float64) *float64 {
	return &value
}

func NonNil(b *float64) float64 {
	if b == nil {
		return 0
	}
	return *b
}

func NonNil32(b *float32) float32 {
	if b == nil {
		return 0
	}
	return *b
}

func IsNil(b *float64) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsNil32(b *float32) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func EvalAsBool(b *float64) *float64 {
	if b == nil {
		return New(0)
	} else if *b == 0 {
		return New(0)
	} else {
		return New(1)
	}
}

func EvalAsBoolOnlyOneIsTrue(b *float64) *float64 {
	if b == nil {
		return New(0)
	} else if *b == 1 {
		return New(1)
	} else {
		return New(0)
	}
}

func Copy(b *float64) *float64 {
	if b == nil {
		return nil
	}
	out := *b
	return &out
}

// RandFloat returns a random float64 within the specified range.
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
