package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func CreateHash(data_to_hash string) string {
	hasher := md5.New()
	hasher.Write([]byte(data_to_hash))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data_to_encrypt string, key string) string {
	if len(key) != 32 {
		key = CreateHash(key)
	}
	block, _ := aes.NewCipher([]byte(key))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data_to_encrypt), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func Decrypt(encrypted_data string, key string) string {
	encrypted_data_bytes, _ := base64.StdEncoding.DecodeString(encrypted_data)
	if len(key) != 32 {
		key = CreateHash(key)
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encrypted_data_bytes[:nonceSize], encrypted_data_bytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext[:])
}
