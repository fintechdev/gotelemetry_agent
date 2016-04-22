package job

import (
	"fmt"
	"time"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
)

// manager instantiates, tracks, and updates all jobs within the Agent
type manager struct {
	jobs                 map[string]*Job
	credentials          gotelemetry.Credentials
	accountStreams       map[string]*gotelemetry.BatchStream
	completionChannel    chan bool
	jobCompletionChannel chan string
	submissionInterval   time.Duration
}

var jobManager *manager

// Init the JobManager with a config file and established error channels
func Init(jobConfig config.Interface, errorChannel chan error, completionChannel chan bool) error {
	jobManager = &manager{
		jobs:                 map[string]*Job{},
		completionChannel:    completionChannel,
		jobCompletionChannel: make(chan string),
	}

	apiToken := jobConfig.APIToken()

	if len(apiToken) == 0 {
		return fmt.Errorf("No API Token found in the configuration file or in the TELEMETRY_API_TOKEN environment variable.")
	}

	credentials, err := gotelemetry.NewCredentials(apiToken, jobConfig.APIURL())

	if err != nil {
		return err
	}

	credentials.SetDebugChannel(errorChannel)

	jobManager.credentials = credentials

	submissionInterval := jobConfig.SubmissionInterval()

	if submissionInterval < time.Second {
		errorChannel <- gotelemetry.NewLogError("Submission interval automatically set to 1s.")
		submissionInterval = time.Second
	} else {
		errorChannel <- gotelemetry.NewLogError("Submission interval set to %ds", submissionInterval/time.Second)
	}

	jobManager.submissionInterval = submissionInterval

	jobManager.accountStreams = map[string]*gotelemetry.BatchStream{}

	// Create each of the jobs listed in the config file
	for _, jobDescription := range jobConfig.Jobs() {

		if err := jobManager.createJob(&jobDescription, false); err != nil {
			return err
		}
		// Job added successfully. Write to database
		if err := database.WriteJob(jobDescription); err != nil {
			return err
		}

	}

	// Fetch jobs located in the database. Do not add jobs already included in the config file
	jobsDatabase, err := database.GetAllJobs()
	for _, jobDescription := range jobsDatabase {
		if _, found := jobManager.jobs[jobDescription.ID]; found {
			continue
		}
		if err := jobManager.createJob(&jobDescription, false); err != nil {
			return err
		}
	}

	if len(jobManager.jobs) == 0 {
		errorChannel <- gotelemetry.NewLogError("No jobs are being scheduled.")
		return nil
	}

	go jobManager.monitorDoneChannel()

	return nil
}

func (m *manager) createJob(jobDescription *config.Job, wait bool) error {
	// Ensure that all jobs have an ID
	if jobDescription.ID == "" {
		if jobDescription.Tag == "" {
			return gotelemetry.NewError(500, "Job ID missing and no `tag` or `id` provided.")
		}
		jobDescription.ID = jobDescription.Tag
	}
	jobID := jobDescription.ID

	if _, found := m.jobs[jobID]; found {
		return gotelemetry.NewError(500, "Duplicate job `"+jobID+"`")
	}

	channelTag := jobDescription.ChannelTag

	accountStream, ok := m.accountStreams[channelTag]

	if !ok {
		var err error

		accountStream, err = gotelemetry.NewBatchStream(m.credentials, channelTag, m.submissionInterval, m.credentials.DebugChannel)

		if err != nil {
			return err
		}

		m.accountStreams[channelTag] = accountStream
	}

	job, err := newJob(m.credentials, accountStream, jobDescription.ID, *jobDescription, m.credentials.DebugChannel, m.jobCompletionChannel, wait)
	if err != nil {
		return err
	}

	m.jobs[job.id] = job
	return nil
}

func (m *manager) monitorDoneChannel() {
	for {
		select {
		case id := <-m.jobCompletionChannel:
			delete(m.jobs, id)

			if len(m.jobs) == 0 {
				for _, accountStream := range m.accountStreams {
					accountStream.Flush()
				}

				if m.completionChannel != nil {
					m.completionChannel <- true
				}
				return
			}
		}
	}
}

// GetJobs returns a list of all jobs being managed
func GetJobs() ([]string, error) {
	var jobsList []string

	if len(jobManager.jobs) == 0 {
		return jobsList, fmt.Errorf("No jobs are scheduled")
	}

	for k := range jobManager.jobs {
		jobsList = append(jobsList, k)
	}

	return jobsList, nil
}

// AddJob triggers the createJob function with a marshaled job config file
func AddJob(jobDescription config.Job) error {
	if err := jobManager.createJob(&jobDescription, false); err != nil {
		return err
	}

	// Job added successfully. Write to database
	err := database.WriteJob(jobDescription)
	return err
}

// GetJobByID searches using an ID string and returns the job with that ID
func GetJobByID(id string) (*config.Job, error) {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return nil, fmt.Errorf("Job not found: %s", id)
	}

	return &foundJob.config, nil
}

// TerminateJob searches for a job by ID string and stops/deletes it
func TerminateJob(id string) error {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return fmt.Errorf("Job not found: %s", id)
	}

	// Delete a script if one exists
	if foundJob.instance.script != nil {
		if err := DeleteScript(foundJob.id); err != nil {
			return err
		}
	}

	delete(jobManager.jobs, id)

	foundJob.instance.terminate()

	err := database.DeleteJob(id)

	return err
}

// ReplaceJob searches for a job by ID string and deletes it and replaces with a new job
func ReplaceJob(jobDescription config.Job) error {
	// Terminate the job if it already exists
	_, found := jobManager.jobs[jobDescription.ID]
	if found {
		if err := TerminateJob(jobDescription.ID); err != nil {
			return err
		}
	}

	if err := AddJob(jobDescription); err != nil {
		return err
	}

	return nil
}

// GetScript gets the source code of a script for the a job by its ID
func GetScript(id string) (string, error) {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return "", fmt.Errorf("Job not found: %s", id)
	}

	if foundJob.instance.script == nil {
		return "", fmt.Errorf("No script set for job: %s", id)
	}

	return foundJob.instance.script.source, nil
}

// AddScript creates or updates a script for a job
func AddScript(id string, scriptSource string) error {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return fmt.Errorf("Job not found: %s", id)
	}

	// Do not add a script if there is an executable
	if foundJob.instance.path != "" {
		return fmt.Errorf("An executable already exists so a script cannot be added to : %s", id)
	}

	database.WriteScript(id, scriptSource)

	// Script already exists. Update
	if foundJob.instance.script != nil {
		foundJob.instance.script.source = scriptSource
		return nil
	}

	// No script has been set. Create a new one
	foundJob.instance.script = newScriptFromSource(scriptSource)
	return nil
}

// DeleteScript removes the script of a job
func DeleteScript(id string) error {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return fmt.Errorf("Job not found: %s", id)
	}

	if foundJob.instance.script == nil {
		return fmt.Errorf("A script has not been set for: %s", id)
	}

	err := database.DeleteScript(id)
	foundJob.instance.script = nil

	return err
}

// RunScriptDebug executes a Lua script and returns the result
func RunScriptDebug(id string) (interface{}, error) {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return nil, fmt.Errorf("Job not found: %s", id)
	}

	if foundJob.instance.script == nil {
		return nil, fmt.Errorf("A script has not been set for: %s", id)
	}

	scriptResult, err := foundJob.instance.script.exec(foundJob)
	if err != nil {
		return nil, err
	}

	return scriptResult, nil
}

// SetScriptState enables or disables the script for a given job ID
func SetScriptState(id string, enableScript bool) error {
	foundJob, found := jobManager.jobs[id]
	if !found {
		return fmt.Errorf("Job not found: %s", id)
	}

	if foundJob.instance.script == nil {
		return fmt.Errorf("A script has not been set for: %s", id)
	}

	if enableScript {
		foundJob.instance.script.enabled = true
		return nil
	}

	foundJob.instance.script.enabled = false
	return nil
}
