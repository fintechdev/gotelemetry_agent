package lua

import (
	"encoding/json"
	"github.com/Shopify/go-lua"
	"github.com/mtabini/goluago/util"
)

var jsonLibrary = []lua.RegistryFunction{
	{
		"encode",
		func(l *lua.State) int {
			lua.CheckAny(l, 1)

			var v interface{}
			var err error

			if l.IsTable(1) {
				if v, err = util.PullTable(l, 1); err != nil {
					lua.Errorf(l, "%s", err)
					panic("unreachable")
				}
			} else {
				v = l.ToValue(1)
			}

			res, err := json.Marshal(v)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, string(res))

			return 1
		},
	},

	{
		"decode",
		func(l *lua.State) int {
			var res interface{}
			err := json.Unmarshal([]byte(lua.CheckString(l, 1)), &res)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, res)

			return 1
		},
	},
}

func openJSONLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, jsonLibrary)
		return 1
	}

	lua.Require(l, "telemetry/json", open, false)
	l.Pop(1)
}
