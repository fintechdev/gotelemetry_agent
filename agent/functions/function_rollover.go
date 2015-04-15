package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("rollover")
	functionHandlers["$rollover"] = rolloverHandler
}

func rolloverHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$rollover", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	counter, err := aggregations.GetCounter(context, data["name"].(string))

	if err != nil {
		return nil, err
	}

	var interval int64 = 0
	var expression interface{} = nil

	if i, ok := data["interval"].(float64); ok {
		interval = int64(i)
	}

	if e, ok := data["expression"].(interface{}); ok {
		expression = e
	}

	counter.SetRolloverMetadata(interval, expression)

	return nil, nil
}
