package hashing

import (
	"sync"
	"sync/atomic"
	"time"
)

// HashStorer is a basic interface for interacting with hash storage
type HashStorer interface {
	SubmitPassword(pass string) SubmitResponse
	GetHash(id int64) GetResponse
}

// SubmitResponse is simple response from submitting a password for hashing
type SubmitResponse struct {
	ID int64 `json:id`
}

// GetResponse is a simple response from getting a hash
type GetResponse struct {
	ID int64   `json:id`
	Hash string `json:hash`
}

// InMemoryHashStore stores hashes an their ids in memory
type InMemoryHashStore struct {
	passwordID int64
	availableHashes map[int64]string
	mapLock sync.Mutex
	wg sync.WaitGroup
}

// NewInMemoryHashStore returns a new InMemoryHashStore instance
func NewInMemoryHashStore() *InMemoryHashStore {
	return &InMemoryHashStore{
		passwordID: 0,
		availableHashes: make(map[int64]string),
		mapLock: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
}

func (h *InMemoryHashStore) getNextPasswordID() int64 {
	return atomic.AddInt64(&h.passwordID, 1)
}

// SubmitPassword accepts new passwords to be hashed, returning the ID so that the hash can be
// retrieved after processing has finished
func (h *InMemoryHashStore) SubmitPassword(pass string) SubmitResponse {
	return SubmitResponse{
		ID: h.waitAndStoreHash(pass, 5 * time.Second),
	}
}

// ForcePassword accepts new passwords without any processing time, inserting them into the store
// immediately
func (h *InMemoryHashStore) ForcePassword(pass string) int64 {
	return h.waitAndStoreHash(pass, 0)
}

func (h *InMemoryHashStore) waitAndStoreHash(pass string, pause time.Duration) int64 {
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
func (h *InMemoryHashStore) GetHash(id int64) GetResponse {
	h.mapLock.Lock()
	defer h.mapLock.Unlock()
	return GetResponse{
		ID: id,
		Hash: h.availableHashes[id],
	}
}

// Flush will block and wait for any processing of in-flight hashing to finish
func (h *InMemoryHashStore) Flush() {
	h.wg.Wait()
}