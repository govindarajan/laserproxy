package monitor

//IMonitor is the interface for the monitor package
type IMonitor interface {
	GetMonitorresults()
	Run()
	Schedule(seconds int)
}
