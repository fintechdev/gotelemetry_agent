package parser

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

type executionContextNotificationProvider interface {
	SendNotification(n gotelemetry.Notification, channelTag string, flowTag string) bool
}

type executionContextJobSpawner interface {
	SpawnJob(id string, plugin string, cfg map[string]interface{}) error
}

type executionContext struct {
	aggregationContext   *aggregations.Context
	variables            map[string]interface{}
	arguments            map[string]interface{}
	output               map[string]interface{}
	notificationProvider executionContextNotificationProvider
	jobSpawner           executionContextJobSpawner
}

func newExecutionContext(ac *aggregations.Context, np executionContextNotificationProvider, js executionContextJobSpawner, args map[string]interface{}) *executionContext {
	return &executionContext{
		aggregationContext:   ac,
		variables:            map[string]interface{}{},
		arguments:            args,
		output:               map[string]interface{}{},
		notificationProvider: np,
		jobSpawner:           js,
	}
}
