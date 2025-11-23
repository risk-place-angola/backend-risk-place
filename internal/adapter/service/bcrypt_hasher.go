package service

import (
	"errors"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type hash struct {
}

func NewBcryptHasher() port.PasswordHasher {
	return &hash{}
}

func (h *hash) Hash(raw string) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *hash) Compare(hash, raw string) bool {
	if len(hash) == 0 || len(raw) == 0 {
		slog.Error("Hash or raw password is empty", "hash", hash, "raw", raw)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			slog.Error("Password comparison failed", "error", err)
			return false
		}
		slog.Error("Unexpected error during password comparison", "error", err)
		return false
	}
	return true
}
