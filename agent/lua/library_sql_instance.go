package lua

import (
	"github.com/jmoiron/sqlx"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var sqlInstanceFunctions = map[string]func(c *sqlx.DB) lua.Function{
	"query": func(db *sqlx.DB) lua.Function {
		return func(l *lua.State) int {
			args := l.Top()

			query := lua.CheckString(l, 1)
			params := make([]interface{}, args-1)

			for i := 2; i <= args; i++ {
				params[i-2] = l.ToValue(i)
			}

			result := []map[string]interface{}{}

			rs, err := db.Queryx(query, params...)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			defer rs.Close()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			for rs.Next() {
				row := make(map[string]interface{})
				err := rs.MapScan(row)

				if err != nil {
					lua.Errorf(l, "%s", err.Error())
				}

				result = append(result, row)
			}

			l.CreateTable(len(result), 0)

			pushArray(l)

			for index, value := range result {
				for i, v := range value {
					if vv, ok := v.([]byte); ok {
						value[i] = string(vv)
					}
				}

				util.DeepPush(l, value)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},
}

func pushSQLInstance(l *lua.State, driverName, dataSourceName string) {
	instance, err := sqlx.Open(driverName, dataSourceName)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
	}

	l.NewTable()

	for name, fn := range sqlInstanceFunctions {
		l.PushGoFunction(fn(instance))
		l.SetField(-2, name)
	}
}
