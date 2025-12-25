package constants

const (
	Magic              = "CRYPTOR2"
	Version            = 2
	FlagDirectory      = 0x01
	DefaultSaltSize    = 32
	DefaultNonceSize   = 12
	MinPasswordLength  = 6
	HeaderSize         = 8 + 1 + 1 + 2 + 4 + 1 + 1 + 1 + 2 + 8
)

type Argon2Params struct {
	TimeCost    uint32
	MemoryCost  uint32
	Parallelism uint8
	HashLen     uint32
}

var DefaultArgon2Params = Argon2Params{
	TimeCost:    4,
	MemoryCost:  128 * 1024,
	Parallelism: 4,
	HashLen:     32,
}


