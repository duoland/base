package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AESEncrypt - AES encryption
func AESEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cryptedData := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptedData, origData)
	return cryptedData, nil
}

// AESDecrypt - AES decryption
func AESDecrypt(cryptedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryptedData))

	blockMode.CryptBlocks(origData, cryptedData)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding - padding algorithm
func PKCS5Padding(cipherData []byte, blockSize int) []byte {
	padding := blockSize - len(cipherData)%blockSize
	paddingData := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherData, paddingData...)
}

// PKCS5UnPadding - unpadding algorithm
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
