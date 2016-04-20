package database

import "github.com/boltdb/bolt"

// WriteConfigParam takes a key/value pair of a config parameter and writes it to the database
func WriteConfigParam(paramName, paramValue string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_config"))

		err := bucket.Put([]byte(paramName), []byte(paramValue))

		return err
	})

	return err
}

// GetConfigParam takes a key string and searches for/returns the corresponding parameter from the database
func GetConfigParam(paramName string) string {
	var paramValue string
	manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_config"))

		val := bucket.Get([]byte(paramName))
		if val != nil {
			paramValue = string(val)
		}

		return nil
	})

	return paramValue
}
