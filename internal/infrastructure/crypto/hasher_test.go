package crypto

import (
	"fmt"
	"testing"

	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

func TestArgon2idHasher_Generate(t *testing.T) {
	password := "Password123"

	argon2id := NewArgon2idHasher(64*1024, 1, 32, 16, 1)

	passwordHash := argon2id.Generate(password)

	fmt.Println(passwordHash)
}

func TestCompare_Roundtrip(t *testing.T) {
	password := "Password123!"
	h := NewArgon2idHasher(64*1024, 1, 32, 16, 1)
	hash := h.Generate(password)
	if err := h.Compare(hash, password); err != nil {
		t.Fatalf("roundtrip failed: %v", err)
	}
}

func TestCompare_WrongPassword(t *testing.T) {
	password := "Password123!"
	h := NewArgon2idHasher(64*1024, 1, 32, 16, 1)
	hash := h.Generate(password)
	if err := h.Compare(hash, "wrong"); err != domainUser.ErrInvalidCredential {
		t.Fatalf("expected ErrInvalidCredential, got %v", err)
	}
}

func TestCompare_ParamMismatch(t *testing.T) {
	password := "Password123!"
	h1 := NewArgon2idHasher(64*1024, 1, 32, 16, 1)
	hash := h1.Generate(password)
	h2 := NewArgon2idHasher(32*1024, 1, 32, 16, 1)
	if err := h2.Compare(hash, password); err == nil {
		t.Fatalf("expected error due to parameter mismatch, got nil")
	}
}

func TestCompare_MalformedInputs(t *testing.T) {
	h := NewArgon2idHasher(64*1024, 1, 32, 16, 1)
	cases := []struct {
		name string
		hash string
	}{
		{"too few parts", "$argon2id$v=19$m=65536,t=1,p=1$saltonly"},
		{"bad algorithm", "$argon2$v=19$m=65536,t=1,p=1$salt$key"},
		{"bad version", "$argon2id$v=abc$m=65536,t=1,p=1$salt$key"},
		{"bad params", "$argon2id$v=19$m=notnum,t=1,p=1$salt$key"},
		{"unsupported version", "$argon2id$v=18$m=65536,t=1,p=1$salt$key"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := h.Compare(c.hash, "pwd"); err == nil {
				t.Fatalf("expected error for case %s, got nil", c.name)
			}
		})
	}
}
