package aggregations

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
)

type Context struct {
	conn          *sqlite3.Conn
	hasError      bool
	inTransaction bool
}

func GetContext() (*Context, error) {
	conn, err := sqlite3.Open(manager.path)

	if err != nil {
		return nil, err
	}

	result := &Context{
		conn: conn,
	}

	conn.Exec("PRAGMA journal_mode = WAL")
	conn.Exec("PRAGMA cache_size = 1000000")
	conn.Exec("PRAGMA busy_timeout = 10000")

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

func (c *Context) fetchRow(query string, values ...interface{}) (sqlite3.RowMap, error) {
	result := sqlite3.RowMap{}

	rs, err := c.conn.Query(query, values...)

	if err != nil {
		return result, err
	}

	if rs != nil {
		defer rs.Close()
	}

	rs.Scan(result)

	return result, nil
}

// Transactions

func (c *Context) Begin() error {
	c.inTransaction = true

	return c.conn.Begin()
}

func (c *Context) Commit() error {
	c.inTransaction = false

	return c.conn.Commit()
}

func (c *Context) Rollback() error {
	c.inTransaction = false

	return c.conn.Rollback()
}

func (c *Context) Close() error {
	if c.inTransaction {
		if c.hasError {
			if err := c.conn.Rollback(); err != nil {
				return err
			}
		} else {
			if err := c.conn.Commit(); err != nil {
				return err
			}
		}
	}

	return c.conn.Close()
}
