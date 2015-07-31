package parser

import (
	"errors"
	"fmt"
)

// Number

type variableExpression struct {
	name string
	l    int
	p    int
}

func newVariableExpression(name string, line, position int) expression {
	result := &variableExpression{
		name: name,
		l:    line,
		p:    position,
	}

	return result
}

func (v *variableExpression) evaluate(c *executionContext) (interface{}, error) {
	if val, ok := c.variables[v.name]; ok {
		return expressionFromInterface(val, v.l, v.p)
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown variable `%s`", v.name))
	}
}

func (v *variableExpression) line() int {
	return v.l
}

func (v *variableExpression) position() int {
	return v.p
}

func (v *variableExpression) String() string {
	return fmt.Sprintf("Variable(%s)", v.name)
}
