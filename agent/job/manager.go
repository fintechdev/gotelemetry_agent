package job

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"time"
)

type JobManager struct {
	credentials          gotelemetry.Credentials
	accountStream        *gotelemetry.BatchStream
	jobs                 map[string]*Job
	completionChannel    chan bool
	jobCompletionChannel chan string
}

func createJob(manager *JobManager, credentials gotelemetry.Credentials, accountStream *gotelemetry.BatchStream, errorChannel chan error, jobDescription config.Job, jobCompletionChannel chan string, wait bool) (*Job, error) {
	pluginFactory, err := GetPlugin(jobDescription.Plugin())

	if err != nil {
		return nil, err
	}

	pluginInstance := pluginFactory()

	return newJob(manager, credentials, accountStream, jobDescription.ID(), jobDescription, pluginInstance, errorChannel, jobCompletionChannel, wait)
}

func NewJobManager(jobConfig config.ConfigInterface, errorChannel chan error, completionChannel chan bool) (*JobManager, error) {
	result := &JobManager{
		jobs:                 map[string]*Job{},
		completionChannel:    completionChannel,
		jobCompletionChannel: make(chan string),
	}

	apiToken, err := jobConfig.APIToken()

	if err != nil {
		return nil, err
	}

	credentials, err := gotelemetry.NewCredentials(apiToken)

	if err != nil {
		return nil, err
	}

	credentials.SetDebugChannel(errorChannel)

	result.credentials = credentials

	submissionInterval := jobConfig.SubmissionInterval()

	if submissionInterval < time.Second {
		errorChannel <- gotelemetry.NewLogError("Submission interval automatically set to 1s. You can change this value by adding a `submission_interval` property to your configuration file.")
		submissionInterval = time.Second
	}

	accountStream, err := gotelemetry.NewBatchStream(credentials, submissionInterval, errorChannel)

	if err != nil {
		return nil, err
	}

	result.accountStream = accountStream

	for _, jobDescription := range jobConfig.Jobs() {
		jobId := jobDescription.ID()

		if jobId == "" {
			return nil, gotelemetry.NewError(500, "Job ID missing and no `flow_tag` provided.")
		}

		if !config.CLIConfig.Filter.MatchString(jobId) {
			continue
		}

		if config.CLIConfig.ForceRunOnce {
			delete(jobDescription, "refresh")
		}

		job, err := createJob(result, credentials, accountStream, errorChannel, jobDescription, result.jobCompletionChannel, false)

		if err != nil {
			return nil, err
		}

		if err := result.addJob(job); err != nil {
			return nil, err
		}
	}

	if len(result.jobs) == 0 {
		errorChannel <- gotelemetry.NewLogError("No jobs are being scheduled. If the server is active, the Agent will run in server mode only.")
		return nil, nil
	}

	go result.monitorDoneChannel()

	return result, nil
}

func (m *JobManager) addJob(job *Job) error {
	if _, found := m.jobs[job.ID]; found {
		return gotelemetry.NewError(500, "Duplicate job `"+job.ID+"`")
	}

	m.jobs[job.ID] = job

	return nil
}

func (m *JobManager) monitorDoneChannel() {
	for {
		select {
		case id := <-m.jobCompletionChannel:
			delete(m.jobs, id)

			if len(m.jobs) == 0 {
				m.accountStream.Flush()

				m.completionChannel <- true
				return
			}
		}
	}
}
