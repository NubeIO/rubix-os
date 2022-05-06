package units

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func Process(value float64, fromUnit, toUint string) (string, UnitVal, error) {
	cmd := fmt.Sprintf("%f %s to %s", value, fromUnit, toUint)
	res, _, ok := convertExpr([]byte(cmd))
	if !ok {
		log.Printf("Invalid command: `%v` %v\n", cmd, res)
		return "", nil, errors.New("invalid unit usage: [amount][from-unit] to [to-unit]")
	}
	c := res.([]interface{})
	cmdFrom := c[0]
	cmdTo := c[2].(string)

	var from UnitVal
	switch uv := cmdFrom.(type) {
	case unparsedUnitVal:
		fromUnit, ok := ParseUnit(uv.unit)
		if !ok {
			msg := fmt.Sprintf("Invalid unit %s", uv.unit)
			return "", nil, errors.New(msg)
		}
		from = fromUnit.FromFloat(uv.val)

	case UnitVal:
		from = uv
	}
	toUnt, ok := ParseUnit(cmdTo)
	if !ok {
		msg := fmt.Sprintf("Invalid unit %s", cmdTo)
		return "", nil, errors.New(msg)
	}
	to, err := from.Convert(toUnt)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf("%s = %s", from, to), to, nil
}

type unparsedUnitVal struct {
	val  float64
	unit string
}

var (
	unitToken = Token(`[A-Za-z+/$€¥£]+`)

	inches     = All(Int, Atom(`"`).Opt()).Map(Index(0))
	feetInches = All(Int, Atom(`'`), inches.Or(0)).Map(mapFeetInches)

	simpleUnitVal = All(Float, unitToken).Map(mapSimpleUnit)

	//currency = All(RuneIn(`$€¥£`), Float).Map(mapCurrency)
	//fromExpr = Any(simpleUnitVal, feetInches, currency)
	fromExpr = Any(simpleUnitVal, feetInches)

	convertExpr = All(fromExpr, Atom(`to`), unitToken)
)

func mapSimpleUnit(v interface{}) interface{} {
	vs := v.([]interface{})
	return unparsedUnitVal{vs[0].(float64), vs[1].(string)}
}

func mapFeetInches(v interface{}) interface{} {
	vs := v.([]interface{})
	feet := vs[0].(int)
	inches := vs[2].(int)
	return FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}
