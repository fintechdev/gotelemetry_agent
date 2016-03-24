// Package job provides an interface to a Telemetry Agent job. A job represents the smallest
// unit of work that the agent recognizes—the specific of what a job does are left to the
// individual plugin
package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"net/http"
)

type Job struct {
	ID                string                   // The ID of the job
	credentials       gotelemetry.Credentials  // The credentials used by the job
	stream            *gotelemetry.BatchStream // The batch stream used by the job.
	instance          *ProcessPlugin           // The process instance
	errorChannel      chan error               // A channel to which all errors are funneled
	config            config.Job               // The configuration associated with the job
	completionChannel chan string              // To be pinged when the job has finished running, so that the manager knows when to quit
	manager           *JobManager              // The manager that owns this job
}

// newJob creates and starts a new Job
func newJob(manager *JobManager, credentials gotelemetry.Credentials, stream *gotelemetry.BatchStream, id string, config config.Job, errorChannel chan error, jobCompletionChannel chan string, wait bool) (*Job, error) {
	result := &Job{
		ID:                id,
		credentials:       credentials,
		stream:            stream,
		errorChannel:      errorChannel,
		config:            config,
		completionChannel: jobCompletionChannel,
		manager:           manager,
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
		j.ReportError(errors.New("Error initializing the job `" + j.ID + "`"))
		j.ReportError(err)
		return

		//TODO Signal failure to manager
	}

	if wait {
		go j.instance.Run(j)
	} else {
		j.instance.Run(j)
		j.completionChannel <- j.ID
	}
}

// Retrieve the configuration data associated with this job
func (j *Job) Config() config.Job {
	return j.config
}

// GetOrCreateBoard either creates a board based on an exported template, or retrieves it
// if a board with the same name already exists.
//
// The template must be passed in JSON format as a string (you can use gotelemetry/boarddump to
// generate a template based on an existing board).
func (j *Job) GetOrCreateBoard(name, prefix string, templateSource string) (*gotelemetry.Board, error) {
	template := &gotelemetry.ExportedBoard{}

	err := json.Unmarshal([]byte(templateSource), &template)

	if err != nil {
		return nil, err
	}

	return gotelemetry.ImportBoard(j.credentials, name, prefix, template)
}

// CreateFlow creates a new flow.
func (j *Job) CreateFlow(tag string, variant, sourceProvider, filter, params string) (*gotelemetry.Flow, error) {
	return gotelemetry.NewFlowWithLayout(j.credentials, tag, variant, sourceProvider, filter, params)
}

// GetFlowLayout returns the layout of a given flow
func (j *Job) GetFlowLayout(id string) (*gotelemetry.Flow, error) {
	return gotelemetry.GetFlowLayout(j.credentials, id)
}

func (j *Job) GetFlowTagLayout(tag string) (*gotelemetry.Flow, error) {
	return gotelemetry.GetFlowLayoutWithTag(j.credentials, tag)
}

func (j *Job) GetOrCreateFlow(tag, variant string, template interface{}) (*gotelemetry.Flow, error) {
	f, err := j.GetFlowTagLayout(tag)

	if err == nil {
		if f.Variant != variant {
			return nil, errors.New("Flow " + f.Id + " is of type " + f.Variant + " instead of the expected " + variant)
		}

		return f, nil
	}

	if template != nil {
		template = config.MapTemplate(template)

		if template, ok := template.(map[string]interface{}); ok {
			f, err = j.CreateFlow(tag, variant, "gotelemetry_agent", "", "")

			if err != nil {
				return nil, err
			}

			err = j.ReadFlow(f)

			if err != nil {
				return nil, err
			}

			err = f.Populate(variant, template)

			if err != nil {
				return nil, err
			}

			err = j.PostImmediateFlowUpdate(f)

			return f, err
		}

		return nil, errors.New("The `template` property is present in the configuration, but is the wrong type.")
	}

	return nil, errors.New(fmt.Sprintf("The flow with the tag `%s` could not be found, and no template was provided to create it. This job will not run.", tag))
}

// ReadFlow populates a flow struct with the data that is currently on the server
// Note that it is not necessary to populate f.Data, as the method will automatically
// initialize a nil value with the appropriate data structure for the flow's variant.
func (j *Job) ReadFlow(f *gotelemetry.Flow) error {
	return f.Read(j.credentials)
}

// PostFlowUpdate queues a flow update. The method returns immediately, but the update
// will most likely be sent to the Telemetry API at a later point based on the configuration
// of the underlying stream
func (j *Job) PostFlowUpdate(flow *gotelemetry.Flow) {
	j.stream.Send(flow)
}

func (j *Job) PostImmediateFlowUpdate(flow *gotelemetry.Flow) error {
	return flow.PostUpdate()
}

// PostDataUpdate queues a data update. The update can contain arbitrary data that is
// sent to the API without any client-side validation.
func (j *Job) QueueDataUpdate(tag string, data interface{}, updateType gotelemetry.BatchType) {
	j.stream.SendData(tag, data, updateType)
}

// ReportError sends a formatted error to the agent's global error log. This should be
// a plugin's preferred error reporting method when running.
func (j *Job) ReportError(err error) {
	actualError := errors.New(j.ID + ": -> " + err.Error())

	if j.errorChannel != nil {
		j.errorChannel <- actualError
	}
}

func (j *Job) SetFlowError(tag string, body interface{}) {
	j.Debugf("Setting error status on flow %s", tag)

	if err := gotelemetry.SetFlowError(j.credentials, tag, body); err != nil {
		j.ReportError(err)
	}
}

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
		j.ReportError(err)
		return true
	}

	return false
}

// Log sends data to the agent's global log. It works like log.Log
func (j *Job) Log(v ...interface{}) {
	for _, val := range v {
		if j.errorChannel != nil {
			if v, ok := val.(string); ok {
				j.errorChannel <- gotelemetry.NewLogError("%s -> %s", j.ID, v)
			} else {
				j.errorChannel <- gotelemetry.NewLogError("%s -> %#v", j.ID, val)
			}
		}
	}
}

// Logf sends a formatted string to the agent's global log. It works like log.Logf
func (j *Job) Logf(format string, v ...interface{}) {
	if j.errorChannel != nil {
		j.errorChannel <- gotelemetry.NewLogError("%s -> %#s", j.ID, fmt.Sprintf(format, v...))
	}
}

// Debugf sends a formatted string to the agent's debug log, if it exists. It works like log.Logf
func (j *Job) Debugf(format string, v ...interface{}) {
	if j.errorChannel != nil {
		j.errorChannel <- gotelemetry.NewDebugError("%s -> %#s", j.ID, fmt.Sprintf(format, v...))
	}
}
