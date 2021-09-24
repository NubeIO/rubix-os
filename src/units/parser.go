package unit

import (
	"bytes"
	"regexp"
	"strconv"
)

var wsPattern = regexp.MustCompile(`^\s*`)

func ws(s []byte) int {
	return len(wsPattern.Find(s))
}

// Parser scans bytes and returns result a result, number of bytes consumed, and whether the parse was successful
type Parser func([]byte) (interface{}, int, bool)

// Token scans for a pattern, skipping leading whitespace if necessary
func Token(pattern string) Parser {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) (interface{}, int, bool) {
		n := ws(s)
		match := regex.Find(s[n:])
		if match == nil {
			return nil, 0, false
		}

		return string(match), n + len(match), true
	}
}

// Sub matches sub expressions in a regex pattern
func Sub(pattern string) Parser {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) (interface{}, int, bool) {
		n := ws(s)
		match := regex.FindSubmatch(s[n:])
		if match == nil {
			return nil, 0, false
		}
		matches := make(map[string]string)
		for i, n := range regex.SubexpNames() {
			if len(n) > 0 {
				matches[n] = string(match[i])
			}
		}
		return matches, n + len(match[0]), true
	}
}

// All matches all parsers in order
// Will result in a slice of all of the parsed values
func All(ps ...Parser) Parser {
	return func(s []byte) (interface{}, int, bool) {
		var (
			vs  []interface{}
			sum int
		)
		for _, p := range ps {
			res, n, ok := p(s[sum:])
			if !ok {
				return nil, 0, false
			}
			vs = append(vs, res)
			sum += n
		}
		return vs, sum, true
	}
}

// Any matches any of the parsers, tested in order
func Any(ps ...Parser) Parser {
	return func(s []byte) (interface{}, int, bool) {
		for _, p := range ps {
			res, n, ok := p(s)
			if ok {
				return res, n, true
			}
		}
		return nil, 0, false
	}
}

// Atom scans for a single atom, skipping whitespace
func Atom(val string) Parser {
	b := []byte(val)
	return func(s []byte) (interface{}, int, bool) {
		if n := ws(s); bytes.HasPrefix(s[n:], b) {
			return nil, n + len(b), true
		}
		return nil, 0, false
	}
}

// None is type representing missing from Opt
type None struct{}

// Opt turns a parser into an optional parser
// Will return a result that's either Missing or the result itself
func (p Parser) Opt() Parser {
	return func(s []byte) (interface{}, int, bool) {
		res, n, ok := p(s)
		if !ok {
			return None{}, 0, true
		}
		return res, n, true
	}
}

// Or is like Opt but returns the value given if the parser is unsuccessful
func (p Parser) Or(or interface{}) Parser {
	return func(s []byte) (interface{}, int, bool) {
		res, n, ok := p(s)
		if !ok {
			return or, 0, true
		}
		return res, n, true
	}
}

// Map maps the result of a parser to a different result
// If the Mapper returns nil, the parser returns as invalid
func (p Parser) Map(f MapperFunc) Parser {
	return func(s []byte) (interface{}, int, bool) {
		if res, n, ok := p(s); ok {
			if v := f(res); v != nil {
				return v, n, true
			}
		}
		return nil, 0, false
	}
}

// Float is a float parser
var Float = Token(`[+-]?\d+([.,]\d*)?([eE][+-]?\d+)?`).Map(mapFloat)

// Int is an integer parser
var Int = Token(`[+-]?\d+`).Map(mapInt)

// MapperFunc is a function for mapping parser results
type MapperFunc = func(interface{}) interface{}

// Index creates a mapper that maps the result to an index in a slice
func Index(i int) MapperFunc {
	return func(v interface{}) interface{} {
		return v.([]interface{})[i]
	}
}

func mapFloat(v interface{}) interface{} {
	f, err := strconv.ParseFloat(v.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func mapInt(v interface{}) interface{} {
	i, err := strconv.Atoi(v.(string))
	if err != nil {
		panic(err)
	}
	return i
}
