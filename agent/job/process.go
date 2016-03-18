package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/lua"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// Struct ProcessPlugin allows the agent to execute an external process and use its
// output as data that can be fed to the Telemetry API.
//
// For configuration parameters, see the Init() function
type ProcessPlugin struct {
	*PluginHelper
	args         []string
	batch        bool
	expiration   time.Duration
	flow         *gotelemetry.Flow
	flowTag      string
	path         string
	scriptArgs   map[string]interface{}
	template     map[string]interface{}
	templateFile string
	url          string
}

func newInstance(job *Job) (*ProcessPlugin, error) {

	p := &ProcessPlugin{
		PluginHelper: NewPluginHelper(),
	}

	c := job.Config()

	job.Debugf("The configuration is %#v", c)

	var ok bool

	p.flowTag = c.Tag
	fmt.Println("sdfsdfsdfghehrjzzz")

	p.batch = c.Batch

	if ok && p.flowTag != "" {
		return nil, errors.New("You cannot specify both `tag` and `batch` properties.")
	}

	exec := c.Exec
	script := c.Script

	if exec != "" && script != "" {
		return nil, errors.New("You cannot specify both `script` and `exec` properties.")
	}

	if exec != "" {
		p.path = exec
	} else if script != "" {
		p.path = script
	}

	p.url = c.Url

	if p.path == "" && p.url == "" {
		return nil, errors.New("You must specify a `script`, `exec`, or `url` property.")
	}

	if p.path != "" && p.url != "" {
		return nil, errors.New("You cannot provide both `script` or `exec` and `url` properties.")
	}

	p.args = []string{}
	p.scriptArgs = map[string]interface{}{}

	if args, ok := c.Args.([]interface{}); ok {
		for _, arg := range args {
			if a, ok := arg.(string); ok {
				p.args = append(p.args, a)
			} else {
				p.args = append(p.args, fmt.Sprintf("%#v", arg))
			}
		}
	} else if args, ok := c.Args.(map[interface{}]interface{}); ok {
		p.scriptArgs = config.MapTemplate(args).(map[string]interface{})
	} else if args, ok := c.Args.(map[string]interface{}); ok {
		p.scriptArgs = args
	}

	if p.path != "" {
		if _, err := os.Stat(p.path); os.IsNotExist(err) {
			return nil, errors.New("File " + p.path + " does not exist.")
		}

		if path.Ext(p.path) == ".lua" {
			p.url = "tpl://" + p.path
			p.path = ""
		} else {
			if len(p.scriptArgs) != 0 {
				return nil, errors.New("You cannot specify an key/value hash of arguments when executing an external process. Provide an array of arguments instead.")
			}
		}
	}

	if p.url != "" {
		if len(p.args) != 0 {
			return nil, errors.New("You cannot specify an array of arguments when executing a template. Provide a key/value hash instead.")
		}

		if strings.HasPrefix(p.url, "tpl://") {
			URL, err := url.Parse(p.url)

			if err != nil {
				return nil, err
			}

			p.templateFile = URL.Host + URL.Path

			if _, err := os.Stat(p.templateFile); os.IsNotExist(err) {
				return nil, errors.New("Template " + p.templateFile + " does not exist.")
			}
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

			p.PluginHelper.AddTaskWithClosure(p.performAllTasks, timeInterval)
		} else {
			return nil, err
		}
	} else {
		p.PluginHelper.AddTaskWithClosure(p.performAllTasks, 0)
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

func (p *ProcessPlugin) performHTTPTask(j *Job) (string, error) {
	j.Debugf("Retrieving expression from URL `%s`", p.url)

	r, err := http.Get(p.url)

	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	out, err := ioutil.ReadAll(r.Body)

	if r.StatusCode > 399 {
		return string(out), gotelemetry.NewErrorWithFormat(r.StatusCode, "HTTP request failed with status %d", nil, r.StatusCode)
	}

	return string(out), nil
}

func (p *ProcessPlugin) performTemplateTaskLua(j *Job) (string, error) {
	source, err := ioutil.ReadFile(p.templateFile)

	if err != nil {
		return "", err
	}

	output, err := lua.Exec(string(source), j, p.scriptArgs)

	if err != nil {
		return "", err
	}

	out, err := json.Marshal(config.MapTemplate(output))

	return string(out), err
}

func (p *ProcessPlugin) performTemplateTask(j *Job) (string, error) {

	if strings.HasSuffix(p.templateFile, ".lua") {
		return p.performTemplateTaskLua(j)
	}

	return "", fmt.Errorf("Unknown script type for file `%s`", p.templateFile)
}

func (p *ProcessPlugin) performAllTasks(j *Job) {
	j.Debugf("Starting process plugin...")

	defer p.PluginHelper.TrackTime(j, time.Now(), "Process plugin completed in %s.")

	var response string
	var err error

	if p.path != "" {
		response, err = p.performScriptTask(j)
	} else if p.templateFile != "" {
		response, err = p.performTemplateTask(j)
	} else if p.url != "" {
		response, err = p.performHTTPTask(j)
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
