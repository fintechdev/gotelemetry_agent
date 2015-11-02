package lua

import (
	"time"

	"github.com/mtabini/go-lua"
	"github.com/mtabini/goluago/util"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

var seriesFunctions = map[string]func(s *aggregations.Series) lua.Function{
	"name": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			util.DeepPush(l, s.Name)

			return 1
		}
	},

	"trimSince": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {

			if l.TypeOf(1) == lua.TypeString {
				duration, err := config.ParseTimeInterval(lua.CheckString(l, 1))

				if err != nil {
					lua.Errorf(l, "%s", err)
					panic("unreachable")
				}

				since := time.Unix(int64(int(time.Now().Unix())-int(duration/time.Second)), 0)

				if err = s.TrimSince(since); err != nil {
					lua.Errorf(l, "%s", err)
					panic("unreachable")
				}

				return 0
			}

			err := s.TrimSince(time.Unix(int64(lua.CheckInteger(l, 1)), 0))

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			return 0
		}
	},

	"trimCount": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {

			count := lua.CheckInteger(l, 1)

			if err := s.TrimCount(count); err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			return 0
		}
	},

	"push": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			value := lua.CheckNumber(l, 1)
			timestamp := time.Unix(int64(lua.OptInteger(l, 2, int(time.Now().Unix()))), 0)

			if err := s.Push(&timestamp, value); err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			return 0
		}
	},

	"pop": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			res, err := s.Pop(l.ToBoolean(1))

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, res)

			return 1
		}
	},

	"last": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			res, err := s.Last()

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, res)

			return 1
		}
	},

	"compute": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			functionType := lua.CheckInteger(l, 1)
			start := time.Unix(int64(lua.CheckInteger(l, 2)), 0)
			end := time.Unix(int64(lua.CheckInteger(l, 3)), 0)

			res, err := s.Compute(aggregations.FunctionType(functionType), &start, &end)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			util.DeepPush(l, res)

			return 1
		}
	},

	"aggregate": func(s *aggregations.Series) lua.Function {
		return func(l *lua.State) int {
			functionType := lua.CheckInteger(l, 1)
			interval := lua.CheckInteger(l, 2)
			count := lua.CheckInteger(l, 3)
			end := time.Unix(int64(lua.CheckInteger(l, 4)), 0)

			res, err := s.Aggregate(aggregations.FunctionType(functionType), interval, count, &end)

			if err != nil {
				lua.Errorf(l, "%s", err)
				panic("unreachable")
			}

			if arr, ok := res.([]interface{}); ok {
				l.CreateTable(len(arr), 0)

				l.NewTable()
				l.PushBoolean(true)
				l.SetField(-2, arrayMarkerField)
				l.SetMetaTable(-2)

				extractor := func(field string) lua.Function {
					return func(l *lua.State) int {
						l.CreateTable(len(arr), 0)

						l.NewTable()
						l.PushBoolean(true)
						l.SetField(-2, arrayMarkerField)
						l.SetMetaTable(-2)

						for index, value := range arr {
							util.DeepPush(l, value.(map[string]interface{})[field])
							l.RawSetInt(-2, index+1)
						}

						return 1
					}
				}

				l.PushGoFunction(extractor("ts"))
				l.SetField(-2, "ts")

				l.PushGoFunction(extractor("value"))
				l.SetField(-2, "values")

				for index, value := range arr {
					util.DeepPush(l, value)
					l.RawSetInt(-2, index+1)
				}
			} else {
				l.PushNil()
			}

			return 1
		}
	},
}

func pushSeries(l *lua.State, name string) {
	series, _, err := aggregations.GetSeries(name)

	if err != nil {
		lua.Errorf(l, "%s", err)
		panic("unreachable")
	}

	l.NewTable()

	for name, fn := range seriesFunctions {
		l.PushGoFunction(fn(series))
		l.SetField(-2, name)
	}
}
