package utils

import (
	"os"
	"testing"
)

var secretStoragePath = "test_secret_store.txt" // Use a temporary test file
var encryptionKey = []byte("abcdefghijklmnop")  // Set a valid 16-byte AES key

func cleanupTestStoragePath(secretStoragePath string) {
	os.Remove(secretStoragePath) // Cleanup after test
}

func TestStoreEncryptedClientSecret(t *testing.T) {
	defer cleanupTestStoragePath(secretStoragePath)

	tests := []struct {
		name      string
		apiKey    string
		wantError bool
	}{
		{"Valid API key", "my-secret-key", false},
		{"Empty API key", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StoreEncryptedClientSecret(tt.apiKey, secretStoragePath, encryptionKey)

			if (err != nil) != tt.wantError {
				t.Errorf("StoreEncryptedClientSecret() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if _, err := os.Stat(secretStoragePath); os.IsNotExist(err) {
					t.Errorf("Expected file to be created but it was not found")
				}
			}
		})
	}
}

func TestGetDecryptedClientSecret(t *testing.T) {
	defer cleanupTestStoragePath(secretStoragePath)

	// Store a secret first
	expectedSecret := "my-secret-key"
	err := StoreEncryptedClientSecret(expectedSecret, secretStoragePath, encryptionKey)
	if err != nil {
		t.Fatalf("Failed to store encrypted secret: %v", err)
	}

	tests := []struct {
		name         string
		setup        func()
		wantError    bool
		expectedText string
	}{
		{"Valid decryption", func() {}, false, expectedSecret},
		{"No file exists", func() { os.Remove(secretStoragePath) }, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			decrypted, err := GetDecryptedClientSecret(secretStoragePath, encryptionKey)

			if (err != nil) != tt.wantError {
				t.Errorf("GetDecryptedClientSecret() error = %v, wantError %v", err, tt.wantError)
			}

			if decrypted != tt.expectedText {
				t.Errorf("Expected %s, got %s", tt.expectedText, decrypted)
			}
		})
	}
}

func TestDeleteClientSecret(t *testing.T) {
	defer cleanupTestStoragePath(secretStoragePath)

	// Store a secret first
	_ = StoreEncryptedClientSecret("my-secret-key", secretStoragePath, encryptionKey)

	tests := []struct {
		name      string
		setup     func()
		wantError bool
	}{
		{"Delete existing file", func() {}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := DeleteClientSecret(secretStoragePath)

			if (err != nil) != tt.wantError {
				t.Errorf("DeleteClientSecret() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if file still exists
			if _, err := os.Stat(secretStoragePath); !os.IsNotExist(err) {
				t.Errorf("Expected file to be deleted, but it still exists")
			}
		})
	}
}
