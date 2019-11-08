package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

func withDumpError(returnVal interface{}, err error) interface{} {
	if err != nil {
		panic("Error : " + err.Error() + "\n")
	}
	return returnVal
}

//if key is not 512 bit, this function will generate a 512 bit string out of the key
func pad256Key(key string) string {
	if len(key) > 64 {
		return key[:32]
	}
	return key + strings.Repeat("0", 64-len(key))
}

func prepare12ByteNonce() []byte {
	buffer := make([]byte, 12)
	io.ReadFull(rand.Reader, buffer)
	return buffer
}

func aesEncrypt(data string, key string) (string, string) {
	key = pad256Key(key)
	hexKey := withDumpError(hex.DecodeString(key))
	plainText := []byte(data)

	aesObject := withDumpError(aes.NewCipher(hexKey.([]byte)))
	aesGCM := withDumpError(cipher.NewGCM(aesObject.(cipher.Block)))
	nonce := prepare12ByteNonce()
	return string(aesGCM.(cipher.AEAD).Seal(nil, nonce, plainText, nil)), string(nonce)
}

func aesDecrypt(data string, nonce string, key string) string {
	key = pad256Key(key)
	hexKey := withDumpError(hex.DecodeString(key))
	cipherText := []byte(data)

	aesObject := withDumpError(aes.NewCipher(hexKey.([]byte)))

	aesGCM := withDumpError(cipher.NewGCM(aesObject.(cipher.Block)))

	plainText := withDumpError(aesGCM.(cipher.AEAD).Open(nil, []byte(nonce), cipherText, nil)).([]byte)
	return string(plainText)
}
