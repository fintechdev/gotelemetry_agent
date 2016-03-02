package job

import (
	"fmt"
	"time"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

type JobManager struct {
	credentials          gotelemetry.Credentials
	accountStreams       map[string]*gotelemetry.BatchStream
	jobs                 map[string]*Job
	completionChannel    chan bool
	jobCompletionChannel chan string
}

var jobManager *JobManager

func Init(jobConfig config.ConfigInterface, errorChannel chan error, completionChannel chan bool) error {
	jobManager = &JobManager{
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

	jobManager.accountStreams = map[string]*gotelemetry.BatchStream{}

	for _, jobDescription := range jobConfig.Jobs() {
		jobId := jobDescription.ID()

		if jobId == "" {
			return gotelemetry.NewError(500, "Job ID missing and no `tag` provided.")
		}

		if !config.CLIConfig.Filter.MatchString(jobId) {
			continue
		}

		channelTag := jobDescription.ChannelTag

		accountStream, ok := jobManager.accountStreams[channelTag]

		if !ok {
			var err error

			accountStream, err = gotelemetry.NewBatchStream(credentials, channelTag, submissionInterval, errorChannel)

			if err != nil {
				return err
			}

			jobManager.accountStreams[channelTag] = accountStream
		}

		job, err := createJob(jobManager, credentials, accountStream, errorChannel, jobDescription, jobManager.jobCompletionChannel, false)

		if err != nil {
			return err
		}

		if err := jobManager.addJob(job); err != nil {
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

func createJob(manager *JobManager, credentials gotelemetry.Credentials, accountStream *gotelemetry.BatchStream, errorChannel chan error, jobDescription config.Job, jobCompletionChannel chan string, wait bool) (*Job, error) {
	pluginFactory, err := GetPlugin(jobDescription.Plugin)

	if err != nil {
		return nil, err
	}

	pluginInstance := pluginFactory()

	return newJob(manager, credentials, accountStream, jobDescription.ID(), jobDescription, pluginInstance, errorChannel, jobCompletionChannel, wait)
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
				for _, accountStream := range m.accountStreams {
					accountStream.Flush()
				}

				m.completionChannel <- true
				return
			}
		}
	}
}

func GetJobsList() {

	for k, _ := range jobManager.jobs {
		fmt.Println("Job ID:", k)
	}
}
