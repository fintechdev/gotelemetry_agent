package parser

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"time"
)

// Number

type globalExpression struct {
	l int
	p int
}

func newGlobalExpression(line, position int) expression {
	result := &globalExpression{line, position}

	return result
}

func (g *globalExpression) evaluate(c *executionContext) (interface{}, error) {
	return g, nil
}

// Properties

type globalProperty func(g *globalExpression) expression

var globalProperties = map[string]globalProperty{
	"counter": func(g *globalExpression) expression {
		return g.counter()
	},
	"now": func(g *globalExpression) expression {
		return g.now()
	},
	"notify": func(g *globalExpression) expression {
		return g.notify()
	},
	"series": func(g *globalExpression) expression {
		return g.series()
	},
}

func (g *globalExpression) now() expression {
	return newCallableExpression(
		"now",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return newNumericExpression(time.Now().Unix(), g.l, g.p), nil
		},
		map[string]callableArgument{},
		g.l,
		g.p,
	)
}

func (g *globalExpression) counter() expression {
	return newCallableExpression(
		"counter",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)
			res, err := aggregations.GetCounter(c.aggregationContext, n)

			if err != nil {
				return nil, err
			}

			return newCounterExpression(n, res, g.p, g.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) series() expression {
	return newCallableExpression(
		"series",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)
			s, err := aggregations.GetSeries(c.aggregationContext, n)

			if err != nil {
				return nil, err
			}

			return newSeriesExpression(n, s, g.p, g.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) notify() expression {
	return newCallableExpression(
		"notify",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			title := args["title"].(string)
			message := args["message"].(string)

			d, err := time.ParseDuration(args["duration"].(string))

			if err != nil {
				return nil, err
			}

			duration := int(d.Seconds())

			if duration < 1 {
				duration = 1
			}

			tag, _ := args["tag"].(string)
			icon, _ := args["icon"].(string)
			sound, ok := args["sound"].(string)

			if !ok {
				sound = "default"
			}

			channelTag := args["channel"].(string)

			notification := gotelemetry.NewNotification(title, message, icon, duration, sound, tag)

			result := c.notificationProvider.SendNotification(notification, channelTag)

			return newBooleanExpression(result, g.l, g.p), nil
		},
		map[string]callableArgument{
			"channel":  callableArgumentString,
			"title":    callableArgumentString,
			"message":  callableArgumentString,
			"duration": callableArgumentString,
			"tag":      callableArgumentOptionalBoolean,
			"icon":     callableArgumentOptionalString,
			"sound":    callableArgumentOptionalString,
		},
		g.l,
		g.p,
	)
}
func (g *globalExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := globalProperties[property]; ok {
		return f(g), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", g, property))
}

func (g *globalExpression) line() int {
	return g.l
}

func (g *globalExpression) position() int {
	return g.p
}

func (g *globalExpression) String() string {
	return fmt.Sprintf("Global()")
}
