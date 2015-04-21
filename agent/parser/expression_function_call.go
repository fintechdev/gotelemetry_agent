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
	if cl, ok := f.target.(callable); ok {
		args := map[string]interface{}{}

		for key, argument := range f.argumentList {
			if res, err := argument.evaluate(c); err == nil {
				args[key] = res
			} else {
				return nil, err
			}
		}

		return cl.call(c, args)
	}

	return nil, errors.New(fmt.Sprintf("%s is not a function", f.target))
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
