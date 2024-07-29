package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func CreateHash(data_to_hash string) string {
	hasher := md5.New()
	hasher.Write([]byte(data_to_hash))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data_to_encrypt string, key string) string {
	block, _ := aes.NewCipher([]byte(CreateHash(key)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data_to_encrypt), nil)
	return string(ciphertext[:])
}

func Decrypt(encrypted_data string, key string) string {
	block, err := aes.NewCipher([]byte(CreateHash(key)))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := []byte(encrypted_data[:nonceSize]), []byte(encrypted_data[nonceSize:])
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext[:])
}
