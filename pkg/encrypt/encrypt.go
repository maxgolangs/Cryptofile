package encrypt

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cryptor/constants"
	"cryptor/internal/crypto"
	"cryptor/internal/file"
	"cryptor/internal/security"
)

func EncryptPath(targetPath, password string, removeOriginal bool) (string, error) {
	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return "", err
	}

	params := constants.DefaultArgon2Params
	salt := make([]byte, constants.DefaultSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	nonce := make([]byte, constants.DefaultNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	nameBytes := []byte(filepath.Base(targetPath))
	flags := uint8(0)
	originalSize := uint64(0)

	if info.IsDir() {
		flags = constants.FlagDirectory
	} else {
		originalSize = uint64(info.Size())
	}

	var plaintext []byte
	if info.IsDir() {
		plaintext, err = file.ArchiveDirectory(targetPath)
		if err != nil {
			return "", err
		}
	} else {
		plaintext, err = os.ReadFile(targetPath)
		if err != nil {
			return "", err
		}
	}

	header := crypto.BuildHeader(flags, params, salt, nonce, nameBytes, originalSize)
	aad := append(append(append(append(header, salt...), nonce...), nameBytes...), []byte{}...)

	key := crypto.DeriveKey(password, salt, params)
	aead, err := crypto.NewAEAD(key)
	if err != nil {
		security.SecureZero(key)
		return "", err
	}

	ciphertext := crypto.EncryptData(aead, nonce, plaintext, aad)
	security.SecureZero(key)

	var outputPath string
	if info.IsDir() {
		outputPath = filepath.Join(filepath.Dir(targetPath), filepath.Base(targetPath)+".encrypted")
	} else {
		suffix := filepath.Ext(targetPath)
		if suffix == "" {
			outputPath = targetPath + ".encrypted"
		} else {
			outputPath = strings.TrimSuffix(targetPath, suffix) + suffix + ".encrypted"
		}
	}

	outputPath = file.EnsureUniquePath(outputPath, false)

	dataToWrite := make([]byte, 0, len(header)+len(salt)+len(nonce)+len(nameBytes)+len(ciphertext))
	dataToWrite = append(dataToWrite, header...)
	dataToWrite = append(dataToWrite, salt...)
	dataToWrite = append(dataToWrite, nonce...)
	dataToWrite = append(dataToWrite, nameBytes...)
	dataToWrite = append(dataToWrite, ciphertext...)

	if err := file.WriteFileBuffered(outputPath, dataToWrite, 0644); err != nil {
		return "", err
	}

	if removeOriginal {
		if info.IsDir() {
			if err := os.RemoveAll(targetPath); err != nil {
				return outputPath, fmt.Errorf("не удалось удалить директорию: %w", err)
			}
		} else {
			if err := file.SecureDelete(targetPath, 1); err != nil {
				return outputPath, fmt.Errorf("не удалось безопасно удалить файл: %w", err)
			}
		}
	}

	return outputPath, nil
}

