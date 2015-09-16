package parser

import (
	"github.com/telemetryapp/gotelemetry"
)

type executionContextNotificationProvider interface {
	SendNotification(n gotelemetry.Notification, channelTag string, flowTag string) bool
}

type executionContextJobSpawner interface {
	SpawnJob(id string, plugin string, cfg map[string]interface{}) error
}

type executionContext struct {
	variables            map[string]interface{}
	arguments            map[string]interface{}
	output               map[string]interface{}
	notificationProvider executionContextNotificationProvider
	jobSpawner           executionContextJobSpawner
}

func newExecutionContext(np executionContextNotificationProvider, js executionContextJobSpawner, args map[string]interface{}) *executionContext {
	return &executionContext{
		variables:            map[string]interface{}{},
		arguments:            args,
		output:               map[string]interface{}{},
		notificationProvider: np,
		jobSpawner:           js,
	}
}
