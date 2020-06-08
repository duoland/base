package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSha256 hash the data using algorithm hmac-sha1
func HmacSha256(data []byte, secret []byte) (output []byte) {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	output = h.Sum(nil)
	return
}

// Sha256HexString creates the hash of the data and return it in hex string format
func Sha256HexString(data []byte) (output string) {
	h := sha256.New()
	h.Write(data)
	output = hex.EncodeToString(h.Sum(nil))
	return
}
