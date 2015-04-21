package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"time"
)

// Number

type seriesExpression struct {
	name   string
	series *aggregations.Series
	l      int
	p      int
}

func newSeriesExpression(name string, series *aggregations.Series, line, position int) expression {
	result := &seriesExpression{
		name:   name,
		series: series,
		l:      line,
		p:      position,
	}

	return result
}

type seriesProperty func(g *seriesExpression) expression

var seriesProperties = map[string]seriesProperty{
	"last": func(s *seriesExpression) expression {
		return s.last()
	},
	"aggregate": func(s *seriesExpression) expression {
		return s.aggregate()
	},
}

func (s *seriesExpression) last() expression {
	return newCallableExpression(
		"last",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			l, err := s.series.Last()

			if err != nil {
				return nil, err
			}

			return newSeriesResultExpression(l["value"].(float64), l["ts"].(int64), s.l, s.p), nil
		},
		map[string]callableArgument{},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) aggregate() expression {
	return newCallableExpression(
		"",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			function, err := aggregations.AggregationFunctionTypeFromName(args["func"].(string))

			if err != nil {
				return nil, err
			}

			intervalDuration, err := time.ParseDuration(args["interval"].(string))

			if err != nil {
				return nil, err
			}

			interval := int(intervalDuration.Seconds())

			if interval < 1 {
				interval = 1
			}

			count := int(args["count"].(float64))

			l, err := s.series.Aggregate(function, interval, count)

			if err != nil {
				return nil, err
			}

			return newSeriesResultAggregateExpression(l.([]interface{}), s.l, s.p), nil
		},
		map[string]callableArgument{
			"func":     callableArgumentString,
			"interval": callableArgumentString,
			"count":    callableArgumentNumeric,
		},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := seriesProperties[property]; ok {
		return f(s), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", s, property))
}

func (s *seriesExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, errors.New(fmt.Sprintf("%s cannot be evaluated", s))
}

func (s *seriesExpression) line() int {
	return s.l
}

func (s *seriesExpression) position() int {
	return s.p
}

func (s *seriesExpression) String() string {
	return fmt.Sprintf("Series(%s)", s.name)
}
