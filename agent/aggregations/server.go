package aggregations

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/telemetryapp/gotelemetry"
	"net/http"
)

func InitServer(addr string, errorChannel chan error) error {
	r := httprouter.New()

	r.GET("/data/:key", getData)
	r.POST("/data/:key", postData)

	go func() {
		http.ListenAndServe(addr, r)
	}()

	errorChannel <- gotelemetry.NewLogError("Data API listening on %s", addr)

	return nil
}

func getData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if out, err := ReadStorage(p.ByName("key")); err == nil {
		if res, err := json.Marshal(out); err == nil {
			w.WriteHeader(200)
			w.Write(res)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func postData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Body == nil {
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	d := json.NewDecoder(r.Body)
	var data interface{}

	if err := d.Decode(&data); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

		return
	}

	if payload, ok := data.(map[string]interface{}); ok {
		if err := WriteStorage(p.ByName("key"), payload); err == nil {
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))

			return
		}
	}

	w.WriteHeader(400)
	w.Write([]byte("Invalid payload; must be a hash"))
}
