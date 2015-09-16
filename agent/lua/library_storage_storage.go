package lua

import (
	"github.com/Shopify/go-lua"
	"github.com/mtabini/goluago/util"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

func pushStorage(l *lua.State) {
	l.PushUserData(map[string]interface{}{})
	l.NewTable()
	l.SetUserValue(-2)

	l.NewTable()
	lua.SetFunctions(l, []lua.RegistryFunction{
		{"__newindex", func(l *lua.State) int {
			var v interface{}
			var err error

			if l.IsTable(3) {
				if v, err = util.PullTable(l, 3); err != nil {
					panic("unreachable")
				}
			} else {
				v = l.ToValue(3)
			}

			if err := aggregations.WriteStorage(lua.CheckString(l, 2), v); err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			return 0
		}},
		{"__index", func(l *lua.State) int {
			if res, err := aggregations.ReadStorage(lua.CheckString(l, 2)); err == nil {
				util.DeepPush(l, res)
			} else {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			return 1
		}},
	}, 0)

	l.SetMetaTable(-2)
}
