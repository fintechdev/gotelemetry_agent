package gotelemetry

import (
	"fmt"
	"net/http"
)

// Notification struct
type Notification struct {
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Duration int    `json:"duration,omitempty"`
	SoundURL string `json:"sound_url,omitempty"`
	FlowTag  string `json:"flow_tag,omitempty"`
}

// NewNotification function
func NewNotification(title, message, icon string, duration int, soundURL string) Notification {
	return Notification{
		Title:    title,
		Message:  message,
		Icon:     icon,
		Duration: duration,
		SoundURL: soundURL,
	}
}

// Channel struct
type Channel struct {
	Tag string
}

// NewChannel function
func NewChannel(tag string) *Channel {
	return &Channel{Tag: tag}
}

// SendNotification function
func (c *Channel) SendNotification(credentials Credentials, notification Notification) error {
	if logger.IsInfo() {
		logger.Info(
			"Sending notification to channel",
			"notification", fmt.Sprintf("%#v", notification),
			"channel", c.Tag,
		)
	}

	req, err := buildRequest(
		"POST",
		credentials,
		"channels/"+c.Tag+"/notifications",
		notification,
	)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(req)

	return err
}

// SendFlowChannelNotification function
func SendFlowChannelNotification(credentials Credentials, flowTag string, notification Notification) error {
	if len(flowTag) == 0 {
		return NewError(http.StatusBadRequest, "flowTag is required")
	}

	if logger.IsInfo() {
		logger.Info(
			"Sending notification to channels of the flow",
			"notification", fmt.Sprintf("%#v", notification),
			"flow", flowTag,
		)
	}

	req, err := buildRequest(
		"POST",
		credentials,
		"flows/"+flowTag+"/notifications",
		notification,
	)
	if err != nil {
		return err
	}

	_, err = sendJSONRequest(req)

	return err
}
