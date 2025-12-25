package decrypt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"cryptor/constants"
	"cryptor/internal/crypto"
	"cryptor/internal/file"
	"cryptor/internal/obfuscation"
	"cryptor/internal/security"
)

func DecryptFile(encryptedPath, password string, removeOriginal bool) (string, error) {
	encryptedPath, err := filepath.Abs(encryptedPath)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(encryptedPath)
	if err != nil {
		return "", err
	}

	if len(data) < constants.HeaderSize {
		return "", errors.New("файл слишком короткий")
	}

	header := data[:constants.HeaderSize]
	version, flags, params, saltLen, nonceLen, nameLen, _, err := crypto.ParseHeader(header)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга заголовка: %w", err)
	}

	expectedVersion := obfuscation.GetVersion()
	if version != expectedVersion {
		return "", fmt.Errorf("неподдерживаемая версия контейнера: %d (требуется версия %d). Файл был создан другой версией программы", version, expectedVersion)
	}

	if err := security.ValidateArgon2Params(params); err != nil {
		return "", fmt.Errorf("небезопасные параметры шифрования в файле: %w", err)
	}

	offset := constants.HeaderSize
	if offset+int(saltLen) > len(data) {
		return "", errors.New("некорректный размер соли")
	}
	salt := data[offset : offset+int(saltLen)]
	offset += int(saltLen)

	if offset+int(nonceLen) > len(data) {
		return "", errors.New("некорректный размер nonce")
	}
	nonce := data[offset : offset+int(nonceLen)]
	offset += int(nonceLen)

	if offset+int(nameLen) > len(data) {
		return "", errors.New("некорректный размер имени")
	}
	nameBytes := data[offset : offset+int(nameLen)]
	offset += int(nameLen)

	ciphertext := data[offset:]

	aad := append(append(append(append(header, salt...), nonce...), nameBytes...), []byte{}...)

	key := crypto.DeriveKey(password, salt, params)
	aead, err := crypto.NewAEAD(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < 16 {
		return "", errors.New("файл повреждён: недостаточно данных")
	}

	plaintext, err := crypto.DecryptData(aead, nonce, ciphertext, aad)
	if err != nil {
		security.SecureZero(key)
		if len(ciphertext) < 16 {
			return "", fmt.Errorf("файл повреждён: недостаточно данных для расшифровки (требуется минимум 16 байт, получено %d)", len(ciphertext))
		}
		return "", fmt.Errorf("неверный пароль или файл повреждён (ошибка AEAD: %v)", err)
	}

	if len(plaintext) == 0 {
		security.SecureZero(key)
		return "", errors.New("файл повреждён: пустые данные")
	}

	name := string(nameBytes)
	destination := filepath.Join(filepath.Dir(encryptedPath), name)

	if flags&constants.FlagDirectory != 0 {
		destination = file.EnsureUniquePath(destination, true)
		if err := file.ExtractDirectory(plaintext, destination); err != nil {
			return "", err
		}
	} else {
		destination = file.EnsureUniquePath(destination, false)
		if err := os.WriteFile(destination, plaintext, 0644); err != nil {
			return "", err
		}
	}

	if removeOriginal {
		if err := file.SecureDelete(encryptedPath, 1); err != nil {
			return destination, fmt.Errorf("не удалось удалить зашифрованный файл: %w", err)
		}
	}

	return destination, nil
}

