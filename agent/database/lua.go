package database

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// GetString will search the database for a string with a given key. If the key
// is not found then the function will return an empty string
func GetString(key string) string {
	var stringVal string

	manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_strings"))
		val := bucket.Get([]byte(key))

		// Return an empty string if the variable cannot be found in the database
		if val == nil {
			stringVal = ""
			return nil
		}

		stringVal = string(val)
		return nil
	})

	return stringVal
}

// SetString will set the value of a string variable under a given key.
// If the user submits an empty string then delete the item from the database
func SetString(key, value string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_strings"))
		if len(value) == 0 {
			return bucket.Delete([]byte(key))
		}
		return bucket.Put([]byte(key), []byte(value))
	})

	return err
}

// GetTable will search the database for a table with a given key. If the key
// is not found then the function will return an empty table
func GetTable(key string) (interface{}, error) {
	var tableVal interface{}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_tables"))
		val := bucket.Get([]byte(key))

		if val == nil {
			val = []byte("{}")
		}

		return json.Unmarshal(val, &tableVal)
	})

	return tableVal, err
}

// SetTable will set the value of a table variable under a given key.
// If the user submits an empty table {} then delete the item from the database
func SetTable(key string, value interface{}) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_tables"))

		data, err := json.Marshal(value)

		if err != nil {
			return err
		}

		if string(data) == "{}" {
			return bucket.Delete([]byte(key))
		}

		return bucket.Put([]byte(key), data)
	})

	return err
}
