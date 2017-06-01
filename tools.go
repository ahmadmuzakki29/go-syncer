package syncer

import "fmt"

var DebugFlag bool

func debugger(msg string) {
	if !DebugFlag {
		return
	}

	fmt.Println(msg)
}
