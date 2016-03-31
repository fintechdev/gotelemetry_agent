package job

import (
	"fmt"
	"time"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// Manager instantiates, tracks, and updates all jobs within the Agent
type Manager struct {
	Jobs                 map[string]*Job
	credentials          gotelemetry.Credentials
	accountStreams       map[string]*gotelemetry.BatchStream
	completionChannel    chan bool
	jobCompletionChannel chan string
	submissionInterval   time.Duration
}

var jobManager *Manager

// Init the JobManager with a config file and established error channels
func Init(jobConfig config.ConfigInterface, errorChannel chan error, completionChannel chan bool) error {
	jobManager = &Manager{
		Jobs:                 map[string]*Job{},
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

	if len(jobManager.Jobs) == 0 {
		errorChannel <- gotelemetry.NewLogError("No jobs are being scheduled.")
		return nil
	}

	go jobManager.monitorDoneChannel()

	return nil
}

func (m *Manager) createJob(jobDescription config.Job, wait bool) error {
	jobID := jobDescription.ID()

	if jobID == "" {
		return gotelemetry.NewError(500, "Job ID missing and no `tag` or `id` provided.")
	}

	if _, found := m.Jobs[jobID]; found {
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

	job, err := newJob(m.credentials, accountStream, jobDescription.ID(), jobDescription, m.credentials.DebugChannel, m.jobCompletionChannel, wait)
	if err != nil {
		return err
	}

	m.Jobs[job.ID] = job
	return nil
}

func (m *Manager) monitorDoneChannel() {
	for {
		select {
		case id := <-m.jobCompletionChannel:
			delete(m.Jobs, id)

			if len(m.Jobs) == 0 {
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
	for k := range jobManager.Jobs {
		fmt.Println("Job ID:", k)
	}
}

// GetJobByID searches using an ID string and returns the job with that ID
func GetJobByID(id string) {
	if foundJob, found := jobManager.Jobs[id]; found {
		fmt.Println(foundJob)
	} else {
		fmt.Println("Job not found: ", id)
	}
}

// TerminateJob searches for a job by ID string and stops/deletes it
func TerminateJob(id string) {
	if foundJob, found := jobManager.Jobs[id]; found {
		foundJob.instance.terminate()
		fmt.Println("Terminated job: ", id)
	} else {
		fmt.Println("Job not found: ", id)
	}
}
