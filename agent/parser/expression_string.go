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

func (s *stringExpression) evaluate(c *executionContext) (interface{}, error) {
	return s.value, s.err
}

func (s *stringExpression) line() int {
	return s.l
}

func (s *stringExpression) position() int {
	return s.p
}

func (s *stringExpression) String() string {
	return fmt.Sprintf("String(%s)", s.value)
}
