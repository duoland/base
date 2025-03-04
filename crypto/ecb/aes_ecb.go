package ecb

import (
	"crypto/aes"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

// AESEncrypt - AES encryption
func AESEncrypt(pt, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	// pad last block of plaintext if block size less than block cipher size
	pt, err = padder.Pad(pt)
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

// AESDecrypt - AES decryption
func AESDecrypt(ct, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	// unpad plaintext after decryption
	pt, err = padder.Unpad(pt)
	if err != nil {
		panic(err.Error())
	}
	return pt
}
