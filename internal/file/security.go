package file

import (
	"crypto/rand"
	"os"

	"golang.org/x/crypto/chacha20poly1305"
)

func SecureDelete(filePath string, passes int) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		return os.RemoveAll(filePath)
	}

	if info.Size() == 0 {
		return os.Remove(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize := info.Size()
	blockSize := int64(1024 * 1024)

	if passes < 1 {
		passes = 1
	}

	for pass := 0; pass < passes; pass++ {
		if _, err := file.Seek(0, 0); err != nil {
			return err
		}

		remaining := fileSize

		for remaining > 0 {
			chunkSize := blockSize
			if chunkSize > remaining {
				chunkSize = remaining
			}

			key := make([]byte, 32)
			nonce := make([]byte, 12)
			if _, err := rand.Read(key); err != nil {
				return err
			}
			if _, err := rand.Read(nonce); err != nil {
				return err
			}

			aead, err := chacha20poly1305.New(key)
			if err != nil {
				return err
			}

			randomData := make([]byte, chunkSize)
			if _, err := rand.Read(randomData); err != nil {
				return err
			}

			ciphertext := aead.Seal(nil, nonce, randomData, nil)

			if len(ciphertext) > int(chunkSize) {
				ciphertext = ciphertext[:chunkSize]
			}

			if _, err := file.Write(ciphertext); err != nil {
				return err
			}

			remaining -= chunkSize
		}

		if err := file.Sync(); err != nil {
			return err
		}
	}

	file.Close()
	return os.Remove(filePath)
}


