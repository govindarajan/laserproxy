package monitor

import (
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
)

var healthyHosts map[int][]CheckResult
var healthyHostsLock sync.RWMutex
var chanCheckResult chan CheckResult

func Init(db *sql.DB) {
	healthyHosts = make(map[int][]CheckResult)
	chanCheckResult = make(chan CheckResult, 10)
	go processCheckResult(chanCheckResult)
	rand.Seed(time.Now().UTC().UnixNano())
	Schedule(db)
}

// GetHealthChecker used to get the HealthChecker interface
func GetHealthChecker(db *sql.DB, fe *store.Frontend) HealthChecker {
	var checker HealthChecker

	healthyHostsLock.RLock()
	checkRes, ok := healthyHosts[fe.Id]
	healthyHostsLock.RUnlock()
	if !ok {
		return nil
	}

	switch fe.Balance {
	case store.BEST:
		checker = &BestRoute{}
	default:
		// Weight based. Order it based on weight.
		checker = &WeightBased{}
		// Should we suffle array???
	}
	checker.Init(checkRes)

	return checker
}

func Schedule(db *sql.DB) {
	interval := time.NewTicker(time.Second)
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	signal.Notify(osChan, syscall.SIGTERM)

	for {
		select {
		case t := <-interval.C:
			ScheduleInternal(db, t.Unix())
		case <-osChan:
			// TODO: graceful shutdown of monitor
			return

		}
	}
}

func ScheduleInternal(db *sql.DB, now int64) {

	bends, err := store.ReadAllBackends(db)
	if err != nil {
		logger.LogError("While getting backend: " + err.Error())
		return
	}
	for _, be := range bends {
		if be.CheckInterval == 0 {
			continue
		}
		if now%int64(be.CheckInterval) != 0 {
			// This is not the right time to do check.
			continue
		}
		go doHealthCheck(be, chanCheckResult)
	}

}

func doHealthCheck(be store.Backend, resChan chan CheckResult) {
	// TODO: Add url check also.
	s := strings.Split(be.Host, ":")
	ip := s[0]
	pingRes, err := GetPingStats(ip)
	if err != nil {
		logger.LogError("Ping error " + err.Error())
	}
	checkRes := CheckResult{be: &be, ping: pingRes}
	resChan <- checkRes
}

func processCheckResult(chanCheckResult chan CheckResult) {
	for {
		// TODO: Add done chan for graceful shutdown
		res := <-chanCheckResult
		healthyHostsLock.RLock()
		checkResults, ok := healthyHosts[res.be.GroupId]
		healthyHostsLock.RUnlock()
		if !ok {
			// First time. Add it to the slice.
			checkResults = append(checkResults, res)
		} else {
			//Replace with existing checkresult
			checkResults = replaceCheckResult(checkResults, res)
		}
		healthyHostsLock.Lock()
		healthyHosts[res.be.GroupId] = checkResults
		healthyHostsLock.Unlock()
	}
}

func replaceCheckResult(existing []CheckResult, res CheckResult) []CheckResult {
	var newRes = make([]CheckResult, 0)
	found := false
	for _, eRes := range existing {
		if res.be.Host == eRes.be.Host {
			found = true
			newRes = append(newRes, res)
		} else {
			newRes = append(newRes, eRes)
		}
	}
	if !found {
		newRes = append(newRes, res)
	}
	return newRes
}
