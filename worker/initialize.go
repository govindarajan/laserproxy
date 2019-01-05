package worker

import (
	"database/sql"
	"os"

	"github.com/govindarajan/laserproxy/helper"
	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
)

var maindb *sql.DB

// Initialize used to initialize requirements for proxy.
func Initialize() {
	initDB()
	// Local IP checker
	go checkAndUpdateIPChange()
}

func initDB() {
	var err error
	maindb, err = store.GetConnection()
	if err != nil {
		logger.LogError("InitDB GetConn: " + err.Error())
		os.Exit(1)
	}
	if err = store.InitMainDB(maindb); err != nil {
		logger.LogError("InitDB InitTable: " + err.Error())
		os.Exit(2)
	}
}

func checkAndUpdateIPChange() {
	changed := make(chan bool, 1)
	if err := updateLocalIP(); err != nil {
		logger.LogError("Update LocalIP Err:" + err.Error())
	}
	// TODO: Get this from config.
	interval := 5
	go helper.WatchNetworkChange(interval, changed)
	for {
		logger.LogInfo("Waiting for Network change")
		<-changed
		logger.LogInfo("Network change detected.")
		if err := updateLocalIP(); err != nil {
			logger.LogError("Update LocalIP Err:" + err.Error())
			continue
		}
	}
}

func updateLocalIP() error {
	IPs, err := helper.GetLocalIPs()
	if err != nil {
		return err
	}

	helper.ConfigureRoute(IPs)

	existingIPs, err := store.ReadLocalRoutes(maindb)
	if err != nil {
		return err
	}

	for _, IP := range IPs {
		lr := &store.LocalRoute{
			Interface: IP.IFace,
			IP:        IP.IP,
			Gateway:   IP.Gateway,
			Weight:    1,
		}
		if err := store.WriteLocalRoute(maindb, lr); err != nil {
			logger.LogError("Write to DB: " + err.Error())
		}
	}

	staleRoutes := findStaleRoutes(existingIPs, IPs)
	for _, sRoute := range staleRoutes {
		if err := store.DeleteLocalRoute(maindb, &sRoute); err != nil {
			logger.LogError("Delete from DB: " + err.Error())
		}
	}
	return nil
}

func findStaleRoutes(existing []store.LocalRoute, cur []helper.LocalIPAddr) []store.LocalRoute {
	var res []store.LocalRoute
	for _, e := range existing {
		found := false
		for _, c := range cur {
			if e.IP.Equal(c.IP) && e.Gateway.Equal(c.Gateway) &&
				e.Interface == c.IFace {
				found = true
				break
			}
		}
		if !found {
			res = append(res, e)
		}
	}
	return res
}
