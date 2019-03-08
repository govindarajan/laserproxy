package monitor

import (
	"errors"
	"time"

	"github.com/govindarajan/laserproxy/store"
	pinger "github.com/sparrc/go-ping"
)

type CheckResult struct {
	be   *store.Backend
	ping *pinger.Statistics
}

type CheckHeap []CheckResult

func (ch CheckHeap) Len() int {
	return len(ch)
}

func (ch CheckHeap) Less(i, j int) bool {
	// TODO: use other parameters also.
	return ch[i].ping.PacketLoss < ch[j].ping.PacketLoss
}

func (ch CheckHeap) Swap(i, j int) {
	ch[i], ch[j] = ch[j], ch[i]
}

func (ch *CheckHeap) Push(v interface{}) {
	val, ok := v.(CheckResult)
	if !ok {
		return
	}
	*ch = append(*ch, val)
}

func (ch *CheckHeap) Pop() interface{} {
	old := *ch
	n := len(old)
	elm := old[n-1]
	*ch = old[0 : n-1]
	return elm
}

const PACKET_COUNT = 10
const TIMEOUT_IN_SEC = 2
const INTERVAL_IN_MS = 100

func GetPingStats(addr string) (*pinger.Statistics, error) {
	if addr == "" {
		return nil, errors.New("Address is empty")
	}
	ping, err := pinger.NewPinger(addr)
	if err != nil {
		return nil, err
	}
	ping.Count = PACKET_COUNT
	ping.Interval = INTERVAL_IN_MS * time.Millisecond
	ping.Timeout = TIMEOUT_IN_SEC * time.Second
	ping.SetPrivileged(true)
	ping.Run()
	return ping.Statistics(), nil
}