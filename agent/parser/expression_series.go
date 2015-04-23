package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"io"
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
	"avg": func(s *seriesExpression) expression {
		return s.compute("avg", aggregations.Avg)
	},
	"min": func(s *seriesExpression) expression {
		return s.compute("min", aggregations.Min)
	},
	"max": func(s *seriesExpression) expression {
		return s.compute("max", aggregations.Max)
	},
	"sum": func(s *seriesExpression) expression {
		return s.compute("sum", aggregations.Sum)
	},
	"count": func(s *seriesExpression) expression {
		return s.compute("count", aggregations.Count)
	},
	"last": func(s *seriesExpression) expression {
		return s.last()
	},
	"aggregate": func(s *seriesExpression) expression {
		return s.aggregate()
	},
	"trim": func(s *seriesExpression) expression {
		return s.trim()
	},
	"push": func(s *seriesExpression) expression {
		return s.push()
	},
	"pop": func(s *seriesExpression) expression {
		return s.pop()
	},
}

func (s *seriesExpression) compute(name string, functionType aggregations.FunctionType) expression {
	return newCallableExpression(
		"name",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			intervalDuration, err := time.ParseDuration(args["interval"].(string))

			if err != nil {
				return nil, err
			}

			end := time.Now()
			start := end.Add(-intervalDuration)

			l, err := s.series.Compute(functionType, &start, &end)

			if err != nil {
				return nil, err
			}

			return newNumericExpression(l, s.l, s.p), nil
		},
		map[string]callableArgument{
			"interval": callableArgumentString,
		},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) last() expression {
	return newCallableExpression(
		"last",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			l, err := s.series.Last()

			if err != nil && err != io.EOF {
				return nil, err
			}

			v, _ := l["value"].(float64)
			ts, _ := l["ts"].(int64)

			return newSeriesResultExpression(v, ts, s.l, s.p), nil
		},
		map[string]callableArgument{},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) push() expression {
	return newCallableExpression(
		"push",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			var timestamp time.Time

			if ts, ok := args["timestamp"].(float64); ok {
				timestamp = time.Unix(int64(ts), 0)
			} else {
				timestamp = time.Now()
			}

			value := args["value"].(float64)

			return nil, s.series.Push(&timestamp, value)
		},
		map[string]callableArgument{
			"timestamp": callableArgumentOptionalNumeric,
			"value":     callableArgumentNumeric,
		},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) pop() expression {
	return newCallableExpression(
		"pop",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			l, err := s.series.Pop(false)

			if err != nil && err != io.EOF {
				return nil, err
			}

			v, _ := l["value"].(float64)
			ts, _ := l["ts"].(int64)

			return newSeriesResultExpression(v, ts, s.l, s.p), nil
		},
		map[string]callableArgument{},
		s.l,
		s.p,
	)
}

func (s *seriesExpression) trim() expression {
	return newCallableExpression(
		"trim",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			since, ok1 := args["since"].(string)
			count, ok2 := args["count"].(float64)

			if ok1 && ok2 {
				return nil, errors.New("A call to trim() cannot contain both `since` and `count` arguments")
			}

			if ok1 {
				d, err := time.ParseDuration(since)

				if err != nil {
					return nil, err
				}

				return nil, s.series.TrimSince(time.Now().Add(d))
			}

			if ok2 {
				return nil, s.series.TrimCount(int(count))
			}

			return nil, errors.New("Either a `since` or `count` argument is required when calling trim()")
		},
		map[string]callableArgument{
			"since": callableArgumentOptionalString,
			"count": callableArgumentOptionalNumeric,
		},
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
