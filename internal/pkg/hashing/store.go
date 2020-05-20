package hashing

import (
	"sync"
	"sync/atomic"
	"time"
)

// HashStore stores hashes an their ids in memory
type HashStore struct {
	passwordID int64
	availableHashes map[int64]string
	mapLock sync.Mutex
}

func NewHashStore() *HashStore {
	return &HashStore{
		passwordID: 0,
		availableHashes: make(map[int64]string),
		mapLock: sync.Mutex{},
	}
}

func (h *HashStore) getNextPasswordID() int64 {
	return atomic.AddInt64(&h.passwordID, 1)
}

func (h *HashStore) SubmitPassword(pass string) int64 {
	id := h.getNextPasswordID()
	go func() {
		time.Sleep(5 * time.Second)
		h.mapLock.Lock()
		defer h.mapLock.Unlock()
		h.availableHashes[id] = GetHash(pass)
	}()
	return id
}

func (h *HashStore) GetHash(id int64) string {
	h.mapLock.Lock()
	defer h.mapLock.Unlock()
	return h.availableHashes[id]
}