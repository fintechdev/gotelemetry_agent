package lua

import (
	"github.com/mtabini/go-lua"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

var counterFunctions = map[string]func(c *aggregations.Counter) lua.Function{
	"value": func(c *aggregations.Counter) lua.Function {
		return func(l *lua.State) int {
			l.PushInteger(int(c.GetValue()))

			return 1
		}
	},

	"set": func(c *aggregations.Counter) lua.Function {
		return func(l *lua.State) int {
			c.SetValue(int64(lua.CheckInteger(l, 1)))

			return 0
		}
	},

	"increment": func(c *aggregations.Counter) lua.Function {
		return func(l *lua.State) int {
			c.Increment(int64(lua.CheckInteger(l, 1)))

			return 0
		}
	},
}

func pushCounter(l *lua.State, name string) {
	counter, _, err := aggregations.GetCounter(name)

	if err != nil {
		lua.Errorf(l, "%s", err)
		panic("unreachable")
	}

	l.NewTable()

	for name, fn := range counterFunctions {
		l.PushGoFunction(fn(counter))
		l.SetField(-2, name)
	}
}
