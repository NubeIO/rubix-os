package units

import (
	"fmt"
	"strings"
	"sync"
)

var (
	unitMap  map[string]UnitType
	unitLock sync.RWMutex
)

func init() {
	unitMap = make(map[string]UnitType)
	updateUnitMap()
}

func updateUnitMap() {
	for unit, aliases := range supportedUnits {
		for _, alias := range aliases {
			alias = strings.ToLower(alias)
			if _, ok := unitMap[alias]; !ok {
				unitMap[alias] = unit
			}
		}
	}
}

// ParseUnit parses a UnitType.
// Lazily loads currency units
func ParseUnit(s string) (UnitType, bool) {
	s = strings.ToLower(s)
	unitLock.RLock()
	defer unitLock.RUnlock()
	u, ok := unitMap[s]
	if !ok {
		unitLock.RUnlock()
		unitLock.RLock()
		u, ok = unitMap[s]
	}
	return u, ok
}

// UnitType represent a single type of unit
type UnitType interface {
	fmt.Stringer
	FromFloat(float64) UnitVal
}

// UnitVal is a value with unit that can be converted to another unit
type UnitVal interface {
	fmt.Stringer
	Convert(to UnitType) (UnitVal, error)
	AsFloat() float64
}

// ErrorConversion occurs when a UnitType cannot be converted to another UnitType
type ErrorConversion struct {
	From, To UnitType
}

func (err ErrorConversion) Error() string {
	return fmt.Sprintf("Can't convert from %s to %s", err.From.String(), err.To.String())
}

type unitCommon string

func (c unitCommon) String() string {
	return string(c)
}

func simpleUnitString(f float64, u UnitType) string {
	return fmt.Sprintf("%.6g %s", f, u.String())
}
