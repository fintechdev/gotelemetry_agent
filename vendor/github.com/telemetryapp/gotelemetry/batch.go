package gotelemetry

import (
	"encoding/json"
)

// BatchResponse struct
type BatchResponse struct {
	Errors  []string `json:"errors"`
	Skipped []string `json:"skipped"`
	Updated []string `json:"updated"`
}

// BatchType type
type BatchType int

const (
	// BatchTypePOST constant
	BatchTypePOST BatchType = iota
	// BatchTypePATCH constant
	BatchTypePATCH BatchType = iota
	// BatchTypeJSONPATCH constant
	BatchTypeJSONPATCH BatchType = iota
)

// Batch type describes a collection of flows that can be submitted simultaneously to the Telemetry servers.
//
// Note the underlying data structure of the batch is a map, and therefore batches are not thread safe
// by default. If you require thread safety, you must mediate access to the batch through some kind
// of synchronization mechansism, like a mutex.
type Batch map[string]interface{}

// SetFlow adds or overwrites a flow to the batch
func (b Batch) SetFlow(f *Flow) {
	b[f.Tag] = f.Data
}

// SetData adds or overwrites data to the batch
func (b Batch) SetData(tag string, data interface{}) {
	b[tag] = data
}

// DeleteFlow deletes a flow from the batch
func (b Batch) DeleteFlow(tag string) {
	delete(b, tag)
}

// Publish submits a batch to the Telemetry API servers, and returns either an instance
// of gotelemetry.Error if a REST error occurs, or errors.Error if any other error occurs.
func (b Batch) Publish(credentials Credentials, channelTag string, submissionType BatchType) error {
	data := map[string]interface{}{}

	for key, submission := range b {
		if logger.IsTrace() {
			payload, _ := json.Marshal(submission)
			logger.Trace(
				"About to post flow with data",
				"flow", key,
				"data", string(payload),
			)
		}

		data[key] = submission
	}

	method := "POST"
	headers := map[string]string{}

	if submissionType != BatchTypePOST {
		method = "PATCH"

		if submissionType == BatchTypeJSONPATCH {
			headers["Content-Type"] = "application/json-patch+json"
		}
	}

	endpoint := "/metrics"

	if channelTag != "" {
		endpoint = "/channels/" + channelTag + "/metrics"
	}

	r, err := buildRequestWithHeaders(
		method,
		credentials,
		endpoint,
		headers,
		map[string]interface{}{
			"data": data,
		},
	)

	if err != nil {
		return err
	}

	response := BatchResponse{}

	err = sendJSONRequestInterface(r, &response)

	if logger.IsWarn() {
		for _, errString := range response.Errors {
			logger.Warn(
				"API Error (Errors)",
				"error", errString,
			)
		}
		for _, skippedFlow := range response.Skipped {
			logger.Warn(
				"API Error (Skipped): The flow was not updated",
				"flow", skippedFlow,
			)
		}
	}

	return err
}
