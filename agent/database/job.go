package database

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// GetAllJobs returns an array of all jobs from the database
func GetAllJobs() ([]config.Job, error) {

	var jobsArray []config.Job

	err := manager.conn.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte("_jobs")).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {

			var fetchedJob config.Job
			if err := json.Unmarshal(v, &fetchedJob); err != nil {
				return nil
			}

			jobsArray = append(jobsArray, fetchedJob)
		}
		return nil
	})

	return jobsArray, err
}

// WriteJob stores a job within the database
func WriteJob(job config.Job) error {
	jobMarshal, err := json.Marshal(job)

	if err != nil {
		return err
	}

	err = manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_jobs"))
		err2 := bucket.Put([]byte(job.ID), jobMarshal)
		return err2
	})

	return err
}

// DeleteJob finds a job by its ID and removes it from the database
func DeleteJob(jobID string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_jobs"))

		// Run a get command first to ensure that the key exists
		if v := bucket.Get([]byte(jobID)); v == nil {
			return fmt.Errorf("Could not find the job in database: %s", jobID)
		}

		err := bucket.Delete([]byte(jobID))

		return err
	})

	return err
}

// WriteScript writes a job's associated Lua source code keyed by the job ID
func WriteScript(jobID, scriptSource string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_scripts"))
		err := bucket.Put([]byte(jobID), []byte(scriptSource))
		return err
	})

	return err
}

// GetScript searches by job ID and returns a string of Lua source code. Returns empty string if not found
func GetScript(jobID string) string {
	var scriptSource string

	manager.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_scripts"))

		val := bucket.Get([]byte(jobID))

		if val != nil {
			scriptSource = string(val)
		}

		return nil
	})

	return scriptSource
}

// DeleteScript searches by job ID for a Lua source code string and deletes the entry. Returns an error if not found
func DeleteScript(jobID string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_scripts"))

		// Run a get command first to ensure that the key exists
		if v := bucket.Get([]byte(jobID)); v == nil {
			return fmt.Errorf("Could not find the script in database for job: %s", jobID)
		}

		err := bucket.Delete([]byte(jobID))

		return err
	})

	return err
}
