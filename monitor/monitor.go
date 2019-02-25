package monitor

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/helper"
	"github.com/govindarajan/laserproxy/logger"
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

//package store

func getLocalInterfaces() []string {
	localIPAddresses, err := helper.GetLocalIPs()
	if err != nil {
		logger.LogInfo("No local interface ips found" + err.Error())
		return nil
	}
	var addresses []string
	for _, ip := range localIPAddresses {
		addresses = append(addresses, ip.IP.String())
	}
	return addresses
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
