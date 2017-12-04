// Package gotelemetry provides bindings for the Telemetry visualization service (http://telemetryapp.com).
//
// In order to use the package, you will need to sign up for an account and obtain your API token from
// https://www.telemetryapp.com/user/documentation/authentication. A full set of documents that explain
// how the Telemetry REST API works is also available at https://www.telemetryapp.com/user/documentation/.
//
// The package supports submitting flow data either individually or in batches to minimize network use.
// Please bear in mind that both the Flow and Batch submission mechanisms are not, by default,
// thread-safe; if thread-safety is desired, it must be provided by the caller.
package gotelemetry

import (
	"encoding/json"
	"reflect"
)

// Flow identifies a flow, defined as the combination of a tag and
// the data associated with it, which must be a pointer to one of the structs declared
// in variants.go
//
// Note that Flow structs are not thread-safe by default, because they store a pointer to
// the underlying variant information. If you require thread-safety, you must mediate access
// to the flow through a synchronization mechanism of some kind, like a mutex.
//
// Flows are designed to be instantiated once and then modified as needed; you can grab
// a pointer to the appropriate underlying data by calling one of the *Data() methods
// of the struct.
type Flow struct {
	credentials    Credentials
	ID             string      `json:"id,omitempty"`
	EmbedID        string      `json:"embed_id,omitempty"`
	Tag            string      `json:"tag"`
	Data           interface{} `json:"data"`
	Variant        string      `json:"variant"`
	SourceProvider string      `json:"source_provider,omitempty"`
	Filter         string      `json:"filter,omitempty"`
	Params         string      `json:"params,omitempty"`
}

// NewFlow creates a new flow. Note that the `data` parameter *must* be a pointer to
// one of the variant structs defined in variant.go. If anything other than a pointer
// is passed, the function panics to prevent the creation of a silently immutable flow.
//
// If the flow is being submitted individually, the tag can be one of:
//
// ** The flow's named tag as entered in the Telemetry admin interface (e.g.: `gauge_1`)
//
// ** The flow's unique ID
//
// ** The flow's embed ID
//
// If, on the other hand, the flow is being submitted as part of a batch, only named
// tags are supported.
func NewFlow(tag string, data interface{}) *Flow {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		panic("NewFlow() expects a pointer to a variant struct")
	}

	return &Flow{Tag: tag, Data: data}
}

// NewFlowWithLayout function
func NewFlowWithLayout(credentials Credentials, tag string, variant, sourceProvider, filter, params string) (*Flow, error) {
	result := &Flow{
		credentials:    credentials,
		Tag:            tag,
		Variant:        variant,
		SourceProvider: sourceProvider,
		Filter:         filter,
		Params:         params,
	}

	err := result.Save()

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetFlowLayoutWithTag function
func GetFlowLayoutWithTag(credentials Credentials, tag string) (*Flow, error) {
	req, err := buildRequest("GET", credentials, "/flows/"+tag, nil)

	if err != nil {
		return nil, err
	}

	result := &Flow{credentials: credentials}

	err = sendJSONRequestInterface(req, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetFlowLayout function
func GetFlowLayout(credentials Credentials, id string) (*Flow, error) {
	req, err := buildRequest("GET", credentials, "/flows/"+id+"/layout", nil)

	if err != nil {
		return nil, err
	}

	result := &Flow{credentials: credentials}

	err = sendJSONRequestInterface(req, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SetFlowError function
func SetFlowError(credentials Credentials, tag string, body interface{}) error {
	req, err := buildRequest("POST", credentials, "/flows/"+tag+"/error", body)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(req)

	return err
}

// Publish sends a flow to the Telemetry API servers. On output, the function return
// nil if the submission was successful, an instance of gotelemetry.Error if a REST
// error occurred, or a errors.Error instance otherwise.
func (f *Flow) Publish(credentials Credentials) error {
	r, err := buildRequest(
		"PUT",
		credentials,
		"/flows/"+f.Tag+"/metrics",
		f.Data,
	)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(r)

	return err
}

// Update function
func (f *Flow) Update(data interface{}) error {
	encoded, err := json.Marshal(data)

	if err != nil {
		return err
	}

	json.Unmarshal(encoded, &f.Data)

	return nil
}

// Populate function
func (f *Flow) Populate(variant string, data interface{}) error {
	encoded, err := json.Marshal(data)

	if err != nil {
		return err
	}

	switch variant {
	case "barchart":
		f.Data = &Barchart{}
	case "box":
		f.Data = &Box{}
	case "bulletchart":
		f.Data = &Bulletchart{}
	case "countdown":
		f.Data = &Countdown{}
	case "custom":
		f.Data = &Custom{}
	case "funnelchart":
		f.Data = &Funnelchart{}
	case "gauge":
		f.Data = &Gauge{}
	case "graph":
		f.Data = &Graph{}
	case "grid":
		f.Data = &Grid{}
	case "histogram":
		f.Data = &Histogram{}
	case "icon":
		f.Data = &Icon{}
	case "image":
		f.Data = &Image{}
	case "log":
		f.Data = &Log{}
	case "map":
		f.Data = &Map{}
	case "multigauge":
		f.Data = &Multigauge{}
	case "multivalue":
		f.Data = &Multivalue{}
	case "piechart":
		f.Data = &Piechart{}
	case "scatterplot":
		f.Data = &Scatterplot{}
	case "servers":
		f.Data = &Servers{}
	case "status":
		f.Data = &Status{}
	case "table":
		f.Data = &Table{}
	case "text":
		f.Data = &Text{}
	case "tickertape":
		f.Data = &Tickertape{}
	case "timeline":
		f.Data = &Timeline{}
	case "timeseries":
		f.Data = &Timeseries{}
	case "upstatus":
		f.Data = &Upstatus{}
	case "value":
		f.Data = &Value{}
	case "video":
		f.Data = &Video{}
	case "waterfall":
		f.Data = &Waterfall{}
	default:
		return NewError(500, "Unknown variant "+variant)
	}

	json.Unmarshal(encoded, &f.Data)

	return nil
}

func (f *Flow) Read(credentials Credentials) error {
	var searchTag string

	if f.EmbedID != "" {
		searchTag = f.EmbedID
	} else if f.ID != "" {
		searchTag = f.ID
	} else {
		searchTag = f.Tag
	}

	req, err := buildRequest(
		"GET",
		credentials,
		"/flows/"+searchTag+"/metrics",
		nil,
	)

	if err != nil {
		return err
	}

	res, err := sendRawRequest(req)

	if err != nil {
		return err
	}

	needsConversion := false

	if f.Data == nil {
		f.Data = &map[string]interface{}{}
		needsConversion = true
	}

	err = readJSONResponseBody(res, f.Data)

	if err != nil {
		return err
	}

	if needsConversion {
		if err = f.Populate(f.Variant, f.Data); err != nil {
			return err
		}
	}

	return err
}

// Save function
func (f *Flow) Save() error {
	request, err := buildRequest("POST", f.credentials, "/flows", f)

	if err != nil {
		return err
	}

	err = sendJSONRequestInterface(request, &f)

	return err
}

// PostUpdate function
func (f *Flow) PostUpdate() error {
	request, err := buildRequest("PUT", f.credentials, "/flows/"+f.EmbedID+"/metrics", f.Data)

	if err != nil {
		return err
	}

	err = sendJSONRequestInterface(request, &f)

	return err
}

// Delete function
func (f *Flow) Delete() error {
	request, err := buildRequest("DELETE", f.credentials, "/flows/"+f.ID, nil)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(request)

	return err
}

// BarchartData function
func (f *Flow) BarchartData() (*Barchart, bool) {
	res, ok := f.Data.(*Barchart)

	return res, ok
}

// BulletchartData function
func (f *Flow) BulletchartData() (*Bulletchart, bool) {
	res, ok := f.Data.(*Bulletchart)

	return res, ok
}

// CountdownData function
func (f *Flow) CountdownData() (*Countdown, bool) {
	res, ok := f.Data.(*Countdown)

	return res, ok
}

// CustomData function
func (f *Flow) CustomData() (*Custom, bool) {
	res, ok := f.Data.(*Custom)

	return res, ok
}

// FunnelchartData function
func (f *Flow) FunnelchartData() (*Funnelchart, bool) {
	res, ok := f.Data.(*Funnelchart)

	return res, ok
}

// GaugeData function
func (f *Flow) GaugeData() (*Gauge, bool) {
	res, ok := f.Data.(*Gauge)

	return res, ok
}

// GraphData function
func (f *Flow) GraphData() (*Graph, bool) {
	res, ok := f.Data.(*Graph)

	return res, ok
}

// GridData function
func (f *Flow) GridData() (*Grid, bool) {
	res, ok := f.Data.(*Grid)

	return res, ok
}

// HistogramData function
func (f *Flow) HistogramData() (*Histogram, bool) {
	res, ok := f.Data.(*Histogram)

	return res, ok
}

// IconData function
func (f *Flow) IconData() (*Icon, bool) {
	res, ok := f.Data.(*Icon)

	return res, ok
}

// ImageData function
func (f *Flow) ImageData() (*Image, bool) {
	res, ok := f.Data.(*Image)

	return res, ok
}

// LogData function
func (f *Flow) LogData() (*Log, bool) {
	res, ok := f.Data.(*Log)

	return res, ok
}

// MapData funcion
func (f *Flow) MapData() (*Map, bool) {
	res, ok := f.Data.(*Map)

	return res, ok
}

// MultigaugeData function
func (f *Flow) MultigaugeData() (*Multigauge, bool) {
	res, ok := f.Data.(*Multigauge)

	return res, ok
}

// MultivalueData function
func (f *Flow) MultivalueData() (*Multivalue, bool) {
	res, ok := f.Data.(*Multivalue)

	return res, ok
}

// PiechartData function
func (f *Flow) PiechartData() (*Piechart, bool) {
	res, ok := f.Data.(*Piechart)

	return res, ok
}

// ScatterplotData function
func (f *Flow) ScatterplotData() (*Scatterplot, bool) {
	res, ok := f.Data.(*Scatterplot)

	return res, ok
}

// ServersData function
func (f *Flow) ServersData() (*Servers, bool) {
	res, ok := f.Data.(*Servers)

	return res, ok
}

// StatusData function
func (f *Flow) StatusData() (*Status, bool) {
	res, ok := f.Data.(*Status)

	return res, ok
}

// TableData function
func (f *Flow) TableData() (*Table, bool) {
	res, ok := f.Data.(*Table)

	return res, ok
}

// TextData function
func (f *Flow) TextData() (*Text, bool) {
	res, ok := f.Data.(*Text)

	return res, ok
}

// TickertapeData function
func (f *Flow) TickertapeData() (*Tickertape, bool) {
	res, ok := f.Data.(*Tickertape)

	return res, ok
}

// TimelineData function
func (f *Flow) TimelineData() (*Timeline, bool) {
	res, ok := f.Data.(*Timeline)

	return res, ok
}

// TimeseriesData function
func (f *Flow) TimeseriesData() (*Timeseries, bool) {
	res, ok := f.Data.(*Timeseries)

	return res, ok
}

// UpstatusData function
func (f *Flow) UpstatusData() (*Upstatus, bool) {
	res, ok := f.Data.(*Upstatus)

	return res, ok
}

// ValueData function
func (f *Flow) ValueData() (*Value, bool) {
	res, ok := f.Data.(*Value)

	return res, ok
}

// WaterfallData function
func (f *Flow) WaterfallData() (*Waterfall, bool) {
	res, ok := f.Data.(*Waterfall)

	return res, ok
}
