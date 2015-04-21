package parser

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

type executionContext struct {
	aggregationContext *aggregations.Context
	variables          map[string]interface{}
	output             map[string]interface{}
}

func newExecutionContext(ac *aggregations.Context) *executionContext {
	return &executionContext{
		aggregationContext: ac,
		variables:          map[string]interface{}{},
		output:             map[string]interface{}{},
	}
}
