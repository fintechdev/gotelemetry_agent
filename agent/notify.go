package agent

import (
	"fmt"
	"net/http"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// ProcessNotificationRequest takes and processes a notification when the Agent is executed in notification mode
func ProcessNotificationRequest(configFile *config.File, errorChannel chan error, completionChannel chan bool, notificationChannel string, notificationFlow string, notification gotelemetry.Notification) {
	errorChannel <- gotelemetry.NewLogError("Notification mode is on.")

	apiToken := configFile.APIToken()

	if len(apiToken) == 0 {
		errorChannel <- fmt.Errorf("No API Token found in the configuration file or in the TELEMETRY_API_TOKEN environment variable.")
		completionChannel <- true
		return
	}

	credentials, err := gotelemetry.NewCredentials(apiToken)

	if err != nil {
		errorChannel <- err
		completionChannel <- true

		return
	}

	//credentials.SetDebugChannel(errorChannel)

	if len(notificationChannel) > 0 {
		channel := gotelemetry.NewChannel(notificationChannel)
		err = channel.SendNotification(credentials, notification)
	} else if len(notificationFlow) > 0 {
		err = gotelemetry.SendFlowChannelNotification(credentials, notificationFlow, notification)
	} else {
		err = gotelemetry.NewError(http.StatusBadRequest, "Either channel or flow is required")
	}

	if err != nil {
		errorChannel <- err
	} else {
		errorChannel <- gotelemetry.NewLogError("Notification sent successfully.")
	}

	completionChannel <- true
}
