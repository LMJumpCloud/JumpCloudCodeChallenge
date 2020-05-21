package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"testing"
	"time"
)

func TestHashStore(t *testing.T) {
	t.Run("hash ID increment", func(t *testing.T) {
		store := hashing.NewInMemoryHashStore()
		test.AssertEqual(t, store.SubmitPassword("first"), hashing.SubmitResponse{ID: 1}, "first hash has id 1")
		test.AssertEqual(t, store.SubmitPassword("second"), hashing.SubmitResponse{ID: 2}, "second hash has id 2")
		test.AssertEqual(t, store.SubmitPassword("third"), hashing.SubmitResponse{ID: 3}, "third hash has id 3")
	})

	t.Run("store returns empty for missing ID", func(t *testing.T) {
		store := hashing.NewInMemoryHashStore()
		test.AssertEqual(t, store.GetHash(2), hashing.GetResponse{
			ID:   2,
			Hash: "",
		}, "empty hash for bad ID")
	})

	t.Run("store returns hash for ID", func(t *testing.T) {
		store := hashing.NewInMemoryHashStore()
		test.AssertEqual(t, store.ForcePassword(input), int64(1), "first hash has id 1")

		// as the inner implementation still uses a goroutine, wait just a moment to let the hash be submitted
		time.Sleep(100 * time.Millisecond)

		test.AssertEqual(t, store.GetHash(1), hashing.GetResponse{
			ID:   1,
			Hash: knownSHA512HashBase64,
		}, "matching hash")
	})
}