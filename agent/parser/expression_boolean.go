package parser

import (
	"errors"
	"fmt"
	"strconv"
)

// Number

type booleanExpression struct {
	value bool
	err   error
	l     int
	p     int
}

var booleanExpressionZero = &booleanExpression{false, nil, 0, 0}

func newBooleanExpression(value interface{}, line, position int) expression {
	result := &booleanExpression{
		l: line,
		p: position,
	}

	switch value.(type) {
	case bool:
		result.value = value.(bool)

	case float64:
		result.value = value.(float64) != 0.0

	case int:
		result.value = float64(value.(int)) != 0

	case int64:
		result.value = float64(value.(int64)) != 0

	case string:
		if v, err := strconv.ParseFloat(value.(string), 64); err == nil {
			result.value = v != 0.0
		} else {
			result.value = value.(string) != ""
		}

	default:
		result.err = errors.New(fmt.Sprintf("Invalid boolean expression `%+v`", value))
	}

	return result
}

func (b *booleanExpression) evaluate(c *executionContext) (interface{}, error) {
	return b.value, b.err
}

func (b *booleanExpression) line() int {
	return b.l
}

func (b *booleanExpression) position() int {
	return b.p
}

func (b *booleanExpression) String() string {
	return fmt.Sprintf("Boolean(%f)", b.value)
}
