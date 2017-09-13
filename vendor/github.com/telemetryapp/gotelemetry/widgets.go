package gotelemetry

type Widget struct {
	credentials Credentials `json:"-"`
	Id          string      `json:"id"`
	Variant     string      `json:"variant"`
	BoardId     string      `json:"board_id"`
	FlowIds     []string    `json:"flow_ids,omitempty"`
	FlowId      string      `json:"flow_id,omitempty"`
	Column      float64     `json:"column"`
	Row         float64     `json:"row"`
	Width       float64     `json:"width"`
	Height      float64     `json:"height"`
	BoardIndex  int         `json:"in_board_index"`
	Background  string      `json:"background"`
}

// Creates a new Widget on Telemetry and binds it to a specific board. Returns the created Widget struct if there are no errros.
func NewWidget(credentials Credentials, board *Board, variant string, column, row, width, height, boardIndex int, background string) (*Widget, error) {
	w := &Widget{
		credentials: credentials,
		Variant:     variant,
		BoardId:     board.Id,
		Column:      float64(column),
		Row:         float64(row),
		Width:       float64(width),
		Height:      float64(height),
		BoardIndex:  boardIndex,
		Background:  background,
	}

	err := w.Save()

	if err != nil {
		return nil, err
	}

	return w, nil
}

// Get a Widget from Telemetry API by ID
func GetWidget(credentials Credentials, id string) (*Widget, error) {
	request, err := buildRequest("GET", credentials, "/widgets/"+id, nil)

	if err != nil {
		return nil, err
	}

	w := &Widget{}

	err = sendJSONRequestInterface(request, &w)

	if err != nil {
		return nil, err
	}

	w.credentials = credentials

	return w, nil
}

// Delete a Widget from Telemetry
func (w *Widget) Delete() error {
	request, err := buildRequest("DELETE", w.credentials, "/widgets/"+w.Id, nil)

	if err != nil {
		return err
	}

	_, err = sendJSONRequest(request)

	return err
}

func (w *Widget) Save() error {
	request, err := buildRequest("POST", w.credentials, "/widgets", w)

	if err != nil {
		return err
	}

	return sendJSONRequestInterface(request, &w)
}
