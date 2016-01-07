package lua

import (
	"bytes"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
	"github.com/telemetryapp/gotelemetry_agent/agent/oauth"
	"io/ioutil"
	"net/http"
)

func oauthRequest(l *lua.State, entryName, method, url string, body string) int {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	res, err := oauth.Do(entryName, req)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	util.DeepPush(l, string(data))

	return 1
}

var oauthLibrary = []lua.RegistryFunction{
	{
		"get",
		func(l *lua.State) int {
			return oauthRequest(
				l,
				lua.CheckString(l, 1),
				"GET",
				lua.CheckString(l, 2),
				"",
			)
		},
	},

	{
		"post",
		func(l *lua.State) int {
			return oauthRequest(
				l,
				lua.CheckString(l, 1),
				"GET",
				lua.CheckString(l, 2),
				lua.CheckString(l, 3),
			)
		},
	},
}

func openOAuthLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, oauthLibrary)
		return 1
	}

	lua.Require(l, "telemetry/oauth", open, false)
	l.Pop(1)
}
