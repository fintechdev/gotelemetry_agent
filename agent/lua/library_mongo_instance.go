package lua

import (
	"github.com/mtabini/go-lua"
	"github.com/mtabini/goluago/util"

	// "github.com/mtabini/goluago/util"
	"gopkg.in/mgo.v2"
)

var mongoCollectionFunctions = map[string]func(c *mgo.Collection) lua.Function{
	"query": func(c *mgo.Collection) lua.Function {
		return func(l *lua.State) int {
			query, err := util.PullTable(l, 1)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			skip := lua.OptInteger(l, 2, 0)
			limit := lua.OptInteger(l, 3, -1)

			result := []map[string]interface{}{}

			q := c.Find(query)

			if skip > 0 {
				q.Skip(skip)
			}

			if limit >= 0 {
				q.Limit(limit)
			}

			err = q.All(&result)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			pushArray(l)

			for index, value := range result {
				util.DeepPush(l, value)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},

	"count": func(c *mgo.Collection) lua.Function {
		return func(l *lua.State) int {
			query, err := util.PullTable(l, 1)

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			skip := lua.OptInteger(l, 2, 0)
			limit := lua.OptInteger(l, 3, -1)

			q := c.Find(query)

			if skip > 0 {
				q.Skip(skip)
			}

			if limit >= 0 {
				q.Limit(limit)
			}

			count, err := q.Count()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
				panic("unreachable")
			}

			l.PushInteger(count)

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

func pushMongoCollection(l *lua.State, db *mgo.Database, name string) {
	c := db.C(name)

	l.NewTable()

	for name, fn := range mongoCollectionFunctions {
		l.PushGoFunction(fn(c))
		l.SetField(-2, name)
	}
}
