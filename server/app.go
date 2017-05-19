package main

import (
	"flag"
	"github.com/ahmadmuzakki29/go-syncer"
)

func main() {
	var debugMode bool
	var port string
	flag.StringVar(&port, "port", "9999", "port to listen")
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.Parse()

	syncer.DebugFlag = debugMode
	syncer.Serve(port)
}
