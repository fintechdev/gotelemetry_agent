package database

import (
	"encoding/json"

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
		err := bucket.Put([]byte(job.ID), jobMarshal)
		return err
	})

	return err
}
