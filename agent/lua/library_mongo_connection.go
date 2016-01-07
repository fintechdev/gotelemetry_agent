package lua

import (
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
	"gopkg.in/mgo.v2"
)

var mongoConnectionFunctions = map[string]func(s *mgo.Session) lua.Function{
	"db": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			pushMongoDatabase(l, s, lua.CheckString(l, 1))

			return 1
		}
	},

	"live_servers": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			pushArray(l)

			for index, server := range s.LiveServers() {
				util.DeepPush(l, server)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},

	"close": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			s.Close()

			return 0
		}
	},
}

func pushGoConnection(l *lua.State, connectionString string) {
	s, err := mgo.Dial(connectionString)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
		panic("unreachable")
	}

	l.NewTable()

	for name, fn := range mongoConnectionFunctions {
		l.PushGoFunction(fn(s))
		l.SetField(-2, name)
	}

	l.CreateTable(0, 1)
	l.PushGoFunction(func(l *lua.State) int {
		s.Close()

		return 0
	})
	l.SetField(-2, "__gc")
	l.SetMetaTable(-2)
}
