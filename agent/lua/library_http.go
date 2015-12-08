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
			url := lua.CheckString(l, 1)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			argIndex := 2
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
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
			url := lua.CheckString(l, 1)
			body := lua.CheckString(l, 2)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			argIndex := 3
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},

	{
		"custom",
		func(l *lua.State) int {
			method := lua.CheckString(l, 1)
			url := lua.CheckString(l, 2)
			body := lua.OptString(l, 3, "")

			if len(method) == 0 {
		    method = "POST"
		  }

			var req *http.Request
			var err error
			if len(body) > 0 {
				req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
			} else {
				req, err = http.NewRequest(method, url, nil)
			}
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			argIndex := 4
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
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
