package aggregations

import (
	"errors"
	"fmt"
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
	//TODO APP-19
	return nil, false, nil
}

func FindSeries(search string) ([]string, error) {
	//TODO APP-19
	return nil, nil
}

func (s *Series) deleteOldData() {
	//TODO APP-19
}

func (s *Series) Push(timestamp *time.Time, value float64) error {
	//TODO APP-19
	return nil
}

func (s *Series) last() (map[string]interface{}, error) {
	//TODO APP-19
	return nil, nil
}

func (s *Series) Last() (map[string]interface{}, error) {
	//TODO APP-19
	return nil, nil
}

func (s *Series) Pop(shouldDelete bool) (map[string]interface{}, error) {
	//TODO APP-19

	return nil, nil
}

func (s *Series) Compute(functionType FunctionType, start, end *time.Time) (float64, error) {
	//TODO APP-19
	var result = 0.0

	return result, nil
}

func (s *Series) Aggregate(functionType FunctionType, interval, count int, endTimePtr *time.Time) (interface{}, error) {
	//TODO APP-19
	output := []interface{}{}

	return interface{}(output), nil
}

func (s *Series) Items(count int) (interface{}, error) {
	//TODO APP-19
	output := []interface{}{}

	return output, nil
}

func (s *Series) TrimSince(since time.Time) error {
	//TODO APP-19
	return nil
}

func (s *Series) TrimCount(count int) error {
	//TODO APP-19
	return nil
}
