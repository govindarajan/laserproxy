package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/govindarajan/laserproxy/admin"
	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/worker"
)

func main() {
	logger.LogInfo("Hello Laser")
	worker.Initialize()

	go admin.StartAdminServer()
	go signalCatcher()
	go worker.StartProxy()
	select {}
}

func signalCatcher() {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt, syscall.SIGTERM)
	<-osChan
	logger.LogInfo("Got ctrl+C signal")
	os.Exit(0)
}
