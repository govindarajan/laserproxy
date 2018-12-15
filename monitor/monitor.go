package monitor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/pinger"
)

var statsChan chan *pinger.Statistics

//Monitor contains monitor funcs and properties
type Monitor struct {
	wtGroup sync.WaitGroup
	Done    chan bool
}

//NewMonitor returns Monitor object
func NewMonitor() (*Monitor, error) {
	return &Monitor{Done: make(chan bool)}, nil
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
			m.Run()
		case <-osChan:
			close(m.Done)
			return
		case <-m.Done:
			m.wtGroup.Wait()
			return

		}
	}
}

//GetMonitorResults returns the monitor statistics
func (m *Monitor) GetMonitorResults() (result Results, err error) {

	return Results{}, nil
}

func (m *Monitor) collectResults() {
	defer m.wtGroup.Done()
	for stat := range statsChan {
		fmt.Println(stat)
	}
}

//Run runs the monitors on the interfaces and target ips
func (m *Monitor) Run() {
	//todo: Get from  package store here
	localRoutes := getLocalInterfaces()
	targetRoutes := getTargetInterfaces()
	statsChan = make(chan *pinger.Statistics, len(localRoutes)+len(targetRoutes))
	var wg sync.WaitGroup
	wg.Add(2)
	m.wtGroup.Add(1)
	go m.collectResults()
	//collect stats  for  local  interface routes
	//or other gateway routes
	go func(wtGrp *sync.WaitGroup, routes []string) {
		defer wtGrp.Done()
		collectStatistics(routes)
	}(&wg, localRoutes)

	//collect  stats  for  target routes
	go func(wtGrp *sync.WaitGroup, routes []string) {
		defer wtGrp.Done()
		collectStatistics(routes)
	}(&wg, targetRoutes)
	wg.Wait()
	close(statsChan)
	m.wtGroup.Wait()

}

func collectStatistics(routes []string) {
	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(wtGrp *sync.WaitGroup, ipRoute string) {
			defer wtGrp.Done()
			fmt.Println("Collecting stats for  ip", ipRoute)
			pong, err := pinger.NewPinger(ipRoute)
			if err != nil {
				log.Println("unable to ping ip ", ipRoute)
				return
			}
			pong.Count = 10
			pong.Timeout = time.Second * 25
			pong.Interval = time.Second / 100
			err = pong.Run()
			if err != nil {
				log.Println("unable to collect stats for ip ", ipRoute)
				return
			}
			stats := pong.Statistics()
			statsChan <- stats
		}(&wg, route)
	}
	wg.Wait()
}

//package store

func getLocalInterfaces() []string {
	return []string{"1.1.1.1", "1.1.4.4"}
}

func getTargetInterfaces() []string {
	return []string{"8.8.8.8", "github.com", "8.8.4.4"}
}
