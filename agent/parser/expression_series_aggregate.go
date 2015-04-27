package parser

import (
	"errors"
	"fmt"
)

// Number

type seriesResultAggregateExpression struct {
	values []interface{}
	l      int
	p      int
}

func newSeriesResultAggregateExpression(values []interface{}, line, position int) expression {
	result := &seriesResultAggregateExpression{
		values: values,
		l:      line,
		p:      position,
	}

	return result
}

func (s *seriesResultAggregateExpression) extract(c *executionContext, property string) (expression, error) {
	switch property {
	case "timestamps":
		ts := make([]interface{}, len(s.values))

		for index, dp := range s.values {
			dataPoint := dp.(map[string]interface{})

			ts[index] = float64(dataPoint["ts"].(int64))
		}

		return newArrayExpression(ts, s.l, s.p), nil

	case "values":
		res := make([]interface{}, len(s.values))

		for index, dp := range s.values {
			dataPoint := dp.(map[string]interface{})

			res[index] = dataPoint["value"].(float64)
		}

		return newArrayExpression(res, s.l, s.p), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", s, property))
}

func (s *seriesResultAggregateExpression) evaluate(c *executionContext) (interface{}, error) {
	ts := make([]interface{}, len(s.values))

	for index, dp := range s.values {
		dataPoint := dp.(map[string]interface{})

		ts[index] = dataPoint["value"].(float64)
	}

	return newArrayExpression(ts, s.l, s.p), nil
}

func (s *seriesResultAggregateExpression) line() int {
	return s.l
}

func (s *seriesResultAggregateExpression) position() int {
	return s.p
}

func (s *seriesResultAggregateExpression) String() string {
	return fmt.Sprintf("SeriesAggregate(%+v)", s.values)
}
