package aggregations

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"math"
	"regexp"
	"strconv"
	"time"
)

type FunctionType int

const (
	None FunctionType = iota
	Sum
	Avg
	Min
	Max
	Count
	StdDev
)

type Series struct {
	Name string
}

func validateSeriesName(name string) error {
	var seriesNameRegex = regexp.MustCompile(`^[A-Za-z\-][A-Za-z0-9_.\-]*$`)
	if seriesNameRegex.MatchString(name) {
		return nil
	}

	return errors.New(fmt.Sprintf("Invalid series name `%s`. Series names must start with a letter or underscore and can only contain letters, underscores, and digits.", name))
}

func GetSeries(name string) (*Series, bool, error) {

	err := validateSeriesName(name)
	if err != nil {
		return nil, false, err
	}

	// Get the requested key
	err = manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))

		_, err := bucket.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	series := &Series{
		Name: name,
	}

	return series, true, nil
}

func (s *Series) Push(timestamp *time.Time, value float64) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))

		if timestamp == nil {
			timestamp = &time.Time{}
			*timestamp = time.Now()
		}

		seriesBucket, err := bucket.CreateBucketIfNotExists([]byte(s.Name))
		if err != nil {
			return err
		}

		err = seriesBucket.Put([]byte(strconv.FormatInt(timestamp.Unix(), 10)), []byte(strconv.FormatFloat(value, 'E', -1, 64)))

		return err
	})

	return err
}

func (s *Series) Last() (map[string]interface{}, error) {

	var output map[string]interface{}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		key, val := bucket.Bucket([]byte(s.Name)).Cursor().Last()
		value, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}

		ts, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			return err
		}

		output = map[string]interface{}{
			"ts":    ts,
			"value": value,
		}

		return nil
	})

	return output, err
}

func (s *Series) Pop() (map[string]interface{}, error) {

	var output map[string]interface{}

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		cursor := bucket.Bucket([]byte(s.Name)).Cursor()
		key, val := cursor.Last()

		value, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}

		ts, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			return err
		}

		output = map[string]interface{}{
			"ts":    ts,
			"value": value,
		}

		err = cursor.Delete()
		if err != nil {
			return err
		}

		return nil
	})

	return output, err
}

func (s *Series) Compute(functionType FunctionType, start, end *time.Time) (float64, error) {

	min := []byte(strconv.FormatInt(start.Unix(), 10))
	max := []byte(strconv.FormatInt(end.Unix(), 10))

	var resultArray []float64

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		c := bucket.Bucket([]byte(s.Name)).Cursor()

		// Iterate over the min/max range
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			valueFloat, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				return err
			}
			resultArray = append(resultArray, valueFloat)

		}

		return nil
	})

	if err != nil {
		return 0.0, err
	}

	var minValue float64
	var maxValue float64
	var sum float64
	count := float64(len(resultArray))

	for i, v := range resultArray {
		if i == 0 {
			maxValue = v
			minValue = v
		}

		sum += v
		if v > maxValue {
			maxValue = v
		}

		if v < minValue {
			minValue = v
		}
	}

	avg := (sum / count)

	switch functionType {
	case Sum:
		return sum, nil
	case Avg:
		return avg, nil
	case Min:
		return minValue, nil
	case Max:
		return maxValue, nil
	case Count:
		return count, nil
	case StdDev:
		var StdDevSum float64
		for _, v := range resultArray {
			StdDevSum += math.Pow((v - avg), 2)
		}
		return math.Sqrt(StdDevSum / (count - 1)), nil
	default:
		return 0.0, errors.New(fmt.Sprintf("Unknown operation %d", functionType))
	}

}

func (s *Series) Aggregate(functionType FunctionType, aggregateInterval int, aggregateCount int, endTimePtr *time.Time) (interface{}, error) {

	interval := int64(aggregateInterval)
	count := int64(aggregateCount)

	output := []interface{}{}

	var startTime int64
	var endTime int64

	if endTimePtr != nil {
		endTime = endTimePtr.Unix()
	} else {
		endTime = time.Now().Unix()
	}

	startTime = endTime - (interval * count)

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		c := bucket.Bucket([]byte(s.Name)).Cursor()

		for i := 0; i < aggregateCount; i++ {

			// Offset the min by 1 so that we are not counting the rollover in each iteration
			min := []byte(strconv.FormatInt(startTime+1, 10))
			startTime += interval
			max := []byte(strconv.FormatInt(startTime, 10))

			var resultArray []float64

			// Iterate over the min/max range
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				valueFloat, err := strconv.ParseFloat(string(v), 64)
				if err != nil {
					return err
				}
				resultArray = append(resultArray, valueFloat)
			}

			var minValue float64
			var maxValue float64
			var sum float64
			count := float64(len(resultArray))

			for i, v := range resultArray {
				if i == 0 {
					maxValue = v
					minValue = v
				}

				sum += v
				if v > maxValue {
					maxValue = v
				}

				if v < minValue {
					minValue = v
				}
			}

			avg := (sum / count)

			var value float64

			switch functionType {
			case Sum:
				value = sum
			case Avg:
				value = avg
			case Min:
				value = minValue
			case Max:
				value = maxValue
			case Count:
				value = count
			case StdDev:
				var StdDevSum float64
				for _, v := range resultArray {
					StdDevSum += math.Pow((v - avg), 2)
				}
				value = math.Sqrt(StdDevSum / (count - 1))
			default:
				return errors.New(fmt.Sprintf("Unknown operation %d", functionType))
			}

			output = append(output, map[string]interface{}{"ts": startTime, "value": value})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return interface{}(output), nil
}

func (s *Series) Items(count int) (interface{}, error) {
	items := []interface{}{}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		cursor := bucket.Bucket([]byte(s.Name)).Cursor()

		key, val := cursor.Last()

		for i := 1; i <= count; i++ {

			if key != nil {
				value, err := strconv.ParseFloat(string(val), 64)
				if err != nil {
					return err
				}

				ts, err := strconv.ParseInt(string(key), 10, 64)
				if err != nil {
					return err
				}

				items = append(items, map[string]interface{}{"ts": ts, "value": value})
			}

			key, val = cursor.Prev()
		}

		return nil
	})

	// Reverse the array since we pushed to it backwards
	output := []interface{}{}
	for i := len(items) - 1; i >= 0; i-- {
		output = append(output, items[i])
	}

	return output, err
}

func (s *Series) TrimSince(since time.Time) error {
	min := []byte(strconv.FormatInt(since.Unix(), 10))
	max := []byte(strconv.FormatInt(time.Now().Unix(), 10))

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		c := bucket.Bucket([]byte(s.Name)).Cursor()

		// Iterate over the min/max range
		for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {

			err := c.Delete()
			if err != nil {
				return err
			}

		}
		return nil
	})

	return err
}

func (s *Series) TrimCount(count int) error {

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_series"))
		cursor := bucket.Bucket([]byte(s.Name)).Cursor()

		var k []byte
		k, _ = cursor.Last()

		for i := 1; i <= count; i++ {
			if k != nil {
				err := cursor.Delete()
				if err != nil {
					return err
				}
			}

			k, _ = cursor.Prev()
		}

		return nil
	})

	return err
}
