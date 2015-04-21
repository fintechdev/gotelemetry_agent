package parser

import (
	"errors"
	"fmt"
)

// Callable expression

type callableClosure func(c *executionContext, argList map[string]interface{}) (expression, error)

type callableArgument int

const (
	callableArgumentString callableArgument = iota
	callableArgumentNumeric
)

type callableExpression struct {
	name         string
	closure      callableClosure
	argumentList map[string]callableArgument
	l            int
	p            int
}

func newCallableExpression(name string, closure callableClosure, argumentList map[string]callableArgument, line, position int) expression {
	result := &callableExpression{
		name:         name,
		closure:      closure,
		argumentList: argumentList,
		l:            line,
		p:            position,
	}

	return result
}

func (e *callableExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, errors.New(fmt.Sprintf("%s is a function call and cannot be evaluated", e.name))
}

func (e *callableExpression) extract(c *executionContext, property string) (expression, error) {
	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", e, property))
}

func (e *callableExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	args := map[string]interface{}{}

	for index, argumentType := range e.argumentList {
		if arg, ok := arguments[index]; ok {
			switch argumentType {
			case callableArgumentString:
				if s, err := newStringExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = s
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s", e, index, err))
				}

			case callableArgumentNumeric:
				if n, err := newNumericExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = n
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s", e, index, err))
				}

			default:
				return nil, errors.New(fmt.Sprintf("%s: unknown argument type `%s` for argument `%s", e, argumentType, index))
			}
		} else {
			return nil, errors.New(fmt.Sprintf("%s: missing argument `%s`", e, index))
		}
	}

	return e.closure(c, args)
}

func (e *callableExpression) line() int {
	return e.l
}

func (e *callableExpression) position() int {
	return e.p
}

func (e *callableExpression) String() string {
	return fmt.Sprintf("%s(%+v)", e.name, e.argumentList)
}
