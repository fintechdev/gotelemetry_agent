package aggregations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"sync"
	"sync/atomic"
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

func GetCounter(context *Context, name string) (*Counter, bool, error) {
	counterLock.Lock()

	defer counterLock.Unlock()

	if counter, ok := counters[name]; ok {
		return counter, false, nil
	}

	if err := context.conn.Exec("INSERT OR IGNORE INTO _counters (name, rollover_last) VALUES (?, ?)", name, time.Now().Unix()); err != nil {
		return nil, false, err
	}

	row, err := context.fetchRow("SELECT value, rollover_last, rollover_interval, rollover_expression FROM _counters WHERE name = ?", name)

	if err == nil {
		v := row["value"].(int64)

		counter := &Counter{
			Name:             name,
			value:            &v,
			rolloverLast:     row["rollover_last"].(int64),
			rolloverInterval: row["rollover_interval"].(int64),
			lock:             &sync.Mutex{},
		}

		if expr, ok := row["rollover_expression"].(string); ok {
			if err := json.Unmarshal([]byte(expr), &counter.rolloverExpression); err != nil {
				return nil, false, err
			}
		}

		counters[name] = counter

		counter.lock.Lock()
		counter.scheduleRollover()
		counter.lock.Unlock()

		return counter, true, nil
	}

	return nil, false, err
}

func (c *Counter) fatal(err error) {
	manager.errorChannel <- errors.New(fmt.Sprintf("Counter %s -> %s", c.Name, err))
}

func (c *Counter) log(format string, data ...interface{}) {
	manager.errorChannel <- gotelemetry.NewLogError("Counter %s -> %s", c.Name, fmt.Sprintf(format, data...))
}

func (c *Counter) debug(format string, data ...interface{}) {
	manager.errorChannel <- gotelemetry.NewLogError("Counter %s -> %s", c.Name, fmt.Sprintf(format, data...))
}

func (c *Counter) scheduleRollover() {
	if c.rolloverTimer != nil {
		c.rolloverTimer.Stop()
	}

	if c.rolloverInterval > 0 {
		var nextTick int64 = -1
		last := c.rolloverLast
		now := time.Now().Unix()

		for nextTick < 0 {
			nextTick = (last/c.rolloverInterval*c.rolloverInterval + c.rolloverInterval) - now

			if nextTick < 0 {
				c.runRollover()
				last = nextTick + now
			}
		}

		c.debug("Scheduling rollover for counter %s in %d second(s); next rollover will be at %s", c.Name, nextTick, time.Now().Add(time.Duration(nextTick)*time.Second))

		c.rolloverTimer = time.AfterFunc(time.Duration(nextTick)*time.Second, c.performRollover)
	}
}

func (c *Counter) runRollover() {
	c.debug("Rolling over counter %s", c.Name)

	var v int64 = 0.0

	if c.rolloverExpression != nil {
		if result, err := Eval(c.rolloverExpression); err == nil {
			c.debug("Counter %s rollover expression evaluates to %#v", c.Name, result)

			switch result.(type) {
			case int:
				v = result.(int64)
			case float64:
				v = int64(result.(float64))
			}
		} else {
			c.fatal(errors.New(fmt.Sprintf("Unable to evaluate rollover expression for counter %s: %s", c.Name, err)))
		}
	}

	c.SetValue(v)
	c.rolloverLast = time.Now().Unix()
}

func (c *Counter) performRollover() {
	c.lock.Lock()

	defer c.lock.Unlock()

	c.runRollover()
	c.scheduleRollover()
}

func (c *Counter) SetRolloverMetadata(interval int64, expression interface{}) {
	c.lock.Lock()

	defer c.lock.Unlock()

	c.debug("Setting rollover metadata for counter `%s`: interval %d, expression %#v", c.Name, interval, expression)

	c.rolloverInterval = interval
	c.rolloverExpression = expression

	c.scheduleRollover()
}

func (c *Counter) GetValue() int64 {
	return atomic.LoadInt64(c.value)
}

func (c *Counter) save() {
	c.lock.Lock()

	defer c.lock.Unlock()

	c.saveTimer = nil

	c.debug("Saving counter %s", c.Name)

	v := c.GetValue()
	context, err := GetContext()

	if err != nil {
		c.fatal(err)
		return
	}

	expr, err := json.Marshal(c.rolloverExpression)

	if err != nil {
		c.fatal(err)
	}

	err = context.conn.Exec("UPDATE _counters SET value = ?, rollover_last = ?, rollover_interval = ?, rollover_expression = ?", v, c.rolloverLast, c.rolloverInterval, expr)

	if err != nil {
		c.fatal(err)
	}
}

func (c *Counter) updateSaveTimer() {
	go func() {
		c.lock.Lock()

		defer c.lock.Unlock()

		if c.saveTimer == nil {
			c.saveTimer = time.AfterFunc(time.Second, c.save)
		}
	}()
}

func (c *Counter) SetValue(newValue int64) {
	atomic.StoreInt64(c.value, newValue)

	c.updateSaveTimer()
}

func (c *Counter) Increment(delta int64) {
	atomic.AddInt64(c.value, delta)

	c.updateSaveTimer()
}
