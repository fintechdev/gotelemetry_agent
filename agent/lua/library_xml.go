package lua

import (
	"github.com/clbanning/mxj"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var xmlLibrary = []lua.RegistryFunction{
	lua.RegistryFunction{
		Name: "encode",
		Function: func(l *lua.State) int {
			lua.CheckAny(l, 1)

			var v interface{}
			var err error

			if l.IsTable(1) {
				if v, err = util.PullTable(l, 1); err != nil {
					lua.Errorf(l, "%s", err.Error())
				}
			} else {
				lua.Errorf(l, "Only tables can be converted to XML")
			}

			res, err := mxj.Map(v.(map[string]interface{})).Xml()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			util.DeepPush(l, string(res))

			return 1
		},
	},
	lua.RegistryFunction{
		Name: "decode",
		Function: func(l *lua.State) int {
			res, err := mxj.NewMapXml([]byte(lua.CheckString(l, 1)))

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
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
