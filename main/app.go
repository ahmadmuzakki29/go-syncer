package main

import (
	"github.com/ahmadmuzakki29/go-syncer"
	"net/http"
)

func main() {
	http.HandleFunc("/lock", syncer.HandlerLock)
	http.HandleFunc("/unlock", syncer.HandlerUnlock)
	http.ListenAndServe(":9001", nil)
}
