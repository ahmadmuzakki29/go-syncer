package syncer

import (
	"net/http"
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

func HandlerLock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lock(id)
	w.WriteHeader(http.StatusOK)
}

func HandlerUnlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	unlock(id)
	w.WriteHeader(http.StatusOK)
}
