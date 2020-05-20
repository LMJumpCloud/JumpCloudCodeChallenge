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

// NewHashStore returns a new HashStore instance
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

// SubmitPassword accepts new passwords to be hashed, returning the ID so that the hash can be
// retrieved after processing has finished
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

// GetHash returns the given has for the provided ID, if one exists
func (h *HashStore) GetHash(id int64) string {
	h.mapLock.Lock()
	defer h.mapLock.Unlock()
	return h.availableHashes[id]
}