package graphite

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"net"
	"testing"
)

var addr = &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 2000}
var s = []byte("ca.tabini.test 123 123\n")
var conn, _ = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})

func testGraphiteMemoryRun(b *testing.B, count int) {
	for index := 0; index < count; index++ {
		_, err := conn.WriteToUDP(s, addr)

		if err != nil {
			b.Errorf("%s", err)
		}
	}
}

func BenchmarkGraphiteMemory(b *testing.B) {
	cfg := config.ConfigFile{}
	cfg.Graphite = config.GraphiteConfig{
		UDPListenPort: ":2000",
	}

	errorChannel := make(chan error, 1)

	go func() {
		for e := range errorChannel {
			b.Logf("%s", e)
		}
	}()

	p := ":2000"
	l := "/tmp/agent.sqlite3"
	ttl := "1h"
	aggregations.Init(&p, &l, &ttl, errorChannel)

	Init(&cfg, errorChannel)

	for n := 0; n < b.N; n++ {
		testGraphiteMemoryRun(b, 1)
	}

	// time.Sleep(time.Second)
}
