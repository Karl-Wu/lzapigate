package main

import (
	"log"
	"os"
)

var (
	appStop = false
)

func main() {
	DEBUG = log.New(os.Stderr, "DEBUG    ", log.Ltime)
	WARN = log.New(os.Stderr, "WARNING  ", log.Ltime)
	CRITICAL = log.New(os.Stderr, "CRITICAL ", log.Ltime)
	ERROR = log.New(os.Stderr, "ERROR    ", log.Ltime)

	login(username, password)

	go serverDetect()
	go servePowerOn()
	go servePositionReq()

	// Block indefinitely
	<-make(chan struct{})
}
