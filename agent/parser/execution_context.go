package parser

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

type executionContext struct {
	aggregationContext   *aggregations.Context
	variables            map[string]interface{}
	output               map[string]interface{}
	notificationProvider executionContextNotificationProvider
}

type executionContextNotificationProvider interface {
	SendNotification(n gotelemetry.Notification, channelTag string) bool
}

func newExecutionContext(ac *aggregations.Context, np executionContextNotificationProvider) *executionContext {
	return &executionContext{
		aggregationContext:   ac,
		variables:            map[string]interface{}{},
		output:               map[string]interface{}{},
		notificationProvider: np,
	}
}
