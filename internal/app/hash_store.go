package app

import "sync/atomic"

type HashStore struct {
	passwordID int64
	availableHashes map[int64]string
}

func NewHashStore() *HashStore {
	return &HashStore{
		passwordID: 0,
		availableHashes: make(map[int64]string),
	}
}

func (h *HashStore) GetNextPasswordID() int64 {
	return atomic.AddInt64(&h.passwordID, 1)
}

func (h *HashStore) StoreHash(id int64, hash string) {
	h.availableHashes[id] = hash
}

func (h *HashStore) GetHash(id int64) string {
	return h.availableHashes[id]
}