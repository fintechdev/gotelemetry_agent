package functions

import (
	"errors"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
	"time"
)

func init() {
	schemas.LoadSchema("trim")
	functionHandlers["$trim"] = trimHandler
}

func trimHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$trim", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	series, _, err := aggregations.GetSeries(context, data["series"].(string))

	if err != nil {
		return nil, err
	}

	if since, ok := data["since"].(float64); ok {
		return nil, series.TrimSince(time.Unix(int64(since), 0))
	}

	if count, ok := data["keep"].(float64); ok {
		return nil, series.TrimCount(int(count))
	}

	return nil, errors.New("Unable to process $trim request")
}
