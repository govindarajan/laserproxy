package pinger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/govindarajan/laserproxy/pinger"
)

//Tests  here  must be run as sudo
//Normal user cannot emit icmp packets
func Test_Pinger(t *testing.T) {
	pong, err := pinger.NewPinger("115.248.131.195")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pong)
	pong.Count = 100
	pong.Interval = time.Second / 100
	pong.Timeout = time.Second * 25
	err = pong.Run()
	if err != nil {
		t.Error(err)
	}
	stats := pong.Statistics()
	if stats.PacketsRecv == 0 || stats.PacketsSent == 0 {
		t.Error("unable to send or receive packets")
	}

	fmt.Println(stats)
}
