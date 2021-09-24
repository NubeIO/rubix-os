package unit

import (
	"fmt"
	"math"

	"github.com/martinlindhe/unit"
)

// LengthUnit is a unit of length
type LengthUnit struct {
	unitCommon
	length unit.Length
	to     func(unit.Length) float64
}

// FromFloat implements SimpleUnit
func (lu *LengthUnit) FromFloat(f float64) UnitVal {
	return LengthVal{unit.Length(f) * lu.length, lu}
}

// Length units
var (
	Meter        = &LengthUnit{"m", unit.Meter, unit.Length.Meters}
	Kilometer    = &LengthUnit{"km", unit.Kilometer, unit.Length.Kilometers}
	Millimeter   = &LengthUnit{"mm", unit.Millimeter, unit.Length.Millimeters}
	Centimeter   = &LengthUnit{"cm", unit.Centimeter, unit.Length.Centimeters}
	Nanometer    = &LengthUnit{"nm", unit.Nanometer, unit.Length.Nanometers}
	Inch         = &LengthUnit{"in", unit.Inch, unit.Length.Inches}
	Foot         = &LengthUnit{"ft", unit.Foot, unit.Length.Feet}
	Yard         = &LengthUnit{"yd", unit.Yard, unit.Length.Yards}
	Mile         = &LengthUnit{"miles", unit.Mile, unit.Length.Miles}
	Furlong      = &LengthUnit{"furlongs", unit.Furlong, unit.Length.Furlongs}
	Lightyear    = &LengthUnit{"ly", unit.LightYear, unit.Length.LightYears}
	NauticalMile = &LengthUnit{"nautical mile", unit.NauticalMile, unit.Length.NauticalMiles}

	FootInch = &FootInchUnit{"feet+inches"}
)

// LengthVal is a length value with unit
type LengthVal struct {
	V unit.Length
	U *LengthUnit
}

func (lv LengthVal) String() string {
	return simpleUnitString(lv.U.to(lv.V), lv.U)
}

func (lv LengthVal) AsFloat() float64 {
	return lv.U.to(lv.V)
}

// Convert implements UnitVal conversion
func (lv LengthVal) Convert(to UnitType) (UnitVal, error) {
	switch to := to.(type) {
	case *LengthUnit:
		lv.U = to
		return lv, nil
	case *FootInchUnit:
		feet, fraction := math.Modf(lv.V.Feet())
		inches := (unit.Length(fraction) * unit.Foot).Inches()
		return FootInchVal{feet, inches}, nil
	default:
		return nil, ErrorConversion{lv.U, to}
	}
}

// FootInchUnit is a unit of both feet + inches
type FootInchUnit struct{ unitCommon }

// FromFloat implements UnitType
func (FootInchUnit) FromFloat(f float64) UnitVal {
	return FootInchVal{Feet: f, Inches: 0}
}

// FootInchVal is a UnitVal of both feet + inches
type FootInchVal struct {
	Feet, Inches float64
}

func (val FootInchVal) String() string {
	if val.Inches == 0 {
		return simpleUnitString(val.Feet, Foot)
	}
	return fmt.Sprintf(`%.f' %.f"`, val.Feet, val.Inches)
}

func (val FootInchVal) AsFloat() float64 {
	return val.Feet
}

// Convert implements UnitVal conversion
func (val FootInchVal) Convert(to UnitType) (UnitVal, error) {
	switch to := to.(type) {
	case *FootInchUnit:
		return val, nil
	case *LengthUnit:
		feet := unit.Length(val.Feet) * unit.Foot
		inches := unit.Length(val.Inches) * unit.Inch
		return LengthVal{feet + inches, to}, nil
	default:
		return nil, ErrorConversion{FootInch, to}
	}
}
