package parser

import (
	"errors"
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

func (s *stringExpression) extract(c *executionContext, property string) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", s, property))
}

func (s *stringExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s is not a function", s))
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
