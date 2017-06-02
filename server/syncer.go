package main

import (
	"flag"
	"github.com/ahmadmuzakki29/go-syncer"
)

func main() {
	var logLevel string
	var port string
	flag.StringVar(&port, "port", "9999", "port to listen")
	flag.StringVar(&logLevel, "log-level", "warning", "[info|warning|error] log level to print. default: warning")
	flag.Parse()

	cfg := syncer.Config{
		Port:     port,
		LogLevel: logLevel,
	}
	syncer.Serve(cfg)
}
