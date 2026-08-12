package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	apperror "github.com/L1mus/backend-ewallet/internal/AppError"
	"github.com/L1mus/backend-ewallet/internal/domain/auth"
	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
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
	const version int8 = 19

	encodedKey := base64.RawStdEncoding.EncodeToString(key)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, a.MemorySize, a.Iterations, a.Parallelism, encodedSalt, encodedKey)
}

func (a *argon2idHasher) Compare(hash, password string) error {
	splitHashPassword := strings.Split(hash, "$")
	if len(splitHashPassword) < 6 {
		return apperror.NewFromError(http.StatusInternalServerError, "INVALID_FORMAT", "internal server error", errors.New("array length after splitting is less than 5"))
	}
	if strings.ToLower(splitHashPassword[1]) != "argon2id" {
		return apperror.NewFromError(http.StatusInternalServerError, "INVALID_FORMAT", "internal server error", errors.New("the algorithm used for password hashing is not argon2id"))
	}

	var version int8
	if scanned, err := fmt.Sscanf(splitHashPassword[2], "v=%d", &version); err != nil || scanned != 1 {
		return apperror.NewFromError(http.StatusInternalServerError, "FAILED_PROCESS", "internal server error", errors.New("failed scan"))
	}
	if version != 19 {
		return apperror.NewFromError(http.StatusInternalServerError, "INVALID_FORMAT", "internal server error", errors.New("wrong version"))
	}

	var memorySize, iteration uint32
	var parallelism uint8
	if scanned, err := fmt.Sscanf(splitHashPassword[3], "m=%d,t=%d,p=%d", &memorySize, &iteration, &parallelism); err != nil || scanned != 3 {
		return apperror.NewFromError(http.StatusInternalServerError, "FAILED_PROCESS", "internal server error", errors.New("failed scan"))
	}

	if memorySize != a.MemorySize || iteration != a.Iterations || parallelism != a.Parallelism {
		return apperror.NewFromError(http.StatusInternalServerError, "INVALID_FORMAT", "internal server error", errors.New("wrong Configuration"))
	}

	saltDecode, err := base64.RawStdEncoding.DecodeString(splitHashPassword[4])
	if err != nil {
		return apperror.NewFromError(http.StatusInternalServerError, "FAILED_PROCESS", "internal server error", errors.New("invalid salt encoding"))
	}

	keyDecode, err := base64.RawStdEncoding.DecodeString(splitHashPassword[5])
	if err != nil {
		return apperror.NewFromError(http.StatusInternalServerError, "FAILED_PROCESS", "internal server error", errors.New("invalid key encoding"))
	}

	if len(saltDecode) == 0 || len(keyDecode) == 0 {
		return apperror.NewFromError(http.StatusInternalServerError, "INVALID_FORMAT", "internal server error", errors.New("salt or key empty"))
	}

	KeyWantToCompare := argon2.IDKey([]byte(password), saltDecode, iteration, memorySize, parallelism, uint32(len(keyDecode)))

	if subtle.ConstantTimeCompare(keyDecode, KeyWantToCompare) == 0 {
		return domainUser.ErrInvalidCredential
	}

	return nil
}

func (a *argon2idHasher) genSalt() []byte {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	return salt
}
