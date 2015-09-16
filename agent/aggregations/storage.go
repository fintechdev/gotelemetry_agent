package aggregations

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

func InitStorage() {
	if err := manager.exec(`CREATE TABLE IF NOT EXISTS "buckets" (key VARCHAR PRIMARY KEY, value VARCHAR)`); err != nil {
		panic(err)
	}
}

func WriteStorage(key string, value interface{}) error {
	key = strings.TrimSpace(key)

	if key == "" {
		return errors.New("Invalid key")
	}

	if _, ok := value.(map[string]interface{}); !ok {
		if _, ok := value.([]interface{}); !ok {
			return errors.New("Only key/value maps or arrays can be written to structured storage.")
		}
	}

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return manager.exec(`UPDATE "buckets" SET value = ? WHERE key = ?; INSERT INTO "buckets" (key, value) SELECT ?, ? WHERE changes() = 0`, data, key, key, data)
}

func ReadStorage(key string) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	err := manager.query(
		func(rs *sql.Rows) error {
			if rs.Next() {
				var s string

				if err := rs.Scan(&s); err != nil {
					return err
				}

				if err := json.Unmarshal([]byte(s), &res); err != nil {
					return err
				}
			}

			return nil
		},

		`SELECT value from "buckets" WHERE key = ?`,

		key,
	)

	return res, err
}
