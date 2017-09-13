package database

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

// WriteOAuthToken writes an OAuth token to the database under a given key string
func WriteOAuthToken(key string, token interface{}) error {
	key = strings.TrimSpace(key)

	if key == "" {
		return errors.New("Invalid key")
	}

	data, err := json.Marshal(token)

	if err != nil {
		return err
	}

	err = manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_oauth"))

		err2 := bucket.Put([]byte(key), data)

		return err2
	})

	return err
}

// ReadOAuthToken takes an OAuth tokens key and returns an un marshaled token object
func ReadOAuthToken(key string, dest interface{}) error {
	key = strings.TrimSpace(key)

	if key == "" {
		return errors.New("Invalid key")
	}

	err := manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_oauth"))

		val := bucket.Get([]byte(key))
		if val != nil {
			if err := json.Unmarshal(val, &dest); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
