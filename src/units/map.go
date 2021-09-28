package unit

var supportedUnits = map[UnitType][]string{
	// Length
	Meter:        {"m", "meter", "meters"},
	Kilometer:    {"km", "kilometer", "kilometers"},
	Millimeter:   {"mm", "millimeter", "millimeters"},
	Centimeter:   {"cm", "centimeter", "centimeters"},
	Nanometer:    {"nm", "nanometer", "nanometers"},
	Inch:         {"in", "inch", "inches"},
	FootInch:     {"ft", "feet+inches", "ftin", "ft+in"},
	Foot:         {"foot", "feet"},
	Yard:         {"yd", "yard", "yards"},
	Mile:         {"mi", "mile", "miles"},
	Furlong:      {"furlong", "furlongs"},
	Lightyear:    {"ly", "lightyear", "lightyears"},
	NauticalMile: {"nmi"},

	// Mass
	Gram:     {"g", "gram", "grams"},
	Kilogram: {"kg", "kilogram", "kilograms"},
	Pound:    {"lb", "lbs", "pound", "pounds"},
	Stone:    {"st", "stone", "stones"},

	// Temperature
	Celsius:    {"c", "celsius"},
	Fahrenheit: {"f", "fahrenheit"},
	Kelvin:     {"k", "kelvin"},

	// Speed
	MilesPerHour:      {"mph"},
	KilometersPerHour: {"kmh", "km/h", "kmph"},

	// Volume
	Liter:      {"l", "liter", "liters"},
	Centiliter: {"cl", "centiliter", "centiliters"},
	Milliliter: {"ml", "milliliter", "milliliters"},
	Gallon:     {"gal", "gals", "gallon", "gallons"},
	Quart:      {"qt", "quart", "quarts"},
	Pint:       {"pt", "pint", "pints"},
	Cup:        {"cup", "cups"},
	FlOunce:    {"oz", "floz", "ounce", "ounces"},
	Tablespoon: {"tbsp", "tablespoon", "tablespoons"},
	Teaspoon:   {"tsp", "teaspoon", "teaspoons"},

	// Duration
	Second: {"s", "sec", "secs", "second", "seconds"},
	Minute: {"min", "mins", "minute", "minutes"},
	Hour:   {"hr", "hrs", "hour", "hours"},
	Day:    {"day", "days"},
	Week:   {"wk", "week", "weeks"},
	Month:  {"month", "months"},
	Year:   {"yr", "year", "years"},
}

const (
	Length = "Length"
	Mass   = "Mass"
)

var UnitsLength = struct {
	Meter      string `json:"meter"`
	Kilometer  string `json:"kilometer"`
	Millimeter string `json:"millimeter"`
	Centimeter string `json:"centimeter"`
	Nanometer  string `json:"nanometer"`
}{
	Meter:      "meter",
	Kilometer:  "kilometer",
	Millimeter: "millimeter",
	Centimeter: "centimeter",
	Nanometer:  "nanometer",
}

var UnitsMass = struct {
	Gram string `json:"gram"`
}{
	Gram: "gram",
}

var UnitsMap = map[string]interface{}{
	"UnitsLength": UnitsLength,
	"UnitsMass":   UnitsMass,
}
