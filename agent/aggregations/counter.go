package aggregations

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/telemetryapp/gotelemetry"
)

// Counter TODO
type Counter struct {
	Name string
}

// GetCounter TODO
func GetCounter(name string) (*Counter, bool, error) {
	isCreated := false

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_counters"))

		val := bucket.Get([]byte(name))
		// Initialize a counter with a value of 0 if one does not already exist
		if val == nil {
			err := bucket.Put([]byte(name), []byte("0"))
			isCreated = true
			return err
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	counter := &Counter{
		Name: name,
	}

	return counter, isCreated, nil
}

func (c *Counter) fatal(err error) {
	manager.errorChannel <- fmt.Errorf("Counter %s -> %s", c.Name, err)
}

func (c *Counter) log(format string, data ...interface{}) {
	manager.errorChannel <- gotelemetry.NewLogError("Counter %s -> %s", c.Name, fmt.Sprintf(format, data...))
}

func (c *Counter) debug(format string, data ...interface{}) {
	manager.errorChannel <- gotelemetry.NewLogError("Counter %s -> %s", c.Name, fmt.Sprintf(format, data...))
}

// GetValue TODO
func (c *Counter) GetValue() int64 {

	var value string
	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_counters"))
		value = string(bucket.Get([]byte(c.Name)))

		return nil
	})

	if err != nil {
		c.fatal(err)
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.fatal(err)
	}

	return valueInt
}

// SetValue TODO
func (c *Counter) SetValue(newValue int64) {

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_counters"))

		err := bucket.Put([]byte(c.Name), []byte(strconv.FormatInt(newValue, 10)))

		return err
	})

	if err != nil {
		c.fatal(err)
	}

}

// Increment TODO
func (c *Counter) Increment(delta int64) {

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_counters"))

		val := string(bucket.Get([]byte(c.Name)))

		valueInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		incremenetedVal := valueInt + delta

		err = bucket.Put([]byte(c.Name), []byte(strconv.FormatInt(incremenetedVal, 10)))

		return err
	})

	if err != nil {
		c.fatal(err)
	}

}
