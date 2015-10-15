package lua

import (
	"bytes"
	"github.com/mtabini/go-lua"
	"github.com/mtabini/goluago/util"
	"io/ioutil"
	"net/http"
)

var httpLibrary = []lua.RegistryFunction{
	{
		"get",
		func(l *lua.State) int {
			res, err := http.Get(lua.CheckString(l, 1))

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},

	{
		"post",
		func(l *lua.State) int {
			res, err := http.Post(lua.CheckString(l, 1), lua.CheckString(l, 2), bytes.NewBuffer([]byte(lua.CheckString(l, 3))))

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},
}

func openHTTPLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, httpLibrary)
		return 1
	}

	lua.Require(l, "telemetry/http", open, false)
	l.Pop(1)
}
