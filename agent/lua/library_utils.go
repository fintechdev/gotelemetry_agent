package lua

import (
	"time"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var utilsLibrary = []lua.RegistryFunction{
	lua.RegistryFunction{
		Name: "nowutc",
		Function: func(l *lua.State) int {
			util.DeepPush(l, time.Now().UTC().Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "now",
		Function: func(l *lua.State) int {
			util.DeepPush(l, time.Now().Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowplusseconds",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, time.Now().Add(time.Duration(seconds)*time.Second).Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowutcplusseconds",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, time.Now().UTC().Add(time.Duration(seconds)*time.Second).Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowminusseconds",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, time.Now().Add(-time.Duration(seconds)*time.Second).Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowutcminusseconds",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, time.Now().UTC().Add(-time.Duration(seconds)*time.Second).Format(time.RFC3339))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowutcepoch",
		Function: func(l *lua.State) int {
			util.DeepPush(l, int(time.Now().UTC().Unix()))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowepoch",
		Function: func(l *lua.State) int {
			util.DeepPush(l, int(time.Now().Unix()))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowplussecondsepoch",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, int(time.Now().Add(time.Duration(seconds)*time.Second).Unix()))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowutcplussecondsepoch",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, int(time.Now().UTC().Add(time.Duration(seconds)*time.Second).Unix()))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowminussecondsepoch",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, int(time.Now().Add(-time.Duration(seconds)*time.Second).Unix()))
			return 1
		},
	},
	lua.RegistryFunction{
		Name: "nowutcminussecondsepoch",
		Function: func(l *lua.State) int {
			seconds := lua.CheckInteger(l, 1)
			util.DeepPush(l, int(time.Now().UTC().Add(-time.Duration(seconds)*time.Second).Unix()))
			return 1
		},
	},
}

func openUtilsLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, utilsLibrary)
		return 1
	}

	lua.Require(l, "telemetry/utils", open, false)
	l.Pop(1)
}
