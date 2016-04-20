package database

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// Manager maintains the status and state of the agent's database
type Manager struct {
	path           string
	ttl            time.Duration
	errorChannel   chan error
	conn           *bolt.DB
	mutex          sync.RWMutex
	cleanupRunning bool
}

var manager *Manager

// Init the aggregation manager instance
func Init(configFile config.Interface, errorChannel chan error) error {
	location := configFile.DatabasePath()
	ttlString := configFile.DataConfig().TTL

	conn, err := bolt.Open(location, 0644, nil)

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

		if _, err = tx.CreateBucketIfNotExists([]byte("_counters")); err != nil {
			return err
		}

		if _, err = tx.CreateBucketIfNotExists([]byte("_oauth")); err != nil {
			return err
		}

		if _, err = tx.CreateBucketIfNotExists([]byte("_jobs")); err != nil {
			return err
		}

		if _, err = tx.CreateBucketIfNotExists([]byte("_config")); err != nil {
			return err
		}

		return nil
	})

	if len(ttlString) > 0 {
		ttl, err2 := config.ParseTimeInterval(ttlString)
		if err2 != nil {
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
		manager.databaseCleanup()

		// Begin the database trim routine
		ticker := time.NewTicker(ttl)
		go func() {
			for {
				select {
				case <-ticker.C:
					if manager.cleanupRunning {
						log.Printf("The database cleanup process is already running. Skipping execution.")
						continue
					}
					manager.cleanupRunning = true
					manager.databaseCleanup()
					manager.cleanupRunning = false
				}
			}
		}()
	} else {
		manager.ttl = 0
	}

	return err
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

// Errorf sends a formatted error string to the agent's error channel
func (m *Manager) Errorf(format string, v ...interface{}) {
	if m.errorChannel != nil {
		m.errorChannel <- fmt.Errorf("Data Manager -> "+format, v...)
	}
}

func (m *Manager) databaseCleanup() {

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

// MergeDatabaseWithConfigFile takes the Agent config file and stores its values
// into the database. If not set then fetch the values from the database and set in config file
func MergeDatabaseWithConfigFile(configFile config.Interface) error {
	// Add jobs from the config file to the database
	for _, jobDescription := range configFile.Jobs() {
		if _, err := WriteJob(jobDescription); err != nil {
			return err
		}
	}

	// Fetch and update the API token
	apiToken := configFile.APIToken()
	if len(apiToken) > 0 {
		WriteConfigParam("api_token", apiToken)
	} else if apiToken = GetConfigParam("api_token"); len(apiToken) > 0 {
		configFile.SetAPIToken(apiToken)
	}

	authToken := configFile.AuthToken()
	if len(authToken) > 0 {
		WriteConfigParam("auth_token", authToken)
	} else if authToken = GetConfigParam("auth_token"); len(authToken) > 0 {
		configFile.SetAuthToken(authToken)
	}

	return nil
}
