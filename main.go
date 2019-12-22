package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DrDelphi/docui/helpers"
)

func main() {
	nodesCount := helpers.InitializeNodes()
	if nodesCount == 0 {
		fmt.Println("Sorry, but I cannot detect any running node")
		return
	}
	helpers.InitDisplay()
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-time.After(time.Millisecond * 500):
			if helpers.AppTerminated {
				return
			}
			helpers.GetNodesInfo()
		case <-sigTerm:
			return
		}
	}
}
