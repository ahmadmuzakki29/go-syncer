package syncer

import (
	"fmt"
	"time"
)

var LogLevel int

const DEFAULT_LOCK_TIMEOUT string = "30s"

const (
	LOG_INFO = iota
	LOG_WARNING
	LOG_ERROR
)

func logger(level int, msgs ...interface{}) {
	if level < LogLevel {
		return
	}
	msgs = append([]interface{}{getPrefix(level)}, msgs...)
	fmt.Println(msgs...)
}

func getPrefix(level int) string {
	switch level {
	case 0:
		return "[INFO]"
	case 1:
		return "[WARNING]"
	default:
		return "[ERROR]"
	}
}

func getLogLevel(level string) int {
	switch level {
	case "info":
		return 0
	case "warning":
		return 1
	default:
		return 2
	}
}

func getLockTimeoutDuration(timeout string) time.Duration {
	locktimeout, err := time.ParseDuration(timeout)
	if err != nil {
		logger(LOG_ERROR, err)
		locktimeout, _ = time.ParseDuration(DEFAULT_LOCK_TIMEOUT)
	}
	return locktimeout
}
