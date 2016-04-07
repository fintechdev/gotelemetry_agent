package job

import (
	"fmt"
	"time"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
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

	apiToken, err := jobConfig.APIToken()

	if err != nil {
		return err
	}

	credentials, err := gotelemetry.NewCredentials(apiToken, jobConfig.APIURL())

	if err != nil {
		return err
	}

	credentials.SetDebugChannel(errorChannel)

	jobManager.credentials = credentials

	submissionInterval := jobConfig.SubmissionInterval()

	if submissionInterval < time.Second {
		errorChannel <- gotelemetry.NewLogError("Submission interval automatically set to 1s. You can change this value by adding a `submission_interval` property to your configuration file.")
		submissionInterval = time.Second
	} else {
		errorChannel <- gotelemetry.NewLogError("Submission interval set to %ds", submissionInterval/time.Second)
	}

	jobManager.submissionInterval = submissionInterval

	jobManager.accountStreams = map[string]*gotelemetry.BatchStream{}

	for _, jobDescription := range jobConfig.Jobs() {

		if err := jobManager.createJob(jobDescription, false); err != nil {
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

func (m *manager) createJob(jobDescription config.Job, wait bool) error {
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

	job, err := newJob(m.credentials, accountStream, jobDescription.ID, jobDescription, m.credentials.DebugChannel, m.jobCompletionChannel, wait)
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

				m.completionChannel <- true
				return
			}
		}
	}
}

// AddJob triggers the createJob function with a marshaled job config file
func AddJob(jobDescription config.Job) error {
	return jobManager.createJob(jobDescription, false)
}

// GetJobs returns a list of all jobs being managed
func GetJobs() {
	for k := range jobManager.jobs {
		fmt.Println("Job ID:", k)
	}
}

// GetJobByID searches using an ID string and returns the job with that ID
func GetJobByID(id string) {
	if foundJob, found := jobManager.jobs[id]; found {
		fmt.Println(foundJob)
	} else {
		fmt.Println("Job not found: ", id)
	}
}

// TerminateJob searches for a job by ID string and stops/deletes it
func TerminateJob(id string) {
	if foundJob, found := jobManager.jobs[id]; found {
		foundJob.instance.terminate()
		fmt.Println("Terminated job: ", id)
	} else {
		fmt.Println("Job not found: ", id)
	}
}

// GetScript gets the source code of a script for the a job by its ID
func GetScript(id string) {
	if foundJob, found := jobManager.jobs[id]; found {
		if foundJob.instance.script != nil {
			fmt.Println("Job script: ", foundJob.instance.script.source)
			return
		}
		fmt.Println("No script set for job: ", id)
	} else {
		fmt.Println("Job not found: ", id)
	}
}

// AddScript creates or updates a script for a job
func AddScript(id string, scriptSource string) {
	if foundJob, found := jobManager.jobs[id]; found {

		// Script already exists. Update
		if foundJob.instance.script != nil {
			foundJob.instance.script.source = scriptSource
			return
		}

		// Do not add a script if there is an executable
		if foundJob.instance.path != "" {
			fmt.Println("An executable already exists so a script cannot be added to : ", id)
			return
		}

		foundJob.instance.script = newScriptFromSource(scriptSource)

		fmt.Println("No script set for job: ", id)
	} else {
		fmt.Println("Job not found: ", id)
	}
}

// DeleteScript removes the script of a job
func DeleteScript(id string) {
	if foundJob, found := jobManager.jobs[id]; found {

		if foundJob.instance.script != nil {
			foundJob.instance.script = nil
			return
		}

	} else {
		fmt.Println("Job not found: ", id)
	}
}

// RunScriptDebug executes a Lua script and returns the result
func RunScriptDebug(id string) {
	if foundJob, found := jobManager.jobs[id]; found {

		if foundJob.instance.script != nil {
			scriptResult, _ := foundJob.instance.script.exec(foundJob)
			fmt.Println(scriptResult)
			return
		}
		fmt.Println("A script has not been set for: ", id)

	} else {
		fmt.Println("Job not found: ", id)
	}
}
