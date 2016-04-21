package database

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/telemetryapp/gotelemetry"
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
func WriteJob(job config.Job) (config.Job, error) {
	// Ensure that all jobs have an ID
	if job.ID == "" {
		if job.Tag == "" {
			return job, gotelemetry.NewError(500, "Job ID missing and no `tag` or `id` provided.")
		}
		job.ID = job.Tag
	}
	// TODO add stronger validation

	jobMarshal, err := json.Marshal(job)

	if err != nil {
		return job, err
	}

	err = manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_jobs"))
		err := bucket.Put([]byte(job.ID), jobMarshal)
		return err
	})

	return job, err
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

func WriteScript(jobID, scriptSource string) error {
	err := manager.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("_scripts"))
		err := bucket.Put([]byte(jobID), []byte(scriptSource))
		return err
	})

	return err
}

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
