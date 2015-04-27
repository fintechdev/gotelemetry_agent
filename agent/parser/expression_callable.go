package parser

import (
	"errors"
	"fmt"
)

// Callable expression

type callableClosure func(c *executionContext, argList map[string]interface{}) (expression, error)

type callableArgument int

func (c callableArgument) isOptional() bool {
	switch c {
	case callableArgumentOptionalBoolean,
		callableArgumentOptionalArray,
		callableArgumentOptionalNumericArray,
		callableArgumentOptionalNumeric,
		callableArgumentOptionalString:

		return true
	}

	return false
}

const (
	callableArgumentString callableArgument = iota
	callableArgumentOptionalString
	callableArgumentNumeric
	callableArgumentOptionalNumeric
	callableArgumentBoolean
	callableArgumentOptionalBoolean
	callableArgumentArray
	callableArgumentOptionalArray
	callableArgumentNumericArray
	callableArgumentOptionalNumericArray
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

func (e *callableExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	// Handle the case of a single unnamed parameter by rewriting the parameter
	// with the correct name

	if len(arguments) == 1 && len(e.argumentList) == 1 {
		if val, ok := arguments[singleUnnamedArgument]; ok {
			for index, _ := range e.argumentList {
				arguments[index] = val
				delete(arguments, singleUnnamedArgument)
			}
		}
	} else {
		if _, ok := arguments[singleUnnamedArgument]; ok {
			return nil, errors.New("Unnamed arguments are not allowed in calls to functions or methods that support multiple arguments")
		}
	}

	args := map[string]interface{}{}

	for index, argumentType := range e.argumentList {
		if arg, ok := arguments[index]; ok {
			switch argumentType {
			case callableArgumentString, callableArgumentOptionalString:
				if s, err := newStringExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = s
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s on line %d:%d", e, index, err, e.l, e.p))
				}

			case callableArgumentNumeric, callableArgumentOptionalNumeric:
				if n, err := newNumericExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = n
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s on line %d:%d", e, index, err, e.l, e.p))
				}

			case callableArgumentBoolean, callableArgumentOptionalBoolean:
				if b, err := newBooleanExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = b
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s on line %d:%d", e, index, err, e.l, e.p))
				}

			case callableArgumentArray, callableArgumentOptionalArray:
				if s, err := newArrayExpression(arg, e.l, e.p).evaluate(c); err == nil {
					args[index] = s
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s on line %d:%d", e, index, err, e.l, e.p))
				}

			case callableArgumentNumericArray, callableArgumentOptionalNumericArray:
				if s, err := newArrayExpression(arg, e.l, e.p).(*arrayExpression).numericArray(c); err == nil {
					args[index] = s
				} else {
					return nil, errors.New(fmt.Sprintf("%s: cannot evaluate argument `%s`: %s on line %d:%d", e, index, err, e.l, e.p))
				}

			default:
				return nil, errors.New(fmt.Sprintf("%s: unknown argument type `%s` for argument `%s` on line %d:%d", e, argumentType, index, e.l, e.p))
			}
		} else {
			if !argumentType.isOptional() {
				return nil, errors.New(fmt.Sprintf("%s: missing argument `%s` on line %d:%d", e, index, e.l, e.p))
			}
		}
	}

	res, err := e.closure(c, args)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s while evaluating %s on line %d:%d", err, e, e.l, e.p))
	}

	return res, nil
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
