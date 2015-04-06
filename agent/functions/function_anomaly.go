package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
	"math"
	"time"
)

func init() {
	schemas.LoadSchema("anomaly")
	functionHandlers["$anomaly"] = anomalyHandler
}

func anomalyHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$anomaly", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	var start, end *time.Time

	if period, ok := data["period"]; ok {
		if times, ok := period.(map[string]interface{}); ok {
			if ts, ok := times["from"].(float64); ok {
				start = &time.Time{}
				*start = time.Unix(int64(ts), 0)
			}

			if ts, ok := times["to"].(float64); ok {
				end = &time.Time{}
				*end = time.Unix(int64(ts), 0)
			}
		} else if duration, ok := data["period"].(float64); ok {
			start = &time.Time{}
			*start = time.Now().Add(-time.Duration(duration) * time.Second)

			end = &time.Time{}
			*end = time.Now()
		}
	}

	value := data["value"].(float64)

	seriesName := data["series"].(string)

	series, err := aggregations.GetSeries(context, seriesName)

	if err != nil {
		return nil, err
	}

	stdDev, err := series.Compute(aggregations.StdDev, start, end)

	if err != nil {
		return nil, err
	}

	mean, err := series.Compute(aggregations.Avg, start, end)

	if err != nil {
		return nil, err
	}

	return math.Abs(value-mean) > 3*stdDev, nil
}
