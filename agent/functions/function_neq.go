package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("neq")
	functionHandlers["$neq"] = neqHandler
}

func neqHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$neq", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	return data["left"].(float64) != data["right"].(float64), nil
}
