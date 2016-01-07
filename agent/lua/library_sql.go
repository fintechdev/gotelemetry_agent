package lua

import (
	"github.com/telemetryapp/go-lua"
)

var sqlLibrary = []lua.RegistryFunction{
	{
		"open",
		func(l *lua.State) int {
			driver := lua.CheckString(l, 1)
			connString := lua.CheckString(l, 2)

			pushSQLInstance(l, driver, connString)

			return 1
		},
	},
}

func openSQLLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, sqlLibrary)
		return 1
	}

	lua.Require(l, "telemetry/sql", open, false)
	l.Pop(1)
}
