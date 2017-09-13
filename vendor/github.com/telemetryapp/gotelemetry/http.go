package gotelemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type TelemetryRequest struct {
	*http.Request
	credentials Credentials
}

var UserAgentString = "Gotelemetry"

func buildRequestWithHeaders(method string, credentials Credentials, fragment string, headers map[string]string, body interface{}, parameters ...map[string]string) (*TelemetryRequest, error) {
	debugChannel := credentials.DebugChannel

	URL := *credentials.ServerURL

	URL.Path = path.Join(URL.Path, fragment)

	if len(parameters) > 0 {
		p := url.Values{}

		for index, value := range parameters[0] {
			p.Add(index, value)
		}

		URL.RawQuery = p.Encode()
	}

	if debugChannel != nil {
		debugChannel <- NewDebugError(fmt.Sprintf("Building request %s %s", method, URL.String()))
	}

	var b []byte
	var err error

	if body == nil {
		b = []byte{}
	} else {
		b, err = json.Marshal(body)

		if err != nil {
			return nil, err
		}

		if debugChannel != nil {
			debugChannel <- NewDebugError(fmt.Sprintf("Request payload: %s", string(b)))
		}
	}

	r, err := http.NewRequest(method, URL.String(), bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	r.Header.Set("user-agent", UserAgentString)
	r.Header.Set("content-type", "application/json")
	r.SetBasicAuth(credentials.APIKey, "")

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	if debugChannel != nil {
		debugChannel <- NewDebugError(fmt.Sprintf("API Key: %s", credentials.APIKey))
	}

	return &TelemetryRequest{r, credentials}, nil
}

func buildRequest(method string, credentials Credentials, fragment string, body interface{}, parameters ...map[string]string) (*TelemetryRequest, error) {
	return buildRequestWithHeaders(method, credentials, fragment, map[string]string{}, body, parameters...)
}

func readJSONResponseBody(r *http.Response, target interface{}, debugChannel chan error) error {
	source, err := ioutil.ReadAll(r.Body)

	if err != nil && err != io.EOF {
		return err
	}

	if debugChannel != nil {
		debugChannel <- NewDebugError(fmt.Sprintf("Response payload: %s", string(source)))
	}

	if len(source) == 0 {
		// Nothing to read
		return nil
	}

	if err := json.Unmarshal(source, target); err != nil {
		return NewError(400, "Invalid JSON response: "+string(source)+" (JSON decode error: "+err.Error()+")")
	}

	return nil
}

func sendRawRequest(request *TelemetryRequest) (*http.Response, error) {
	r, err := http.DefaultClient.Do(request.Request)

	return r, err
}

func sendJSONRequestInterface(request *TelemetryRequest, target interface{}) error {
	debugChannel := request.credentials.DebugChannel

	r, err := sendRawRequest(request)

	if err != nil {
		return err
	}

	if debugChannel != nil {
		debugChannel <- NewDebugError(fmt.Sprintf("Response status code: %d", r.StatusCode))

		for key, value := range r.Header {
			debugChannel <- NewDebugError(fmt.Sprintf("Response header %s: %s", key, value))
		}
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.StatusCode > 399 {
		v, _ := ioutil.ReadAll(r.Body)

		if len(v) > 0 && debugChannel != nil {
			debugChannel <- NewDebugError(fmt.Sprintf("Response payload: %s", string(v)))
		}

		return NewErrorWithData(r.StatusCode, r.Status, v)
	}

	return readJSONResponseBody(r, target, request.credentials.DebugChannel)
}

func sendJSONRequest(request *TelemetryRequest) (interface{}, error) {
	var body interface{}

	err := sendJSONRequestInterface(request, &body)

	return body, err
}
