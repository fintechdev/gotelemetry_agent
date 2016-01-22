package aggregations

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"sync"
	"time"
)

type Manager struct {
	path         string
	ttl          time.Duration
	errorChannel chan error
	conn         string
	mutex        sync.RWMutex
}

var manager *Manager = nil

func Init(listen, location, ttlString *string, errorChannel chan error) error {
	if location != nil {

		manager = &Manager{
			errorChannel: errorChannel,
			conn:         "",
			mutex:        sync.RWMutex{},
		}

		if ttlString != nil && len(*ttlString) > 0 {
			ttl, err := config.ParseTimeInterval(*ttlString)
			if err != nil {
				return err
			}
			manager.ttl = ttl
		} else {
			manager.ttl = 0
		}

		manager.Debugf("Writing data layer database to %s", manager.path)

		InitStorage()

		if listen != nil {
			InitServer(*listen, errorChannel)
		}

		return nil
	}

	manager.Errorf("Data Manager -> No `data.path` property provided. The Data Manager will not run.")

	return nil
}

// Logf sends a formatted string to the agent's global log. It works like log.Logf
func (m *Manager) Logf(format string, v ...interface{}) {
	if m.errorChannel != nil {
		m.errorChannel <- gotelemetry.NewLogError("Data Manager -> %#s", fmt.Sprintf(format, v...))
	}
}

// Debugf sends a formatted string to the agent's debug log, if it exists. It works like log.Logf
func (m *Manager) Debugf(format string, v ...interface{}) {
	if m.errorChannel != nil {
		m.errorChannel <- gotelemetry.NewDebugError("Data Manager -> %#s", fmt.Sprintf(format, v...))
	}
}

func (m *Manager) Errorf(format string, v ...interface{}) {
	if m.errorChannel != nil {
		m.errorChannel <- errors.New(fmt.Sprintf("Data Manager -> "+format, v...))
	}
}

// Data

func (m *Manager) exec(bucket string, key string, value string) error {

	// Store some data


	return nil
}

func (m *Manager) query(bucket string, key string) ([]byte, error) {

	var val []byte

	// retrieve the data


	return val, nil

}
