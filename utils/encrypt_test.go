package utils

import (
	"testing"
)

func setupTestKey() []byte {
	return []byte("abcdefghijklmnop") // Set a valid 16-byte key for testing
}

func TestEncrypt(t *testing.T) {
	key := setupTestKey()

	tests := []struct {
		name      string
		plaintext string
		wantError bool
	}{
		{"Valid encryption", "Hello, Golang!", false},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.plaintext, key)

			if (err != nil) != tt.wantError {
				t.Errorf("Encrypt() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && encrypted == "" {
				t.Errorf("Expected encrypted text, got empty string")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	key := setupTestKey()

	encrypted, _ := Encrypt("Secret Message", key) // Ensure valid ciphertext

	tests := []struct {
		name        string
		ciphertext  string
		wantError   bool
		expectedMsg string
	}{
		{"Valid decryption", encrypted, false, "Secret Message"},
		{"Invalid hex", "invalidhex", true, ""},
		{"Short ciphertext", "1234", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decrypted, err := Decrypt(tt.ciphertext, key)

			if (err != nil) != tt.wantError {
				t.Errorf("Decrypt() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && decrypted != tt.expectedMsg {
				t.Errorf("Expected %s, got %s", tt.expectedMsg, decrypted)
			}
		})
	}
}
