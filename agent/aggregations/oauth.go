package aggregations

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

func InitOAuthStorage() {
	if err := manager.exec(`CREATE TABLE IF NOT EXISTS "oauth_buckets" (key VARCHAR PRIMARY KEY, value VARCHAR)`); err != nil {
		panic(err)
	}
}

func WriteOAuthToken(key string, token interface{}) error {
	key = strings.TrimSpace(key)

	if key == "" {
		return errors.New("Invalid key")
	}

	data, err := json.Marshal(token)

	if err != nil {
		return err
	}

	return manager.exec(`UPDATE "oauth_buckets" SET value = ? WHERE key = ?; INSERT INTO "oauth_buckets" (key, value) SELECT ?, ? WHERE changes() = 0`, data, key, key, data)
}

func ReadOAuthToken(key string, dest interface{}) error {
	return manager.query(
		func(rs *sql.Rows) error {
			if rs.Next() {
				var s string

				if err := rs.Scan(&s); err != nil {
					return err
				}

				if err := json.Unmarshal([]byte(s), &dest); err != nil {
					return err
				}
			}

			return nil
		},

		`SELECT value from "oauth_buckets" WHERE key = ?`,

		key,
	)
}
