package security

import (
	"cryptor/constants"
	"encoding/binary"
	"errors"
)

func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

func ValidateFileFormat(data []byte) error {
	if len(data) < constants.HeaderSize {
		return errors.New("файл повреждён или неполный")
	}

	header := data[:constants.HeaderSize]

	magicStr := constants.Magic
	if string(header[0:8]) != magicStr {
		return errors.New("неверный формат файла")
	}

	if header[8] != constants.Version {
		return errors.New("несовместимая версия формата")
	}

	saltLen := uint16(header[17])
	nonceLen := uint16(header[18])
	nameLen := binary.BigEndian.Uint16(header[19:21])

	if saltLen == 0 || saltLen > 64 {
		return errors.New("некорректный формат файла")
	}
	if nonceLen != 12 {
		return errors.New("некорректный формат файла")
	}
	if nameLen > 4096 {
		return errors.New("некорректный формат файла")
	}

	minSize := constants.HeaderSize + int(saltLen) + int(nonceLen) + int(nameLen) + 16
	if len(data) < minSize {
		return errors.New("файл повреждён или обрезан")
	}

	return nil
}

func ValidateArgon2Params(params constants.Argon2Params) error {
	if params.TimeCost < 3 {
		return errors.New("небезопасные параметры шифрования")
	}
	if params.MemoryCost < 64*1024 {
		return errors.New("небезопасные параметры шифрования")
	}
	if params.Parallelism == 0 || params.Parallelism > 16 {
		return errors.New("небезопасные параметры шифрования")
	}
	if params.HashLen != 32 {
		return errors.New("небезопасные параметры шифрования")
	}
	return nil
}

func VerifyProgramIntegrity() bool {
	return true
}


