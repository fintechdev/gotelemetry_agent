package lua

import (
	"errors"
	"github.com/mtabini/go-lua"
	"github.com/mtabini/goluago"
	"github.com/mtabini/goluago/util"
)

const arrayMarkerField = "_is_array"

func Exec(source string, np notificationProvider, args map[string]interface{}) (map[string]interface{}, error) {
	l := lua.NewState()

	lua.OpenLibraries(l)
	goluago.Open(l)

	openJSONLibrary(l)
	openHTTPLibrary(l)
	openStorageLibrary(l)
	openNotificationsLibrary(l, np)

	util.DeepPush(l, args)

	l.SetGlobal("args")

	util.DeepPush(l, map[string]interface{}{})

	l.SetGlobal("output")

	err := lua.DoString(l, source)

	l.Global("output")

	defer l.Pop(1)

	table, err := util.PullTable(l, 1)

	if err != nil {
		return nil, err
	}

	if output, ok := table.(map[string]interface{}); ok {
		return output, err
	} else {
		if err == nil {
			return nil, errors.New("The output global has been overwritten with something other than a table.")
		}

		return nil, err
	}
}
