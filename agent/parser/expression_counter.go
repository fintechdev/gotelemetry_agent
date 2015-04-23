package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

// Number

type counterExpression struct {
	name    string
	counter *aggregations.Counter
	l       int
	p       int
}

func newCounterExpression(name string, counter *aggregations.Counter, line, position int) expression {
	result := &counterExpression{
		name:    name,
		counter: counter,
		l:       line,
		p:       position,
	}

	return result
}

type counterProperty func(x *counterExpression) expression

var counterProperties = map[string]counterProperty{
	"set": func(x *counterExpression) expression {
		return x.set()
	},
	"increment": func(x *counterExpression) expression {
		return x.increment()
	},
	"reset": func(x *counterExpression) expression {
		return x.reset()
	},
}

func (x *counterExpression) set() expression {
	return newCallableExpression(
		"set",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			x.counter.SetValue(int64(args["value"].(float64)))

			return nil, nil
		},
		map[string]callableArgument{
			"value": callableArgumentNumeric,
		},
		x.l,
		x.p,
	)
}

func (x *counterExpression) increment() expression {
	return newCallableExpression(
		"increment",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			x.counter.Increment(int64(args["delta"].(float64)))

			return nil, nil
		},
		map[string]callableArgument{
			"delta": callableArgumentNumeric,
		},
		x.l,
		x.p,
	)
}

func (x *counterExpression) reset() expression {
	return newCallableExpression(
		"reset",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			x.counter.SetValue(0.0)

			return nil, nil
		},
		map[string]callableArgument{},
		x.l,
		x.p,
	)
}

func (x *counterExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := counterProperties[property]; ok {
		return f(x), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", x, property))
}

func (x *counterExpression) evaluate(c *executionContext) (interface{}, error) {
	return newNumericExpression(float64(x.counter.GetValue()), x.l, x.p), nil
}

func (x *counterExpression) line() int {
	return x.l
}

func (x *counterExpression) position() int {
	return x.p
}

func (x *counterExpression) String() string {
	return fmt.Sprintf("Series(%s)", x.name)
}
