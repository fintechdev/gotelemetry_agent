package gotelemetry

import (
	"strings"
)

// Board struct
type Board struct {
	credentials      Credentials
	Prefix           string    `json:"-"`
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Theme            string    `json:"theme,omitempty"`
	DisplayName      bool      `json:"display_board_name,omitempty"`
	AspectRatio      string    `json:"aspect_ratio,omitempty"`
	FontFamily       string    `json:"font_family,omitempty"`
	FontSize         string    `json:"font_size,omitempty"`
	WidgetBackground string    `json:"widget_background,omitempty"`
	WidgetMargins    int64     `json:"widget_margins,omitempty"`
	WidgetPadding    int64     `json:"widget_padding,omitempty"`
	Widgets          []*Widget `json:"widgets,omitempty"`
	ChannelIds       []string  `json:"channel_ids,omitempty"`
}

// NewBoard creates a new board in Telemetry by informing the basic parameters and return a new *Board instance
func NewBoard(credentials Credentials, name, theme string, displayName bool, aspectRatio string) (*Board, error) {
	result := &Board{
		credentials: credentials,
		Name:        name,
		Theme:       theme,
		DisplayName: displayName,
		AspectRatio: aspectRatio,
	}

	result.Save()

	return result, nil
}

// GetBoard returns a board from Telemetry API by its ID
func GetBoard(credentials Credentials, id string) (*Board, error) {
	request, err := buildRequest("GET", credentials, "/boards/"+id, nil)

	if err != nil {
		return nil, err
	}

	b := &Board{}

	err = sendJSONRequestInterface(request, b)

	if err != nil {
		return nil, err
	}

	b.credentials = credentials

	return b, err
}

// GetBoardByName function
func GetBoardByName(credentials Credentials, name string) (*Board, error) {
	request, err := buildRequest("GET", credentials, "/boards", nil, map[string]string{"name": name})

	if err != nil {
		return nil, err
	}

	b := []*Board{}

	err = sendJSONRequestInterface(request, &b)

	if err != nil {
		return nil, err
	}

	if len(b) < 1 {
		return nil, NewError(404, "Board not found")
	}

	result := b[0]
	result.credentials = credentials

	return result, err
}

// Save function
func (b *Board) Save() error {
	request, err := buildRequest("POST", b.credentials, "/boards", b)

	if err != nil {
		return err
	}

	err = sendJSONRequestInterface(request, b)

	return err
}

// Delete deletes a board from Telemetry
func (b *Board) Delete() error {
	request, err := buildRequest("DELETE", b.credentials, "/boards/"+b.ID, nil)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(request)
	return err
}

// MapWidgetsToFlows function
func (b *Board) MapWidgetsToFlows() (map[string]*Flow, error) {
	var err error

	if b.Widgets == nil {
		if b, err = GetBoard(b.credentials, b.ID); err != nil {
			return nil, err
		}
	}

	result := map[string]*Flow{}

	for _, widget := range b.Widgets {
		if len(widget.FlowIds) == 0 {
			continue
		}

		f, err := GetFlowLayout(b.credentials, widget.FlowIds[0])

		if err != nil {
			return nil, err
		}

		err = f.Read(b.credentials)

		if err != nil {
			return nil, err
		}

		result[strings.TrimPrefix(f.Tag, b.Prefix)] = f
	}

	return result, nil
}
