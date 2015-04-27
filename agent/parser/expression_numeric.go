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
	case *numericExpression:
		v := value.(*numericExpression)

		result.value = v.value
		result.err = v.err

	case float64:
		result.value = value.(float64)

	case int:
		result.value = float64(value.(int))

	case int64:
		result.value = float64(value.(int64))

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

func (n *numericExpression) evaluate(c *executionContext) (interface{}, error) {
	return n.value, n.err
}

func (n *numericExpression) line() int {
	return n.l
}

func (n *numericExpression) position() int {
	return n.p
}

func (n *numericExpression) String() string {
	return fmt.Sprintf("Number(%f)", n.value)
}
