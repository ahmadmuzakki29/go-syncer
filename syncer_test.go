package syncer_test

import (
	"fmt"
	"github.com/ahmadmuzakki29/go-syncer"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

// run this test after the server running
func TestSyncer(t *testing.T) {
	// we have 4 process with the same ID
	id := "sameid"
	address := fmt.Sprint("localhost", syncer.PORT)
	syncer.Init(address)

	processCount := 20

	res := make(chan string, processCount*2)
	var i int
	for i < processCount {
		go func() {
			syncer.Lock(id)
			d := getRandomDuration()
			res <- "start process"
			// simulate random duration process
			time.Sleep(time.Duration(d))
			res <- "finish process"
			syncer.Unlock(id)
		}()
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
