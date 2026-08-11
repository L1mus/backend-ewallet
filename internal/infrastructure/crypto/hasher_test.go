package crypto

import (
	"fmt"
	"testing"
)

func TestArgon2idHasher_Generate(t *testing.T) {
	password := "Password123"

	argon2id := NewArgon2idHasher(64*1024, 1, 32, 16, 1)

	passwordHash := argon2id.Generate(password)

	fmt.Println(passwordHash)
}
