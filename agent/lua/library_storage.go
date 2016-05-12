package lua

import (
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
)

var storageLibrary = []lua.RegistryFunction{
	{
		"series",
		func(l *lua.State) int {
			pushSeries(l, lua.CheckString(l, 1))

			return 1
		},
	},
	{
		"series_find",
		func(l *lua.State) int {
			res, err := database.FindSeries(lua.CheckString(l, 1))

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			pushArray(l)

			for index, name := range res {
				pushSeries(l, name)
				l.RawSetInt(-2, index+1)
			}

			return 1
		},
	},
	{
		"counter",
		func(l *lua.State) int {
			pushCounter(l, lua.CheckString(l, 1))

			return 1
		},
	},
}

func openStorageLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, storageLibrary)

		l.NewTable()

		l.PushInteger(int(database.Sum))
		l.SetField(-2, "SUM")

		l.PushInteger(int(database.Avg))
		l.SetField(-2, "AVG")

		l.PushInteger(int(database.Min))
		l.SetField(-2, "MIN")

		l.PushInteger(int(database.Max))
		l.SetField(-2, "MAX")

		l.PushInteger(int(database.Count))
		l.SetField(-2, "COUNT")

		l.PushInteger(int(database.StdDev))
		l.SetField(-2, "STDDEV")

		l.SetField(-2, "Functions")

		return 1
	}

	lua.Require(l, "telemetry/storage", open, false)
	l.Pop(1)
}
