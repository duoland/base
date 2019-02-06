package crypto

import (
	"encoding/base64"
	"testing"
)

var (
	secret256  = "aaaabbbbccccdddd"
	samples256 = [][]string{
		[]string{"hello", "chjcslJQ8vbrwAC7dB4f3A=="},
		[]string{"world", "GLUAwYNhSPIol7r4xXxZXw=="},
		[]string{"golang programming", "7yyWtwJkDHF1vHHn7OStOnl5WDdM/asmsUvch+6OroE="},
	}
)

func TestEncryptUsingAES256(t *testing.T) {
	for _, sample := range samples256 {
		origData := []byte(sample[0])
		cryptedData, err := AESEncrypt(origData, []byte(secret256))
		if err != nil {
			t.Fatal(err)
			return
		}

		hexedCryptedData := base64.StdEncoding.EncodeToString(cryptedData)
		if hexedCryptedData != sample[1] {
			t.Fatal("invalid aes256 algorithm in encryption")
		}
	}
}

func TestDecryptUsingAES256(t *testing.T) {
	for _, sample := range samples256 {
		cryptedData, _ := base64.StdEncoding.DecodeString(sample[1])
		origData, err := AESDecrypt(cryptedData, []byte(secret256))
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(origData) != sample[0] {
			t.Fatal("invalid aes256 algorithm in decryption")
		}
	}
}
