package parser

import (
	"errors"
	"fmt"
	"strconv"
)

// Number

type numericExpression struct {
	value float64
	err   error
	l     int
	p     int
}

var numericExpressionZero = &numericExpression{0.0, nil, 0, 0}

func newNumericExpression(value interface{}, line, position int) expression {
	result := &numericExpression{
		l: line,
		p: position,
	}

	switch value.(type) {
	case float64:
		result.value = value.(float64)

	case int:
		result.value = float64(value.(int))

	case string:
		if v, err := strconv.ParseFloat(value.(string), 64); err == nil {
			result.value = v
		} else {
			result.err = err
		}

	default:
		result.err = errors.New(fmt.Sprintf("Invalid numeric expression `%+v`", value))
	}

	return result
}

func (e *numericExpression) evaluate(c *executionContext) (interface{}, error) {
	return e.value, e.err
}

func (n *numericExpression) line() int {
	return n.l
}

func (n *numericExpression) position() int {
	return n.p
}
