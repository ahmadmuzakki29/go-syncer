package main

import (
	"flag"
	"github.com/ahmadmuzakki29/go-syncer"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.Parse()

	syncer.DebugFlag = debugMode
	syncer.Serve()
}
