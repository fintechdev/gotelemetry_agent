package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanphx/json-patch"
	"github.com/tealeg/xlsx"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"strings"
	"time"
)

// init() registers this plugin with the Plugin Manager.
// The plugin provides an ExcelPluginFactory that the manager calls whenever it needs
// to create a new job
func init() {
	job.RegisterPlugin("com.telemetryapp.excel", ExcelPluginFactory)
}

// Func ExcelPluginFactory generates a blank plugin instance
func ExcelPluginFactory() job.PluginInstance {
	return &ExcelPlugin{
		PluginHelper: job.NewPluginHelper(),
	}
}

// Struct ExcelPlugin allows you to extract data from an Excel sheet and use it to
// populate a Telemetry flow.
type ExcelPlugin struct {
	*job.PluginHelper
	filePath string
	refresh  time.Duration
	cells    []spreadsheetCellReference
	patch    string
	flowTag  string
	variant  string
	flow     *gotelemetry.Flow
}

// Init initializes the plugin.
//
// The required configuration parameters are:
//
// - path                         The path to Excel file
//
// - observe                      Whether the plugin should observe the file for changes, and run whenever changes are detected
//
// - source                       The data to be extracted; a comma-separated list of one or more cells (e.g.: “A12”) or one-dimensional cell ranges (e.g.: “A1:A14”). The plugin supports both string and numeric values.
//
// - flow_tag                     The tag of the flow to populate
//
// - variant                      The variant of the flow
//
// - template                     A template that will be used to populate the flow when it is created
//
// - patch                        A JSON Patch payload that describes how the data extracted from the database must be applied to the flow
//
// The patch is executed only once per update. You can use $$# as a placeholder that will
// be replaced by the data extracted from your Excel sheet at runtime. You can also
// use $$n as a placeholder that will be replaced by an individual value extracted from the
// sheet.
func (p *ExcelPlugin) Init(job *job.Job) error {
	var err error

	c := job.Config()

	p.filePath = c["path"].(string)
	p.flowTag = c["flow_tag"].(string)
	p.variant = c["variant"].(string)

	p.cells, err = parseRange(c["source"].(string))

	if err != nil {
		job.ReportError(err)
		return err
	}

	patch, err := json.Marshal(config.MapFromYaml(c["patch"]))

	if err != nil {
		job.ReportError(err)
		return err
	}

	p.patch = string(patch)

	p.flow, err = job.GetOrCreateFlow(p.flowTag, p.variant, c["template"])

	if err != nil {
		return err
	}

	if refresh, ok := c["refresh"]; ok {
		p.PluginHelper.AddTaskWithClosure(p.performAllTasks, time.Duration(refresh.(int))*time.Second)
	} else if ok, observe := c["observe"].(bool); ok && observe {
		p.PluginHelper.AddTaskWithFileObservation(p.performAllTasks, p.filePath)
	} else {
		p.PluginHelper.AddTaskWithClosure(p.performAllTasks, 0)
	}

	return nil
}

func (p *ExcelPlugin) performAllTasks(j *job.Job) {
	j.Log("Starting Excel plugin...")

	defer p.PluginHelper.TrackTime(j, time.Now(), "Excel plugin completed in %s.")

	f, err := xlsx.OpenFile(p.filePath)

	if err != nil {
		j.ReportError(err)
		return
	}

	data := []interface{}{}

	sheet := f.Sheets[0]

	for _, cell := range p.cells {
		c := sheet.Cell(cell.Row, cell.Column)

		switch c.Type() {
		case xlsx.CellTypeBool:
			data = append(data, c.Bool())

		case xlsx.CellTypeNumeric:
			if v, err := c.Float(); err == nil {
				data = append(data, v)
			} else {
				j.ReportError(err)
				return
			}

		case xlsx.CellTypeString:
			data = append(data, c.String())

		default:
			j.ReportError(errors.New(fmt.Sprintf("Unable to handle value of type %d", c.Type())))
		}
	}

	if err := j.ReadFlow(p.flow); err != nil {
		j.ReportError(err)
		return
	}

	doc, err := json.Marshal(p.flow.Data)

	if err != nil {
		j.ReportError(err)
		return
	}

	marshalled, err := json.Marshal(data)

	if err != nil {
		j.ReportError(err)
		return
	}

	patchSource := strings.Replace(p.patch, "$$#", string(marshalled), -1)

	for index, value := range data {
		marshalled, err := json.Marshal(value)

		if err != nil {
			j.ReportError(err)
			return
		}

		patchSource = strings.Replace(
			patchSource,
			fmt.Sprintf("$$%d", index),
			string(marshalled),
			-1,
		)
	}

	patch, err := jsonpatch.DecodePatch([]byte(patchSource))

	doc, err = patch.Apply(doc)

	if err != nil {
		j.ReportError(err)
		return
	}

	err = json.Unmarshal(doc, &p.flow.Data)

	if err != nil {
		j.ReportError(err)
	}

	j.Logf("Posting flow (%s) %s", p.flowTag, p.flow.Id)

	j.PostFlowUpdate(p.flow)
}
