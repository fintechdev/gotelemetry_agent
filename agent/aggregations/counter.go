package aggregations

import (
	"sync"
	"time"
)

var Eval func(expr interface{}) (interface{}, error)

type Counter struct {
	Name               string
	value              *int64
	rolloverLast       int64
	rolloverInterval   int64
	rolloverExpression interface{}
	rolloverTimer      *time.Timer
	lock               *sync.Mutex
	saveTimer          *time.Timer
}

var counters = map[string]*Counter{}
var counterLock = sync.Mutex{}

func GetCounter(name string) (*Counter, bool, error) {
	//TODO APP-19
	return nil, true, nil
}

func (c *Counter) fatal(err error) {
	//TODO APP-19
}

func (c *Counter) log(format string, data ...interface{}) {
	//TODO APP-19
}

func (c *Counter) debug(format string, data ...interface{}) {
	//TODO APP-19
}

func (c *Counter) scheduleRollover() {
	//TODO APP-19
}

func (c *Counter) runRollover() {
	//TODO APP-19
}

func (c *Counter) performRollover() {
	//TODO APP-19
}

func (c *Counter) SetRolloverMetadata(interval int64, expression interface{}) {
	//TODO APP-19
}

func (c *Counter) GetValue() int64 {
	//TODO APP-19
	return 0
}

func (c *Counter) save() {
	//TODO APP-19
}

func (c *Counter) updateSaveTimer() {
	//TODO APP-19
}

func (c *Counter) SetValue(newValue int64) {
	//TODO APP-19
}

func (c *Counter) Increment(delta int64) {
	//TODO APP-19
}
