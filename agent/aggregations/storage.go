package aggregations

import (
	"encoding/json"
	"errors"
	"strings"
)

func InitStorage() {
	c, err := GetContext()

	if err != nil {
		panic(err)
	}

	defer c.Close()

	if err := c.Exec(`CREATE TABLE IF NOT EXISTS "buckets" (key VARCHAR PRIMARY KEY, value VARCHAR)`); err != nil {
		panic(err)
	}
}

func WriteStorage(key string, value interface{}) error {
	key = strings.TrimSpace(key)

	if key == "" {
		return errors.New("Invalid key")
	}

	c, err := GetContext()

	if err != nil {
		return err
	}

	defer c.Close()

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.Exec(`UPDATE "buckets" SET value = ? WHERE key = ?; INSERT INTO "buckets" (key, value) SELECT ?, ? WHERE changes() = 0`, data, key, key, data)
}

func ReadStorage(key string) (map[string]interface{}, error) {
	c, err := GetContext()

	if err != nil {
		return nil, err
	}

	defer c.Close()

	rs, err := c.query(`SELECT value from "buckets" WHERE key = ?`, key)

	if err != nil {
		return nil, err
	}

	if rs != nil {
		defer rs.Close()
	}

	res := map[string]interface{}{}

	if rs.Next() {
		var s string

		if err = rs.Scan(&s); err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(s), &res); err != nil {
			return nil, err
		}
	}

	return res, nil
}
