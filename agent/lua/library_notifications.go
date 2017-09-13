package lua

import (
	"time"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/gotelemetry"
)

type notificationProvider interface {
	SendNotification(n gotelemetry.Notification, channelTag string, flowTag string) bool
}

func openNotificationsLibrary(l *lua.State, p notificationProvider) {
	var notificationsLibrary = []lua.RegistryFunction{
		lua.RegistryFunction{
			Name: "post",
			Function: func(l *lua.State) int {
				channelTag := lua.CheckString(l, 1)
				flowTag := lua.CheckString(l, 2)
				title := lua.CheckString(l, 3)
				message := lua.CheckString(l, 4)

				var duration int

				if l.IsNumber(5) {
					duration = lua.CheckInteger(l, 5)
				} else if l.IsString(5) {
					d := lua.CheckString(l, 5)
					dd, err := time.ParseDuration(d)

					if err != nil {
						lua.Errorf(l, "%s", err.Error())
					}

					duration = int(dd.Seconds())
				} else {
					lua.Errorf(l, "Invalid time duration %v", l.ToValue(5))
				}

				if duration < 1 {
					duration = 1
				}

				icon := lua.OptString(l, 6, "")
				sound := lua.OptString(l, 7, "default")

				notification := gotelemetry.NewNotification(title, message, icon, duration, sound)
				l.PushBoolean(p.SendNotification(notification, channelTag, flowTag))

				return 1
			},
		},
	}

	open := func(l *lua.State) int {
		lua.NewLibrary(l, notificationsLibrary)
		return 1
	}

	lua.Require(l, "telemetry/notifications", open, false)
	l.Pop(1)
}
