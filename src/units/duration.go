package unit

import "github.com/martinlindhe/unit"

// DurationUnit is a unit of time
type DurationUnit struct {
	unitCommon
	time unit.Duration
	to   func(unit.Duration) float64
}

// FromFloat implements SimpleUnit
func (du *DurationUnit) FromFloat(f float64) UnitVal {
	return DurationVal{unit.Duration(f) * du.time, du}
}

// Time units
var (
	Second = &DurationUnit{"s", unit.Second, unit.Duration.Seconds}
	Minute = &DurationUnit{"min", unit.Minute, unit.Duration.Minutes}
	Hour   = &DurationUnit{"hr", unit.Hour, unit.Duration.Hours}
	Day    = &DurationUnit{"days", unit.Day, unit.Duration.Days}
	Week   = &DurationUnit{"weeks", unit.Week, unit.Duration.Weeks}
	Month  = &DurationUnit{"months", unit.ThirtyDayMonth, unit.Duration.ThirtyDayMonths}
	Year   = &DurationUnit{"years", unit.JulianYear, unit.Duration.JulianYears}
)

// DurationVal is a time value with unit
type DurationVal struct {
	V unit.Duration
	U *DurationUnit
}

func (dv DurationVal) String() string {
	return simpleUnitString(dv.U.to(dv.V), dv.U)
}

func (dv DurationVal) AsFloat() float64 {
	return dv.U.to(dv.V)
}

// Convert implements UnitVal conversion
func (dv DurationVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*DurationUnit); ok {
		dv.U = to
		return dv, nil
	}
	return nil, ErrorConversion{dv.U, to}
}
