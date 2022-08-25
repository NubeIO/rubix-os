package integer

import (
	"math/rand"
	"time"
)

func New(value int) *int {
	return &value
}

func NewUint(value uint) *uint {
	return &value
}

func NonNil(b *int) int {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func IsNil(b *int) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsUnitNil(b *uint) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsUnit32Nil(b *uint32) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

// RandInt returns a random int within the specified range.
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
