package job

import (
	"time"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

type JobManager struct {
	Jobs                 map[string]*Job
	credentials          gotelemetry.Credentials
	accountStreams       map[string]*gotelemetry.BatchStream
	completionChannel    chan bool
	jobCompletionChannel chan string
	submissionInterval   time.Duration

}

var jobManager *JobManager

func Init(jobConfig config.ConfigInterface, errorChannel chan error, completionChannel chan bool) error {
	jobManager = &JobManager{
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

		if err := jobManager.CreateJob(jobDescription, false); err != nil {
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

func (m *JobManager) CreateJob(jobDescription config.Job, wait bool) error {
	jobId := jobDescription.ID()

	if jobId == "" {
		return gotelemetry.NewError(500, "Job ID missing and no `tag` or `id` provided.")
	}

	if _, found := m.Jobs[jobId]; found {
		return gotelemetry.NewError(500, "Duplicate job `"+jobId+"`")
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

	job, err := newJob(m, m.credentials, accountStream, jobDescription.ID(), jobDescription, m.credentials.DebugChannel, m.jobCompletionChannel, wait)
	if err != nil {
		return err
	}

	m.Jobs[job.ID] = job
	return nil
}

func (m *JobManager) monitorDoneChannel() {
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

func GetJobManager() JobManager {
	return *jobManager

}

func TerminateJob(id string) {
	if foundJob, found := jobManager.Jobs[id]; found {
		foundJob.instance.Terminate(foundJob)
		fmt.Println("Terminated job: ", id)
	} else {
		fmt.Println("Job not found: ", id)
	}
}
