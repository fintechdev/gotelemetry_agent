package gotelemetry

type BarchartBar struct {
	Color string      `json:"color,omitempty"`
	Label string      `json:"label,omitempty"`
	Value interface{} `json:"value"`
}

type Barchart struct {
	ExpiresAt int64         `json:"expires_at,omitempty"`
	Opacity   *float64      `json:"opacity,omitempty"`
	Title     string        `json:"title,omitempty"`
	Priority  int           `json:"priority,omitempty"`
	Min       interface{}   `json:"min,omitempty"`
	Max       interface{}   `json:"max,omitempty"`
	Sort      interface{}   `json:"sort,omitempty"`
	Bars      []BarchartBar `json:"bars"`
}

type Box struct {
}

type BulletchartChart struct {
	Colors     []string `json:"colors,omitempty"`
	Label      string   `json:"label,omitempty"`
	Marker     int      `json:"marker,omitempty"`
	Max        int      `json:"max"`
	Thresholds []int    `json:"thresholds,omitempty"`
	Value      int      `json:"value"`
	ValueType  string   `json:"value_type,omitempty"`
}

type Bulletchart struct {
	ExpiresAt    int64              `json:"expires_at,omitempty"`
	Opacity      *float64           `json:"opacity,omitempty"`
	Title        string             `json:"title,omitempty"`
	Bulletcharts []BulletchartChart `json:"bulletcharts"`
}

type Countdown struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Message   string   `json:"message"`
	Time      int64    `json:"time"`
}

type Custom struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
}

type FunnelchartChart struct {
	Color string  `json:"color,omitempty"`
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type Funnelchart struct {
	ExpiresAt int64              `json:"expires_at,omitempty"`
	Opacity   *float64           `json:"opacity,omitempty"`
	Title     string             `json:"title,omitempty"`
	Priority  int                `json:"priority,omitempty"`
	Values    []FunnelchartChart `json:"values"`
}

type Gauge struct {
	ExpiresAt   int64    `json:"expires_at,omitempty"`
	Opacity     *float64 `json:"opacity,omitempty"`
	Title       string   `json:"title,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Value       float64  `json:"value"`
	ValueColor  string   `json:"value_color,omitempty"`
	ValueType   string   `json:"value_type,omitempty"`
	GaugeColor  string   `json:"gauge_color,omitempty"`
	Min         float64  `json:"min,omitempty"`
	Max         float64  `json:"max,omitempty"`
	Range       int      `json:"range,omitempty"`
	Value2      float64  `json:"value_2,omitempty"`
	Value2Color string   `json:"value_2_color,omitempty"`
	Value2Label string   `json:"value_2_label,omitempty"`
	Icon        string   `json:"icon,omitempty"`
}

type GraphSeries struct {
	Color  string    `json:"color,omitempty"`
	Label  string    `json:"label,omitempty"`
	Values []float64 `json:"values"`
}

type Graph struct {
	ExpiresAt int64         `json:"expires_at,omitempty"`
	Opacity   *float64      `json:"opacity,omitempty"`
	Title     string        `json:"title,omitempty"`
	Priority  int           `json:"priority,omitempty"`
	Series    []GraphSeries `json:"series"`
	Baseline  string        `json:"basline,omitempty"`
	EndTime   int64         `json:"end_time,omitempty"`
	StartTime int64         `json:"start_time,omitempty"`
	Label1    string        `json:"label_1,omitempty"`
	Label2    string        `json:"label_2,omitempty"`
	Label3    string        `json:"label_3,omitempty"`
	MinScale  float64       `json:"min_scale,omitempty"`
	Renderer  string        `json:"renderer,omitempty"`
	Unstack   bool          `json:"unstack,omitempty"`
	ValueType string        `json:"value_type,omitempty"`
	XLabels   []string      `json:"x_labels,omitempty"`
}

type GridData struct {
	Fill      int    `json:"fill"`
	Label     string `json:"label"`
	FillColor string `json:"fill_color,omitempty"`
	BGColor   string `json:"bg_color,omitempty"`
	Color     string `json:"color,omitempty"`
}

type Grid struct {
	ExpiresAt int64        `json:"expires_at,omitempty"`
	Opacity   *float64     `json:"opacity,omitempty"`
	Title     string       `json:"title,omitempty"`
	Priority  int          `json:"priority,omitempty"`
	Data      [][]GridData `json:"data"`
}

type Histogram struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
}

type IconIcon struct {
	Color string `json:"color"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type Icon struct {
	ExpiresAt int64      `json:"expires_at,omitempty"`
	Opacity   *float64   `json:"opacity,omitempty"`
	Title     string     `json:"title,omitempty"`
	Priority  int        `json:"priority,omitempty"`
	Icons     []IconIcon `json:"icons"`
}

type Image struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Link      string   `json:"link,omitempty"`
	Mode      string   `json:"mode,omitempty"`
	URL       string   `json:"url"`
}

type LogMessage struct {
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
	Color     string `json:"color,omitempty"`
}

type Log struct {
	ExpiresAt int64        `json:"expires_at,omitempty"`
	Opacity   *float64     `json:"opacity,omitempty"`
	Title     string       `json:"title,omitempty"`
	Priority  int          `json:"priority,omitempty"`
	Messages  []LogMessage `json:"messages"`
}

type MapCoord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type MapCoordWithZoom struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Zoom int     `json:"zoom"`
}

type MapCircle struct {
	Center      MapCoord `json:"center"`
	FillColor   string   `json:"fill_color,omitempty"`
	Label       string   `json:"label"`
	LineWidth   int      `json:"line_width"`
	Radius      int      `json:"radius"`
	StrokeColor string   `json:"stroke_color"`
}

type MapMarker struct {
	Color  string   `json:"color"`
	Coords MapCoord `json:"coords"`
	Icon   string   `json:"icon"`
	Label  string   `json:"label"`
}

type MapPolygon struct {
	FillColor   string     `json:"fill_color"`
	Label       string     `json:"label"`
	LineWidth   int        `json:"line_width"`
	StrokeColor string     `json:"stroke_color"`
	Vertices    []MapCoord `json:"vertices"`
}

type MapPolyline struct {
	Label       string     `json:"label"`
	LineWidth   int        `json:"line_width"`
	StrokeColor string     `json:"stroke_color"`
	Vertices    []MapCoord `json:"vertices"`
}

type Map struct {
	ExpiresAt int64            `json:"expires_at,omitempty"`
	Opacity   *float64         `json:"opacity,omitempty"`
	Title     string           `json:"title,omitempty"`
	Circles   []MapCircle      `json:"circles,omitempty"`
	Coords    MapCoordWithZoom `json:"coords"`
	MapboxId  string           `json:"mapbox_id,omitempty"`
	Markers   []MapMarker      `json:"markers,omitempty"`
	Polygons  []MapPolygon     `json:"polygons,omitempty"`
	Polylines []MapPolyline    `json:"polylines,omitempty"`
	Type      string           `json:"type,omitempty"`
}

type MultigaugeGauge struct {
	Label     string  `json:"label"`
	Value     float64 `json:"value"`
	Max       float64 `json:"max,omitempty"`
	ValueType string  `json:"value_type,omitempty"`
	Icon      string  `json:"icon,omitempty"`
}

type Multigauge struct {
	GaugeColor string            `json:"gauge_color,omitempty"`
	ExpiresAt  int64             `json:"expires_at,omitempty"`
	Opacity    *float64          `json:"opacity,omitempty"`
	Title      string            `json:"title,omitempty"`
	Priority   int               `json:"priority,omitempty"`
	Layout     string            `json:"layout,omitempty"`
	Gauges     []MultigaugeGauge `json:"gauges"`
}

type MultivalueValue struct {
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Color      string  `json:"color,omitempty"`
	ValueType  string  `json:"value_type,omitempty"`
	Abbreviate bool    `json:"abbreviate,omitempty"`
	Rounding   int     `json:"rounding,omitempty"`
	Icon       string  `json:"icon,omitempty"`
	LabelColor string  `json:"label_color,omitempty"`
}

type Multivalue struct {
	ExpiresAt int64             `json:"expires_at,omitempty"`
	Opacity   *float64          `json:"opacity,omitempty"`
	Title     string            `json:"title,omitempty"`
	Priority  int               `json:"priority,omitempty"`
	Values    []MultivalueValue `json:"values"`
}

type Piechart struct {
	ExpiresAt int64     `json:"expires_at,omitempty"`
	Opacity   *float64  `json:"opacity,omitempty"`
	Title     string    `json:"title,omitempty"`
	Priority  int       `json:"priority,omitempty"`
	Colors    []string  `json:"colors,omitempty"`
	Labels    []string  `json:"labels"`
	Renderer  string    `json:"renderer,omitempty"`
	Values    []float64 `json:"values"`
}

type Scatterplot struct {
	ExpiresAt int64     `json:"expires_at,omitempty"`
	Opacity   *float64  `json:"opacity,omitempty"`
	Title     string    `json:"title,omitempty"`
	Priority  int       `json:"priority,omitempty"`
	Values    []float64 `json:"values"`
	XLabel    string    `json:"x_label,omitempty"`
	YLabel    string    `json:"y_label,omitempty"`
}

type Server struct {
	Labels []string  `json:"labels,omitempty"`
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

type Servers struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Orange    float64  `json:"name,omitempty"`
	Red       float64  `json:"name,omitempty"`
	Servers   []Server `json:"servers"`
}

type StatusItem struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type Status struct {
	ExpiresAt int64        `json:"expires_at,omitempty"`
	Opacity   *float64     `json:"opacity,omitempty"`
	Title     string       `json:"title,omitempty"`
	Priority  int          `json:"priority,omitempty"`
	Statuses  []StatusItem `json:"statuses"`
}

type TableCell struct {
	Value     interface{} `json:"value,omitempty"`
	Color     string      `json:"color,omitempty"`
	Alignment string      `json:"alignment,omitempty"`
	Icon      string      `json:"icon,omitempty"`
	ValueType string      `json:"value_type,omitempty"`
	Sparkline []float64   `json:"sparkline,omitempty"`
}

type Table struct {
	ExpiresAt int64         `json:"expires_at,omitempty"`
	Opacity   *float64      `json:"opacity,omitempty"`
	Title     string        `json:"title,omitempty"`
	Priority  int           `json:"priority,omitempty"`
	Headers   []string      `json:"headers,omitempty"`
	Cells     [][]TableCell `json:"cells"`
}

type Text struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Alignment string   `json:"alignment,omitempty"`
	Text      string   `json:"text"`
}

type Tickertape struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Messages  []string `json:"messages"`
}

type TimelineMessage struct {
	From      string `json:"from"`
	IconURL   string `json:"icon_url,omitempty"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type Timeline struct {
	ExpiresAt int64             `json:"expires_at,omitempty"`
	Opacity   *float64          `json:"opacity,omitempty"`
	Title     string            `json:"title,omitempty"`
	Priority  int               `json:"priority,omitempty"`
	Messages  []TimelineMessage `json:"messages"`
}

type TimeseriesSeriesMetadata struct {
	Aggregation string `json:"aggregation"`
	Label       string `json:"label,omitempty"`
	Color       string `json:"color,omitempty"`
	ValueType   string `json:"value_type,omitempty"`
	Interpolate bool   `json:"interpolate,omitempty"`
}

type Timeseries struct {
	ExpiresAt      int64                      `json:"expires_at,omitempty"`
	Opacity        *float64                   `json:"opacity,omitempty"`
	Title          string                     `json:"title,omitempty"`
	Renderer       string                     `json:"renderer,omitempty"`
	Baseline       string                     `json:"baseline,omitempty"`
	Interval       string                     `json:"interval"`
	IntervalCount  int                        `json:"interval_count"`
	SeriesMetadata []TimeseriesSeriesMetadata `json:"series_metadata"`
	Values         []float64                  `json:"values"`
}

type Upstatus struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Down      []string `json:"down,omitempty"`
	Up        []string `json:"up,omitempty"`
	LastDown  int64    `json:"last_down,omitempty"`
	Uptime    float64  `json:"uptime,omitempty"`
}

type Video struct {
	ExpiresAt int64    `json:"expires_at,omitempty"`
	Opacity   *float64 `json:"opacity,omitempty"`
	Title     string   `json:"title,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Mode      string   `json:"mode,omitempty"`
	MP4       string   `json:"mp4,omitempty"`
	Muted     bool     `json:"muted,omitempty"`
	OGG       string   `json:"ogg,omitempty"`
	Poster    string   `json:"poster,omitempty"`
	WebM      string   `json:"webm,omitempty"`
}

type Value struct {
	ExpiresAt  int64       `json:"expires_at,omitempty"`
	Opacity    *float64    `json:"opacity,omitempty"`
	Title      string      `json:"title,omitempty"`
	Priority   int         `json:"priority,omitempty"`
	Color      string      `json:"color,omitempty"`
	Delta      float64     `json:"delta,omitempty"`
	DeltaType  string      `json:"delta_type,omitempty"`
	Label      string      `json:"label,omitempty"`
	Sparkline  []float64   `json:"sparkline,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	ValueType  string      `json:"value_type,omitempty"`
	Abbreviate bool        `json:"abbreviate,omitempty"`
	Rounding   int         `json:"rounding,omitempty"`
	Icon       string      `json:"icon,omitempty"`
	LabelColor string      `json:"label_color,omitempty"`
}

type WaterfallData struct {
	Serial int           `json:"serial"`
	Values []interface{} `json:"values"`
}

type Waterfall struct {
	ExpiresAt int64           `json:"expires_at,omitempty"`
	Opacity   *float64        `json:"opacity,omitempty"`
	Title     string          `json:"title,omitempty"`
	Priority  int             `json:"priority,omitempty"`
	Color     string          `json:"color,omitempty"`
	Direction string          `json:"direction,omitempty"`
	Spread    int             `json:"spread,omitempty"`
	ValueType string          `json:"value_type,omitempty"`
	Data      []WaterfallData `json:"data"`
}
