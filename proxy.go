package main

import (
	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/worker"
)

func main() {
	logger.LogInfo("Hello Laser")
	worker.Initialize()

	go worker.StartProxy()
	select {}
}
