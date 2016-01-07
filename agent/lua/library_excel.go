package lua

import (
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
	"github.com/tealeg/xlsx"
)

var excelLibrary = []lua.RegistryFunction{
	{
		"import",
		func(l *lua.State) int {

			path := lua.CheckString(l, 1)

			res, err := xlsx.FileToSlice(path)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, res)
			return 1
		},
	},
}

func openExcelLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, excelLibrary)
		return 1
	}

	lua.Require(l, "telemetry/excel", open, false)
	l.Pop(1)
}
