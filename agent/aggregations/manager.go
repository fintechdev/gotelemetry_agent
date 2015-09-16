package aggregations

import (
	"database/sql"
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
	conn         *sql.DB
	mutex        sync.RWMutex
}

var manager *Manager = nil

func Init(listen, location, ttlString *string, errorChannel chan error) error {
	if location != nil {

		conn, err := sql.Open("sqlite3", *location)

		if err != nil {
			return err
		}

		conn.Exec(`
			PRAGMA busy_timeout = 10;
	    PRAGMA automatic_index = ON;
	    PRAGMA cache_size = 32768;
	    PRAGMA cache_spill = OFF;
	    PRAGMA foreign_keys = ON;
	    PRAGMA journal_size_limit = 67110000;
	    PRAGMA locking_mode = NORMAL;
	    PRAGMA page_size = 4096;
	    PRAGMA recursive_triggers = ON;
	    PRAGMA secure_delete = ON;
	    PRAGMA synchronous = NORMAL;
	    PRAGMA temp_store = MEMORY;
	    PRAGMA journal_mode = WAL;
	    PRAGMA wal_autocheckpoint = 16384;
		`)

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

		if err := manager.exec("CREATE TABLE IF NOT EXISTS _counters (name VARCHAR NOT NULL PRIMARY KEY, value INT NOT NULL DEFAULT(0), rollover_last INT NOT NULL, rollover_interval INT NOT NULL DEFAULT(0), rollover_expression VARCHAR)"); err != nil {
			return err
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

func (m *Manager) exec(query string, values ...interface{}) error {
	// m.mutex.Lock()

	// defer m.mutex.Unlock()

	_, err := m.conn.Exec(query, values...)

	return err
}

type queryClosure func(*sql.Rows) error

func (m *Manager) query(closure queryClosure, query string, values ...interface{}) error {
	// m.mutex.RLock()

	// defer m.mutex.RUnlock()

	rs, err := m.conn.Query(query, values...)

	if err != nil {
		return err
	}

	defer rs.Close()

	return closure(rs)
}
