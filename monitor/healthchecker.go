package monitor

import (
	"math/rand"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
)

type HealthChecker interface {
	Init(interface{})
	GetNext() *store.Backend
}

type WeightBased struct {
	totalWeight int
	backends    []CheckResult
}

func (r *WeightBased) GetNext() *store.Backend {
	if len(r.backends) <= 0 {
		return nil
	}
	be := r.backends[0].be
	r.backends = r.backends[1:]
	return be
}

func (r *WeightBased) Init(val interface{}) {
	var ok bool
	if r.backends, ok = val.([]CheckResult); !ok {
		logger.LogError("Conversion to []Backend failed")
	}

	weight := 0
	for _, ber := range r.backends {
		weight += ber.be.Weight
	}
	r.totalWeight = weight
	r.fix()
}

func (r *WeightBased) fix() {
	if r.totalWeight <= 0 {
		logger.LogDebug("Weight is 0")
		return
	}
	var res []CheckResult
	rand := randInt(1, r.totalWeight+1)
	curMin := 1
	for i, ber := range r.backends {
		curMax := curMin + ber.be.Weight - 1
		if rand >= curMin && rand <= curMax {
			res = append(res, ber)
			res = append(res, r.backends[:i]...)
			res = append(res, r.backends[i+1:]...)
			break
		} else {
			curMin += ber.be.Weight
		}
	}
	r.backends = res
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
