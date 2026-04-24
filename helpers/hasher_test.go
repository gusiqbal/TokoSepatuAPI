package helpers

import (
	"strings"
	"testing"
)

func TestHashPassword_ReturnsHash(t *testing.T) {
	hash, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if hash == "mypassword" {
		t.Error("hash must differ from input")
	}
	if !strings.HasPrefix(hash, "$2a$") {
		t.Errorf("expected bcrypt hash prefix, got %q", hash[:4])
	}
}

func TestHashPassword_DifferentCallsDifferentHashes(t *testing.T) {
	h1, _ := HashPassword("samepassword")
	h2, _ := HashPassword("samepassword")
	if h1 == h2 {
		t.Error("expected different hashes for same password (bcrypt uses random salt)")
	}
}

func TestHashPassword_EmptyString(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("unexpected error hashing empty string: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash for empty password")
	}
}

func TestVerifyPassword_Correct(t *testing.T) {
	hash, _ := HashPassword("correctpassword")
	err := VerifyPassword(hash, "correctpassword")
	if err != nil {
		t.Errorf("expected no error for correct password, got %v", err)
	}
}

func TestVerifyPassword_Wrong(t *testing.T) {
	hash, _ := HashPassword("correctpassword")
	err := VerifyPassword(hash, "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestVerifyPassword_EmptyInput(t *testing.T) {
	hash, _ := HashPassword("somepassword")
	err := VerifyPassword(hash, "")
	if err == nil {
		t.Error("expected error when verifying empty password against non-empty hash")
	}
}
