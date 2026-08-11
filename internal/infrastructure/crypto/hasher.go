package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/L1mus/backend-ewallet/internal/domain/auth"
	"golang.org/x/crypto/argon2"
)

const (
	Byte = 1 << (10 * iota) // 1 Byte
	KiB                     // 1024 Bytes (Kibibyte)
	MiB                     // 1048576 Bytes (Mebibyte)
)

type argon2idHasher struct {
	MemorySize  uint32
	Iterations  uint32
	Parallelism uint8
	KeyLen      uint32
	SaltLen     uint32
}

func NewArgon2idHasher(memory, iteration, keyLen, saltLen uint32, parallelism uint8) auth.HasherPassword {
	return &argon2idHasher{
		MemorySize:  memory,
		Iterations:  iteration,
		Parallelism: parallelism,
		KeyLen:      keyLen,
		SaltLen:     saltLen,
	}
}

// base on OWASP recommendation
func (a *argon2idHasher) useRecommendation() {
	a.MemorySize = 65536 * KiB
	a.Iterations = 1
	a.Parallelism = 1
	a.KeyLen = 32
	a.SaltLen = 16
}

func (a *argon2idHasher) Generate(password string) string {
	salt := a.genSalt()

	key := argon2.IDKey([]byte(password), salt, a.Iterations, a.MemorySize, a.Parallelism, a.KeyLen)

	encodedKey := base64.RawStdEncoding.EncodeToString(key)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.MemorySize, a.Iterations, a.Parallelism, encodedSalt, encodedKey)
}

func (a *argon2idHasher) Compare(hash, password string) error {
	return nil
}

func (a *argon2idHasher) genSalt() []byte {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	return salt
}
