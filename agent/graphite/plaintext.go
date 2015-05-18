package graphite

import (
	"bufio"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Init(cfg config.ConfigInterface, errorChannel chan error) error {
	graphiteConfig := cfg.GraphiteConfig()

	if graphiteConfig.TCPListenPort != "" {
		go setupTCPListener(graphiteConfig.TCPListenPort, errorChannel)
	}

	if graphiteConfig.UDPListenPort != "" {
		go setupUDPListener(graphiteConfig.UDPListenPort, errorChannel)
	}

	return nil
}

func setupTCPListener(listen string, errorChannel chan error) {
	l, err := net.Listen("tcp", listen)

	if err != nil {
		errorChannel <- err
		return
	}

	defer l.Close()

	errorChannel <- gotelemetry.NewLogError("Graphite => Listening for TCP plaintext connections on %s", l.Addr())

	for {
		conn, err := l.Accept()

		if err != nil {
			errorChannel <- err
			return
		}

		go handleTCPRequest(conn, errorChannel, true)
	}
}

func setupUDPListener(listen string, errorChannel chan error) {
	addr, err := net.ResolveUDPAddr("udp", listen)

	if err != nil {
		errorChannel <- err
		return
	}

	conn, err := net.ListenUDP("udp", addr)

	errorChannel <- gotelemetry.NewLogError("Graphite => Listening for UDP plaintext messages on %s", conn.LocalAddr())

	if err != nil {
		errorChannel <- err
		return
	}

	context, err := aggregations.GetContext()

	defer context.Close()

	if err != nil {
		errorChannel <- gotelemetry.NewErrorWithFormat(
			500, "Graphite => [UDP] Unable to obtain data context: %s",
			nil,
			err.Error(),
		)

		return
	}

	buf := make([]byte, 2048)

	for {
		if n, addr, err := conn.ReadFromUDP(buf); err == nil {
			remoteAddress := addr.String() + ", UDP"

			if err := parseRequest(context, remoteAddress, string(buf[0:n]), errorChannel); err != nil {
				errorChannel <- gotelemetry.NewErrorWithFormat(400, "Graphite => [%s, UDP] Error %s while receving data", nil, addr, err)
			}
		} else {
			errorChannel <- gotelemetry.NewErrorWithFormat(400, "Graphite => [%s, UDP] Error %s while receving data", nil, addr, err)
		}
	}
}

var splitter = regexp.MustCompile(" +")

func parseTraditionalRequest(context *aggregations.Context, remoteAddress string, line []string, errorChannel chan error) error {
	seriesName := line[0]

	value, err := strconv.ParseFloat(line[1], 64)

	if err != nil {
		return gotelemetry.NewErrorWithFormat(
			400, "Graphite => [%s] Invalid value %s: %s",
			nil,
			remoteAddress,
			line[1],
			err.Error(),
		)
	}

	timestamp, err := strconv.ParseInt(line[2], 10, 64)

	if err != nil {
		return gotelemetry.NewErrorWithFormat(
			400, "Graphite => [%s] Invalid timestamp %s: %s",
			nil,
			remoteAddress,
			line[2],
			err.Error(),
		)
	}

	series, isCreated, err := aggregations.GetSeries(context, seriesName)

	if err != nil {
		return gotelemetry.NewErrorWithFormat(
			500, "Graphite => [%s] Unable to get series %s: %s",
			nil,
			remoteAddress,
			seriesName,
			err.Error(),
		)
	}

	if isCreated {
		errorChannel <- gotelemetry.NewLogError("Graphite => Started receiving graphite data for '%s'", seriesName)
	}

	ts := time.Unix(timestamp, 0)

	if err := series.Push(&ts, value); err != nil {
		return gotelemetry.NewErrorWithFormat(
			500, "Graphite => [%s] Unable to push value %f with timestamp %s to series %s: %s",
			nil,
			remoteAddress,
			value,
			ts,
			seriesName,
			err.Error(),
		)
	}

	errorChannel <- gotelemetry.NewDebugError(
		"Graphite => [%s] Pushed value %f to series %s at time %s",
		remoteAddress,
		value,
		seriesName,
		ts,
	)

	return nil
}

func parseCounterRequest(context *aggregations.Context, remoteAddress string, line []string, errorChannel chan error) error {
	counterName := line[0]

	value, err := strconv.ParseInt(line[1], 10, 64)

	if err != nil {
		return gotelemetry.NewErrorWithFormat(
			400, "Graphite => [%s] Invalid value %s: %s",
			nil,
			remoteAddress,
			line[1],
			err.Error(),
		)
	}

	counter, isCreated, err := aggregations.GetCounter(context, counterName)

	if isCreated {
		errorChannel <- gotelemetry.NewLogError("Graphite => Started receiving graphite data for '%s'", counterName)
	}

	if err != nil {
		return gotelemetry.NewErrorWithFormat(
			500, "Graphite => [%s] Unable to get counter %s: %s",
			nil,
			remoteAddress,
			counterName,
			err.Error(),
		)
	}

	counter.Increment(value)

	return nil
}

func parseRequest(context *aggregations.Context, remoteAddress, request string, errorChannel chan error) error {
	request = strings.TrimSpace(request)

	if request == "" {
		return nil
	}

	line := splitter.Split(request, -1)

	context.Begin()

	defer context.Commit()

	switch len(line) {
	case 2:
		return parseCounterRequest(context, remoteAddress, line, errorChannel)

	case 3:
		return parseTraditionalRequest(context, remoteAddress, line, errorChannel)

	default:
		return gotelemetry.NewErrorWithFormat(
			400, "Graphite => [%s] Unable to parse request %s",
			nil,
			remoteAddress,
			request,
		)
	}
}

func handleTCPRequest(conn net.Conn, errorChannel chan error, closeOnError bool) {
	defer conn.Close()

	remoteAddress := conn.RemoteAddr().String()

	errorChannel <- gotelemetry.NewDebugError("Graphite => New connection from %s", remoteAddress)

	r := bufio.NewScanner(conn)

	context, err := aggregations.GetContext()

	defer context.Close()

	if err != nil {
		errorChannel <- gotelemetry.NewErrorWithFormat(
			500, "Graphite => [%s] Unable to obtain data context: %s",
			nil,
			remoteAddress,
			err.Error(),
		)

		if closeOnError {
			return
		}
	}

	for r.Scan() {
		if err := parseRequest(context, remoteAddress, r.Text(), errorChannel); err != nil {
			errorChannel <- err

			if closeOnError {
				return
			}
		}
	}

	if err := r.Err(); err != nil && err != io.EOF {
		errorChannel <- err
	}

	errorChannel <- gotelemetry.NewDebugError("Graphite => Connection from %s terminated.", remoteAddress)
}
