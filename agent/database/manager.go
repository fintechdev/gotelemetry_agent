package database

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
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

		if _, err = tx.CreateBucketIfNotExists([]byte("_scripts")); err != nil {
			return err
		}

		return nil
	})

	ttlString := configFile.DatabaseTTL()
	if len(ttlString) == 0 {
		if ttlString = GetConfigParam("ttl"); len(ttlString) > 0 {
			configFile.SetDatabaseTTL(ttlString)
		}
	}

	if len(ttlString) > 0 {
		var ttl time.Duration
		ttl, err = config.ParseTimeInterval(ttlString)
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

			if bytes.HasPrefix(name, []byte("_")) {
				return nil
			}

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

	// Fetch and update the API token
	apiToken := configFile.APIToken()
	if len(apiToken) == 0 {
		if apiToken = GetConfigParam("api_token"); len(apiToken) > 0 {
			configFile.SetAPIToken(apiToken)
		}
	}

	authKey := configFile.AuthKey()
	if len(authKey) == 0 {
		if authKey = GetConfigParam("auth_key"); len(authKey) > 0 {
			configFile.SetAuthKey(authKey)
		} else if len(apiToken) == 0 {
			// If both the API token and Auth token have not been set then generate a random auth token
			rand.Seed(time.Now().UnixNano())
			const alphanum = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
			bytes := make([]byte, 16)
			for i := range bytes {
				bytes[i] = alphanum[rand.Intn(len(alphanum))]
			}

			authKey = string(bytes)
			configFile.SetAuthKey(authKey)
		}
	}

	portNumber := configFile.Listen()
	if len(portNumber) == 0 {
		if portNumber = GetConfigParam("listen"); len(portNumber) > 0 {
			configFile.SetListen(portNumber)
		}
	}

	UDPListenPort := configFile.GraphiteConfig().UDPListenPort
	if len(UDPListenPort) == 0 {
		if UDPListenPort = GetConfigParam("listen_udp"); len(UDPListenPort) > 0 {
			configFile.SetListen(UDPListenPort)
		} else { // Only fetch a TCP Port number if there is no UDP port
			TCPListenPort := configFile.GraphiteConfig().TCPListenPort
			if len(TCPListenPort) == 0 {
				if TCPListenPort = GetConfigParam("listen_tcp"); len(TCPListenPort) > 0 {
					configFile.SetListen(TCPListenPort)
				}
			}
		}
	}

	return nil
}
