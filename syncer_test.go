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
	// we have process with the same ID
	id := "sameid"
	address := fmt.Sprint("localhost:9999")
	syncer.Init(address)

	processCount := 10

	res := make(chan string, processCount*2)
	var i int
	for i < processCount {
		go func(a int) {
			syncer.Lock(id)
			d := getRandomDuration()
			res <- fmt.Sprint("start process ", a)
			// simulate random duration process
			time.Sleep(time.Duration(d))
			res <- fmt.Sprint("finish process ", a)
			if i < 2 {
				// deliberately not unlocking the first 2 process
				return
			}
			syncer.Unlock(id)
		}(i)
		i += 1
	}

	i = 0
	var msg string
	for i < processCount {
		// this expect the result will be synchronous
		msg = <-res
		assert.Equal(t, fmt.Sprint("start process ", i), msg)
		fmt.Println(msg)
		msg = <-res
		assert.Equal(t, fmt.Sprint("finish process ", i), msg)
		fmt.Println(msg)
		i += 1
	}
}

func getRandomDuration() time.Duration {
	ms := rand.Float32() * 1000 // to make it ms
	return time.Duration(ms) * time.Millisecond
}
