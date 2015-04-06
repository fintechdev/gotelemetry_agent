package functions

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("if")
	functionHandlers["$if"] = ifHandler
}

func ifHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$if", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})

	if data["condition"].(bool) {
		return data["then"], nil
	} else {
		return data["else"], nil
	}
}
