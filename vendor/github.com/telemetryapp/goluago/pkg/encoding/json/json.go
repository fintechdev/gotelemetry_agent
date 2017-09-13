package json

import (
	"encoding/json"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

func Open(l *lua.State) {
	jsonOpen := func(l *lua.State) int {
		lua.NewLibrary(l, jsonLibrary)
		return 1
	}
	lua.Require(l, "goluago/encoding/json", jsonOpen, false)
	l.Pop(1)
}

var jsonLibrary = []lua.RegistryFunction{
	{"unmarshal", unmarshal},
}

func unmarshal(l *lua.State) int {
	payload := lua.CheckString(l, 1)

	var output interface{}

	if err := json.Unmarshal([]byte(payload), &output); err != nil {
		lua.Errorf(l, err.Error())
		panic("unreachable")
	}
	return util.DeepPush(l, output)
}
