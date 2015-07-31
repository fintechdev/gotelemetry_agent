package parser

import (
	"errors"
	"fmt"
)

// Number

type nullExpression struct {
	l int
	p int
}

func newNullExpression(value interface{}, line, position int) expression {
	result := &nullExpression{
		l: line,
		p: position,
	}

	return result
}

func (a *nullExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, nil
}

func (a *nullExpression) resolve(c *executionContext) (interface{}, error) {
	return nil, nil
}

func (a *nullExpression) extract(c *executionContext, property string) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", a, property))
}

func (a *nullExpression) line() int {
	return a.l
}

func (a *nullExpression) position() int {
	return a.p
}

func (a *nullExpression) String() string {
	return fmt.Sprintf("null")
}
