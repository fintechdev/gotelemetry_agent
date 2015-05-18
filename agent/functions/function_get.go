package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("get")
	functionHandlers["$get"] = getHandler
}

func getHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$get", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	counter, _, err := aggregations.GetCounter(context, data["name"].(string))

	if err != nil {
		return nil, err
	}

	return counter.GetValue(), nil
}
