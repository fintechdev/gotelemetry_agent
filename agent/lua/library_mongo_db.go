package lua

import (
	"encoding/json"

	"github.com/globalsign/mgo"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var mongoDBFunctions = map[string]func(db *mgo.Database) lua.Function{
	"collections": func(db *mgo.Database) lua.Function {
		return func(l *lua.State) int {
			names, err := db.CollectionNames()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			pushArray(l)

			for index, value := range names {
				util.DeepPush(l, value)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},

	"collection": func(db *mgo.Database) lua.Function {
		return func(l *lua.State) int {
			pushMongoCollection(l, db, lua.CheckString(l, 1))

			return 1
		}
	},

	"command": func(db *mgo.Database) lua.Function {
		return func(l *lua.State) int {
			var result interface{}

			cmd, err := util.PullTable(l, 1)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			err = db.Run(cmd, &result)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			bytes, err := json.Marshal(&result)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			var resultJSON interface{}
			json.Unmarshal(bytes, &resultJSON)

			util.DeepPush(l, resultJSON)

			return 1
		}
	},
}

func pushMongoDatabase(l *lua.State, s *mgo.Session, dbName string) {
	db := s.DB(dbName)

	l.NewTable()

	for name, fn := range mongoDBFunctions {
		l.PushGoFunction(fn(db))
		l.SetField(-2, name)
	}
}
