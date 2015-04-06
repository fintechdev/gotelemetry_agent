package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("mul")
	functionHandlers["$mul"] = mulHandler
}

func mulHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$mul", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	return data["left"].(float64) * data["right"].(float64), nil
}
