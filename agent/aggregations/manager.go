package aggregations

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

type Manager struct {
	path         string
	ttl          int
	errorChannel chan error
}

var manager *Manager = nil

func Init(cfg config.ConfigInterface, errorChannel chan error) error {
	dataConfig := cfg.DataConfig()

	if dataConfig.DataLocation != nil {
		var ttl int

		if dataConfig.DefaultTTL != nil {
			ttl = *dataConfig.DefaultTTL
		} else {
			ttl = -1
		}

		manager = &Manager{
			path:         *dataConfig.DataLocation,
			ttl:          ttl,
			errorChannel: errorChannel,
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
		c.Debugf("Default data layer TTL is set to %d", manager.ttl)

		return nil
	}

	errorChannel <- gotelemetry.NewLogError("Data Manager -> No `data.path` property provided. The Data Manager will not run.")

	return nil
}
