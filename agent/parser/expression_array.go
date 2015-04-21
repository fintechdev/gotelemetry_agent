package parser

import (
	"fmt"
)

// Number

type arrayExpression struct {
	values []interface{}
	l      int
	p      int
}

func newArrayExpression(values []interface{}, line, position int) expression {
	result := &arrayExpression{
		values: values,
		l:      line,
		p:      position,
	}

	return result
}

func (a *arrayExpression) evaluate(c *executionContext) (interface{}, error) {
	return a.values, nil
}

func (a *arrayExpression) line() int {
	return a.l
}

func (a *arrayExpression) position() int {
	return a.p
}

func (a *arrayExpression) String() string {
	return fmt.Sprintf("Array(%v)", a.values)
}
