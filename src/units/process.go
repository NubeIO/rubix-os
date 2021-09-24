package unit

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func Process(value float64, fromUnit, toUint string) (string, UnitVal) {
	cmd := fmt.Sprintf("%f %s to %s", value, fromUnit, toUint)
	res, _, ok := convertExpr([]byte(cmd))
	if !ok {
		log.Printf("Invalid command: `%v` %v\n", cmd, res)
		return "Usage: !conv [amount][from-unit] to [to-unit]", nil
	}
	c := res.([]interface{})
	cmdFrom := c[0]
	cmdTo := c[2].(string)

	var from UnitVal
	switch uv := cmdFrom.(type) {
	case unparsedUnitVal:
		fromUnit, ok := ParseUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit), nil
		}
		from = fromUnit.FromFloat(uv.val)

	case UnitVal:
		from = uv
	}
	toUnt, ok := ParseUnit(cmdTo)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", cmdTo), nil
	}
	to, err := from.Convert(toUnt)
	if err != nil {
		return err.Error(), nil
	}
	return fmt.Sprintf("%s = %s", from, to), to
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

//func mapCurrency(v interface{}) interface{} {
//	c := v.([]interface{})
//	u, ok := ParseUnit(string(c[0].(rune)))
//	if !ok {
//		return nil
//	}
//	return CurrencyVal{V: c[1].(float64), U: u.(*CurrencyUnit)}
//}
