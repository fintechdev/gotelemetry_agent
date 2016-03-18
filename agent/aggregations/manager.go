package aggregations

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"log"
	"strconv"
	"sync"
	"time"
)

type Manager struct {
	path           string
	ttl            time.Duration
	errorChannel   chan error
	conn           *bolt.DB
	mutex          sync.RWMutex
	cleanupRunning bool
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

		// Create default buckets
		err = conn.Update(func(tx *bolt.Tx) error {

			if _, err := tx.CreateBucketIfNotExists([]byte("_counters")); err != nil {
				return err
			}

			if _, err := tx.CreateBucketIfNotExists([]byte("_oauth")); err != nil {
				return err
			}

			return nil
		})

		if ttlString != nil && len(*ttlString) > 0 {
			ttl, err := config.ParseTimeInterval(*ttlString)
			if err != nil {
				return err
			}
			manager.ttl = ttl

			// The cleanup job should run at least once every 24 hours
			timeInterval := ttl
			oneDayInterval, _ := config.ParseTimeInterval("24h")
			if timeInterval > oneDayInterval {
				timeInterval = oneDayInterval
			}

			// Run once initially
			manager.DatabaseCleanup()

			// Begin the database trim routine
			ticker := time.NewTicker(ttl)
			go func() {
				for {
					 select {
						case <- ticker.C:
							if manager.cleanupRunning {
								log.Printf("The database cleanup process is already running. Skipping execution.")
								continue
							}
							manager.cleanupRunning = true
							manager.DatabaseCleanup()
							manager.cleanupRunning = false
						}
					}
			 }()
		} else {
			manager.ttl = 0
		}

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

func (m *Manager) DatabaseCleanup() {
	fmt.Printf("db cleaned")
	since := time.Now().Add(-m.ttl)
	max := []byte(strconv.FormatInt(since.Unix(), 10))

	err := m.conn.Update(func(tx *bolt.Tx) error {

		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {

			if string(string(name)[0]) != "_" {

				cursor := b.Cursor()

				// Start by finding the closest value to our trim target
				cursor.Seek(max)
				// Step backwards since we do not want to remove the target
				k, _ := cursor.Prev()

				// Delete all items that take place before this point
				for k != nil {
					err := cursor.Delete()
					if err != nil {
						return err
					}
					k, _ = cursor.Prev()
				}

			}

			return nil
		})

		return err
	})

	if err != nil {
		m.Errorf("Database Cleanup Error: %s", err)
	}

}
