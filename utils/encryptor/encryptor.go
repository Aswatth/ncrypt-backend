package encryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

// Create a hash using SHA-256
func CreateHash(data_to_hash string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data_to_hash))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrpyt plain text using SHA-256 hash as key
func Encrypt(plaintext string, keyStr string) (string, error) {
	key, err := hex.DecodeString(CreateHash(keyStr))

	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad plaintext to be multiple of block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	paddedPlaintext := append([]byte(plaintext), padtext...)

	// Encrypt the plaintext
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Combine IV and ciphertext for output
	combined := append(iv, ciphertext...)
	return hex.EncodeToString(combined), nil
}

// Decrpyt encrypted text using hsa-256 hash
func Decrypt(ciphertextHex string, keyStr string) (string, error) {
	key, err := hex.DecodeString(CreateHash(keyStr))

	if err != nil {
		return "", err
	}

	// Decode the combined IV and ciphertext from hex
	combined, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	if len(combined) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := combined[:aes.BlockSize]
	ciphertext := combined[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext
	mode := cipher.NewCBCDecrypter(block, iv)
	paddedPlaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(paddedPlaintext, ciphertext)

	// Remove padding
	padding := int(paddedPlaintext[len(paddedPlaintext)-1])
	plaintext := paddedPlaintext[:len(paddedPlaintext)-padding]
	return string(plaintext), nil
}
