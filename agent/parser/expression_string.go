package parser

import (
	"fmt"
)

// Number

type stringExpression struct {
	value string
	err   error
	l     int
	p     int
}

func newStringExpression(value interface{}, line, position int) expression {
	result := &stringExpression{
		value: fmt.Sprintf("%v", value),
		l:     line,
		p:     position,
	}

	return result
}

func (e *stringExpression) evaluate(c *executionContext) (interface{}, error) {
	return e.value, e.err
}

func (n *stringExpression) line() int {
	return n.l
}

func (n *stringExpression) position() int {
	return n.p
}
