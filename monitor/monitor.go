package monitor

import (
	"database/sql"
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
	pinger "github.com/sparrc/go-ping"
)

var statsChan chan *pinger.Statistics
var localIfc map[string]bool
var fastestlocalIP string

//Monitor contains monitor funcs and properties
type Monitor struct {
	wtGroup sync.WaitGroup
	done    chan bool
}

//NewMonitor returns Monitor object
func NewMonitor() (*Monitor, error) {
	return &Monitor{done: make(chan bool)}, nil
}

//Schedule runs Run() on the intreval seconds passed
func (m *Monitor) Schedule(seconds int) {
	interval := time.NewTicker(time.Duration(seconds))
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	signal.Notify(osChan, syscall.SIGTERM)

	for {
		select {
		case <-interval.C:
			// TODO: Schedule check
		case <-osChan:
			close(m.done)
			return
		case <-m.done:
			m.wtGroup.Wait()
			return

		}
	}
}

//GetMonitorResults returns the monitor statistics
func GetMonitorResults() (result Results, err error) {

	return Results{Interfaces: []string{fastestlocalIP}, TargetIPs: []string{}}, nil
}

func getTargetInterfaces() []string {
	return []string{"8.8.8.8", "github.com", "8.8.4.4"}
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

type HealthyBackends struct {
	totalWeight int
	backends    []store.Backend
}

var healthyhosts map[int]HealthyBackends
var lock sync.RWMutex

func Init() {
	healthyhosts = make(map[int]HealthyBackends)
	rand.Seed(time.Now().UTC().UnixNano())
}

func DoHealthCheck(db *sql.DB) error {
	fends, err := store.ReadFrontends(db)
	if err != nil {
		return err
	}

	for _, fe := range fends {

		logger.LogDebug(strconv.Itoa(fe.Id))
		// // get the backends for this frontend
		// bends, err = store.ReadBackends(db, fe.Id)
		// if err != nil {
		// 	logger.LogError("While getting backend for FE:" + strconv.Itoa(fe.Id))
		// 	continue
		// }

	}
	return nil
}

// GetHealthyBackends used to get the ordered list of
// backend servers for a given frontend.
func GetHealthyBackends(db *sql.DB, fe *store.Frontend) []store.Backend {

	lock.RLock()
	hHost, ok := healthyhosts[fe.Id]
	lock.RUnlock()
	if !ok {
		return nil
	}
	switch fe.Balance {
	case store.BEST:
		return hHost.backends
	default:
		// Weight based. Order it based on weight.
		// should we suffle array???
		return orderByWeight(hHost)
	}

}

func orderByWeight(hHosts HealthyBackends) []store.Backend {
	if hHosts.totalWeight <= 0 {
		logger.LogDebug("Weight is 0")
		return nil
	}
	var res []store.Backend
	rand := randInt(1, hHosts.totalWeight+1)
	curMin := 1
	for i, be := range hHosts.backends {
		curMax := curMin + be.Weight - 1
		if rand >= curMin && rand <= curMax {
			res = append(res, be)
			res = append(res, hHosts.backends[:i]...)
			res = append(res, hHosts.backends[i+1:]...)
			break
		} else {
			curMin += be.Weight
		}
	}
	return res
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
