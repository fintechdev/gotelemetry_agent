package gotelemetry

import (
	"flag"
)

var testKey = flag.String("api-key", "", "Telemetry API key used for testing.")

func getTestKey() string {
	if *testKey == "" {
		panic("Missing test key. Run tests with `go test -api-key=<your api key> ./...`")
	}

	return *testKey
}
