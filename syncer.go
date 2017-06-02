package syncer

import (
	"context"
	"fmt"
	"github.com/ahmadmuzakki29/go-syncer/pb"
	"github.com/orcaman/concurrent-map"
	"sync"
)

var syncer = cmap.New()

type Syncer struct {
	m      *sync.Mutex
	buffer int
	done   context.CancelFunc
}

func lock(req *pb.LockRequest) {
	id := req.Id
	locktimeout := getLockTimeoutDuration(req.Locktimeout)

	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)
		lock.buffer += 1

		syncer.Set(id, lock)
		logger(LOG_WARNING, fmt.Sprint("suspending "+id, " buffer:", lock.buffer))
		lock.m.Lock()
	} else {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, locktimeout)

		syncer.Set(id, Syncer{
			m:      &sync.Mutex{},
			buffer: 1,
			done:   cancel,
		})

		// will unlock when timeout comes or unlock() called
		// which one first
		unlocker(req, ctx.Done())

		logger(LOG_INFO, "locking "+id, " buffer:", 1)
		l, _ := syncer.Get(id)
		l.(Syncer).m.Lock()
	}
}

func doUnlock(req *pb.LockRequest) {
	id := req.Id
	locktimeout := getLockTimeoutDuration(req.Locktimeout)

	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)

		lock.m.Unlock()
		lock.buffer -= 1

		logger(LOG_INFO, "releasing "+id, " buffer:", lock.buffer)
		if lock.buffer == 0 {
			syncer.Remove(id)
			return
		}

		// refreshing the timeout
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, locktimeout)
		lock.done = cancel
		unlocker(req, ctx.Done())

		syncer.Set(id, lock)
	}
}

func unlock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)
		lock.done()
	}
}

func unlocker(req *pb.LockRequest, done <-chan struct{}) {
	go func(req *pb.LockRequest, done <-chan struct{}) {
		defer func() {
			if r := recover(); r != nil {
				logger(LOG_ERROR, r)
			}
		}()
		<-done
		doUnlock(req)
	}(req, done)
}
