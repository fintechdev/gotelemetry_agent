package lua

import (
	"github.com/clbanning/mxj"
	"github.com/mtabini/go-lua"
	"github.com/mtabini/goluago/util"
)

var xmlLibrary = []lua.RegistryFunction{
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
				lua.Errorf(l, "Only tables can be converted to XML")
				panic("unreachable")
			}

			res, err := mxj.Map(v.(map[string]interface{})).Xml()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			util.DeepPush(l, string(res))

			return 1
		},
	},

	{
		"decode",
		func(l *lua.State) int {
			res, err := mxj.NewMapXml([]byte(lua.CheckString(l, 1)))

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			util.DeepPush(l, res)

			return 1
		},
	},
}

func openXMLLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, xmlLibrary)
		return 1
	}

	lua.Require(l, "telemetry/xml", open, false)
	l.Pop(1)
}
