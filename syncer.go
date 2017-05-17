package syncer

import (
	"sync"
)

var syncer = make(map[string]Syncer)

type Syncer struct {
	m      *sync.Mutex
	buffer int
}

func lock(id string) {
	if lock, ok := syncer[id]; ok {
		lock.buffer += 1
		syncer[id] = lock
		lock.m.Lock()
	} else {
		syncer[id] = Syncer{
			m:      &sync.Mutex{},
			buffer: 1,
		}
		syncer[id].m.Lock()
	}
}

func unlock(id string) {
	if lock, ok := syncer[id]; ok {
		if lock.buffer == 1 {
			delete(syncer, id)
		}

		lock.m.Unlock()
		lock.buffer -= 1
		syncer[id] = lock
	}
}
