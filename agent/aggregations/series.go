package aggregations

import (
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
	context *Context
	Name    string
}

var cachedSeries = map[string]*Series{}

var seriesNameRegex = regexp.MustCompile(`^[A-Za-z\-][A-Za-z0-9_.\-]*$`)

func GetSeries(context *Context, name string) (*Series, error) {
	if err := validateSeriesName(name); err != nil {
		return nil, err
	}

	result := &Series{
		context: context,
		Name:    name,
	}

	if _, ok := cachedSeries[name]; !ok {
		if err := createSeries(context, name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *Series) Push(timestamp *time.Time, value float64) error {
	if timestamp == nil {
		timestamp = &time.Time{}
		*timestamp = time.Now()
	}

	return s.exec("INSERT INTO ?? (ts, value) VALUES (?, ?)", *timestamp, value)
}

func (s *Series) last() (map[string]interface{}, error) {
	return s.fetchRow("SELECT rowid, ts, value FROM ?? ORDER BY ts DESC LIMIT 1")
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

	row, err := s.fetchRow("SELECT CAST("+operation+" AS FLOAT) AS result FROM ?? WHERE ts BETWEEN ? AND ?", *start, *end)

	if err != nil {
		return 0.0, err
	}

	result, ok := row["result"].(float64)

	if !ok {
		result = 0.0
	}

	if functionType == StdDev {
		return math.Sqrt(result), nil
	} else {
		return result, nil
	}
}

func (s *Series) Aggregate(functionType FunctionType, interval, count int) (interface{}, error) {
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

	start := int(time.Now().Add(-time.Duration(interval*count)*time.Second).Unix()) / interval * interval

	rs, err := s.query("SELECT (ts - ?) / ? * ? AS interval, "+operation+"(value) AS result FROM ?? WHERE ts >= ? GROUP BY interval", start, interval, interval, start)

	if err != nil && err != io.EOF {
		return nil, err
	}

	rows := map[int]float64{}

	if err != io.EOF {
		for rs.Next() == nil {
			var index int
			var value float64

			err := rs.Scan(&index, &value)

			if err != nil {
				return nil, err
			}

			rows[index] = value
		}
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

func (s *Series) TrimSince(since time.Time) error {
	return s.exec("DELETE FROM ?? WHERE ts < ?", since.Unix())
}

func (s *Series) TrimCount(count int) error {
	return s.exec("DELETE FROM ?? WHERE rowid <= (SELECT rowid FROM ?? ORDER BY ts DESC LIMIT 1 OFFSET ?)", count)
}
