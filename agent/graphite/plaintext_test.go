package graphite

import (
	"net"
	"testing"

	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
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
	cfg := config.File{}
	cfg.Graphite = config.GraphiteConfig{
		UDPListenPort: ":2000",
	}
	cfg.Data = config.DataConfig{
		TTL:          "1h",
		DataLocation: "/tmp/agent.db",
	}

	errorChannel := make(chan error, 1)

	go func() {
		for e := range errorChannel {
			b.Logf("%s", e)
		}
	}()

	database.Init(&cfg, errorChannel)

	Init(&cfg, errorChannel)

	for n := 0; n < b.N; n++ {
		testGraphiteMemoryRun(b, 1)
	}

}
