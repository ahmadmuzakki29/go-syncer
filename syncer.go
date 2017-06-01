package syncer

import (
	"context"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"sync"
	"time"
)

var TIMEOUT time.Duration

var syncer = cmap.New()

type Syncer struct {
	m      *sync.Mutex
	buffer int
	done   context.CancelFunc
}

func lock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)
		lock.buffer += 1

		syncer.Set(id, lock)
		debugger(fmt.Sprint("suspending "+id, " buffer:", lock.buffer))
		lock.m.Lock()
	} else {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, TIMEOUT)

		syncer.Set(id, Syncer{
			m:      &sync.Mutex{},
			buffer: 1,
			done:   cancel,
		})

		// will unlock when timeout comes or unlock() called
		// which one first
		unlocker(id, ctx.Done())

		debugger(fmt.Sprint("locking "+id, " buffer:", 1))
		l, _ := syncer.Get(id)
		l.(Syncer).m.Lock()
	}
}

func doUnlock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)

		lock.m.Unlock()
		lock.buffer -= 1

		debugger(fmt.Sprint("releasing "+id, " buffer:", lock.buffer))
		if lock.buffer == 0 {
			syncer.Remove(id)
			return
		}

		// refreshing the timeout
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
		lock.done = cancel
		unlocker(id, ctx.Done())

		syncer.Set(id, lock)
	}
}

func unlock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)
		lock.done()
	}
}

func unlocker(id string, done <-chan struct{}) {
	go func(id string, done <-chan struct{}) {
		<-done
		doUnlock(id)
	}(id, done)
}
