package syncer_test

import (
	"fmt"
	"github.com/ahmadmuzakki29/go-syncer/client"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

// run this test after the server running
func TestSyncer(t *testing.T) {
	// we have process with the same ID
	id := "sameid"
	address := fmt.Sprint("localhost:9999")
	cfg := client.Config{
		EndPoint:    address,
		LockTimeout: time.Duration(10) * time.Second,
	}
	client.Init(cfg)

	processCount := 10

	res := make(chan string, processCount*2)
	var i int
	for i < processCount {
		go func(a int) {
			client.Lock(id)
			d := getRandomDuration()
			res <- "start process"
			// simulate random duration process
			time.Sleep(time.Duration(d))
			res <- "finish process"
			if i > 7 {
				// deliberately not unlocking the last 2 process
				return
			}
			client.Unlock(id)
		}(i)
		i += 1
	}

	i = 0
	var msg string
	for i < processCount {
		// this expect the result will be synchronous
		msg = <-res
		assert.Equal(t, "start process", msg)
		fmt.Println(msg)
		msg = <-res
		assert.Equal(t, "finish process", msg)
		fmt.Println(msg)
		i += 1
	}
}

func getRandomDuration() time.Duration {
	ms := rand.Float32() * 1000 // to make it ms
	return time.Duration(ms) * time.Millisecond
}
