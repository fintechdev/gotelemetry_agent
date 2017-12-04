package gotelemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// TelemetryRequest struct
type TelemetryRequest struct {
	*http.Request
	credentials Credentials
}

// UserAgentString variable
var UserAgentString = "Gotelemetry"

func buildRequestWithHeaders(method string, credentials Credentials, fragment string, headers map[string]string, body interface{}, parameters ...map[string]string) (*TelemetryRequest, error) {
	URL := *credentials.ServerURL

	URL.Path = path.Join(URL.Path, fragment)

	if len(parameters) > 0 {
		p := url.Values{}

		for index, value := range parameters[0] {
			p.Add(index, value)
		}

		URL.RawQuery = p.Encode()
	}

	if logger.IsDebug() {
		logger.Debug(
			"Building request",
			"method", method,
			"url", URL.String(),
		)
	}

	var err error
	var r *http.Request

	if body == nil {
		r, err = http.NewRequest(method, URL.String(), nil)
	} else {
		var b []byte
		b, err = json.Marshal(body)

		if err != nil {
			return nil, err
		}

		if logger.IsTrace() {
			logger.Trace(
				"Request payload",
				"payload", string(b),
				"payload_size", len(b),
			)
		}

		r, err = http.NewRequest(method, URL.String(), bytes.NewReader(b))
	}

	if err != nil {
		return nil, err
	}

	r.Header.Set("user-agent", UserAgentString)
	r.Header.Set("content-type", "application/json")
	r.SetBasicAuth(credentials.APIKey, "")

	for key, value := range headers {
		if logger.IsTrace() {
			logger.Trace(
				"Request header",
				"header_name", key,
				"header_value", value,
			)
		}
		r.Header.Set(key, value)
	}

	if logger.IsTrace() {
		logger.Trace(
			"API Key",
			"key", strings.Repeat("*", len(credentials.APIKey)),
		)
	}

	return &TelemetryRequest{r, credentials}, nil
}

func buildRequest(method string, credentials Credentials, fragment string, body interface{}, parameters ...map[string]string) (*TelemetryRequest, error) {
	return buildRequestWithHeaders(method, credentials, fragment, map[string]string{}, body, parameters...)
}

func readJSONResponseBody(r *http.Response, target interface{}) error {
	defer func() {
		closeErr := r.Body.Close()
		if closeErr != nil {
			if logger.IsWarn() {
				logger.Warn(
					"failed to close response body reader",
					"error", closeErr.Error(),
				)
			}
		}
	}()

	source, err := ioutil.ReadAll(r.Body)

	if err != nil && err != io.EOF {
		return err
	}

	if logger.IsTrace() {
		logger.Trace(
			"Response payload",
			"payload", string(source),
		)
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

var client *http.Client

func sendRawRequest(request *TelemetryRequest) (*http.Response, error) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:       120 * time.Second,
				TLSHandshakeTimeout:   15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 30 * time.Second,
		}
	}
	return client.Do(request.Request)
}

func sendJSONRequestInterface(request *TelemetryRequest, target interface{}) error {
	r, err := sendRawRequest(request)
	if err != nil {
		return err
	}

	if logger.IsDebug() && r.StatusCode < 400 {
		logger.Debug(
			"Response status",
			"status_code", r.StatusCode,
			"status", r.Status,
		)

		for k, v := range r.Header {
			logger.Debug(
				"Response header",
				"header_name", k,
				"header_values", v,
			)
		}
	}

	if r.StatusCode > 399 {
		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				if logger.IsWarn() {
					logger.Warn(
						"failed to close response body reader",
						"error", closeErr.Error(),
					)
				}
			}
		}()

		if logger.IsWarn() {
			logger.Warn(
				"Response status",
				"status_code", r.StatusCode,
				"status", r.Status,
			)

			for k, v := range r.Header {
				logger.Warn(
					"Response header",
					"header_name", k,
					"header_values", v,
				)
			}
		}

		var v []byte
		v, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		if len(v) > 0 && logger.IsTrace() {
			logger.Trace(
				"Response payload",
				"payload", string(v),
			)
		}

		return NewErrorWithData(r.StatusCode, r.Status, v)
	}

	return readJSONResponseBody(r, target)
}

func sendJSONRequest(request *TelemetryRequest) (interface{}, error) {
	var body interface{}

	err := sendJSONRequestInterface(request, &body)

	return body, err
}
