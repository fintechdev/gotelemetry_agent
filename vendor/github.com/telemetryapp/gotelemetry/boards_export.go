package gotelemetry

import (
	"encoding/json"
)

// ExportedFlow struct
type ExportedFlow struct {
	Tag  string      `json:"tag"`
	Data interface{} `json:"data"`
}

// ExportedWidget struct
type ExportedWidget struct {
	Flow       *ExportedFlow `json:"flow"`
	Variant    string        `json:"variant"`
	Column     float64       `json:"column"`
	Row        float64       `json:"row"`
	Width      float64       `json:"width"`
	Height     float64       `json:"height"`
	BoardIndex int           `json:"in_board_index"`
	Background string        `json:"background"`
}

// ExportedBoard struct
type ExportedBoard struct {
	Name             string            `json:"name,omitempty"`
	Theme            string            `json:"theme,omitempty"`
	DisplayName      bool              `json:"display_board_name,omitempty"`
	AspectRatio      string            `json:"aspect_ratio,omitempty"`
	FontFamily       string            `json:"font_family,omitempty"`
	FontSize         string            `json:"font_size,omitempty"`
	WidgetBackground string            `json:"widget_background,omitempty"`
	WidgetMargins    int64             `json:"widget_margins,omitempty"`
	WidgetPadding    int64             `json:"widget_padding,omitempty"`
	Widgets          []*ExportedWidget `json:"widgets"`
}

// ImportBoard function
func ImportBoard(credentials Credentials, name string, prefix string, board *ExportedBoard) (*Board, error) {
	// First, make sure the board doesn't already exist

	result, err := GetBoardByName(credentials, name)

	if err == nil && result != nil {
		result, err = GetBoard(credentials, result.ID)

		if err != nil {
			return nil, err
		}

		result.Prefix = prefix

		//TODO: Make sure that the board retrieved from the API matches
		// the makeup of the board template we're trying to import.

		return result, nil
	}

	result = &Board{credentials: credentials, Prefix: prefix}

	encoded, err := json.Marshal(board)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(encoded, result)

	if err != nil {
		return nil, err
	}

	result.Widgets = []*Widget{}

	result.Name = name

	err = result.Save()

	if err != nil {
		return nil, err
	}

	for _, exportedWidget := range board.Widgets {
		var flow *Flow
		var err error

		if exportedWidget.Flow != nil {
			flow, err = NewFlowWithLayout(
				credentials,
				prefix+exportedWidget.Flow.Tag,
				exportedWidget.Variant,
				"",
				"",
				"",
			)

			if err != nil {
				result.Delete()
				return nil, err
			}

			flow.Data = exportedWidget.Flow.Data
			err = flow.Publish(credentials)

			if err != nil {
				flow.Delete()
				result.Delete()
				return nil, err
			}
		}

		widget := &Widget{credentials: credentials}

		encoded, err = json.Marshal(exportedWidget)

		if err != nil {
			if flow != nil {
				flow.Delete()
			}

			result.Delete()
			return nil, err
		}

		err = json.Unmarshal(encoded, &widget)

		if err != nil {
			if flow != nil {
				flow.Delete()
			}

			result.Delete()
			return nil, err
		}

		widget.BoardID = result.ID

		if flow != nil {
			widget.FlowID = flow.ID
		}

		err = widget.Save()

		if err != nil {
			if flow != nil {
				flow.Delete()
			}

			result.Delete()
			return nil, err
		}

		result.Widgets = append(result.Widgets, widget)
	}

	return result, nil
}

// Export function
func (b *Board) Export() (*ExportedBoard, error) {
	result := &ExportedBoard{}

	encoded, err := json.Marshal(b)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(encoded, &result)

	if err != nil {
		return nil, err
	}

	for index, widget := range b.Widgets {
		if len(widget.FlowIds) > 0 {
			flow, err := GetFlowLayout(b.credentials, widget.FlowIds[0])

			if err != nil {
				return nil, err
			}

			err = flow.Read(b.credentials)

			if err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(flow)

			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(encoded, &result.Widgets[index].Flow)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
