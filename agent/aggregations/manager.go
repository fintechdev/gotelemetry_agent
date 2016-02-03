package aggregations

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"sync"
	"time"
	"log"
	"github.com/boltdb/bolt"
)

type Manager struct {
	path         string
	ttl          time.Duration
	errorChannel chan error
	conn         *bolt.DB
	mutex        sync.RWMutex
}

var manager *Manager = nil

func Init(listen, location, ttlString *string, errorChannel chan error) error {
	if location != nil {

		conn, err := bolt.Open(*location, 0644, nil)

		if err != nil {
			return err
		}

		manager = &Manager{
			errorChannel: errorChannel,
			conn:         conn,
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

		// Create default buckets
		err = conn.Update(func(tx *bolt.Tx) error {

			if _, err := tx.CreateBucketIfNotExists([]byte("_counters")); err != nil {
				return err
			}

			if _, err := tx.CreateBucketIfNotExists([]byte("_series")); err != nil {
				return err
			}

			if _, err := tx.CreateBucketIfNotExists([]byte("_oauth")); err != nil {
				return err
			}

			return nil
		})

		return err
	}
	log.Printf("Error: %s", "Data Manager -> No `data.path` property provided. The Data Manager will not run.")

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
