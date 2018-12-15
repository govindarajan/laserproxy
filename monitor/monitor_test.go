package monitor_test

import (
	"testing"

	"github.com/govindarajan/laserproxy/monitor"
)

func Test_MonitorRun(t *testing.T) {
	monitor, err := monitor.NewMonitor()
	if err != nil {
		t.Error(err)
	}
	monitor.Schedule(100)
}
