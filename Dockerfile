FROM golang:1.9
ADD . /go/src/github.com/telemetryapp/gotelemetry_agent
CMD [ "/go/src/github.com/telemetryapp/gotelemetry_agent/build.sh" ]
