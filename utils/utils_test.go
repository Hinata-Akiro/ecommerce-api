package utils

import (
	"testing"
)

// TestHashPassword verifies that hashing a password produces a non-empty string and no errors.
func TestHashPassword(t *testing.T) {
	plaintext := "mySecurePassword123"

	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash, got empty string")
	}
}

// TestComparePasswords verifies that comparing the correct plaintext password with a hash returns no error.
func TestComparePasswords(t *testing.T) {
	plaintext := "mySecurePassword123"

	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = ComparePasswords(plaintext, hash)
	if err != nil {
		t.Errorf("expected passwords to match, got error: %v", err)
	}

	wrongPassword := "incorrectPassword123"
	err = ComparePasswords(wrongPassword, hash)
	if err == nil {
		t.Error("expected an error for incorrect password comparison, got nil")
	}
}
