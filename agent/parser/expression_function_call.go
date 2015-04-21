package parser

import (
	"errors"
	"fmt"
)

// Number

type functionCallExpression struct {
	target       expression
	argumentList map[string]expression
	l            int
	p            int
}

func newFunctionCallExpression(target expression, argumentList map[string]expression, line, position int) expression {
	result := &functionCallExpression{
		target:       target,
		argumentList: argumentList,
		l:            line,
		p:            position,
	}

	return result
}

func (f *functionCallExpression) evaluate(c *executionContext) (interface{}, error) {
	args := map[string]interface{}{}

	for key, argument := range f.argumentList {
		if res, err := argument.evaluate(c); err == nil {
			args[key] = res
		} else {
			return nil, err
		}
	}

	return f.target.call(c, args)
}

func (f *functionCallExpression) extract(c *executionContext, property string) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", f, property))
}

func (f *functionCallExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s is not a function", f))
}

func (f *functionCallExpression) line() int {
	return f.l
}

func (f *functionCallExpression) position() int {
	return f.p
}

func (f *functionCallExpression) String() string {
	return fmt.Sprintf("FunctionCall(%s(%+v))", f.target, f.argumentList)
}
