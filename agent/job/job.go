// Package job provides an interface to a Telemetry Agent job. A job represents the smallest
// unit of work that the agent recognizes
package job

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// Job is a unit of work that the Agent manages. Jobs manage their own
// processes that interact with the Telemetry API
type Job struct {
	id                string                   // The ID of the job
	credentials       gotelemetry.Credentials  // The credentials used by the job
	stream            *gotelemetry.BatchStream // The batch stream used by the job.
	instance          *processPlugin           // The process instance
	errorChannel      chan error               // A channel to which all errors are funneled
	config            config.Job               // The configuration associated with the job
	completionChannel chan string              // To be pinged when the job has finished running so that the manager knows when to quit
}

// newJob creates and starts a new Job
func newJob(credentials gotelemetry.Credentials, stream *gotelemetry.BatchStream, id string, config config.Job, errorChannel chan error, jobCompletionChannel chan string, wait bool) (*Job, error) {
	result := &Job{
		id:                id,
		credentials:       credentials,
		stream:            stream,
		errorChannel:      errorChannel,
		config:            config,
		completionChannel: jobCompletionChannel,
	}

	if wait {
		result.start(true)
	} else {
		go result.start(false)
	}

	return result, nil
}

// start starts a job. It must be executed asychronously in its own goroutine
func (j *Job) start(wait bool) {
	var err error
	j.instance, err = newInstance(j)
	if err != nil {
		j.reportError(errors.New("Error initializing the job `" + j.id + "`"))
		j.reportError(err)
		return

		//TODO Signal failure to manager
	}

	if wait {
		go j.instance.run(j)
	} else {
		j.instance.run(j)
		if j.completionChannel != nil {
			j.completionChannel <- j.id
		}
	}
}

// getConfig retrieves the configuration data associated with this job.
func (j *Job) getConfig() config.Job {
	return j.config
}

// getOrCreateFlow checks with the Telemetry API to see if a flow tag exists.
// Returns the flow layout if present or creates a new tag defined by optional
// template if provided.
func (j *Job) getOrCreateFlow(tag, variant string, template interface{}) (*gotelemetry.Flow, error) {
	f, err := gotelemetry.GetFlowLayoutWithTag(j.credentials, tag)

	if err == nil {
		if f.Variant != variant {
			return nil, errors.New("Flow " + f.Id + " is of type " + f.Variant + " instead of the expected " + variant)
		}

		return f, nil
	}

	if template != nil {
		template = config.MapTemplate(template)

		if template, ok := template.(map[string]interface{}); ok {
			f, err = gotelemetry.NewFlowWithLayout(j.credentials, tag, variant, "gotelemetry_agent", "", "")

			if err != nil {
				return nil, err
			}

			// populates a flow struct with the data that is currently on the server
			// Note that it is not necessary to populate f.Data, as the method will automatically
			// initialize a nil value with the appropriate data structure for the flow's variant.
			err = f.Read(j.credentials)

			if err != nil {
				return nil, err
			}

			err = f.Populate(variant, template)

			if err != nil {
				return nil, err
			}

			err = f.PostUpdate()

			return f, err
		}

		return nil, errors.New("The `template` property is present in the configuration, but is the wrong type.")
	}

	return nil, fmt.Errorf("The flow with the tag `%s` could not be found, and no template was provided to create it. This job will not run.", tag)
}

// queueDataUpdate queues a data update. The update can contain arbitrary data that is
// sent to the API without any client-side validation.
func (j *Job) queueDataUpdate(tag string, data interface{}, updateType gotelemetry.BatchType) {
	j.stream.SendData(tag, data, updateType)
}

// reportError sends a formatted error to the agent's global error log. This should be
// a plugin's preferred error reporting method when running.
func (j *Job) reportError(err error) {
	actualError := errors.New(j.id + ": -> " + err.Error())

	if j.errorChannel != nil {
		j.errorChannel <- actualError
	}
}

// setFlowError sets a given flow to the error state
func (j *Job) setFlowError(tag string, body interface{}) {
	j.debugf("Setting error status on flow %s", tag)

	if err := gotelemetry.SetFlowError(j.credentials, tag, body); err != nil {
		j.reportError(err)
	}
}

// SendNotification pings the Telemetry API with a notification to a particular flow or channel
func (j *Job) SendNotification(notification gotelemetry.Notification, channelTag string, flowTag string) bool {
	var err error

	if len(channelTag) > 0 {
		channel := gotelemetry.NewChannel(channelTag)
		err = channel.SendNotification(j.credentials, notification)
	} else if len(flowTag) > 0 {
		err = gotelemetry.SendFlowChannelNotification(j.credentials, flowTag, notification)
	} else {
		err = gotelemetry.NewError(http.StatusBadRequest, "Either channel or flow is required")
	}

	if err != nil {
		j.reportError(err)
		return true
	}

	return false
}

// log sends data to the agent's global log. It works like log.Log
func (j *Job) log(v ...interface{}) {
	for _, val := range v {
		if j.errorChannel != nil {
			if v, ok := val.(string); ok {
				j.errorChannel <- gotelemetry.NewLogError("%s -> %s", j.id, v)
			} else {
				j.errorChannel <- gotelemetry.NewLogError("%s -> %#v", j.id, val)
			}
		}
	}
}

// logf sends a formatted string to the agent's global log. It works like log.Logf
func (j *Job) logf(format string, v ...interface{}) {
	if j.errorChannel != nil {
		j.errorChannel <- gotelemetry.NewLogError("%s -> %#s", j.id, fmt.Sprintf(format, v...))
	}
}

// debugf sends a formatted string to the agent's debug log, if it exists. It works like log.Logf
func (j *Job) debugf(format string, v ...interface{}) {
	if j.errorChannel != nil {
		j.errorChannel <- gotelemetry.NewDebugError("%s -> %#s", j.id, fmt.Sprintf(format, v...))
	}
}
