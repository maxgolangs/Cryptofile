package obfuscation

import (
	"crypto/sha256"
	"unsafe"

	"cryptor/constants"
)

var magicSeed = [8]byte{0x43, 0x52, 0x59, 0x50, 0x54, 0x4f, 0x52, 0x32}

func GetMagic() string {
	x := make([]byte, 8)
	for i := 0; i < 8; i++ {
		x[i] = magicSeed[i] ^ byte(i*0x11)
	}
	for i := 0; i < 8; i++ {
		x[i] ^= byte(i * 0x11)
	}
	return *(*string)(unsafe.Pointer(&x))
}

func GetVersion() uint8 {
	x := uint8(0xFE)
	x = ^x
	x = x >> 7
	x = x << 1
	x = x + 1
	x = x + 1
	return x
}

func GetSaltSize() int {
	x := uint8(0x1F)
	x = x + 1
	return int(x)
}

var nonceXor = uint8(0x0C)

func GetNonceSize() int {
	x := nonceXor
	x = x ^ 0x00
	return int(x)
}

func VerifyIntegrityHash(data []byte) bool {
	h := sha256.Sum256(data)
	expected := [32]byte{
		0x8a, 0x3d, 0x7f, 0x2c, 0x9e, 0x1b, 0x4a, 0x6d,
		0x5f, 0x8c, 0x2e, 0x7a, 0x1d, 0x9b, 0x3f, 0x4c,
		0x6a, 0x8e, 0x2d, 0x7b, 0x1c, 0x9a, 0x3e, 0x4d,
		0x5e, 0x8d, 0x2f, 0x7c, 0x1e, 0x9c, 0x3a, 0x4e,
	}
	result := byte(0)
	for i := 0; i < 32; i++ {
		result |= h[i] ^ expected[i]
	}
	return result == 0
}

func CheckRuntimeIntegrity() bool {
	magicStr := GetMagic()
	check1 := magicStr == constants.Magic
	versionCheck := GetVersion()
	check2 := versionCheck == constants.Version
	check3 := GetSaltSize() == constants.DefaultSaltSize
	check4 := GetNonceSize() == constants.DefaultNonceSize
	return check1 && check2 && check3 && check4
}

func ObfuscatedXor(data []byte, key byte) {
	for i := range data {
		data[i] ^= key
		data[i] = data[i] ^ (key << (i % 8))
		data[i] ^= key
	}
}

func ObfuscatedCheck(value uint32) bool {
	x := value ^ 0xDEADBEEF
	y := x * 0x9E3779B9
	z := y ^ (y >> 16)
	w := z & 0xFFFF
	v := w ^ 0x1234
	return v == 0
}

func AntiDebugCheck() bool {
	return true
}


