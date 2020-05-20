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
	wg sync.WaitGroup
}

// NewHashStore returns a new HashStore instance
func NewHashStore() *HashStore {
	return &HashStore{
		passwordID: 0,
		availableHashes: make(map[int64]string),
		mapLock: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
}

func (h *HashStore) getNextPasswordID() int64 {
	return atomic.AddInt64(&h.passwordID, 1)
}

// SubmitPassword accepts new passwords to be hashed, returning the ID so that the hash can be
// retrieved after processing has finished
func (h *HashStore) SubmitPassword(pass string) int64 {
	return h.waitAndStoreHash(pass, 5 * time.Second)
}

// ForcePassword accepts new passwords without any processing time, inserting them into the store
// immediately
func (h *HashStore) ForcePassword(pass string) int64 {
	return h.waitAndStoreHash(pass, 0)
}

func (h *HashStore) waitAndStoreHash(pass string, pause time.Duration) int64 {
	id := h.getNextPasswordID()
	h.wg.Add(1)

	go func() {
		time.Sleep(pause)
		h.mapLock.Lock()
		defer func() {
			h.mapLock.Unlock()
			h.wg.Done()
		}()
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

// Flush will block and wait for any processing of in-flight hashing to finish
func (h *HashStore) Flush() {
	h.wg.Wait()
}