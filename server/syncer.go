package main

import (
	"flag"
	"github.com/ahmadmuzakki29/go-syncer"
	"log"
	"time"
)

func main() {
	var logLevel string
	var port string
	flag.StringVar(&port, "port", "9999", "port to listen")
	flag.StringVar(&logLevel, "log-level", "warning", "[info|warning|error] log level to print. default: warning")
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "30s", "time until unlocked automatically. default: 30s")
	flag.Parse()

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal(err)
	}

	cfg := syncer.Config{
		Port:     port,
		Timeout:  timeout,
		LogLevel: logLevel,
	}
	syncer.Serve(cfg)
}
