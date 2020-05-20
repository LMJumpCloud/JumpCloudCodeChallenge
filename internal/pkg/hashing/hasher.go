package hashing

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// GetHash will generate the SHA512 hash and return a base64 encoded string of the hash
func GetHash(str string) string {
	hashBytes := sha512.Sum512([]byte(str))
	asString := fmt.Sprintf("%x", hashBytes)
	return base64.StdEncoding.EncodeToString([]byte(asString))
}