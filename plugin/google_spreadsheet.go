package plugin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/evanphx/json-patch"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// init() registers this plugin with the Plugin Manager.
// The plugin provides an GoogleSpreadsheetPluginFactory that the manager calls whenever it needs
// to create a new job
func init() {
	job.RegisterPlugin("com.telemetryapp.google_spreadsheet", GoogleSpreadsheetPluginFactory)
}

// Func GoogleSpreadsheetPluginFactory generates a blank plugin instance
func GoogleSpreadsheetPluginFactory() job.PluginInstance {
	return &GoogleSpreadsheetPlugin{
		PluginHelper: job.NewPluginHelper(),
	}
}

// Struct GoogleSpreadsheetPlugin allows you to extract data from an GoogleSpreadsheet sheet and use it to
// populate a Telemetry flow.
type GoogleSpreadsheetPlugin struct {
	*job.PluginHelper
	spreadsheetId  string
	worksheetIndex int
	refresh        time.Duration
	cells          []spreadsheetCellReference
	patch          string
	flowTag        string
	variant        string
	flow           *gotelemetry.Flow
}

type cellField struct {
	XMLName      xml.Name `xml:"cell"`
	NumericValue string   `xml:"numericValue,attr"`
	Value        string   `xml:",innerxml"`
}

type entryField struct {
	XMLName xml.Name  `xml:"entry"`
	Id      string    `xml:"id"`
	Cell    cellField `xml:"cell"`
}

type atomFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Entries []entryField `xml:"entry"`
}

// Init initializes the plugin.
//
// The required configuration parameters are:
//
// - spreadsheet_id               The path to Google Spreadsheet Id, e.g. 1gtvewgRXv3YpcduvJQXiD0Db0L3q05113jGKYJv4De8
//
// - worksheet_index              The index of the worksheet. The first worksheet is indexed as 0.
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
// be replaced by the data extracted from your GoogleSpreadsheet sheet at runtime. You can also
// use $$n as a placeholder that will be replaced by an individual value extracted from the
// sheet.
func (p *GoogleSpreadsheetPlugin) Init(job *job.Job) error {
	var err error

	c := job.Config()

	p.spreadsheetId = c["spreadsheet_id"].(string)

	if worksheetIndex, ok := c["worksheet_index"].(int64); ok {
		p.worksheetIndex = int(worksheetIndex)
	}
	if worksheetIndex, ok := c["worksheet_index"].(int); ok {
		p.worksheetIndex = worksheetIndex
	}

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
	} else {
		p.PluginHelper.AddTaskWithClosure(p.performAllTasks, 0)
	}

	return nil
}

func (p *GoogleSpreadsheetPlugin) performAllTasks(j *job.Job) {
	j.Log("Starting GoogleSpreadsheet plugin...")

	defer p.PluginHelper.TrackTime(j, time.Now(), "GoogleSpreadsheet plugin completed in %s.")

	res, err := http.Get(fmt.Sprintf("https://spreadsheets.google.com/feeds/worksheets/%s/public/full", p.spreadsheetId))
	if err != nil {
		j.ReportError(err)
		return
	}

	rawXML, err := ioutil.ReadAll(res.Body)
	if err != nil {
		j.ReportError(err)
		return
	}

	spreadsheetFeed := atomFeed{}
	err = xml.Unmarshal(rawXML, &spreadsheetFeed)
	if err != nil {
		j.ReportError(err)
		return
	}

	worksheetURL := spreadsheetFeed.Entries[p.worksheetIndex].Id
	worksheetId := strings.Replace(worksheetURL, fmt.Sprintf("https://spreadsheets.google.com/feeds/worksheets/%s/public/full/", p.spreadsheetId), "", 1)

	data := []interface{}{}

	for _, cell := range p.cells {
		cellURL := fmt.Sprintf("https://spreadsheets.google.com/feeds/cells/%s/%s/public/full?min-row=%d&max-row=%d&min-col=%d&max-col=%d ", p.spreadsheetId, worksheetId, cell.Row, cell.Row, cell.Column+1, cell.Column+1)

		res, err := http.Get(cellURL)
		if err != nil {
			j.ReportError(err)
			return
		}

		rawXML, err = ioutil.ReadAll(res.Body)
		if err != nil {
			j.ReportError(err)
			return
		}

		cellFeed := atomFeed{}
		err = xml.Unmarshal(rawXML, &cellFeed)
		if err != nil {
			j.ReportError(err)
			return
		}

		c := cellFeed.Entries[0].Cell

		if c.Value == "FALSE" || c.Value == "TRUE" {
			data = append(data, (c.Value == "TRUE"))
		} else if len(c.NumericValue) > 0 {
			numericValue, err := strconv.ParseFloat(c.NumericValue, 64)
			if err != nil {
				j.ReportError(err)
				return
			}
			data = append(data, numericValue)
		} else {
			data = append(data, c.Value)
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

	if err != nil {
		j.ReportError(err)
		return
	}

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
