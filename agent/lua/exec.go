package lua

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago"
	"github.com/telemetryapp/goluago/util"
)

const arrayMarkerField = "_is_array"

var errorRegex = regexp.MustCompile(`:([^:]+)+:(.+)$`)

// Exec takes a Lua source code string and set of arguments and executes the code using the go-lua interpreter
func Exec(source string, np notificationProvider, args map[string]interface{}) (map[string]interface{}, error) {
	l := lua.NewState()

	lua.OpenLibraries(l)
	goluago.Open(l)

	openOAuthLibrary(l)
	openJSONLibrary(l)
	openHTTPLibrary(l)
	openStorageLibrary(l)
	openExcelLibrary(l)
	openNotificationsLibrary(l, np)
	openSQLLibrary(l)
	openMongoLibrary(l)
	openXMLLibrary(l)

	util.DeepPush(l, args)

	l.SetGlobal("args")

	util.DeepPush(l, map[string]interface{}{})

	l.SetGlobal("output")

	err := lua.LoadString(l, source)

	if err != nil {
		matches := errorRegex.FindStringSubmatch(lua.CheckString(l, -1))

		if len(matches) != 3 {
			return nil, err
		}

		return nil, fmt.Errorf("Parse error on line %s: %s", matches[1], matches[2])
	}

	err = l.ProtectedCall(0, 0, 0)

	if err != nil {
		matches := errorRegex.FindStringSubmatch(lua.CheckString(l, -1))

		if len(matches) != 3 {
			return nil, err
		}

		return nil, fmt.Errorf("Runtime error on line %s: %s", matches[1], matches[2])
	}

	l.Global("output")

	defer l.Pop(1)

	table, err := util.PullTable(l, 1)

	if err != nil {
		return nil, err
	}

	if output, ok := table.(map[string]interface{}); ok {
		return output, err
	}

	if err == nil {
		return nil, errors.New("The output global has been overwritten with something other than a table.")
	}

	return nil, err
}
