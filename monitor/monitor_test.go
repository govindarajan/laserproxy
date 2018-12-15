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
	//monitor.Schedule(100)
	monitor.Run()
}

// func Test_Monitoring(t *testing.T) {
// 	//out, _ := exec.Command("ping", "-4", "8.8.8.8", "-c 5", "-i 0.1", "-w 10").Output()
// 	out, _ := exec.Command("ping 8.8.8.8 -c 5 -I 192.168.168.132").CombinedOutput()
// 	fmt.Println(string(out))
// 	if strings.Contains(string(out), "Destination Host Unreachable") {
// 		fmt.Println("TANGO DOWN")
// 	} else {
// 		fmt.Println("IT'S ALIVEEE")
// 	}

// }
