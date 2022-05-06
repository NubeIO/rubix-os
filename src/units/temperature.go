package units

import (
	"github.com/martinlindhe/unit"
)

// TemperatureUnit is a unit of temperature
type TemperatureUnit struct {
	unitCommon
	from func(float64) unit.Temperature
	to   func(unit.Temperature) float64
}

// FromFloat implements SimpleUnit
func (tu *TemperatureUnit) FromFloat(f float64) UnitVal {
	return TemperatureVal{tu.from(f), tu}
}

// Temperature units
var (
	Celsius    = &TemperatureUnit{"°C", unit.FromCelsius, unit.Temperature.Celsius}
	Fahrenheit = &TemperatureUnit{"°F", unit.FromFahrenheit, unit.Temperature.Fahrenheit}
	Kelvin     = &TemperatureUnit{"K", unit.FromKelvin, unit.Temperature.Kelvin}
)

// TemperatureVal is a temperature value with unit
type TemperatureVal struct {
	V unit.Temperature
	U *TemperatureUnit
}

func (tv TemperatureVal) String() string {
	return simpleUnitString(tv.U.to(tv.V), tv.U)
}

func (tv TemperatureVal) AsFloat() float64 {
	return tv.U.to(tv.V)
}

// Convert implements UnitVal conversion
func (tv TemperatureVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*TemperatureUnit); ok {
		tv.U = to
		return tv, nil
	}
	return nil, ErrorConversion{tv.U, to}
}
