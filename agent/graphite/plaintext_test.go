package graphite

import (
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"net"
	"runtime"
	"testing"
	"time"
)

func testGraphiteMemoryRun(t *testing.T, count int) uint64 {
	var before, after runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&before)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})

	if err != nil {
		t.Errorf("%s", err)
	}

	for index := 0; index < count; index++ {
		s := fmt.Sprintf("ca.tabini.test 123 123\n")

		_, err = conn.WriteToUDP([]byte(s), &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 2000})

		if err != nil {
			t.Errorf("%s", err)
		}
	}

	time.Sleep(time.Second)

	runtime.GC()
	runtime.ReadMemStats(&after)

	return after.HeapAlloc - before.HeapAlloc
}

func TestGraphiteMemory(t *testing.T) {
	cfg := config.ConfigFile{}
	cfg.Graphite = config.GraphiteConfig{
		UDPListenPort: ":2000",
	}

	errorChannel := make(chan error, 1)

	go func() {
		for e := range errorChannel {
			t.Logf("%s", e)
		}
	}()

	l := "/tmp/agent.sqlite3"
	ttl := "1h"
	aggregations.Init(&l, &ttl, errorChannel)

	Init(&cfg, errorChannel)

	run1 := testGraphiteMemoryRun(t, 100)
	run2 := testGraphiteMemoryRun(t, 100)
	run3 := testGraphiteMemoryRun(t, 100)
	run4 := testGraphiteMemoryRun(t, 100)

	t.Logf("%d - %d - %d - %d", run1, run2, run3, run4)
}
