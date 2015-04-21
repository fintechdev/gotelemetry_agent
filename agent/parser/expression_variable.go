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
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown variable `%s`", v.name))
	}
}

func (v *variableExpression) extract(c *executionContext, property string) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", v, property))
}

func (v *variableExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s is not a function", v))
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
