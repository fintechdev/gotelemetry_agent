package parser

import (
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

func resolveExpression(c *executionContext, e expression) (interface{}, error) {
	var v interface{} = e
	var err error

	for {
		if vv, ok := v.(expression); ok {
			v, err = vv.evaluate(c)

			if err != nil {
				return nil, err
			}
		} else {
			return v, nil
		}
	}
}
