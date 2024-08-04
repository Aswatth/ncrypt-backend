package encryptor

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := "some text"
	plain_text := "this is a secret"

	encrypted_text, _ := Encrypt(plain_text, key)

	if encrypted_text == "" {
		t.Errorf("Encryption failed")
	}
}

func TestDecrypt(t *testing.T) {
	key := "some text"
	plain_text := "this is a secret"

	encrypted_text, _ := Encrypt(plain_text, key)
	decrypted_text, _ := Decrypt(encrypted_text, key)

	if plain_text != decrypted_text {
		t.Errorf("Decryption failed")
	}
}
