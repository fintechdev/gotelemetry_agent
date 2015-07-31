package parser

import (
	"errors"
	"fmt"
)

type callable interface {
	call(c *executionContext, arguments map[string]interface{}) (expression, error)
}

type extractable interface {
	extract(c *executionContext, property string) (expression, error)
}

type expression interface {
	fmt.Stringer
	evaluate(c *executionContext) (interface{}, error)
	line() int
	position() int
}

func expressionFromInterface(v interface{}, line, position int) (expression, error) {
	if _, ok := v.(expression); ok {
		return v.(expression), nil
	}

	switch v.(type) {
	case nil:
		return newNullExpression(nil, line, position), nil

	case map[string]interface{}:
		return newMapExpression(v.(map[string]interface{}), line, position), nil

	case []interface{}:
		return newArrayExpression(v.([]interface{}), line, position), nil

	case int, int64, float64:
		return newNumericExpression(v, line, position), nil

	case string:
		return newStringExpression(v, line, position), nil

	case bool:
		return newBooleanExpression(v, line, position), nil

	default:
		return nil, errors.New(fmt.Sprintf("Unable to create an expression with a value of type %T", v))
	}
}

type resolvable interface {
	resolve(*executionContext) (interface{}, error)
}

func resolveExpression(c *executionContext, e expression) (interface{}, error) {
	var v interface{} = e
	var err error

	for {
		switch v.(type) {
		case resolvable:
			v, err = v.(resolvable).resolve(c)
			break

		case expression:
			v, err = v.(expression).evaluate(c)

		default:
			return v, nil
		}

		if err != nil {
			return nil, err
		}
	}
}
