package monitor

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/govindarajan/laserproxy/helper"
	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/pinger"
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
			m.Run()
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

func (m *Monitor) collectResults() {
	minPktLoss := 1000.00
	defer m.wtGroup.Done()
	for stat := range statsChan {
		if _, ok := localIfc[stat.Addr]; ok {
			//this is a local interface
			if stat.PacketsLoss < minPktLoss {
				minPktLoss = stat.PacketsLoss
				fastestlocalIP = stat.Addr
			}
		}
		//fmt.Println(stat)
	}
	fmt.Println(fastestlocalIP, "is the fastest local IP")
}

//Run runs the monitors on the interfaces and target ips
func (m *Monitor) Run() {
	localIfc = make(map[string]bool)
	//todo: Get from  package store here
	localRoutes := getLocalInterfaces()
	fmt.Println(localRoutes)
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
		collectStatistics(routes, true)
	}(&wg, localRoutes)

	//collect  stats  for  target routes
	go func(wtGrp *sync.WaitGroup, routes []string) {
		defer wtGrp.Done()
		collectStatistics(routes, false)
	}(&wg, targetRoutes)
	wg.Wait()
	close(statsChan)
	m.wtGroup.Wait()

}

func collectStatistics(routes []string, isLocalInterface bool) {
	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(wtGrp *sync.WaitGroup, ipRoute string) {
			defer wtGrp.Done()
			pong, err := pinger.NewPinger(ipRoute)
			if err != nil {
				log.Println("unable to ping ip ", ipRoute)
				return
			}
			if isLocalInterface {
				pong.Source = ipRoute
				ipaddr, _ := net.ResolveIPAddr("ip", "8.8.8.8")
				pong.IPAddr = ipaddr
				localIfc[ipRoute] = true
			}
			pong.Count = 100
			pong.Timeout = time.Second * 25
			pong.Interval = time.Second / 100
			err = pong.Run()
			if err != nil {
				log.Println("unable to collect stats for ip ", ipRoute, err)
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
