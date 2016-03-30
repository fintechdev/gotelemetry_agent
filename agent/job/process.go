package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// A simple task closure
type PluginHelperClosure func(job *Job)

type pluginHelperTask func(job *Job, doneChannel chan bool)

// Struct ProcessPlugin allows the agent to execute an external process and use its
// output as data that can be fed to the Telemetry API.
type ProcessPlugin struct {
	args            []string
	batch           bool
	expiration      time.Duration
	flow            *gotelemetry.Flow
	flowTag         string
	path            string
	script          *Script
	template        map[string]interface{}
	tasks           []pluginHelperTask
	closures        []PluginHelperClosure
	jobDoneChannel  chan bool
	taskDoneChannel chan bool
	waitGroup       *sync.WaitGroup
	isRunning       bool
}

func newInstance(job *Job) (*ProcessPlugin, error) {

	p := &ProcessPlugin{
		tasks:           []pluginHelperTask{},
		jobDoneChannel:  make(chan bool, 0),
		taskDoneChannel: make(chan bool, 2),
		waitGroup:       &sync.WaitGroup{},
	}

	c := job.Config()

	job.Debugf("The configuration is %#v", c)

	var ok bool

	p.flowTag = c.Tag
	p.batch = c.Batch

	if ok && p.flowTag != "" {
		return nil, errors.New("You cannot specify both `tag` and `batch` properties.")
	}

	exec := c.Exec
	script := c.Script

	if exec != "" && script != "" {
		return nil, errors.New("You cannot specify both `script` and `exec` properties.")
	}

	if exec == "" && script == "" {
		return nil, errors.New("You must specify a `script` or `exec` property.")
	}

	p.args = []string{}
	scriptArgs := map[string]interface{}{}

	if args, ok := c.Args.([]interface{}); ok {
		for _, arg := range args {
			if a, ok := arg.(string); ok {
				p.args = append(p.args, a)
			} else {
				p.args = append(p.args, fmt.Sprintf("%#v", arg))
			}
		}
	} else if args, ok := c.Args.(map[interface{}]interface{}); ok {
		scriptArgs = config.MapTemplate(args).(map[string]interface{})
	} else if args, ok := c.Args.(map[string]interface{}); ok {
		scriptArgs = args
	}

	if exec != "" {
		p.path = exec

		if _, err := os.Stat(p.path); os.IsNotExist(err) {
			return nil, errors.New("File " + p.path + " does not exist.")
		}

		if len(scriptArgs) != 0 {
			return nil, errors.New("You cannot specify an key/value hash of arguments when executing an external process. Provide an array of arguments instead.")
		}

	} else if script != "" {
		var err error
		p.script, err = newScript(c.Script, scriptArgs)

		if err != nil {
			return nil, err
		}

	}

	template := c.Template
	variant := c.Variant

	if variant != "" && template != nil {
		if p.flowTag == "" {
			return nil, errors.New("The required `tag` property (`string`) is either missing or of the wrong type.")
		}

		if f, err := job.GetOrCreateFlow(p.flowTag, variant, template); err != nil {
			return nil, err
		} else {
			p.flow = f
		}

	}

	if c.Interval != "" {
		if timeInterval, err := config.ParseTimeInterval(c.Interval); err == nil {
			if p.expiration == 0 {
				p.expiration = timeInterval * 3.0
			}

			p.AddTaskWithClosure(p.performAllTasks, timeInterval)
		} else {
			return nil, err
		}
	} else {
		p.AddTaskWithClosure(p.performAllTasks, 0)
	}

	if p.expiration > 0 && p.expiration < time.Second*60 {
		p.expiration = time.Second * 60
	}

	switch c.Expiration.(type) {
	case int:
		p.expiration = time.Duration(c.Expiration.(int)) * time.Second

	case int64:
		p.expiration = time.Duration(c.Expiration.(int64)) * time.Second

	case string:
		if timeInterval, err := config.ParseTimeInterval(c.Expiration.(string)); err == nil {
			p.expiration = timeInterval
		} else {
			return nil, errors.New("Invalid expiration value. Must be either a number of seconds or a time interval string.")
		}
	}

	if p.expiration < 0 {
		return nil, errors.New("Invalid expiration time")
	}

	if p.expiration > 0 {
		job.Debugf("Expiration is set to %dµs", p.expiration)
	} else {
		job.Debugf("Expiration is off.")
	}

	return p, nil
}

func (p *ProcessPlugin) performDataUpdate(j *Job, flowTag string, isReplace bool, data map[string]interface{}) {

	if config.CLIConfig.DebugMode == true {
		// Debug Mode. Print data dump. Do not send API update
		jsonOutput, err := json.MarshalIndent(data, "", "  ")

		if err != nil {
			return
		}

		fmt.Printf("\nPrinting the output results of \"%s\":\n", flowTag)
		fmt.Println(string(jsonOutput))
		return
	}

	if isReplace {
		if p.expiration > 0 {
			newExpiration := time.Now().Add(p.expiration)
			newUnixExpiration := newExpiration.Unix()

			j.Debugf("Forcing expiration to %d (%s)", newUnixExpiration, newExpiration)

			data["expires_at"] = newUnixExpiration
		}

		j.QueueDataUpdate(flowTag, data, gotelemetry.BatchTypePOST)
	} else {
		if p.expiration > 0 {
			newExpiration := time.Now().Add(p.expiration)
			newUnixExpiration := newExpiration.Unix()

			j.Debugf("Forcing expiration to %d (%s)", newUnixExpiration, newExpiration)

			data["expires_at"] = newUnixExpiration
		}

		j.QueueDataUpdate(flowTag, data, gotelemetry.BatchTypePATCH)
	}
}

func (p *ProcessPlugin) analyzeAndSubmitProcessResponse(j *Job, response string) error {
	isReplace := false

	if strings.HasPrefix(response, "REPLACE\n") {
		isReplace = true
		response = strings.TrimPrefix(response, "REPLACE\n")
	}

	data := map[string]interface{}{}

	for _, command := range strings.Split(response, "\n") {
		command = strings.TrimSpace(command)

		if command == "" {
			continue
		}

		if err := json.Unmarshal([]byte(command), &data); err != nil {
			return err
		}
	}

	if p.batch {
		for key, value := range data {
			if valueMap, ok := value.(map[string]interface{}); ok {
				p.performDataUpdate(j, key, isReplace, valueMap)
			} else {
				return errors.New(fmt.Sprintf("Invalid data for flow %s", key))
			}
		}

		return nil
	}

	if p.flowTag == "" {
		if j.ID != "" {
			// Flow-less job
			return nil
		}
		return errors.New("The required `tag` property (`string`) is either missing or of the wrong type.")
	}

	p.performDataUpdate(j, p.flowTag, isReplace, data)

	return nil
}

func (p *ProcessPlugin) performScriptTask(j *Job) (string, error) {
	if len(p.args) > 0 {
		j.Debugf("Executing `%s` with arguments %#v", p.path, p.args)
	} else {
		j.Debugf("Executing `%s` with no arguments", p.path)
	}

	out, err := exec.Command(p.path, p.args...).Output()

	return string(out), err
}

func (p *ProcessPlugin) performAllTasks(j *Job) {
	j.Debugf("Starting process plugin...")

	defer p.TrackTime(j, time.Now(), "Process plugin completed in %s.")

	var response string
	var err error

	if p.path != "" {
		response, err = p.performScriptTask(j)
	} else if p.script != nil {
		response, err = p.script.exec(j)
	} else {
		err = errors.New("Nothing to do!")
	}

	if err != nil {
		if p.flowTag != "" {
			res := err.Error() + " : " + strings.TrimSpace(string(response))

			if res == "" {
				res = "No output detected."
			}

			j.SetFlowError(p.flowTag, map[string]interface{}{"message": res})
		}

		j.ReportError(err)
		return
	}

	j.Debugf("Process output: %s", strings.Replace(response, "\n", "\\n", -1))

	if p.flowTag != "" {
		j.Debugf("Posting flow %s", p.flowTag)
	}

	if err := p.analyzeAndSubmitProcessResponse(j, response); err != nil {
		j.ReportError(errors.New("Unable to analyze process output: " + err.Error()))
	}
}

// Adds a task to the plugin. The task will be run automarically after the duration specified by
// the interval parameter. Note that interval is measured starting from the end of the last
// execution; therefore, you do not need to worry about conditions like slow networking causing
// successive iterations of a task to “execute over each other.”
func (p *ProcessPlugin) AddTaskWithClosure(c PluginHelperClosure, interval time.Duration) {
	var t pluginHelperTask = nil

	runJob := func(j *Job) {
		p.isRunning = true

		go func(j *Job) {
			c(j)

			p.isRunning = false
		}(j)
	}

	if interval > 0 {
		t = func(job *Job, doneChannel chan bool) {
			runJob(job)

			t := time.NewTicker(interval)

			for {
				select {
				case <-t.C:
					if p.isRunning {
						job.Log("The previous instance of the job is still running; skipping this execution.")
						continue
					}

					runJob(job)

					break
				case <-doneChannel:
					t.Stop()
					return
				}
			}
		}
	}
	p.addTask(t, c)
}

func (p *ProcessPlugin) addTask(t pluginHelperTask, c PluginHelperClosure) {
	if t != nil {
		p.tasks = append(p.tasks, t)
	}

	p.closures = append(p.closures, c)
}

// Run method satisfies the requirements of the PluginInstance interface,
// executing all the tasks asynchronously.
func (p *ProcessPlugin) Run(job *Job) {
	if len(p.tasks) == 0 {
		// Since there are no scheduled tasks, we just run everything once and
		// exit. This makes it possible to schedule a run of the agent through
		// some external mechanism like cron.

		p.RunOnce(job)
		return
	}

	for _, t := range p.tasks {
		p.waitGroup.Add(1)

		go func(t pluginHelperTask) {
			t(job, p.taskDoneChannel)
			p.waitGroup.Done()
		}(t)
	}
	select {
	case <-p.jobDoneChannel:
		return
	}
}

func (p *ProcessPlugin) RunOnce(job *Job) {
	for _, c := range p.closures {
		c(job)
	}
}

// Terminate waits for all outstanding tasks to be completed and then returns.
func (p *ProcessPlugin) Terminate() {
	p.taskDoneChannel <- true
	p.jobDoneChannel <- true
	p.waitGroup.Wait()
}

// TrackTime can be used in a deferred call near the beginning of a function
// to automatically determine how long that function runs for.
func (p *ProcessPlugin) TrackTime(job *Job, start time.Time, template string) {
	job.Logf(template, time.Since(start))
}
