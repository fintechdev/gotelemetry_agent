package database

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// FunctionType TODO
type FunctionType int

const (
	// None TODO
	None FunctionType = iota
	// Sum TODO
	Sum
	// Avg TODO
	Avg
	// Min TODO
	Min
	// Max TODO
	Max
	// Count TODO
	Count
	// StdDev TODO
	StdDev
)

// Series TODO
type Series struct {
	Name string
}

// validateSeriesName TODO
func validateSeriesName(name string) error {
	var seriesNameRegex = regexp.MustCompile(`^[A-Za-z\-][A-Za-z0-9_.\-]*$`)
	if seriesNameRegex.MatchString(name) {
		return nil
	}

	return fmt.Errorf("Invalid series name `%s`. Series names must start with a letter or underscore and can only contain letters, underscores, and digits.", name)
}

// GetSeries TODO
func GetSeries(name string) (*Series, bool, error) {
	isCreated := false

	err := validateSeriesName(name)
	if err != nil {
		return nil, isCreated, err
	}

	// Get the requested key
	err = manager.conn.Update(func(tx *bolt.Tx) error {

		if tx.Bucket([]byte(name)) == nil {
			_, err = tx.CreateBucket([]byte(name))
			if err != nil {
				return err
			}
			isCreated = true
		}

		return nil
	})

	if err != nil {
		return nil, isCreated, err
	}

	series := &Series{
		Name: name,
	}

	return series, isCreated, nil
}

// Push TODO
func (s *Series) Push(timestamp *time.Time, value float64) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {

		if timestamp == nil {
			timestamp = &time.Time{}
			*timestamp = time.Now()
		}

		seriesBucket, err := tx.CreateBucketIfNotExists([]byte(s.Name))
		if err != nil {
			return err
		}

		err = seriesBucket.Put([]byte(strconv.FormatInt(timestamp.Unix(), 10)), []byte(strconv.FormatFloat(value, 'E', -1, 64)))

		return err
	})

	return err
}

// Last TODO
func (s *Series) Last() (map[string]interface{}, error) {

	var output map[string]interface{}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		key, val := tx.Bucket([]byte(s.Name)).Cursor().Last()
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

// Pop TODO
func (s *Series) Pop() (map[string]interface{}, error) {

	var output map[string]interface{}

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(s.Name)).Cursor()
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

// Compute TODO
func (s *Series) Compute(functionType FunctionType, start, end *time.Time) (float64, error) {

	min := []byte(strconv.FormatInt(start.Unix(), 10))
	max := []byte(strconv.FormatInt(end.Unix(), 10))

	var resultArray []float64

	err := manager.conn.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(s.Name)).Cursor()

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

	// Do not compute if there are not any items
	count := float64(len(resultArray))
	if count < 1 {
		return 0.0, nil
	}

	var minValue float64
	var maxValue float64
	var sum float64

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
		// Standard deviation formula requies at least two values
		if count < 2 {
			return 0.0, nil
		}
		var StdDevSum float64
		for _, v := range resultArray {
			StdDevSum += math.Pow((v - avg), 2)
		}
		return math.Sqrt(StdDevSum / (count - 1)), nil
	default:
		return 0.0, fmt.Errorf("Unknown operation %d", functionType)
	}

}

// Aggregate TODO
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
		c := tx.Bucket([]byte(s.Name)).Cursor()

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

			var value float64

			// Only aggregate if there are items in the array
			count := float64(len(resultArray))
			if count > 0 {
				var minValue float64
				var maxValue float64
				var sum float64

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
					// Standard deviation formula requies at least two values
					if count > 1 {
						var StdDevSum float64
						for _, v := range resultArray {
							StdDevSum += math.Pow((v - avg), 2)
						}
						value = math.Sqrt(StdDevSum / (count - 1))
					} else {
						value = 0.0
					}
				default:
					return fmt.Errorf("Unknown operation %d", functionType)
				}
			} else {
				value = 0.0
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

// Items TODO
func (s *Series) Items(count int) (interface{}, error) {
	items := []interface{}{}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(s.Name)).Cursor()

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

// TrimSince TODO
func (s *Series) TrimSince(since time.Time) error {
	max := []byte(strconv.FormatInt(since.Unix(), 10))

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(s.Name)).Cursor()

		// Start by finding the closest value to our trim target
		cursor.Seek(max)
		// Step backwards since we do not want to remove the target
		k, _ := cursor.Prev()

		// Delete all items that take place before this point
		for k != nil {
			err := cursor.Delete()
			if err != nil {
				return err
			}
			k, _ = cursor.Prev()
		}

		return nil
	})

	return err
}

// TrimCount TODO
func (s *Series) TrimCount(count int) error {

	err := manager.conn.Update(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(s.Name)).Cursor()

		k, _ := cursor.Last()

		for i := 1; i <= count; i++ {
			k, _ = cursor.Prev()
			// We do nothing if we hit a nil value before the trim point
			if k != nil {
				return nil
			}
		}

		// Delete all items before the trim point
		for k != nil {
			err := cursor.Delete()
			if err != nil {
				return err
			}
			k, _ = cursor.Prev()
		}

		return nil
	})

	return err
}
