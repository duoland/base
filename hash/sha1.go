package hash

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// HmacSha1 hash the data using algorithm hmac-sha1
func HmacSha1(data []byte, secret []byte) (output []byte) {
	h := hmac.New(sha1.New, secret)
	h.Write(data)
	output = h.Sum(nil)
	return
}

// Sha1HexString creates the hash of the data and return it in hex string format
func Sha1HexString(data []byte) (output string) {
	h := sha1.New()
	h.Write(data)
	output = hex.EncodeToString(h.Sum(nil))
	return
}
