package lua

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

				// Convert basic types into BSON types
				if queryMap, ok := query.(map[string]interface{}); ok {
					query, err = convertTypes(queryMap)
					if err != nil {
						lua.Errorf(l, "%s", err.Error())
						panic("unreachable")
					}
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

func convertTypes(query map[string]interface{}) (interface{}, error) {
	for key, value := range query {
		switch valueType := value.(type) {
		case map[string]interface{}:
			// Recursively call this function if the value is also a map
			var err error
			query[key], err = convertTypes(valueType)
			if err != nil {
				return nil, err
			}
		case string:
			// We are looking for values with a type (#) prefix
			if strings.HasPrefix(valueType, "#") {
				if strings.HasPrefix(valueType, "#Timestamp=") {
					timestampString := strings.TrimPrefix(valueType, "#Timestamp=")
					timestampInt, err := strconv.ParseInt(timestampString, 10, 64)
					if err != nil {
						return nil, err
					}
					query[key] = time.Unix(timestampInt, 0)
				} else if strings.HasPrefix(valueType, "#ObjectId=") {
					objectIDString := strings.TrimPrefix(valueType, "#ObjectId=")
					objectIDBytes, err := hex.DecodeString(objectIDString)
					if err != nil {
						return nil, err
					}
					query[key] = bson.ObjectId(objectIDBytes)
				}
			}
		}

	}
	return query, nil
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
