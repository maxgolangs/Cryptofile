package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"cryptor/constants"

	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string, salt []byte, params constants.Argon2Params) []byte {
	base64Password := base64.StdEncoding.EncodeToString([]byte(password))
	hash := sha256.Sum256([]byte(base64Password))
	return argon2.IDKey(hash[:], salt, params.TimeCost, params.MemoryCost, params.Parallelism, params.HashLen)
}

