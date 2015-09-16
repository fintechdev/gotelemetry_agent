package aggregations

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
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

func AggregationFunctionTypeFromName(name string) (FunctionType, error) {
	switch strings.ToLower(name) {
	case "avg", "average":
		return Avg, nil

	case "sum":
		return Sum, nil

	case "min", "minimum":
		return Min, nil

	case "max", "maximum":
		return Max, nil

	case "count":
		return Count, nil

	default:
		return None, errors.New(fmt.Sprintf("Unknown aggregation function `%s`", name))
	}
}

type Series struct {
	Name string
}

var cachedSeries = map[string]*Series{}

var seriesNameRegex = regexp.MustCompile(`^[A-Za-z\-][A-Za-z0-9_.\-]*$`)

func GetSeries(name string) (*Series, bool, error) {
	if err := validateSeriesName(name); err != nil {
		return nil, false, err
	}

	result := &Series{
		Name: name,
	}

	created := false

	if _, ok := cachedSeries[name]; !ok {
		if err := createSeries(name); err != nil {
			return nil, false, err
		}

		// if manager.ttl > 0 {
		// 	ticker := time.Tick(time.Second)
		// 	go func() {
		// 		for range ticker {
		// 			result.deleteOldData()
		// 		}
		// 	}()
		// }

		created = true
	}

	return result, created, nil
}

func FindSeries(search string) ([]string, error) {
	res := []string{}

	err := manager.query(
		func(rs *sql.Rows) error {
			var name string

			for rs.Next() {
				if err := rs.Scan(&name); err != nil {
					return err
				}

				res = append(res, strings.TrimSuffix(name, "_series"))
			}

			return rs.Err()
		},

		"SELECT name FROM sqlite_master WHERE type='table' AND name LIKE ?",

		search+"_series",
	)

	return res, err
}

func (s *Series) deleteOldData() {
	s.exec("DELETE FROM ?? WHERE ts < ?", int(time.Now().Add(-manager.ttl).Unix()))
}

func (s *Series) Push(timestamp *time.Time, value float64) error {
	if timestamp == nil {
		timestamp = &time.Time{}
		*timestamp = time.Now()
	}

	return s.exec("INSERT INTO ?? (ts, value) VALUES (?, ?)", (*timestamp).Unix(), value)
}

func (s *Series) last() (map[string]interface{}, error) {
	var result map[string]interface{}

	err := s.query(
		func(rs *sql.Rows) error {
			if rs.Next() {
				var rowid, ts int
				var value float64

				err := rs.Scan(&rowid, &ts, &value)

				if err != nil {
					return err
				}

				result = map[string]interface{}{
					"rowid": rowid,
					"ts":    ts,
					"value": value,
				}
			}

			return nil
		},

		"SELECT rowid, ts, value FROM ?? ORDER BY ts DESC LIMIT 1",
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Series) Last() (map[string]interface{}, error) {
	row, err := s.last()

	if err != nil {
		return nil, err
	}

	delete(row, "rowid")

	return row, nil
}

func (s *Series) Pop(shouldDelete bool) (map[string]interface{}, error) {
	row, err := s.last()

	if err != nil {
		return nil, err
	}

	rowId := row["rowid"]

	delete(row, "rowid")

	if shouldDelete {
		s.exec("DELETE FROM ?? WHERE rowid = ?", rowId)
	}

	return row, nil
}

func (s *Series) Compute(functionType FunctionType, start, end *time.Time) (float64, error) {
	var operation string

	switch functionType {
	case Sum:
		operation = "TOTAL(value)"

	case Avg:
		operation = "AVG(value)"

	case Min:
		operation = "MIN(value)"

	case Max:
		operation = "MAX(value)"

	case Count:
		operation = "COUNT(*)"

	case StdDev:
		avg, err := s.Compute(Avg, start, end)

		if err != nil {
			return 0.0, err
		}

		operation = fmt.Sprintf("AVG((value - %f) * (value - %f))", avg, avg)

	default:
		return 0.0, errors.New(fmt.Sprintf("Unknown operation %d", functionType))
	}

	if start == nil {
		start = &time.Time{}
		*start = time.Unix(0, 0)
	}

	if end == nil {
		end = &time.Time{}
		*end = time.Now()
	}

	var result = 0.0

	err := s.query(
		func(rs *sql.Rows) error {
			if rs.Next() {
				if err := rs.Scan(&result); err != nil {
					return err
				}

				if functionType == StdDev {
					result = math.Sqrt(result)
				}
			}

			return nil
		},

		"SELECT COALESCE(CAST("+operation+" AS FLOAT), 0.0) AS result FROM ?? WHERE ts BETWEEN ? AND ?",

		(*start).Unix(),
		(*end).Unix(),
	)

	return result, err
}

func (s *Series) Aggregate(functionType FunctionType, interval, count int, endTimePtr *time.Time) (interface{}, error) {
	var operation string

	switch functionType {
	case Sum:
		operation = "TOTAL"

	case Avg:
		operation = "AVG"

	case Min:
		operation = "MIN"

	case Max:
		operation = "MAX"

	case Count:
		operation = "COUNT"

	default:
		return nil, errors.New(fmt.Sprintf("Unknown operation %d", functionType))
	}

	var endTime time.Time
	if endTimePtr != nil {
		endTime = *endTimePtr
	} else {
		endTime = time.Now()
	}

	start := int(endTime.Add(-time.Duration(interval*count)*time.Second).Unix()) / interval * interval

	rows := map[int]float64{}

	err := s.query(

		func(rs *sql.Rows) error {
			for rs.Next() {
				var index int
				var value float64

				err := rs.Scan(&index, &value)

				if err != nil {
					return err
				}

				rows[index] = value
			}

			return nil
		},

		"SELECT (ts - ?) / ? * ? AS interval, "+operation+"(value) AS result FROM ?? WHERE ts >= ? GROUP BY interval",

		start,
		interval,
		interval,
		start,
	)

	if err != nil && err != io.EOF {
		return nil, err
	}

	output := []interface{}{}

	for index := 0; index < count; index++ {
		t := index * interval
		ts := start + t

		if value, ok := rows[t]; ok {
			output = append(output, map[string]interface{}{"ts": ts, "value": value})
		} else {
			output = append(output, map[string]interface{}{"ts": ts, "value": value})
		}
	}

	return interface{}(output), nil
}

func (s *Series) Items(count int) (interface{}, error) {
	output := []interface{}{}

	err := s.query(

		func(rs *sql.Rows) error {
			for rs.Next() {
				var ts int
				var value float64

				err := rs.Scan(&ts, &value)

				if err != nil {
					return err
				}

				output = append(output, map[string]interface{}{"ts": ts, "value": value})
			}

			return nil
		},

		"SELECT ts, value FROM ?? WHERE rowid IN (SELECT rowid FROM ?? ORDER BY ts DESC LIMIT ?) ORDER BY ts",

		count,
	)

	if err != nil && err != io.EOF {
		return nil, err
	}

	return output, nil
}

func (s *Series) TrimSince(since time.Time) error {
	return s.exec("DELETE FROM ?? WHERE ts < ?", since.Unix())
}

func (s *Series) TrimCount(count int) error {
	return s.exec("DELETE FROM ?? WHERE rowid <= (SELECT rowid FROM ?? ORDER BY ts DESC LIMIT 1 OFFSET ?)", count)
}
