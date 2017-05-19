package syncer

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"sync"
)

var syncer = cmap.New() //make(map[string]Syncer)

type Syncer struct {
	m      *sync.Mutex
	buffer int
}

func lock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)
		lock.buffer += 1
		syncer.Set(id, lock)
		debugger("suspending "+id, lock.buffer)
		lock.m.Lock()
	} else {
		syncer.Set(id, Syncer{
			m:      &sync.Mutex{},
			buffer: 1,
		})

		debugger("locking "+id, 1)
		l, _ := syncer.Get(id)
		l.(Syncer).m.Lock()
	}
}

func unlock(id string) {
	if tmp, ok := syncer.Get(id); ok {
		lock := tmp.(Syncer)

		lock.m.Unlock()
		lock.buffer -= 1

		debugger("releasing "+id, lock.buffer)
		if lock.buffer == 1 {
			syncer.Remove(id)
			return
		}

		syncer.Set(id, lock)
	}
}

var DebugFlag bool

func debugger(msg string, buffer int) {
	if !DebugFlag {
		return
	}

	fmt.Println(msg, " buffer:", buffer)
}
