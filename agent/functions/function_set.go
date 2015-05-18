package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("set")
	functionHandlers["$set"] = setHandler
}

func setHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$set", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	counter, _, err := aggregations.GetCounter(context, data["name"].(string))

	if err != nil {
		return nil, err
	}

	counter.SetValue(int64(data["value"].(float64)))

	return nil, nil
}
