package main

import (
	"flag"
	"github.com/ahmadmuzakki29/go-syncer"
	"github.com/prometheus/common/log"
	"time"
)

func main() {
	var debugMode bool
	var port string
	flag.StringVar(&port, "port", "9999", "port to listen")
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "30s", "time until unlocked automatically")
	flag.Parse()

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal(err)
	}

	syncer.DebugFlag = debugMode

	cfg := syncer.Config{
		Port:    port,
		Timeout: timeout,
	}
	syncer.Serve(cfg)
}
