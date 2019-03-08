package monitor

import (
	"math/rand"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
)

type HealthChecker interface {
	Init([]store.Backend)
	GetNext() *store.Backend
}

type WeightBased struct {
	totalWeight int
	backends    []store.Backend
}

func (r *WeightBased) GetNext() *store.Backend {
	if len(r.backends) <= 0 {
		return nil
	}
	be := r.backends[0]
	r.backends = r.backends[1:]
	return &be
}

func (r *WeightBased) Init(bes []store.Backend) {
	r.backends = bes
	weight := 0
	for _, be := range bes {
		weight += be.Weight
	}
	r.totalWeight = weight
	r.fix()
}

func (r *WeightBased) fix() {
	if r.totalWeight <= 0 {
		logger.LogDebug("Weight is 0")
		return
	}
	var res []store.Backend
	rand := randInt(1, r.totalWeight+1)
	curMin := 1
	for i, be := range r.backends {
		curMax := curMin + be.Weight - 1
		if rand >= curMin && rand <= curMax {
			res = append(res, be)
			res = append(res, r.backends[:i]...)
			res = append(res, r.backends[i+1:]...)
			break
		} else {
			curMin += be.Weight
		}
	}
	r.backends = res
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
