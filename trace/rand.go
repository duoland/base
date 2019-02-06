package trace

import (
	"math/rand"
	"time"
)

var randSeed = "abcdefghijklmnopqrstuvwxyz0123456789"

// GenRandBytes creates a fixed number of random bytes
func GenRandBytes(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		// create a random string
		seedLen := len(randSeed)
		for i := 0; i < size; i++ {
			buf[i] = randSeed[rand.Intn(seedLen)]
		}
		return buf
	}

	return buf
}
