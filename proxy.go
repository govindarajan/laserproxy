package main

import (
	"./helper"
	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/worker"
)

func main() {
	logger.LogInfo("Hello Laser")
	helper.ConfigureRoute()

	go worker.StartProxy()
	select {}
}
