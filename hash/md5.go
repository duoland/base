package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5HexString returns md5 hash of data in hex format
func Md5HexString(from []byte) string {
	md5Hasher := md5.New()
	md5Hasher.Write(from)
	return hex.EncodeToString(md5Hasher.Sum(nil))
}
