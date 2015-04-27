package parser

import (
	"errors"
	"fmt"
	"time"
)

// Number

type seriesResultExpression struct {
	value     float64
	timestamp int64
	l         int
	p         int
}

func newSeriesResultExpression(value float64, timestamp int64, line, position int) expression {
	result := &seriesResultExpression{
		value:     value,
		timestamp: timestamp,
		l:         line,
		p:         position,
	}

	return result
}

func (s *seriesResultExpression) extract(c *executionContext, property string) (expression, error) {
	switch property {
	case "timestamp":
		return newNumericExpression(float64(s.timestamp), s.l, s.p), nil

	case "value":
		return newNumericExpression(float64(s.value), s.l, s.p), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", s, property))
}

func (s *seriesResultExpression) evaluate(c *executionContext) (interface{}, error) {
	return newNumericExpression(s.value, s.l, s.p), nil
}

func (s *seriesResultExpression) line() int {
	return s.l
}

func (s *seriesResultExpression) position() int {
	return s.p
}

func (s *seriesResultExpression) String() string {
	return fmt.Sprintf("SeriesResult(%f, %d (%s))", s.value, s.timestamp, time.Unix(s.timestamp, 0))
}
