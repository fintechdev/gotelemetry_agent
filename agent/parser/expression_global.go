package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"time"
)

// Number

type globalExpression struct {
	l int
	p int
}

func newGlobalExpression(line, position int) expression {
	result := &globalExpression{line, position}

	return result
}

func (g *globalExpression) evaluate(c *executionContext) (interface{}, error) {
	return g, nil
}

// Properties

type globalProperty func(g *globalExpression) expression

var globalProperties = map[string]globalProperty{
	"now": func(g *globalExpression) expression {
		return g.now()
	},
	"series": func(g *globalExpression) expression {
		return g.series()
	},
}

func (g *globalExpression) now() expression {
	return newCallableExpression(
		"now",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return newNumericExpression(time.Now().Unix(), g.l, g.p), nil
		},
		map[string]callableArgument{},
		g.l,
		g.p,
	)
}

func (g *globalExpression) series() expression {
	return newCallableExpression(
		"series",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)
			s, err := aggregations.GetSeries(c.aggregationContext, n)

			if err != nil {
				return nil, err
			}

			return newSeriesExpression(n, s, g.p, g.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := globalProperties[property]; ok {
		return f(g), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", g, property))
}

func (g *globalExpression) line() int {
	return g.l
}

func (g *globalExpression) position() int {
	return g.p
}

func (g *globalExpression) String() string {
	return fmt.Sprintf("Global()")
}
