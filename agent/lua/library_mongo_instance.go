package lua

import (
	"encoding/json"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"

	"gopkg.in/mgo.v2"
)

var mongoCollectionFunctions = map[string]func(c *mgo.Collection) lua.Function{
	"find": func(collection *mgo.Collection) lua.Function {
		return func(l *lua.State) int {

			var err error
			var query interface{}

			if l.IsTable(1) {
				query, err = util.PullTable(l, 1)
				if err != nil {
					lua.Errorf(l, "%s", err.Error())
					panic("unreachable")
				}

			}

			skip := lua.OptInteger(l, 2, 0)
			limit := lua.OptInteger(l, 3, -1)

			queryResult := collection.Find(query)

			if skip > 0 {
				queryResult.Skip(skip)
			}

			if limit >= 0 {
				queryResult.Limit(limit)
			}

			pushMongoQuery(l, queryResult)

			return 1
		}
	},
	"name": func(c *mgo.Collection) lua.Function {
		return func(l *lua.State) int {
			l.PushString(c.Name)

			return 1
		}
	},
}

var mongoQueryFunctions = map[string]func(query *mgo.Query) lua.Function{
	"all": func(query *mgo.Query) lua.Function {
		return func(l *lua.State) int {
			var result []interface{}

			err := query.All(&result)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			pushMongoResult(l, &result)

			return 1
		}
	},
	"count": func(query *mgo.Query) lua.Function {
		return func(l *lua.State) int {
			count, err := query.Count()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			l.PushInteger(count)

			return 1
		}
	},
	"distinct": func(query *mgo.Query) lua.Function {
		return func(l *lua.State) int {
			var result []interface{}

			key := lua.CheckString(l, 1)
			err := query.Distinct(key, &result)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			pushMongoResult(l, &result)

			return 1
		}
	},
}

func pushMongoResult(l *lua.State, result *[]interface{}) {
	// Lua only supports basic types. Several BSON types will throw an error
	// so we marshal the query results in JSON format
	bytes, err := json.Marshal(&result)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	var resultJSON []interface{}
	json.Unmarshal(bytes, &resultJSON)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	pushArray(l)

	for index, value := range resultJSON {
		util.DeepPush(l, value)
		l.RawSetInt(-2, index+1)
	}
}

func pushMongoCollection(l *lua.State, db *mgo.Database, name string) {
	c := db.C(name)

	l.NewTable()

	for name, fn := range mongoCollectionFunctions {
		l.PushGoFunction(fn(c))
		l.SetField(-2, name)
	}
}

func pushMongoQuery(l *lua.State, query *mgo.Query) {
	l.NewTable()

	for name, fn := range mongoQueryFunctions {
		l.PushGoFunction(fn(query))
		l.SetField(-2, name)
	}
}
