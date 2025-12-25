package crypto

import (
	"cryptor/constants"
	"cryptor/internal/obfuscation"
	"encoding/binary"
	"errors"
)

func BuildHeader(flags uint8, params constants.Argon2Params, salt, nonce, nameBytes []byte, originalSize uint64) []byte {
	buf := make([]byte, constants.HeaderSize)
	copy(buf[0:8], []byte(obfuscation.GetMagic()))
	buf[8] = obfuscation.GetVersion()
	buf[9] = flags
	binary.BigEndian.PutUint16(buf[10:12], uint16(params.TimeCost))
	binary.BigEndian.PutUint32(buf[12:16], params.MemoryCost)
	buf[16] = params.Parallelism
	buf[17] = uint8(len(salt))
	buf[18] = uint8(len(nonce))
	binary.BigEndian.PutUint16(buf[19:21], uint16(len(nameBytes)))
	binary.BigEndian.PutUint64(buf[21:29], originalSize)
	return buf
}

func ParseHeader(header []byte) (version uint8, flags uint8, params constants.Argon2Params, saltLen, nonceLen, nameLen uint16, originalSize uint64, err error) {
	if len(header) != constants.HeaderSize {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректный заголовок файла")
	}

	magicStr := obfuscation.GetMagic()
	if string(header[0:8]) != magicStr {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("файл не принадлежит формату CryptoFile или создан другой версией программы")
	}

	version = header[8]
	if version != obfuscation.GetVersion() {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("неверная версия формата")
	}
	flags = header[9]
	params.TimeCost = uint32(binary.BigEndian.Uint16(header[10:12]))
	params.MemoryCost = binary.BigEndian.Uint32(header[12:16])
	params.Parallelism = header[16]
	params.HashLen = 32
	saltLen = uint16(header[17])
	nonceLen = uint16(header[18])
	nameLen = binary.BigEndian.Uint16(header[19:21])
	originalSize = binary.BigEndian.Uint64(header[21:29])

	if params.TimeCost == 0 || params.TimeCost > 100 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}
	if params.MemoryCost == 0 || params.MemoryCost > 2048*1024*1024 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}
	if params.Parallelism == 0 || params.Parallelism > 16 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}
	if saltLen == 0 || saltLen > 64 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}
	if nonceLen != 12 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}
	if nameLen > 4096 {
		return 0, 0, constants.Argon2Params{}, 0, 0, 0, 0, errors.New("некорректные параметры шифрования")
	}

	return version, flags, params, saltLen, nonceLen, nameLen, originalSize, nil
}

