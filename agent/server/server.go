package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions"
	"io"
	"net/http"
	"time"
)

func Init(configFile *config.ConfigFile, errorChannel chan error) error {
	var err error

	apiKey, err := configFile.Accounts()[0].GetAPIKey()

	if err != nil {
		return err
	}

	credentials, err := gotelemetry.NewCredentials(apiKey)

	if err != nil {
		return err
	}

	credentials.SetDebugChannel(errorChannel)

	stream, err := gotelemetry.NewBatchStream(credentials, time.Second, errorChannel)

	if err != nil {
		return err
	}

	errorChannel <- gotelemetry.NewDebugError("Eval server listening on %s", configFile.ListenAddress())
	go startup(configFile.ListenAddress(), errorChannel, stream)

	return nil
}

func startup(listen string, errorChannel chan error, stream *gotelemetry.BatchStream) {
	router := httprouter.New()

	router.POST("/eval", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handleEval(w, r, errorChannel)
	})

	router.POST("/flow/:tag", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handleFlow(w, r, errorChannel, gotelemetry.BatchTypePOST, params.ByName("tag"), stream)
	})

	router.PATCH("/flow/:tag", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handleFlow(w, r, errorChannel, gotelemetry.BatchTypePOST, params.ByName("tag"), stream)
	})

	http.ListenAndServe(listen, router)
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
	debug(r, errorChannel, gotelemetry.NewDebugError("Received eval request via HTTP from %s", r.RemoteAddr))

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
	} else if expression, ok := payload.(map[string]interface{}); ok {
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
	} else {
		err = gotelemetry.NewErrorWithFormat(http.StatusBadRequest, "Unable to process payload: `%#v`", nil, payload)
		writeError(w, r, errorChannel, err)
	}

}

func handleFlow(w http.ResponseWriter, r *http.Request, errorChannel chan error, submissionType gotelemetry.BatchType, tag string, stream *gotelemetry.BatchStream) {
	debug(r, errorChannel, gotelemetry.NewDebugError("Received flow request via HTTP from %s", r.RemoteAddr))

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

	if expression, ok := payload.(map[string]interface{}); ok {
		var result map[string]interface{}

		output, err := functions.Parse(context, expression)

		if err == nil {
			if output != nil {
				var submissionType gotelemetry.BatchType

				if r.Method == "POST" {
					submissionType = gotelemetry.BatchTypePOST
				} else {
					submissionType = gotelemetry.BatchTypePATCH
				}

				stream.SendData(tag, output, submissionType)
			} else {
				err = gotelemetry.NewErrorWithFormat(http.StatusBadRequest, "The payload processed into a nil value.", nil, payload)
				writeError(w, r, errorChannel, err)
			}
		} else {
			writeError(w, r, errorChannel, err)
		}

		if err := json.NewEncoder(w).Encode([]interface{}{result}); err != nil {
			writeError(w, r, errorChannel, err)
			return
		}
	} else {
		err = gotelemetry.NewErrorWithFormat(http.StatusBadRequest, "Unable to process payload: `%#v`", nil, payload)
		writeError(w, r, errorChannel, err)
	}
}
