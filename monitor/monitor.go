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

//Schedule runs Run() on the intreval seconds passed
func Schedule(seconds int) {
	interval := time.NewTicker(time.Duration(seconds))
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	signal.Notify(osChan, syscall.SIGTERM)

	for {
		select {
		case <-interval.C:
			// TODO: Schedule check
		case <-osChan:
			// TODO: graceful shutdown of monitor
			return

		}
	}
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

var healthyhostsWeight map[int][]store.Backend
var healthyhostsWeightLock sync.RWMutex

func Init() {
	healthyhostsWeight = make(map[int][]store.Backend)
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

// GetHealthChecker used to get the HealthChecker interface
func GetHealthChecker(db *sql.DB, fe *store.Frontend) HealthChecker {

	var checker HealthChecker

	switch fe.Balance {
	case store.BEST:
		//return hHost.backends
	default:
		// Weight based. Order it based on weight.
		healthyhostsWeightLock.RLock()
		hHosts, ok := healthyhostsWeight[fe.Id]
		healthyhostsWeightLock.RUnlock()
		if !ok {
			return nil
		}
		checker = &WeightBased{}
		// Should we suffle array???
		checker.Init(hHosts)
	}

	return checker
}
