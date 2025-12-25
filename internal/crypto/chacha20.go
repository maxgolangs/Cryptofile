package crypto

import (
	"crypto/cipher"
	"golang.org/x/crypto/chacha20poly1305"
)

func NewAEAD(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.New(key)
}

func EncryptData(aead cipher.AEAD, nonce, plaintext, aad []byte) []byte {
	return aead.Seal(nil, nonce, plaintext, aad)
}

func DecryptData(aead cipher.AEAD, nonce, ciphertext, aad []byte) ([]byte, error) {
	return aead.Open(nil, nonce, ciphertext, aad)
}
