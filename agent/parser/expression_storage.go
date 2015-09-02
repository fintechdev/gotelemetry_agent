package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

// Number

type storageExpression struct {
	l int
	p int
}

func newStorageExpression(line, position int) expression {
	result := &storageExpression{
		l: line,
		p: position,
	}

	return result
}

type storageProperty func(g *storageExpression) expression

var storageProperties = map[string]storageProperty{
	"get": func(s *storageExpression) expression {
		return s.get()
	},
	"set": func(s *storageExpression) expression {
		return s.set()
	},
}

func (s *storageExpression) get() expression {
	return newCallableExpression(
		"get",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			res, err := aggregations.ReadStorage(args["key"].(string))

			if err != nil {
				return nil, err
			}

			return newMapExpression(res, s.l, s.p), nil
		},
		map[string]callableArgument{
			"key": callableArgumentString,
		},
		s.l,
		s.p,
	)
}

func (s *storageExpression) set() expression {
	return newCallableExpression(
		"set",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			err := aggregations.WriteStorage(args["key"].(string), args["value"].(map[string]interface{}))

			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		map[string]callableArgument{
			"key":   callableArgumentString,
			"value": callableArgumentMap,
		},
		s.l,
		s.p,
	)
}

func (s *storageExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := storageProperties[property]; ok {
		return f(s), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", s, property))
}

func (s *storageExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, errors.New(fmt.Sprintf("%s cannot be evaluated", s))
}

func (s *storageExpression) line() int {
	return s.l
}

func (s *storageExpression) position() int {
	return s.p
}

func (s *storageExpression) String() string {
	return fmt.Sprintf("Storage")
}
