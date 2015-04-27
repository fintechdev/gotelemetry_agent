package parser

import (
	"errors"
	"fmt"
)

// Number

type arrayExpression struct {
	values []interface{}
	l      int
	p      int
}

func newArrayExpression(value interface{}, line, position int) expression {
	var values []interface{}
	var ok bool

	if values, ok = value.([]interface{}); !ok {
		values = []interface{}{value}
	}

	result := &arrayExpression{
		values: values,
		l:      line,
		p:      position,
	}

	return result
}

func (a *arrayExpression) evaluate(c *executionContext) (interface{}, error) {
	res := make([]interface{}, len(a.values))

	for index, e := range a.values {
		ex, err := expressionFromInterface(e, a.l, a.p)

		if err != nil {
			return nil, errors.New(fmt.Sprintf("[Evaluating item at index %d]: %s", index, err))
		}

		if v, err := resolveExpression(c, ex); err == nil {
			res[index] = v
		} else {
			return nil, errors.New(fmt.Sprintf("[Evaluating item at index %d]: %s", index, err))
		}
	}

	return res, nil
}

// Properties

type arrayProperty func(a *arrayExpression) expression

var arrayProperties = map[string]arrayProperty{
	"item": func(a *arrayExpression) expression {
		return a.item()
	},
	"count": func(a *arrayExpression) expression {
		return a.count()
	},
	"sum": func(a *arrayExpression) expression {
		return a.sum()
	},
	"min": func(a *arrayExpression) expression {
		return a.min()
	},
	"max": func(a *arrayExpression) expression {
		return a.max()
	},
	"avg": func(a *arrayExpression) expression {
		return a.avg()
	},
	"average": func(a *arrayExpression) expression {
		return a.avg()
	},
	"stddev": func(a *arrayExpression) expression {
		return a.stddev()
	},
}

func (a *arrayExpression) count() expression {
	return newCallableExpression(
		"count",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return newNumericExpression(len(a.values), a.p, a.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) sum() expression {
	return newCallableExpression(
		"sum",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if result, err := a.numericArray(c); err == nil {
				return newNumericExpression(result.sum(), a.p, a.p), nil
			} else {
				return nil, err
			}
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) min() expression {
	return newCallableExpression(
		"min",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if result, err := a.numericArray(c); err == nil {
				return newNumericExpression(result.min(), a.p, a.p), nil
			} else {
				return nil, err
			}
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) max() expression {
	return newCallableExpression(
		"max",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if result, err := a.numericArray(c); err == nil {
				return newNumericExpression(result.max(), a.p, a.p), nil
			} else {
				return nil, err
			}
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) avg() expression {
	return newCallableExpression(
		"avg",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if result, err := a.numericArray(c); err == nil {
				return newNumericExpression(result.avg(), a.p, a.p), nil
			} else {
				return nil, err
			}
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) stddev() expression {
	return newCallableExpression(
		"avg",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if result, err := a.numericArray(c); err == nil {
				return newNumericExpression(result.stddev(), a.p, a.p), nil
			} else {
				return nil, err
			}
		},
		map[string]callableArgument{},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) item() expression {
	return newCallableExpression(
		"item",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			index := int(args["index"].(float64))

			if index < 0 || index > len(a.values)-1 {
				return nil, errors.New(fmt.Sprintf("Invalid index %d", index))
			}

			return expressionFromInterface(a.values[index], a.l, a.p)
		},
		map[string]callableArgument{"index": callableArgumentNumeric},
		a.l,
		a.p,
	)
}

func (a *arrayExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := arrayProperties[property]; ok {
		return f(a), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", a, property))
}

func (a *arrayExpression) line() int {
	return a.l
}

func (a *arrayExpression) position() int {
	return a.p
}

func (a *arrayExpression) String() string {
	return fmt.Sprintf("Array(%v)", a.values)
}
