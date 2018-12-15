package monitor

import "sync"

//Monitor contains monitor funcs and properties
type Monitor struct {
}

//GetMonitorResults returns the monitor statistics
func (m *Monitor) GetMonitorResults() (result Results, err error) {
	return Results{}, nil
}

//Run runs the monitors on the interfaces and target ips
func (m *Monitor) Run() {
	//todo: Get from  package store here
	localRoutes := getLocalInterfaces()
	targetRoutes := getTargetInterfaces()
	var wg sync.WaitGroup
	wg.Add(2)

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
}

func collectStatistics(routes []string) {

}

//package store

func getLocalInterfaces() []string {
	return []string{"182.76.143.194", "182.76.143.175"}
}

func getTargetInterfaces() []string {
	return []string{"14.143.69.107", "github.com"}
}
