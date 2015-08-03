package parser

import (
	"errors"
	"fmt"
)

// Number

type mapExpression struct {
	values map[string]interface{}
	l      int
	p      int
}

func newMapExpression(value interface{}, line, position int) expression {
	var values map[string]interface{}

	switch value.(type) {
	case *mapExpression:
		values = value.(*mapExpression).values

	case map[string]interface{}:
		values = value.(map[string]interface{})

	case map[string]expression:
		source := value.(map[string]expression)

		values = make(map[string]interface{})

		for index, value := range source {
			values[index] = value
		}

	default:
		return nil
	}

	result := &mapExpression{
		values: values,
		l:      line,
		p:      position,
	}

	return result
}

func (a *mapExpression) evaluate(c *executionContext) (interface{}, error) {
	result := map[string]interface{}{}

	for key, v := range a.values {
		if vv, ok := v.(resolvable); ok {
			if vvv, err := vv.resolve(c); err == nil {
				result[key] = vvv
			} else {
				return nil, err
			}

		} else {
			result[key] = v
		}
	}

	return newMapExpression(result, a.l, a.p), nil
}

func (a *mapExpression) resolve(c *executionContext) (interface{}, error) {
	res := make(map[string]interface{})

	for index, e := range a.values {
		ex, err := expressionFromInterface(e, a.l, a.p)

		if err != nil {
			return nil, errors.New(fmt.Sprintf("[Evaluating item at index %s]: %s", index, err))
		}

		if v, err := resolveExpression(c, ex); err == nil {
			res[index] = v
		} else {
			return nil, errors.New(fmt.Sprintf("[Evaluating item at index %s]: %s", index, err))
		}
	}

	return res, nil
}

// Properties

type mapProperty func(a *mapExpression) expression

var mapProperties = map[string]mapProperty{
	"item": func(a *mapExpression) expression {
		return a.item()
	},
	"count": func(a *mapExpression) expression {
		return a.count()
	},
	"set": func(a *mapExpression) expression {
		return a.set()
	},
}

func (a *mapExpression) count() expression {
	return newCallableExpression(
		"count",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return newNumericExpression(len(a.values), a.p, a.p), nil
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *mapExpression) item() expression {
	return newCallableExpression(
		"item",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			index := args["index"].(string)

			if v, ok := a.values[index]; ok {
				return expressionFromInterface(v, a.l, a.p)
			}
			return nil, errors.New(fmt.Sprintf("Invalid index %s", index))
		},
		map[string]callableArgument{"index": callableArgumentString},
		a.l,
		a.p,
	)
}

func (a *mapExpression) set() expression {
	return newCallableExpression(
		"set",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			index := args["index"].(string)
			value := args["value"]

			if vv, ok := value.(resolvable); ok {
				if vvv, err := vv.resolve(c); err == nil {
					value = vvv
				} else {
					return nil, err
				}
			}

			a.values[index] = value

			return nil, nil
		},
		map[string]callableArgument{
			"index": callableArgumentString,
			"value": callableArgumentInterface,
		},
		a.l,
		a.p,
	)
}

func (a *mapExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := mapProperties[property]; ok {
		return f(a), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", a, property))
}

func (a *mapExpression) line() int {
	return a.l
}

func (a *mapExpression) position() int {
	return a.p
}

func (a *mapExpression) String() string {
	return fmt.Sprintf("Map(%v)", a.values)
}
