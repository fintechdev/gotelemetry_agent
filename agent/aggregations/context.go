package aggregations

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/telemetryapp/gotelemetry"
)

type Context struct {
	conn     *sql.DB
	hasError bool
}

func GetContext() (*Context, error) {
	if manager == nil {
		return nil, errors.New("No data context is available. Did you set the `data.path` property in the Agent's configuration file?")
	}

	conn, err := sql.Open("sqlite3", manager.path)

	if err != nil {
		return nil, err
	}

	result := &Context{
		conn: conn,
	}

	conn.Exec(`
		PRAGMA busy_timeout = 10
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

	return result, nil
}

// Logf sends a formatted string to the agent's global log. It works like log.Logf
func (c *Context) Logf(format string, v ...interface{}) {
	if manager.errorChannel != nil {
		manager.errorChannel <- gotelemetry.NewLogError("Data Manager -> %#s", fmt.Sprintf(format, v...))
	}
}

// Debugf sends a formatted string to the agent's debug log, if it exists. It works like log.Logf
func (c *Context) Debugf(format string, v ...interface{}) {
	if manager.errorChannel != nil {
		manager.errorChannel <- gotelemetry.NewDebugError("Data Manager -> %#s", fmt.Sprintf(format, v...))
	}
}

func (c *Context) Errorf(format string, v ...interface{}) {
	if manager.errorChannel != nil {
		manager.errorChannel <- errors.New(fmt.Sprintf("Data Manager -> "+format, v...))
	}
}

func (c *Context) SetError() {
	c.hasError = true
}

// Data

func (c *Context) Exec(query string, values ...interface{}) error {
	_, err := c.conn.Exec(query, values...)

	return err
}

func (c *Context) query(query string, values ...interface{}) (*sql.Rows, error) {
	return c.conn.Query(query, values...)
}

func (c *Context) Close() error {
	return c.conn.Close()
}
