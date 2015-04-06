package server

import (
	"encoding/json"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions"
	"io"
	"net/http"
)

func Init(listen string, errorChannel chan error) {
	errorChannel <- gotelemetry.NewDebugError("Eval server listening on %s", listen)
	go startup(listen, errorChannel)
}

func startup(listen string, errorChannel chan error) {
	http.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		handleEval(w, r, errorChannel)
	})
	http.ListenAndServe(listen, nil)
}

func log(r *http.Request, errorChannel chan error, err interface{}) {
	switch err.(type) {
	case error:
		e := err.(error)

		errorChannel <- gotelemetry.NewLogError("Eval %s -> %#s", r.RemoteAddr, e.Error())

	default:
		errorChannel <- gotelemetry.NewLogError("Eval %s -> %#v", r.RemoteAddr, err)
	}
}

func debug(r *http.Request, errorChannel chan error, err interface{}) {
	switch err.(type) {
	case error:
		e := err.(error)

		errorChannel <- gotelemetry.NewDebugError("Eval %s -> %#s", r.RemoteAddr, e.Error())

	default:
		errorChannel <- gotelemetry.NewDebugError("Eval %s -> %#v", r.RemoteAddr, err)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, errorChannel chan error, err interface{}) {
	switch err.(type) {
	case *gotelemetry.Error:
		e := err.(*gotelemetry.Error)

		errorChannel <- e

		w.WriteHeader(e.StatusCode)

		if e.StatusCode > 499 {
			io.WriteString(w, "A server error has occurred.")
		} else {
			io.WriteString(w, e.Error())
		}

	case error:
		e := err.(error)

		errorChannel <- e

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "A server error has occurred.")

	default:
		errorChannel <- gotelemetry.NewErrorWithFormat(500, "%#v", nil, err)

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "A server error has occurred.")
	}
}

func handleEval(w http.ResponseWriter, r *http.Request, errorChannel chan error) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	errorChannel <- gotelemetry.NewDebugError("Received eval request via HTTP from %s", r.RemoteAddr)

	var payload interface{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		err = gotelemetry.NewErrorWithFormat(http.StatusBadRequest, "Unable to parse JSON payload: `%s`", nil, err.Error())
		writeError(w, r, errorChannel, err)
		return
	}

	context, err := aggregations.GetContext()

	if err != nil {
		writeError(w, r, errorChannel, err)
		return
	}

	defer context.Close()

	log(r, errorChannel, payload)

	if expressions, ok := payload.([]interface{}); ok {
		result := []interface{}{}

		for _, expression := range expressions {
			output, err := functions.Parse(context, expression)

			if err != nil {
				context.SetError()

				result = append(result, map[string]interface{}{
					"err": err.Error(),
					"out": output,
				})
			} else {
				result = append(result, map[string]interface{}{
					"err": nil,
					"out": output,
				})
			}
		}

		debug(r, errorChannel, fmt.Sprintf("Result is %#v", result))

		if err := json.NewEncoder(w).Encode(result); err != nil {
			writeError(w, r, errorChannel, err)
			return
		}
	}

	if expression, ok := payload.(map[string]interface{}); ok {
		var result map[string]interface{}

		output, err := functions.Parse(context, expression)

		if err != nil {
			result = map[string]interface{}{
				"err": err.Error(),
				"out": output,
			}
		} else {
			result = map[string]interface{}{
				"err": nil,
				"out": output,
			}
		}

		if err := json.NewEncoder(w).Encode([]interface{}{result}); err != nil {
			writeError(w, r, errorChannel, err)
			return
		}
	}
}
