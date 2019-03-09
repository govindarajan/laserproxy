package monitor

import (
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
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

var healthyhostsWeight map[int][]store.Backend
var healthyhostsWeightLock sync.RWMutex
var healthyhostsBest map[int][]CheckResult
var healthyhostsBestLock sync.RWMutex
var chanCheckResult chan CheckResult

func Init() {
	healthyhostsWeight = make(map[int][]store.Backend)
	healthyhostsBest = make(map[int][]CheckResult)
	chanCheckResult = make(chan CheckResult, 10)
	go processCheckResult(chanCheckResult)
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
		healthyhostsBestLock.RLock()
		checkRes, ok := healthyhostsBest[fe.Id]
		healthyhostsBestLock.RUnlock()
		if !ok {
			return nil
		}
		checker = &BestRoute{}
		checker.Init(checkRes)
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
