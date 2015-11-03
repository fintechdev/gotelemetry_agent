package aggregations

import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func validateSeriesName(name string) error {
	if seriesNameRegex.MatchString(name) {
		return nil
	}

	return errors.New(fmt.Sprintf("Invalid series name `%s`. Series names must start with a letter or underscore and can only contain letters, underscores, and digits.", name))
}

func createSeries(name string) error {
	result := &Series{
		Name: name,
	}

	if err := result.createTable(); err != nil {
		return err
	}

	cachedSeries[name] = result

	return nil
}

func (s *Series) prepQuery(query string) string {
	return strings.Replace(query, "??", `"`+s.Name+`_series"`, -1)
}

func (s *Series) exec(query string, values ...interface{}) error {
	if manager == nil {
		return errors.New("The storage subsystem has not been initialized. Please specify a storage location in your configuration file.")
	}

	return manager.exec(s.prepQuery(query), values...)
}

func (s *Series) query(closure queryClosure, query string, values ...interface{}) error {
	if manager == nil {
		return errors.New("The storage subsystem has not been initialized. Please specify a storage location in your configuration file.")
	}

	return manager.query(closure, s.prepQuery(query), values...)
}

func (s *Series) createTable() error {
	if err := s.exec("CREATE TABLE IF NOT EXISTS ?? (ts INT NOT NULL, value FLOAT)"); err != nil {
		return err
	}

	n := `"` + s.Name + `_index"`

	if err := s.exec("CREATE INDEX IF NOT EXISTS " + n + " ON ?? (ts)"); err != nil {
		return err
	}

	return nil
}
