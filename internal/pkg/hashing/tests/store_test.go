package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"testing"
)

const input = "password"

// raw hash: "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86"
const knownSHA512HashBase64 = "YjEwOWYzYmJiYzI0NGViODI0NDE5MTdlZDA2ZDYxOGI5MDA4ZGQwOWIzYmVmZDFiNWUwNzM5NGM3MDZhOGJiOTgwYjFkNzc4NWU1OTc2ZWMwNDliNDZkZjVmMTMyNmFmNWEyZWE2ZDEwM2ZkMDdjOTUzODVmZmFiMGNhY2JjODY="

func TestHasher(t *testing.T) {
	t.Run("generate hash", func(t *testing.T) {
		test.AssertEqual(t, hashing.GetHash(input), knownSHA512HashBase64, "correctly generates SHA512 hash")
	})
}