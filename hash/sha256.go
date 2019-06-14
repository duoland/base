package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256HexString creates the hash of the data and return it in hex string format
func Sha256HexString(data []byte) (output string) {
	h := sha256.New()
	h.Write(data)
	output = hex.EncodeToString(h.Sum(nil))
	return
}
