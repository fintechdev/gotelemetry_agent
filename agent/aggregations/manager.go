package aggregations

import (
	"github.com/telemetryapp/gotelemetry"
)

type Manager struct {
	path         string
	ttl          int
	errorChannel chan error
}

var manager *Manager = nil

func Init(location *string, ttl *int, errorChannel chan error) error {
	if location != nil {
		manager = &Manager{
			path:         *location,
			errorChannel: errorChannel,
		}

		if ttl != nil && *ttl > 0 {
			manager.ttl = *ttl
		} else {
			manager.ttl = 0
		}

		c, err := GetContext()

		if err != nil {
			return err
		}

		defer c.Close()

		if err := c.conn.Exec("CREATE TABLE IF NOT EXISTS _counters (name VARCHAR NOT NULL PRIMARY KEY, value INT NOT NULL DEFAULT(0), rollover_last INT NOT NULL, rollover_interval INT NOT NULL DEFAULT(0), rollover_expression VARCHAR)"); err != nil {
			return err
		}

		c.Debugf("Writing data layer database to %s", manager.path)

		return nil
	}

	errorChannel <- gotelemetry.NewLogError("Data Manager -> No `data.path` property provided. The Data Manager will not run.")

	return nil
}
