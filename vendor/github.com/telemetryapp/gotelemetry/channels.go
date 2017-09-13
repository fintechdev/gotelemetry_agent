package gotelemetry

import (
	"net/http"
)

type Notification struct {
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Duration int    `json:"duration,omitempty"`
	SoundURL string `json:"sound_url,omitempty"`
	FlowTag  string `json:"flow_tag,omitempty"`
}

func NewNotification(title, message, icon string, duration int, soundUrl string) Notification {
	return Notification{
		Title:    title,
		Message:  message,
		Icon:     icon,
		Duration: duration,
		SoundURL: soundUrl,
	}
}

type Channel struct {
	Tag string
}

func NewChannel(tag string) *Channel {
	return &Channel{Tag: tag}
}

func (c *Channel) SendNotification(credentials Credentials, notification Notification) error {
	if credentials.DebugChannel != nil {
		credentials.DebugChannel <- NewDebugError("Sending notification %#v to channel %s", notification, c.Tag)
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

func SendFlowChannelNotification(credentials Credentials, flowTag string, notification Notification) error {
	if len(flowTag) == 0 {
		return NewError(http.StatusBadRequest, "flowTag is required")
	}

	if credentials.DebugChannel != nil {
		credentials.DebugChannel <- NewDebugError("Sending notification %#v to channels of the flow %s", notification, flowTag)
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
