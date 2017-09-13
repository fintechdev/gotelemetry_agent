package lua

import (
	"github.com/telemetryapp/go-lua"
)

var mongoLibrary = []lua.RegistryFunction{
	lua.RegistryFunction{
		Name: "open",
		Function: func(l *lua.State) int {
			connString := lua.CheckString(l, 1)

			pushGoConnection(l, connString)

			return 1
		},
	},
}

func openMongoLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, mongoLibrary)
		return 1
	}

	lua.Require(l, "telemetry/mongodb", open, false)
	l.Pop(1)
}
